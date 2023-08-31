package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	db "github.com/A1sal/AvitoTech/database"
	sg "github.com/A1sal/AvitoTech/segment"
	usseg "github.com/A1sal/AvitoTech/usersegment"
)

type createSegmentBody struct {
	AudienceCvg int `json:"audience_cvg"`
}

type userSegmentsModifyBody struct {
	UserId           int      `json:"user_id"`
	SegmentsToAdd    []string `json:"segments_add"`
	SegmentsToRemove []string `json:"segments_remove"`
}

type userSegmentsResponseBody struct {
	Segments []string `json:"segments"`
}

type userModifyErrorResponse struct {
	Message          string   `json:"message"`
	SegmentsToRemove []string `json:"segments_remove,omitempty"`
	SegmentsToAdd    []string `json:"segments_add,omitempty"`
}

func (u userSegmentsModifyBody) String() string {
	return fmt.Sprintf("\n========\nUSER %dto add: %v\nto remove%v\n========\n", u.UserId, u.SegmentsToAdd, u.SegmentsToRemove)
}

func helloRootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func createSegmentHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody createSegmentBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		writeResponse(w, []byte("incorrect body format"), 400)
		return
	}
	if reqBody.AudienceCvg < 0 || reqBody.AudienceCvg > 100 {
		writeResponse(w, []byte("audience_cvg must be in [0, 100] integers"), 400)
		return
	}
	segmentName := chi.URLParam(r, "segmentName")
	segment := sg.NewSegment(segmentName, reqBody.AudienceCvg)
	if err := serviceRepo.SegmentDb.CreateObject(segment); err != nil {
		switch err.(type) {
		case db.ErrInternal:
			writeResponse(w, []byte("internal error"), 500)
			return
		default:
			writeResponse(w, []byte("segment with such name already exists"), 400)
			return
		}
	}
	if err := serviceRepo.SetRandomSegmentAuditory(segment); err != nil {
		writeResponse(w, []byte("internal error"), 500)
		return
	}
	writeResponse(w, []byte("OK"), 200)
}

func deleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	segmentName := chi.URLParam(r, "segmentName")
	res := "OK"
	logStatus := "SUCCESS"
	statusCode := 200
	if err := serviceRepo.SegmentDb.DeleteByName(segmentName); err != nil {
		res = "internal error"
		statusCode = 500
		logStatus = "DENIED"
	}
	fmt.Printf("%s %s ==> delete segment %s | %s\n", r.Method, r.URL.Path, segmentName, logStatus)
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, res)
}

func modifyUserSegments(w http.ResponseWriter, r *http.Request) {
	var reqBody userSegmentsModifyBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&reqBody); err != nil {
		// add case of incorrect body format
		writeResponse(w, []byte("incorrect body format"), 400)
		return
	}
	toAdd := xorStringArrays(reqBody.SegmentsToAdd, reqBody.SegmentsToRemove)
	toRm := xorStringArrays(reqBody.SegmentsToRemove, reqBody.SegmentsToAdd)
	userId := reqBody.UserId

	// it must be like some kind of a transaction
	// so if one value is incorrect, the others will be ignored
	unableToRm, removable := serviceRepo.CheckNonExistantSegments(toRm)
	unableToAdd, addable := serviceRepo.CheckNonExistantSegments(toAdd)
	var errorResponse userModifyErrorResponse = userModifyErrorResponse{}
	if !(len(unableToAdd) == 0 && len(unableToRm) == 0) {
		errorResponse.Message = "objects with these values were not found"
		if len(unableToAdd) > 0 {
			errorResponse.SegmentsToAdd = unableToAdd
		}
		if len(unableToRm) > 0 {
			errorResponse.SegmentsToRemove = unableToRm
		}
		resp, _ := json.Marshal(errorResponse)
		writeResponse(w, resp, 400)
		return
	}
	for _, v := range removable {
		err := serviceRepo.UserSegmentDb.SetUserSegmentInactive(userId, v)
		if err != nil {
			writeResponse(w, []byte("internal error"), 500)
			return
		}
	}
	for _, v := range addable {
		userSegment := usseg.NewUserSegment(userId, v)
		err := serviceRepo.UserSegmentDb.CreateObject(userSegment)
		if err != nil {
			var resp []byte
			var statusCode int
			switch err.(type) {
			case db.ErrUniqueConstraintFailed:
				resp = []byte("user already has such segments: " + err.Error())
				statusCode = 400
			default:
				resp = []byte("internal error")
				statusCode = 500
			}
			writeResponse(w, resp, statusCode)
			return
		}
	}
	writeResponse(w, []byte("OK"), 200)
}

func getUserSegments(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")
	userSegments := userSegmentsResponseBody{}

	// we ignore error, returned by atoi
	// because our router checks if the
	// value contains digits only
	userId, _ := strconv.Atoi(userIdStr)
	segmentNames := serviceRepo.GetUserActiveSegments(userId)
	if segmentNames == nil {
		writeResponse(w, []byte("user not found"), 404)
		return
	}
	userSegments.Segments = segmentNames
	res, err := json.Marshal(userSegments)
	if err != nil {
		writeResponse(w, []byte("internal error"), 500)
		return
	}
	writeResponse(w, res, 200)
}

func getUserSegmentsInPeriod(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(chi.URLParam(r, "userId"))
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))

	currentYear := time.Now().Year()
	currentMonth := getMonthNumber(time.Now().Month().String())
	if currentYear < year || currentYear == year && currentMonth < month {
		writeResponse(w, []byte("date %d/%d is in future - cannot calculate"), 400)
		return
	}
	if year <= 1970 {
		writeResponse(w, []byte("date %d/%d is in very past - cannot calculate"), 400)
		return
	}
	segments := serviceRepo.UserSegmentDb.GetUserSegmentActionsInPeriod(userId, year, month)
	if segments == nil {
		writeResponse(w, []byte("user not found"), 404)
		return
	}
	filename, err := createUserReport(segments)
	if err != nil {
		writeResponse(w, []byte("internal error"), 500)
		return
	}
	writeResponse(w, []byte(fmt.Sprintf("http://%s/user-report/%s", r.Host, filename)), 200)
}

func downloadUserReport(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	filepath := fmt.Sprintf("./data/%s.csv", filename)
	file, err := os.Open(filepath)
	if err != nil && !os.IsNotExist(err) {
		writeResponse(w, []byte("internal error"), 500)
		return
	} else if err != nil && os.IsNotExist(err) {
		writeResponse(w, []byte("user report does not exist"), 404)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename+".csv")
	http.ServeFile(w, r, filepath)
	file.Close()
	if err = os.Remove(filepath); err != nil {
		fmt.Println(err)
	}
}

func getMonthNumber(month string) int {
	switch month {
	case "January":
		return 1
	case "February":
		return 2
	case "March":
		return 3
	case "April":
		return 4
	case "May":
		return 5
	case "June":
		return 6
	case "July":
		return 7
	case "August":
		return 8
	case "September":
		return 9
	case "October":
		return 10
	case "November":
		return 11
	case "December":
		return 12
	default:
		return -1
	}
}

func xorStringArrays(str1, str2 []string) []string {
	checker := make(map[string]bool, len(str1)+len(str2))
	res := make([]string, 0, len(str1))
	for _, s := range str2 {
		checker[s] = true
	}
	for _, s := range str1 {
		if _, ok := checker[s]; !ok {
			res = append(res, s)
			checker[s] = true
		}
	}
	return res
}

func writeResponse(w http.ResponseWriter, mesg []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(mesg)
}

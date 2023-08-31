package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"

	usg "github.com/A1sal/AvitoTech"
)

func createUserReport(data []usg.UserSegmentActions) (string, error) {
	filename, er := uuid.NewRandom()

	if er != nil {
		fmt.Println(er)
		return "", er
	}
	filepath := fmt.Sprintf("./data/%s.csv", filename.String())
	file, er := os.Create(filepath)
	if er != nil {
		fmt.Println(er)
		return "", er
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	er = csvWriter.Write([]string{"user_id", "segment_name", "date", "operation"})
	if er != nil {
		fmt.Println(er)
		return "", er
	}
	for _, v := range data {
		er = csvWriter.Write([]string{strconv.Itoa(v.UserId), v.SegmentName, v.Date.String(), v.OperationType})
		if er != nil {
			fmt.Println(er)
			return "", er
		}
	}
	csvWriter.Flush()
	return filename.String(), nil
}

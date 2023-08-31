package repository

import (
	"fmt"

	db "github.com/A1sal/AvitoTech/database"
	seg "github.com/A1sal/AvitoTech/segment"
	u "github.com/A1sal/AvitoTech/user"
	usseg "github.com/A1sal/AvitoTech/usersegment"
)

type ServiceMockRepository struct {
	SegmentDb     seg.SegmentDatabase
	UserSegmentDb usseg.UserSegmentDatabase
	UserDb        u.UserDatabase
}

func NewServiceMockRepository(segmentDb seg.SegmentDatabase, usgDb usseg.UserSegmentDatabase, userDb u.UserDatabase) *ServiceMockRepository {
	return &ServiceMockRepository{
		SegmentDb:     segmentDb,
		UserSegmentDb: usgDb,
		UserDb:        userDb,
	}
}

func (r *ServiceMockRepository) GetSegmentsByUserId(id int) ([]seg.Segment, error) {
	usgs := r.UserSegmentDb.GetByUserId(id)
	res := make([]seg.Segment, 0)
	for _, userSegment := range usgs {
		v, _ := r.SegmentDb.GetByName(userSegment.GetSegmentName())
		res = append(res, v)
	}
	return res, nil
}

func (r *ServiceMockRepository) GetUserIdsBySegmentName(name string) ([]int, error) {
	usgs := r.UserSegmentDb.GetBySegmentName(name)
	res := make([]int, 0)
	for _, userSegment := range usgs {
		res = append(res, userSegment.GetUserId())
	}
	return res, nil
}

func (r *ServiceMockRepository) CheckNonExistantSegments(segmentNames []string) ([]string, []string) {
	existing := make([]string, 0, len(segmentNames))
	nonExisting := make([]string, 0, len(segmentNames))
	for _, v := range segmentNames {
		s, err := r.SegmentDb.GetByName(v)
		if err == nil {
			existing = append(existing, s.GetName())
			continue
		}
		switch err.(type) {
		case db.ErrObjNotFound:
			nonExisting = append(nonExisting, v)
		default:
			return nil, nil
		}
	}
	return nonExisting, existing
}

func (r *ServiceMockRepository) GetUserActiveSegments(user_id int) []string {
	var names []string
	res := r.UserSegmentDb.GetUserActiveSegments(user_id)
	if res == nil {
		return nil
	}
	names = make([]string, 0, len(res))
	for _, v := range res {
		names = append(names, v.GetSegmentName())
	}
	return names
}

func (r *ServiceMockRepository) SetRandomSegmentAuditory(s seg.Segment) error {
	users := r.UserDb.GetRandomUsersByPercent(s.GetAudienceCvg())
	if users == nil {
		return db.ErrInternal{}
	}
	for _, user := range users {
		userSegment := usseg.NewUserSegment(user.Id, s.Name)
		if err := r.UserSegmentDb.CreateObject(userSegment); err != nil {
			fmt.Println(err)
			return db.ErrInternal{}
		}
	}
	return nil
}

package usersegment

import (
	"time"
)

type UserSegmentActions struct {
	UserId        int       `ksql:"user_id"`
	SegmentName   string    `ksql:"segment_name"`
	Date          time.Time `ksql:"date"`
	OperationType string    `ksql:"operation"`
}

type UserSegment struct {
	UserId      int    `ksql:"user_id"`
	SegmentName string `ksql:"segment_name"`
	IsActive    bool   `ksql:"is_active"`
}

func NewUserSegment(userId int, segmentName string) UserSegment {
	return UserSegment{
		UserId:      userId,
		SegmentName: segmentName,
		IsActive:    true,
	}
}

func (u UserSegment) GetUserId() int {
	return u.UserId
}

func (u UserSegment) GetStatus() bool {
	return u.IsActive
}

func (u UserSegment) GetSegmentName() string {
	return u.SegmentName
}



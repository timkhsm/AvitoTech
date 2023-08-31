package usersegment

import (
	datab "github.com/A1sal/AvitoTech/database"
)

type UserSegmentDatabase interface {
	datab.Database[UserSegment]
	GetByUserId(int) []UserSegment
	GetBySegmentName(string) []UserSegment
	GetUserActiveSegments(int) []UserSegment
	GetUserSegmentActionsInPeriod(int, int, int) []UserSegmentActions
	DeleteBySegmentName(string) error
	DeleteByUserId(int) error
	SetUserSegmentInactive(int, string) error
}

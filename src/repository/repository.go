package repository


import (
	seg "github.com/A1sal/AvitoTech/segment"
)

type ServiceRepository interface {
	GetUserIdsBySegmentName(string) ([]int, error)
	GetSegmentsByUserId(int) ([]seg.Segment, error)
	GetUserActiveSegments(int) ([]string, error)
	CheckNonExistantSegments([]string) ([]string, []string)
	SetRandomSegmentAuditory(seg.Segment) error
}

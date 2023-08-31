package segment

import (
	datab "github.com/A1sal/AvitoTech/database"
)

type SegmentDatabase interface {
	datab.Database[Segment]
	GetByName(string) (Segment, error)
	DeleteByName(string) error
}

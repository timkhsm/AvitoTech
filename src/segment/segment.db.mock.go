package segment

import (
	db "github.com/A1sal/AvitoTech/database"
)

type SegmentMockDatabase struct {
	storage map[string]Segment
}

func NewSegmentMockDatabase() *SegmentMockDatabase {
	return &SegmentMockDatabase{
		storage: make(map[string]Segment),
	}
}

func (d *SegmentMockDatabase) GetByName(name string) (Segment, error) {
	if v, ok := d.storage[name]; !ok {
		return v, db.ErrObjNotFound{}
	} else {
		return v, nil
	}
}

func (d *SegmentMockDatabase) CreateObject(s Segment) error {
	if _, ok := d.storage[s.Name]; ok {
		return db.ErrUniqueConstraintFailed{Field: "name", Value: s.Name}
	}
	d.storage[s.Name] = s
	return nil
}

func (d *SegmentMockDatabase) DeleteObject(s Segment) error {
	delete(d.storage, s.Name)
	return nil
}

func (d *SegmentMockDatabase) DeleteByName(name string) error {
	delete(d.storage, name)
	return nil
}

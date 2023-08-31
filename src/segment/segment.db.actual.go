package segment

import (
	"context"
	"fmt"
    "github.com/vingarcia/ksql"
	db_ "github.com/A1sal/AvitoTech/database"
)

type SegmentActualDatabase struct {
	db    ksql.DB
	table ksql.Table
}

func NewSegmentActualDatabase(db ksql.DB) *SegmentActualDatabase {
	return &SegmentActualDatabase{
		db:    db,
		table: ksql.NewTable("segments"),
	}
}

func (d *SegmentActualDatabase) CreateObject(s Segment) error {
	_, err := d.db.Exec(context.Background(), "insert into segments(name, audience_cvg) values($1, $2)", s.Name, s.AudienceCvg)
	if err != nil {
		fmt.Println(err)
		err = db_.ErrUniqueConstraintFailed{Field: "name", Value: s.Name}
	}
	return err
}

func (d *SegmentActualDatabase) DeleteObject(s Segment) error {
	err := d.db.Delete(context.Background(), d.table, &s)
	if err != nil {
		fmt.Println(err)
		err = db_.ErrObjNotFound{}
	}
	return err
}

func (d *SegmentActualDatabase) DeleteByName(name string) error {
	_, err := d.db.Exec(context.Background(), "delete from segments where name=$1", name)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *SegmentActualDatabase) GetByName(name string) (Segment, error) {
	var res Segment
	err := d.db.QueryOne(context.Background(), &res, "select * from segments where name=$1", name)
	if err != nil {
		fmt.Println(err)
		err = db_.ErrObjNotFound{}
	}
	return res, err
}

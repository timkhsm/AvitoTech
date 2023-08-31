package usersegment

import (
	"context"
	"fmt"
	datab_ "github.com/A1sal/AvitoTech/database"
	"github.com/vingarcia/ksql"
)

type UserSegmentActualDatabase struct {
	db    ksql.DB
	table ksql.Table
}

func NewUserSegmentActualDatabase(db ksql.DB) *UserSegmentActualDatabase {
	return &UserSegmentActualDatabase{
		db:    db,
		table: ksql.NewTable("user_segments"),
	}
}

func (d *UserSegmentActualDatabase) CreateObject(s UserSegment) error {
	_, err := d.db.Exec(context.Background(), "insert into user_segments(user_id, segment_name) values($1, $2)", s.UserId, s.SegmentName)
	if err != nil {
		fmt.Println(err)
		err = datab_.ErrUniqueConstraintFailed{Field: "user_id&segment_name", Value: fmt.Sprintf("%d&%s", s.UserId, s.SegmentName)}
	}
	return err
}

func (d *UserSegmentActualDatabase) GetByUserId(id int) []UserSegment {
	var res []UserSegment
	err := d.db.Query(context.Background(), &res, "select * from user_segments where user_id=$1", id)
	if err != nil {
		fmt.Println(err)
		res = nil
	}
	return res
}
func (d *UserSegmentActualDatabase) DeleteObject(s UserSegment) error {
	err := d.db.Delete(context.Background(), d.table, &s)
	if err != nil {
		err = datab_.ErrObjNotFound{}
	}
	return err
}
func (d *UserSegmentActualDatabase) GetUserActiveSegments(id int) []UserSegment {
	res := make([]UserSegment, 0)
	err := d.db.Query(
		context.Background(),
		&res,
		`select
			*
		from
			user_segments
		where
			user_segments.user_id=$1
		and
			user_segments.is_active='true'`,
		id,
	)
	if err != nil {
		fmt.Println(err)
		res = nil
	}
	return res
}

func (d *UserSegmentActualDatabase) GetUserSegmentActionsInPeriod(userId, year, month int) []UserSegmentActions {
	res := make([]UserSegmentActions, 0)
	err := d.db.Query(
		context.Background(),
		&res,
		`select
			user_id, segment_name, added_at as date, 'added' as operation
		from
			user_segments
		where
			user_id=$3
		and
			extract(year from added_at)=$1
		and
			extract(month from added_at)=$2
		union all
		select
			user_id, segment_name, removed_at as date, 'removed' as operation
		from
			user_segments
		where
			user_id=$3
		and
			extract(year from removed_at)=$1
		and
			extract(month from removed_at)=$2
			`,
		year,
		month,
		userId,
	)
	if err != nil {
		fmt.Println(err)
		res = nil
	}
	return res
}

func (d *UserSegmentActualDatabase) GetBySegmentName(name string) []UserSegment {
	res := make([]UserSegment, 0)
	err := d.db.Query(context.Background(), &res, "select * from user_segments where segment_name=$1", name)
	if err != nil {
		fmt.Println(err)
		res = nil
	}
	return res
}

func (d *UserSegmentActualDatabase) DeleteByUserId(id int) error {
	_, err := d.db.Exec(context.Background(), "delete from user_segments where user_id=$1", id)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *UserSegmentActualDatabase) SetUserSegmentInactive(user_id int, segment_name string) error {
	_, err := d.db.Exec(context.Background(), "update user_segments set removed_at=now(), is_active='false' where user_id=$1 and segment_name=$2", user_id, segment_name)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *UserSegmentActualDatabase) DeleteBySegmentName(name string) error {
	_, err := d.db.Exec(context.Background(), "delete from user_segments where segment_name=$1", name)
	if err != nil {
		fmt.Println(err)
	}
	return err
}



package user

import (
	"context"
	"fmt"
	"github.com/vingarcia/ksql"
)

type UserActualDatabase struct {
	db    ksql.DB
	table ksql.Table
}

func NewUserActualDatabase(db ksql.DB) *UserActualDatabase {
	return &UserActualDatabase{
		db:    db,
		table: ksql.NewTable("users"),
	}
}

func (d *UserActualDatabase) DeleteObject(u User) error {
	err := d.db.Delete(context.Background(), d.table, u)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *UserActualDatabase) CreateObject(u User) error {
	_, err := d.db.Exec(context.Background(), "insert into users default values")
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (d *UserActualDatabase) GetRandomUsersByPercent(percent int) []User {
	var res []User
	err := d.db.Query(context.Background(), &res, "select * from users order by random() limit ((select count(*) from users) * $1/100)", percent)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (d *UserActualDatabase) GetUserById(id int) (User, error) {
	var res User
	err := d.db.QueryOne(context.Background(), &res, "select * from users where id=$1", id)
	if err != nil {
		fmt.Println(err)
	}
	return res, err
}



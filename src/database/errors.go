package database

import (
	"fmt"
)


type ErrObjAlreadyExists struct {
	Id int
}

func (err ErrObjAlreadyExists) Error() string {
	return fmt.Sprintf("object with id %d already exists\n", err.Id)
}

type ErrObjNotFound struct {
}

func (err ErrObjNotFound) Error() string {
	return "object not found"
}

type ErrInternal struct {
}

func (err ErrInternal) Error() string {
	return "internal error"
}

type ErrUniqueConstraintFailed struct {
	Field, Value string
}

func (err ErrUniqueConstraintFailed) Error() string {
	return fmt.Sprintf("field(s) %s value(s) %s\n", err.Field, err.Value)
}

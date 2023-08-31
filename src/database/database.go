package database



type Database[T any] interface {
	CreateObject(obj T) error
	DeleteObject(obj T) error
}

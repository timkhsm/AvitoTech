package user

import (
	datab_ "github.com/A1sal/AvitoTech/database"
)



type UserDatabase interface {
	datab_.Database[User]
	GetUserById(int) (User, error)
	GetRandomUsersByPercent(int) []User
}

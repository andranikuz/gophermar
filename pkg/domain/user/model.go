package user

import "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

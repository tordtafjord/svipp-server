package models

type Role string

const (
	RoleUser   Role = "user"
	RoleDriver Role = "driver"
	RoleAdmin  Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

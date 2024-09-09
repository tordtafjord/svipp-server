package models

type Role string

const (
	RoleUser   Role = "user"
	RoleDriver Role = "driver"
	RoleAdmin  Role = "admin"
)

// No need for String() method

// IsValid checks if the role is valid
func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleDriver, RoleAdmin:
		return true
	}
	return false
}

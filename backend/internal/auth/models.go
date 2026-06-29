package auth

import "github.com/golang-jwt/jwt/v5"

type UserType string

const (
	UserHuman  UserType = "human"
	UserDevice UserType = "device"
)

type Permissions string

const (
	PermissionsAdmin   Permissions = "admin"
	PermissionsDefault Permissions = "default"
)

type User struct {
	Type        UserType    `json:"type"`
	Id          string      `json:"id"`
	Permissions Permissions `json:"permissions"`
}

func (u *User) IsAuthorized(p Permissions) bool {
	switch u.Permissions {
	case PermissionsAdmin:
		return true
	case PermissionsDefault:
		return p != PermissionsAdmin
	default:
		return false
	}
}

type jwtClaims struct {
	User
	jwt.RegisteredClaims
}

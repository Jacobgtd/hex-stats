package authn

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

type jwtClaims struct {
	User
	jwt.RegisteredClaims
}

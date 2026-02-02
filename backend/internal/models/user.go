package models

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

type User struct {
	ID        uint `gormm:"primaryKey"`
	Name      string
	Email     string `gormm:"uniqueIndex"`
	Password  string
	Role      UserRole `gorm:"type:varchar(10)"`
	CreatedAt time.Time
}

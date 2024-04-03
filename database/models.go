package database

import "gorm.io/gorm"

type Order struct {
	ID string

	gorm.Model
}

package model

import (
	"time"
)

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	CreationDate time.Time `db:"creation_date"`
	UpdateDate   time.Time `db:"update_date"`
}

package model

import (
	"time"
)

type Wallet struct {
	ID      		int			`db:"id"`
	UserId     		int			`db:"user_id"`
	Balance			float64 	`db:"balance"`
	CreationDate	time.Time	`db:"creation_date"`
	UpdateDate		time.Time	`db:"update_date"`
	DeletionDate	*time.Time 	`db:"deletion_date"` 
}
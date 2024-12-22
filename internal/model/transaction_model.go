package model

import (
	"time"
)

type Transaction struct {
	ID               int        `db:"id"`
	SenderUserId     int        `db:"sender_user_id"`
	SenderWalletId   int        `db:"sender_wallet_id"`
	ReceiverUserId   int        `db:"receiver_user_id"`
	ReceiverWalletId int        `db:"receiver_wallet_id"`
	Amount           float64    `db:"amount"`
	Type             string     `db:"type"`
	CreationDate     time.Time  `db:"creation_date"`
	UpdateDate       time.Time  `db:"update_date"`
	DeletionDate     *time.Time `db:"deletion_date"`
}

type TransactionData struct {
	Type            string
	Amount          float64
	Receiver        string
	TransactionDate time.Time
}

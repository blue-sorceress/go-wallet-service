package object

type TransferParams struct {
	UserId         int     `json:"user_id"`
	ReceiverUserId int     `json:"receiver_user_id"`
	Amount         float64 `json:"amount" binding:"required,min=0,max=10000000"`
}

type UserParams struct {
	UserId int `json:"user_id"`
}

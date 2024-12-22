package object

type OAuthParams struct {
	UserId int `json:"user_id" binding:"required"`
}

type BalanceParams struct {
	UserId int `json:"user_id" binding:"required"`
}

type DepositParams struct {
	UserId int     `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,min=0,max=10000000"`
}

type WithdrawParams struct {
	UserId int     `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,min=0,max=10000000"`
}

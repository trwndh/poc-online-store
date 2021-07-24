package model

type RegisterPayment struct {
	ID        int64
	UserID    int64
	CartID    int64
	Products  string
	Status    string
	CreatedAt string
	ExpiredAt string
}

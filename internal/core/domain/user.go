package domain

type User struct {
	ID          int    `gorm:"column:uid" json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number" gorm:" column:phone_number"`
	EmailID     string `json:"email_id" gorm:"column:email"`
}

type Summary struct {
	ID      int     `gorm:"column:uid" json:"uid"`
	Balance float32 `json:"Balance" gorm:"column:amount"`
}

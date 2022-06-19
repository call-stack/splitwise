package domain

type Expense struct {
	ID               int     `json:"id"`
	UserID           int     `json:"paid_by"`
	Category         string  `json:"category"`
	IsExpenseSettled bool    `json:"is_expense_settled"`
	Amount           float32 `json:"amount"`
	PaidTo           []int   `json:"paid_to,omitempty"`
	Gid              int     `json:"gid,omitempty"`
}

type ExpenseEntity struct {
	ID               int    `gorm:"column:eid"`
	UserID           int    `gorm:"column:uid"`
	Category         string `gorm:"column:category"`
	IsExpenseSettled bool   `gorm:"column:is_expense_settled"`
	Amount           float32
	PeopleInvolved   int
}

type UserExpenseEntity struct {
	UID       int     `gorm:"column:paid_by"`
	ExpenseID int     `gorm:"column:eid"`
	Amount    float32 `gorm:"column:amount"`
	OldAmount float32 `gorm:"-"`
}

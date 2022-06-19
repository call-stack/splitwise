package port

import "splitwise/internal/core/domain"

type SplitwiseService interface {
	AddUser(name, phoneNumber, email string) (domain.User, error)
	AddGroup(name string, uid []int) (domain.Group, error)
	NewExpense(paidByUserID int, splitUserID []int, category string, amount float32) (domain.Expense, error)
	UpdateExpense(txnID int, amount float32) (domain.Expense, error)
	SettleExpense(txnID int) (domain.Expense, error)
	ViewExpense(expenseID int) (domain.Expense, error)
	Summary(uid int) (domain.Summary, error)
	NewGroupExpense(paidByUserID int, groupID int, category string, amount float32) (domain.Expense, error)
	GetAllUnSettledExpense() ([]*domain.Expense, error)
}

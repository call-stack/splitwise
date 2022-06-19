package port

import "splitwise/internal/core/domain"

type SplitwiseRepo interface {
	AddUser(user domain.User) (domain.User, error)
	AddGroup(group domain.Group) (domain.GroupEntity, error)
	NewExpense(transaction domain.ExpenseEntity, userExpense []domain.UserExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error)
	GetExpense(txnID int) (domain.ExpenseEntity, error)
	GetUIDs(txnID int) ([]domain.UserExpenseEntity, error)
	UpdateExpense(expense domain.ExpenseEntity, userExpense []domain.UserExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error)
	GetSummary(uid int) (domain.Summary, error)
	GetGroupUsers(gid int) ([]int, error)
	SettleExpense(transaction domain.ExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error)
	GetAllUnSettledExpense() ([]*domain.ExpenseEntity, error)
}

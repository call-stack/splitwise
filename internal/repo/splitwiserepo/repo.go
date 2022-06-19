package splitwiserepo

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"splitwise/internal/core"
	"splitwise/internal/core/domain"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo() *Repo {
	username, password, dbname, host := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("DATABASE_HOST")
	fmt.Println(username, password, dbname, host)
	dbURL := fmt.Sprintf("postgresql://%v:%v@%v:5432/%v", username, password, host, dbname)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return &Repo{db}
}

func (r *Repo) AddUser(user domain.User) (domain.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(core.UserTable).Create(&user).Error; err != nil {
			return err
		}
		summary := &domain.Summary{ID: user.ID}
		if err := tx.Table(core.Summary).Create(&summary).Error; err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (r *Repo) AddGroup(group domain.Group) (domain.GroupEntity, error) {
	grpEntity := domain.GroupEntity{Name: group.Name}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(core.GroupTable).Create(&grpEntity).Error; err != nil {
			return err
		}

		for _, id := range group.UserID {
			userGroup := &domain.UserGroup{GID: grpEntity.ID, UID: id}
			err := tx.Table(core.UserGroup).Create(&userGroup).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return grpEntity, err
}

func (r *Repo) NewExpense(expense domain.ExpenseEntity, userExpense []domain.UserExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(core.ExpenseTable).Create(&expense).Error; err != nil {
			return err
		}

		n := len(userExpense)
		for i := 0; i < n; i++ {
			userExpense[i].ExpenseID = expense.ID

			if err := tx.Table(core.UserExpenseTable).Create(&userExpense[i]).Error; err != nil {
				return err
			}

			err := tx.Table(core.Summary).Where("uid=?", summary[i].ID).UpdateColumn("amount", gorm.Expr("amount + ?", summary[i].Balance)).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return expense, err
}

func (r *Repo) ModifyTransaction(transaction domain.Expense, userExpense []domain.UserExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbExpense domain.Expense
		if err := tx.Table(core.ExpenseTable).Model(&dbExpense).Where("eid = ?", transaction.ID).Update("amount", transaction.Amount).Error; err != nil {
			return err
		}
		n := len(userExpense)
		for i := 0; i < n; i++ {
			if err := tx.Table(core.UserExpenseTable).Where("eid = ? and paid_by = ?", transaction.ID, userExpense[i].UID).Update("amount", userExpense[i].Amount).Error; err != nil {
				return err
			}

			if err := tx.Table(core.Summary).Where("uid = ?", summary[i].ID).Update("amount", gorm.Expr("amount  + ?", summary[i].Balance)).Error; err != nil {
				return err
			}
		}

		return nil

	})
	return domain.ExpenseEntity{}, err
}

func (r *Repo) GetExpense(txnID int) (domain.ExpenseEntity, error) {
	expense := domain.ExpenseEntity{ID: txnID}
	if err := r.db.Table(core.ExpenseTable).Where("eid = ?", txnID).Take(&expense).Error; err != nil {
		return expense, err
	}
	return expense, nil
}

func (r *Repo) GetUIDs(txnID int) ([]domain.UserExpenseEntity, error) {
	var users []domain.UserExpenseEntity
	if err := r.db.Table(core.UserExpenseTable).Where("eid = ?", txnID).Find(&users).Error; err != nil {
		return []domain.UserExpenseEntity{}, err
	}

	return users, nil
}

func (r *Repo) UpdateExpense(expense domain.ExpenseEntity, userExpense []domain.UserExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Table(core.ExpenseTable).Model(&domain.ExpenseEntity{}).Where("eid = ?", expense.ID).Update("amount", expense.Amount).Error; err != nil {
			return err
		}

		n := len(userExpense)
		for i := 0; i < n; i++ {

			if err := tx.Table(core.UserExpenseTable).Where("eid = ? and paid_by = ?", expense.ID, userExpense[i].UID).Update("amount", userExpense[i].Amount).Error; err != nil {
				return err
			}

			if err := tx.Table(core.Summary).Where("uid = ?", summary[i].ID).Update("amount", gorm.Expr("amount  + ?", summary[i].Balance)).Error; err != nil {
				return err
			}
		}

		return nil

	})
	return expense, err
}

func (r *Repo) GetSummary(userID int) (domain.Summary, error) {
	var summary domain.Summary
	if err := r.db.Table(core.Summary).Where("uid = ?", userID).First(&summary).Error; err != nil {
		return domain.Summary{}, err
	}

	return summary, nil
}

func (r *Repo) GetGroupUsers(gid int) ([]int, error) {
	var users []int
	if err := r.db.Table(core.UserGroup).Where("gid = ?", gid).Select("uid").Scan(&users).Error; err != nil {
		return []int{}, err
	}

	return users, nil
}

func (r *Repo) SettleExpense(expense domain.ExpenseEntity, summary []domain.Summary) (domain.ExpenseEntity, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Table(core.ExpenseTable).Model(&domain.ExpenseEntity{}).Where("eid = ?", expense.ID).Update("is_expense_settled", true).Error; err != nil {
			return err
		}

		n := len(summary)
		for i := 0; i < n; i++ {
			if err := tx.Table(core.Summary).Where("uid = ?", summary[i].ID).Update("amount", gorm.Expr("amount  + ?", summary[i].Balance)).Error; err != nil {
				return err
			}
		}

		return nil

	})
	return expense, err
}
func (r *Repo) GetAllUnSettledExpense() ([]*domain.ExpenseEntity, error) {
	var expenses []*domain.ExpenseEntity
	if err := r.db.Table(core.ExpenseTable).Where("is_expense_settled", false).Find(&expenses).Error; err != nil {
		return expenses, err
	}

	return expenses, nil

}

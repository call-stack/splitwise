package services

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"splitwise/internal/core"
	"splitwise/internal/core/domain"
	"splitwise/internal/core/port"
)

type Service struct {
	repo port.SplitwiseRepo
}

func (s *Service) AddUser(name, phoneNumber, email string) (domain.User, error) {
	err := validation.Errors{
		"phone_number": validation.Validate(phoneNumber, validation.Required, validation.By(validatePhoneNumber())),
		"email":        validation.Validate(email, validation.Required, is.Email),
	}.Filter()

	if err != nil {
		return domain.User{}, err
	}

	usr := domain.User{
		Name:        name,
		PhoneNumber: phoneNumber,
		EmailID:     email,
	}
	usr, err = s.repo.AddUser(usr)
	if err != nil {
		return domain.User{}, err
	}

	return usr, nil

}

func (s *Service) AddGroup(name string, uid []int) (domain.Group, error) {

	if len(uid) <= 1 {
		return domain.Group{}, fmt.Errorf(core.InvalidGroupSize)
	}

	grp := domain.Group{
		Name:   name,
		UserID: uid,
	}
	grpE, err := s.repo.AddGroup(grp)
	if err != nil {
		return domain.Group{}, err
	}

	return domain.Group{
		ID:     grpE.ID,
		Name:   grpE.Name,
		UserID: uid,
	}, nil
}

func (s *Service) NewExpense(paidByUserID int, paidToIds []int, category string, amount float32) (domain.Expense, error) {

	err := validation.Errors{
		"amount":  validation.Validate(amount, validation.By(validateAmount())),
		"paid_to": validation.Validate(paidToIds, validation.Required, validation.By(validateLength()), validation.By(validateEntries(paidByUserID))),
	}.Filter()

	if err != nil {
		return domain.Expense{}, err
	}

	userInvolved := len(paidToIds) + 1
	expense := domain.ExpenseEntity{UserID: paidByUserID, Category: category, IsExpenseSettled: false, Amount: amount, PeopleInvolved: userInvolved}
	userExpenses := make([]domain.UserExpenseEntity, userInvolved)
	summary := make([]domain.Summary, userInvolved)
	var amountSplit float32
	splitsIds := append(paidToIds, paidByUserID)
	for idx, uid := range splitsIds {
		if uid != paidByUserID {
			amountSplit = -1 * amount / float32(userInvolved)
			summary[idx] = domain.Summary{ID: uid, Balance: amountSplit}
		} else {
			amountSplit = amount / float32(userInvolved)
			summary[idx] = domain.Summary{ID: uid, Balance: amount - amountSplit}
		}
		userExpenses[idx] = domain.UserExpenseEntity{
			UID:    uid,
			Amount: amountSplit,
		}

	}

	expenseEntity, err := s.repo.NewExpense(expense, userExpenses, summary)
	if err != nil {
		return domain.Expense{}, err
	}

	return domain.Expense{
		ID:               expenseEntity.ID,
		UserID:           expenseEntity.UserID,
		Category:         expenseEntity.Category,
		IsExpenseSettled: expenseEntity.IsExpenseSettled,
		Amount:           expenseEntity.Amount,
		PaidTo:           paidToIds,
	}, err
}

func (s *Service) UpdateExpense(txnID int, amount float32) (domain.Expense, error) {
	err := validation.Errors{
		"amount": validation.Validate(amount, validation.By(validateAmount())),
	}.Filter()

	if err != nil {
		return domain.Expense{}, err
	}

	expenseEntity, err := s.repo.GetExpense(txnID)
	if err != nil {
		return domain.Expense{}, err
	}

	if expenseEntity.IsExpenseSettled {
		return domain.Expense{}, fmt.Errorf(core.CannotModifySettledExpense)
	}

	userExpense, err := s.repo.GetUIDs(expenseEntity.ID)
	if err != nil {
		return domain.Expense{}, err
	}

	var amountSplit float32

	splitUserIDs := make([]int, expenseEntity.PeopleInvolved-1)
	summary := make([]domain.Summary, expenseEntity.PeopleInvolved)
	for idx, ue := range userExpense {

		if ue.UID != expenseEntity.UserID {
			splitUserIDs[idx] = ue.UID
			amountSplit = -1 * amount / float32(expenseEntity.PeopleInvolved)

			summary[idx] = domain.Summary{ID: ue.UID, Balance: amountSplit + expenseEntity.Amount/float32(expenseEntity.PeopleInvolved)}

		} else {
			amountSplit = amount / float32(expenseEntity.PeopleInvolved)
			summary[idx] = domain.Summary{ID: ue.UID, Balance: (amount - amountSplit) - (expenseEntity.Amount - expenseEntity.Amount/float32(expenseEntity.PeopleInvolved))}
		}
		userExpense[idx].OldAmount = ue.Amount
		userExpense[idx].Amount = amountSplit

	}

	expenseEntity.Amount = amount
	expenseEntity, err = s.repo.UpdateExpense(expenseEntity, userExpense, summary)
	if err != nil {
		return domain.Expense{}, nil
	}

	return domain.Expense{
		ID:               expenseEntity.ID,
		UserID:           expenseEntity.UserID,
		Category:         expenseEntity.Category,
		IsExpenseSettled: expenseEntity.IsExpenseSettled,
		Amount:           expenseEntity.Amount,
		PaidTo:           splitUserIDs,
	}, nil
}

func (s *Service) SettleExpense(txnID int) (domain.Expense, error) {
	expenseEntity, err := s.repo.GetExpense(txnID)
	if err != nil {
		return domain.Expense{}, err
	}

	if expenseEntity.IsExpenseSettled {
		return domain.Expense{}, fmt.Errorf(core.ExpenseIsAlreadySettled)
	}
	userExpense, err := s.repo.GetUIDs(expenseEntity.ID)
	if err != nil {
		return domain.Expense{}, err
	}

	var amountSplit float32
	splitUserIDs := make([]int, expenseEntity.PeopleInvolved-1)
	summary := make([]domain.Summary, expenseEntity.PeopleInvolved)
	var amount float32
	for idx, ue := range userExpense {

		if ue.UID != expenseEntity.UserID {
			splitUserIDs[idx] = ue.UID
			amountSplit = -1 * amount / float32(expenseEntity.PeopleInvolved)

			summary[idx] = domain.Summary{ID: ue.UID, Balance: amountSplit + expenseEntity.Amount/float32(expenseEntity.PeopleInvolved)}

		} else {
			amountSplit = amount / float32(expenseEntity.PeopleInvolved)
			summary[idx] = domain.Summary{ID: ue.UID, Balance: (amount - amountSplit) - (expenseEntity.Amount - expenseEntity.Amount/float32(expenseEntity.PeopleInvolved))}
		}
		userExpense[idx].OldAmount = ue.Amount
		userExpense[idx].Amount = amountSplit

	}

	expenseEntity.IsExpenseSettled = true
	expenseEntity, err = s.repo.SettleExpense(expenseEntity, summary)
	if err != nil {
		return domain.Expense{}, nil
	}

	return domain.Expense{
		ID:               expenseEntity.ID,
		UserID:           expenseEntity.UserID,
		Category:         expenseEntity.Category,
		IsExpenseSettled: expenseEntity.IsExpenseSettled,
		Amount:           expenseEntity.Amount,
		PaidTo:           splitUserIDs,
	}, nil

}

func (s *Service) Summary(uid int) (domain.Summary, error) {
	summery, err := s.repo.GetSummary(uid)
	return summery, err
}

func (s *Service) ViewExpense(txnID int) (domain.Expense, error) {
	exp, err := s.repo.GetExpense(txnID)
	if err != nil {
		return domain.Expense{}, err
	}

	userIds, err := s.repo.GetUIDs(exp.ID)
	if err != nil {
		return domain.Expense{}, err
	}

	splitsIDs := make([]int, 0)
	for _, u := range userIds {
		if u.UID != exp.UserID {
			splitsIDs = append(splitsIDs, u.UID)
		}
	}

	return domain.Expense{
		ID:               exp.ID,
		UserID:           exp.UserID,
		Category:         exp.Category,
		IsExpenseSettled: exp.IsExpenseSettled,
		Amount:           exp.Amount,
		PaidTo:           splitsIDs,
	}, nil
}

func New(repo port.SplitwiseRepo) *Service {
	return &Service{repo}
}

func (s *Service) NewGroupExpense(paidByUserID int, groupID int, category string, amount float32) (domain.Expense, error) {
	userIDs, err := s.repo.GetGroupUsers(groupID)
	if err != nil {
		return domain.Expense{}, err
	}

	if len(userIDs) == 0 {
		return domain.Expense{}, fmt.Errorf(core.GroupNotFound)
	}

	users := make([]int, 0)
	isPayingUserPartOfGroup := false
	fmt.Println(userIDs, paidByUserID)
	for _, id := range userIDs {
		if id != paidByUserID {
			users = append(users, id)
		}

		if id == paidByUserID {
			isPayingUserPartOfGroup = true
		}
	}

	if !isPayingUserPartOfGroup {
		return domain.Expense{}, fmt.Errorf(core.PayingUserNotPartOfGroup)
	}

	expense, err := s.NewExpense(paidByUserID, users, category, amount)
	if err != nil {
		return domain.Expense{}, err
	}

	return expense, nil
}

func (s *Service) GetAllUnSettledExpense() ([]*domain.Expense, error) {
	expense, err := s.repo.GetAllUnSettledExpense()
	var exp []*domain.Expense
	if err != nil {
		return exp, err
	}
	for _, e := range expense {
		ex := &domain.Expense{
			ID:               e.ID,
			UserID:           e.UserID,
			Category:         e.Category,
			IsExpenseSettled: e.IsExpenseSettled,
			Amount:           e.Amount,
		}

		exp = append(exp, ex)
	}
	return exp, err
}

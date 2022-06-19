package services

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"splitwise/internal/core"
)

func validateAmount() validation.RuleFunc {
	return func(value interface{}) error {
		val, _ := value.(float32)
		if val <= 0 {
			return errors.New(core.InvalidAmount)
		}

		return nil
	}
}

func validateLength() validation.RuleFunc {
	return func(value interface{}) error {
		v, _ := value.([]int)
		if len(v) == 0 {
			return errors.New(core.InvalidPaidTo)
		}

		return nil
	}
}

func validateEntries(paidByUserID int) validation.RuleFunc {
	return func(value interface{}) error {
		v, _ := value.([]int)
		for _, u := range v {
			if u == paidByUserID {
				return errors.New(core.PaidToContainPaidBy)
			}
		}
		return nil
	}
}
func validatePhoneNumber() validation.RuleFunc {
	return func(value interface{}) error {
		val, _ := value.(string)
		if len(val) != 10 {
			return errors.New(core.InvalidPhoneNumber)
		}

		return nil
	}
}

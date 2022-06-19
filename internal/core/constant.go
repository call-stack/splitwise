package core

const (
	UserTable        = "users"
	GroupTable       = "groups"
	ExpenseTable     = "expense"
	UserExpenseTable = "user_expense"
	Summary          = "summary"
	UserGroup        = "user_group"
)

const (
	PayingUserNotPartOfGroup = "paying user does not belong to given group"
	GroupNotFound            = "group not found"
	ExpenseIsAlreadySettled  = "expense is already settled"
	InvalidGroupSize         = "there should be atlease 2 user belong to a group"
	InvalidPaidTo            = "there should be at-least 1 user in paid to"
	PaidToContainPaidBy      = "paid_to should not contain paid_by user id"
	InvalidPhoneNumber       = "incorrect phone number. Please add phone number without +91 or 0"
	InvalidAmount            = "amount should be greater than 0"
)

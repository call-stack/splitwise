package splithdl

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"splitwise/internal/core/domain"
	"splitwise/internal/core/port"
	"strconv"
)

type HttpHandler struct {
	Srv port.SplitwiseService
}

func NewHttpHandlers(splitWiseService port.SplitwiseService) *HttpHandler {
	return &HttpHandler{splitWiseService}
}

func (h *HttpHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr, err := h.Srv.AddUser(user.Name, user.PhoneNumber, user.EmailID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *HttpHandler) AddGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var group domain.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grp, err := h.Srv.AddGroup(group.Name, group.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(grp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *HttpHandler) NewExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var expense domain.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exp, err := h.Srv.NewExpense(expense.UserID, expense.PaidTo, expense.Category, expense.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(exp)
	if err != nil {
		return
	}

}
func (h *HttpHandler) ModifyTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var expense domain.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(expense)
	exp, err := h.Srv.UpdateExpense(expense.ID, expense.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *HttpHandler) SettleExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var expense domain.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(expense)
	expense, err = h.Srv.SettleExpense(expense.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		return
	}
}

func (h *HttpHandler) GroupExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var expense domain.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exp, err := h.Srv.NewGroupExpense(expense.UserID, expense.Gid, expense.Category, expense.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(exp)
	if err != nil {
		return
	}
}

func (h *HttpHandler) ViewExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return
	}
	exp, err := h.Srv.ViewExpense(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(exp)
	if err != nil {
		return
	}
}

func (h *HttpHandler) Summary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return
	}

	summ, err := h.Srv.Summary(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(summ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
}

func (h *HttpHandler) GetAllUnsettledExpenses(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	expense, err := h.Srv.GetAllUnSettledExpense()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"splitwise/internal/core/services"
	"splitwise/internal/handlers/splithdl"
	"splitwise/internal/repo/splitwiserepo"
)

func main() {
	repo := splitwiserepo.NewRepo()
	srv := services.New(repo)
	hdl := splithdl.NewHttpHandlers(srv)
	r := mux.NewRouter()
	r.HandleFunc("/user", hdl.AddUser).Methods("POST")
	r.HandleFunc("/group", hdl.AddGroup).Methods("POST")
	r.HandleFunc("/expense", hdl.NewExpense).Methods("POST")
	r.HandleFunc("/expense", hdl.ModifyTransaction).Methods("PATCH")
	r.HandleFunc("/settle", hdl.SettleExpense).Methods("POST")
	r.HandleFunc("/group_expense", hdl.GroupExpense).Methods("POST")
	r.HandleFunc("/view_expense", hdl.ViewExpense).Methods("GET")
	r.HandleFunc("/summary", hdl.Summary).Methods("GET")
	r.HandleFunc("/all_unsettled", hdl.GetAllUnsettledExpenses).Methods("GET")
	http.ListenAndServe("127.0.0.1:8080", r)

}

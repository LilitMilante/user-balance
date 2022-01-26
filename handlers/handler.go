package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user-balance/domain/entity"
	"user-balance/domain/service"

	"github.com/gorilla/mux"
)

type Server struct {
	r  *mux.Router
	bs *service.BalanceService
}

func NewServer(bs *service.BalanceService) *Server {
	r := mux.NewRouter()

	s := Server{
		r:  r,
		bs: bs,
	}

	return &s
}

func (s *Server) Start(port string) error {
	s.r.HandleFunc("/balance", s.MoneyTransactionHandler).Methods(http.MethodPatch)
	s.r.HandleFunc("/transfer", s.TransferMoneyHandler).Methods(http.MethodPatch)
	s.r.HandleFunc("/users/{id}/balance", s.CheckBalanceHandler).Methods(http.MethodGet)

	return http.ListenAndServe(":"+port, s.r)
}

func (s *Server) MoneyTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var b entity.Balance

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if b.TypeOp != entity.Plus && b.TypeOp != entity.Minus {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	b, err = s.bs.BalanceOperations(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *Server) TransferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	var t entity.Transfer

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if t.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = s.bs.TransferringFunds(t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *Server) CheckBalanceHandler(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	b, err := s.bs.UserBalance(id, currency)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (s *Server) UserTransactionsList(w http.ResponseWriter, r *http.Request) {

}

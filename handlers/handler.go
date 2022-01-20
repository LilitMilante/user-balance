package handlers

import (
	"encoding/json"
	"net/http"
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
	s.r.HandleFunc("/balance", s.PlusBalanceHandler).Methods(http.MethodPatch)
	s.r.HandleFunc("/balance", s.MinusBalanceHandler).Methods(http.MethodPatch)
	s.r.HandleFunc("/transfer", s.TransferMoneyHandler).Methods(http.MethodPatch)
	s.r.HandleFunc("/users/{id}/balance", s.CheckBalanceHandler).Methods(http.MethodGet)

	return http.ListenAndServe(":"+port, s.r)
}

func (s *Server) PlusBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var b entity.Balance

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	b, err = s.bs.CreditingFunds(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) MinusBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Списания средств с баланса"))
}

func (s *Server) TransferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Перевода средств от пользователя к пользователю"))
}

func (s *Server) CheckBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Получения текущего баланса пользователя"))
}

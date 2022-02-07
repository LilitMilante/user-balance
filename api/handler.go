package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"user-balance/domain/entity"
)

type Servicer interface {
	TransferMoney(t entity.Transaction) (entity.Transaction, error)
	Transactions(uID int64) ([]entity.Transaction, error)
	UserBalance(id int64, currency string) (entity.Balance, error)
}

type server struct {
	r  *mux.Router
	bs Servicer
}

func NewServer(bs Servicer) *server {
	r := mux.NewRouter()
	s := server{
		r:  r,
		bs: bs,
	}

	return &s
}

func (s *server) Start(port string) error {
	s.r.HandleFunc("/transactions", s.MoneyTransactionHandler).Methods(http.MethodPost)
	s.r.HandleFunc("/users/{id}/balance", s.CheckBalanceHandler).Methods(http.MethodGet)
	s.r.HandleFunc("/users/{id}/transactions", s.UserTransactions).Methods(http.MethodGet)

	return http.ListenAndServe(":"+port, s.r)
}

func (s *server) MoneyTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var t entity.Transaction

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if t.Event != entity.EventCrediting && t.Event != entity.EventWriteOffs && t.Event != entity.EventTransfer {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	t, err = s.bs.TransferMoney(t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (s *server) CheckBalanceHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *server) UserTransactions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	t, err := s.bs.Transactions(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

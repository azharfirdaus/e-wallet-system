package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/users", createUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{currency_code}", createWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/topup", topUpWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/pay", payWithWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transfer", transferWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/suspend", suspendWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}", getWalletHandler).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func createUserHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func createWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func topUpWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func payWithWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func transferWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func suspendWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func getWalletHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

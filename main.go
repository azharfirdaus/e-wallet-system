package main

import (
	"log"
	"net/http"

	"github.com/azharfirdaus/e-wallet-system/config"
	"github.com/azharfirdaus/e-wallet-system/handler"
	"github.com/gorilla/mux"
)

func main() {
	postgresConfig, err := config.LoadPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/users", handler.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{currency_code}", handler.CreateWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/topup", handler.TopUpWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/pay", handler.PayWithWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transfer", handler.TransferWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/suspend", handler.SuspendWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}", handler.GetWalletHandler).Methods(http.MethodGet)

	log.Printf(
		"postgres config loaded: host=%s port=%s database=%s user=%s sslmode=%s",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Database,
		postgresConfig.User,
		postgresConfig.SSLMode,
	)

	log.Fatal(http.ListenAndServe(":8080", router))
}

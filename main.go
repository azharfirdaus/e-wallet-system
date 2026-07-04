package main

import (
	"log"
	"net/http"

	"github.com/azharfirdaus/e-wallet-system/config"
	postgresdb "github.com/azharfirdaus/e-wallet-system/db"
	"github.com/azharfirdaus/e-wallet-system/handler"
	"github.com/gorilla/mux"
)

func main() {
	postgresConfig, err := config.LoadPostgresConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := postgresdb.NewPostgres(postgresConfig.DSN())
	defer db.Close()

	userHandler := handler.NewUserHandler(db)
	walletHandler := handler.NewWalletHandler(db)

	router := mux.NewRouter()

	router.HandleFunc("/users", userHandler.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets", walletHandler.CreateWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/topup", walletHandler.TopUpWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/pay", walletHandler.PayWithWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transfer/{country_code}", walletHandler.TransferWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/suspend", walletHandler.SuspendWalletHandler).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}", walletHandler.GetWalletHandler).Methods(http.MethodGet)
	router.HandleFunc("/wallet/update_close_balance", walletHandler.UpdateCloseBalanceHandler).Methods(http.MethodPost)

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

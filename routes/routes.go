package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maickmachado/blockchain-go/blockchain/controllers"
)

func HandleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/blockchain", controllers.CreateGenesisBlock).Methods("POST")
	myRouter.HandleFunc("/blockchain", controllers.GetAllData).Methods("GET")
	myRouter.HandleFunc("/blockchain/balance/{address}", controllers.GetBalanceAddress).Methods("GET")
	myRouter.HandleFunc("/blockchain/transactions", controllers.NewTransaction).Methods("POST")
	// myRouter.HandleFunc("/api/pokemons/{id}", controllers.GetPokemonById).Methods("GET")
	// myRouter.HandleFunc("/api/pokemons/{id}", controllers.UpdateDataPokemon).Methods("PUT")
	// myRouter.HandleFunc("/api/pokemons/{id}", controllers.DeleteDataPokemon).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
	//log.Fatal(http.ListenAndServe(config.AppConfig.Port, myRouter))
	//2022/08/01 18:14:44 listen tcp: address 8080: missing port in address
	//exit status 1
}

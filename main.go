package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){
	

	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)


	log.Fatal(http.ListenAndServe(":4242", nil))
}

// http handlers
func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request){
	fmt.Println("Create Payment Intent")
	responseObj := []string{"Hello", "World"}
	w.Write(responseObj, )
}

// Timestamp: 2.17.38
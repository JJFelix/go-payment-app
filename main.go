package main

import (
	"bytes"
	"encoding/json"
	"io"

	// "errors"
	"fmt"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

var STRIPE_API_KEY = os.Getenv("STRIPE_SECRET_KEY")

func main(){
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	http.HandleFunc("/create-payment-intent", handleCreatePaymentIntent)
	http.HandleFunc("/health", handleHealth)

	log.Println("listening on http://localhost:4242...")
	log.Fatal(http.ListenAndServe(":4242", nil))
}


// http handlers
func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request){

	if r.Method != "POST"{ // request method validation
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// fmt.Println("Create Payment Intent")

	var req struct{
		ProductId 	string `json:"product_id"`
		FirstName 	string `json:"first_name"`
		LastName 	string `json:"last_name"`
		Address1 	string `json:"address_1"`
		Address2 	string `json:"address_2"`
		City 		string `json:"city"`
		State 		string `json:"state"`
		Zip 		string `json:"zip"`
		Country 	string `json:"country"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		// log.Println(err) // DEBUG
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println(req.FirstName, req.LastName, req.City) // testing the request Body data

	// params to be passed to the stripe server to create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(calculateOrderAmount(req.ProductId)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	paymentIntent, err := paymentintent.New(params) // sends a request to stripe's API
	if err != nil{
		// log.Println(err) // DEBUG
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println(paymentIntent.ClientSecret) // DEBUG
	var response struct{
		ClientSecret string `json:"clientSecret"`
	}

	response.ClientSecret = paymentIntent.ClientSecret

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(response)
	if err != nil{
		// log.Println(err) // DEBUG
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, &buf)
	if err != nil{
		fmt.Println(err) // DEBUG
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
	}

}

func calculateOrderAmount(productId string) int64 {
	switch productId{
	case "Forever Pants" :
		return 26000
	case "Forever Shirt" :
		return 15500
	case "Forever Shorts" :
		return 30000
	}
	return 0
}

func handleHealth(writer http.ResponseWriter, request *http.Request){
	fmt.Println("Up and Running")
	// response is a slice of bytes
	response := []byte("Server is up and running") // convert a string message into slice of bytes

	_, err := writer.Write(response)
	if err != nil{
		fmt.Println(err)
	}
}
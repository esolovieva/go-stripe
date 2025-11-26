package main

import (
	"encoding/json"
	"net/http"
)

// stripePayload represents the expected incoming JSON payload
// for creating a Stripe payment intent, containing the currency
// and amount provided by the client.
type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// jsonResponse defines the structure of the JSON response returned
// by the server. It contains information about the request status,
// an optional message, additional content, and an optional resource ID.
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Content string `json:"content"`
	ID      int    `json:"id"`
}

// GetPaymentIntent handles HTTP requests for creating or retrieving
// a payment intent. It returns a JSON response indicating whether the
// request was successful. Currently, it returns a basic placeholder
// response without interacting with any payment provider.
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	j := jsonResponse{
		OK: true,
	}

	out, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

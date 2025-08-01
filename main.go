package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
	"github.com/stripe/stripe-go/v79/charge"
)

func main() {
	chargeID := flag.String("charge", "", "Stripe charge ID (e.g., ch_123abc)")
	flag.Parse()

	if *chargeID == "" {
		fmt.Println("Usage: stripefees-cli -charge <charge_id>")
		os.Exit(1)
	}

	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		log.Fatal("Error: STRIPE_SECRET_KEY environment variable is not set.")
	}
	stripe.Key = apiKey

	ch, err := charge.Get(*chargeID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve charge: %v", err)
	}

	bt, err := balancetransaction.Get(ch.BalanceTransaction.ID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve balance transaction: %v", err)
	}

	fmt.Printf("Charge ID: %s\n", ch.ID)
	fmt.Printf("Amount Charged: $%.2f\n", float64(ch.Amount)/100)
	fmt.Printf("Stripe Fee: $%.2f\n", float64(bt.Fee)/100)
	fmt.Printf("Net Amount: $%.2f\n", float64(bt.Net)/100)

	fmt.Println("\nFee Breakdown:")
	for _, fee := range bt.FeeDetails {
		fmt.Printf("- %s: $%.2f\n", fee.Type, float64(fee.Amount)/100)
	}
}
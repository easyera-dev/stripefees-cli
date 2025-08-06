package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
	"github.com/stripe/stripe-go/v79/charge"
)

type FeeDetail struct {
	Type   string
	Amount int64
}

type ChargeInfo struct {
	ID     string
	Amount int64
	Fee    int64
	Net    int64
	FeeDetails []*FeeDetail
}

func getCharge(chargeID string) (*stripe.Charge, error) {
	if chargeID == "" {
		return getLatestCharge()
	}
	return charge.Get(chargeID, nil)
}

func getLatestCharge() (*stripe.Charge, error) {
	params := &stripe.ChargeListParams{}
	params.Limit = stripe.Int64(1)
	i := charge.List(params)
	
	if i.Next() {
		ch := i.Charge()
		fmt.Printf("No charge ID provided, using latest charge: %s\n", ch.ID)
		return ch, nil
	}
	
	if i.Err() != nil {
		return nil, i.Err()
	}
	
	return nil, errors.New("no charges found in your Stripe account")
}

func getChargeInfo(ch *stripe.Charge) (*ChargeInfo, error) {
	bt, err := balancetransaction.Get(ch.BalanceTransaction.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve balance transaction: %v", err)
	}

	feeDetails := make([]*FeeDetail, len(bt.FeeDetails))
	for i, fee := range bt.FeeDetails {
		feeDetails[i] = &FeeDetail{
			Type:   fee.Type,
			Amount: fee.Amount,
		}
	}

	return &ChargeInfo{
		ID:         ch.ID,
		Amount:     ch.Amount,
		Fee:        bt.Fee,
		Net:        bt.Net,
		FeeDetails: feeDetails,
	}, nil
}

func displayChargeInfo(info *ChargeInfo) {
	fmt.Printf("Charge ID: %s\n", info.ID)
	fmt.Printf("Amount Charged: $%.2f\n", float64(info.Amount)/100)
	fmt.Printf("Stripe Fee: $%.2f\n", float64(info.Fee)/100)
	fmt.Printf("Net Amount: $%.2f\n", float64(info.Net)/100)

	fmt.Println("\nFee Breakdown:")
	for _, fee := range info.FeeDetails {
		fmt.Printf("- %s: $%.2f\n", fee.Type, float64(fee.Amount)/100)
	}
}

func main() {
	chargeID := flag.String("charge", "", "Stripe charge ID (e.g., ch_123abc). If not provided, uses the latest charge.")
	flag.Parse()

	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		log.Fatal("Error: STRIPE_SECRET_KEY environment variable is not set.")
	}
	stripe.Key = apiKey

	ch, err := getCharge(*chargeID)
	if err != nil {
		log.Fatalf("Failed to retrieve charge: %v", err)
	}

	info, err := getChargeInfo(ch)
	if err != nil {
		log.Fatalf("Failed to get charge info: %v", err)
	}

	displayChargeInfo(info)
}
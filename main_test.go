package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestChargeInfo_Creation(t *testing.T) {
	feeDetails := []*FeeDetail{
		{Type: "stripe_fee", Amount: 30},
		{Type: "application_fee", Amount: 29},
	}

	info := &ChargeInfo{
		ID:         "ch_test123",
		Amount:     1000,
		Fee:        59,
		Net:        941,
		FeeDetails: feeDetails,
	}

	if info.ID != "ch_test123" {
		t.Errorf("Expected ID to be 'ch_test123', got %s", info.ID)
	}
	if info.Amount != 1000 {
		t.Errorf("Expected Amount to be 1000, got %d", info.Amount)
	}
	if info.Fee != 59 {
		t.Errorf("Expected Fee to be 59, got %d", info.Fee)
	}
	if info.Net != 941 {
		t.Errorf("Expected Net to be 941, got %d", info.Net)
	}
	if len(info.FeeDetails) != 2 {
		t.Errorf("Expected 2 fee details, got %d", len(info.FeeDetails))
	}
}

func TestDisplayChargeInfo(t *testing.T) {
	feeDetails := []*FeeDetail{
		{Type: "stripe_fee", Amount: 30},
		{Type: "application_fee", Amount: 29},
	}

	info := &ChargeInfo{
		ID:         "ch_test123",
		Amount:     1000,
		Fee:        59,
		Net:        941,
		FeeDetails: feeDetails,
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayChargeInfo(info)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedOutputs := []string{
		"Charge ID: ch_test123",
		"Amount Charged: $10.00",
		"Stripe Fee: $0.59",
		"Net Amount: $9.41",
		"Fee Breakdown:",
		"- stripe_fee: $0.30",
		"- application_fee: $0.29",
	}

	for _, expected := range expectedOutputs {
		if !contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it didn't. Full output:\n%s", expected, output)
		}
	}
}

func TestAmountFormatting(t *testing.T) {
	testCases := []struct {
		amount   int64
		expected string
	}{
		{1000, "$10.00"},
		{59, "$0.59"},
		{941, "$9.41"},
		{0, "$0.00"},
		{1, "$0.01"},
		{100000, "$1000.00"},
	}

	for _, tc := range testCases {
		result := fmt.Sprintf("$%.2f", float64(tc.amount)/100)
		if result != tc.expected {
			t.Errorf("For amount %d, expected %s, got %s", tc.amount, tc.expected, result)
		}
	}
}

func TestChargeInfoStructure(t *testing.T) {
	info := &ChargeInfo{}
	
	if info.ID != "" {
		t.Errorf("Expected empty ID for new ChargeInfo, got %s", info.ID)
	}
	if info.Amount != 0 {
		t.Errorf("Expected zero Amount for new ChargeInfo, got %d", info.Amount)
	}
	if info.Fee != 0 {
		t.Errorf("Expected zero Fee for new ChargeInfo, got %d", info.Fee)
	}
	if info.Net != 0 {
		t.Errorf("Expected zero Net for new ChargeInfo, got %d", info.Net)
	}
	if info.FeeDetails != nil {
		t.Errorf("Expected nil FeeDetails for new ChargeInfo, got %v", info.FeeDetails)
	}
}

func TestFeeDetailsHandling(t *testing.T) {
	info := &ChargeInfo{
		ID:         "ch_test",
		Amount:     1000,
		Fee:        30,
		Net:        970,
		FeeDetails: []*FeeDetail{},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayChargeInfo(info)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !contains(output, "Fee Breakdown:") {
		t.Error("Expected 'Fee Breakdown:' to be displayed even with empty fee details")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkDisplayChargeInfo(b *testing.B) {
	feeDetails := []*FeeDetail{
		{Type: "stripe_fee", Amount: 30},
		{Type: "application_fee", Amount: 29},
	}

	info := &ChargeInfo{
		ID:         "ch_benchmark",
		Amount:     1000,
		Fee:        59,
		Net:        941,
		FeeDetails: feeDetails,
	}

	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		displayChargeInfo(info)
	}
}
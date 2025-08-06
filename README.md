# Stripe Fees CLI

A command-line tool to retrieve and display Stripe fees for a given charge ID.

## Features

- Fetch charge details from Stripe
- Display the original charge amount, Stripe fees, and net amount
- Show detailed fee breakdown by type

## Prerequisites

- Go 1.23.1 or higher
- Stripe secret key

## Installation

1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd stripefees-cli
   ```

2. Build the application:
   ```bash
   go build -o stripefees-cli main.go
   ```

## Testing

Run the test suite:

```bash
go test
```

Run tests with verbose output:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

The test suite includes:
- Unit tests for the `ChargeInfo` struct
- Tests for amount formatting and display logic
- Tests for fee details handling
- Benchmarks for performance testing

## Configuration

Set your Stripe secret key as an environment variable:

```bash
export STRIPE_SECRET_KEY="sk_test_your_stripe_secret_key_here"
```

## Usage

```bash
./stripefees-cli [-charge <charge_id>]
```

If no charge ID is provided, the tool will automatically use the latest charge from your Stripe account.

### Examples

With a specific charge ID:
```bash
./stripefees-cli -charge ch_1234567890abcdef
```

Using the latest charge:
```bash
./stripefees-cli
```

Output:
```
Charge ID: ch_1234567890abcdef
Amount Charged: $10.00
Stripe Fee: $0.59
Net Amount: $9.41

Fee Breakdown:
- stripe_fee: $0.30
- application_fee: $0.29
```

## Parameters

- `-charge`: Optional. The Stripe charge ID to analyze (format: `ch_xxxxxx`). If not provided, uses the latest charge.

## Error Handling

The tool will exit with an error if:
- The `STRIPE_SECRET_KEY` environment variable is not set
- The charge ID is invalid or not found
- No charges exist in your Stripe account (when using default behavior)
- There's an issue connecting to the Stripe API

## Dependencies

- [stripe-go](https://github.com/stripe/stripe-go) v79.12.0
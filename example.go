package solanatracker

import (
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	startTime := time.Now()
	privateKey := "YOUR_SECRET_KEY" // replace with your base58 private key

	// Create a keypair from the secret key
	keypair, err := solana.PrivateKeyFromBase58(privateKey)
	if err != nil {
		log.Fatalf("Error creating keypair: %v", err)
	}

	rpcUrl := "https://solana-rpc.publicnode.com"

	// Initialize a new Solana tracker with the keypair and RPC endpoint
	tracker := NewSolanaTracker(keypair, rpcUrl)

	priorityFee := 0.00005 // priorityFee requires a pointer, thus we store it in a variable

	// Get the swap instructions for the specified tokens, amounts, and other parameters
	swapResponse, err := tracker.GetSwapInstructions(
		"So11111111111111111111111111111111111111112",
		"4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
		0.0001,                       // Amount to swap
		30,                           // Slippage
		keypair.PublicKey().String(), // Payer public key
		&priorityFee,                 // Priority fee (Recommended while network is congested)
		false,
	)
	if err != nil {
		// Log and exit if there's an error getting the swap instructions
		log.Fatalf("Error getting swap instructions: %v", err)
	}

	maxRetries := uint(5) // maxRetries requires a pointer, thus we store it in a variable

	// Define the options for the swap transaction
	options := SwapOptions{
		SendOptions: rpc.TransactionOpts{
			SkipPreflight: true,
			MaxRetries:    &maxRetries,
		},
		ConfirmationRetries:        50,
		ConfirmationRetryTimeout:   1000 * time.Millisecond,
		LastValidBlockHeightBuffer: 200,
		Commitment:                 rpc.CommitmentProcessed,
		ResendInterval:             1500 * time.Millisecond,
		ConfirmationCheckInterval:  100 * time.Millisecond,
		SkipConfirmationCheck:      false,
	}

	// Perform the swap transaction with the specified options
	sendTime := time.Now()
	txid, err := tracker.PerformSwap(swapResponse, options)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime).Seconds()
	if err != nil {
		fmt.Printf("Swap failed: %v\n", err)
		fmt.Printf("Time elapsed before failure: %.2f seconds\n", elapsedTime)
		// Add retries or additional error handling as needed
	} else {
		fmt.Printf("Transaction ID: %s\n", txid)
		fmt.Printf("Transaction URL: https://solscan.io/tx/%s\n", txid)
		fmt.Printf("Swap completed in %.2f seconds\n", elapsedTime)
		fmt.Printf("Transaction finished in %.2f seconds\n", endTime.Sub(sendTime).Seconds())
	}
}

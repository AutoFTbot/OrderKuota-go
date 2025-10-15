package qris

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PaymentStatus stores the payment status information.
// PaymentStatus menyimpan informasi status pembayaran.
type PaymentStatus struct {
	Status    string // Payment status (PAID/UNPAID) / Status pembayaran (PAID/UNPAID)
	Amount    int64  // Payment amount / Nominal pembayaran
	Reference string // Payment reference / Referensi pembayaran
	Date      string // Payment date (if PAID) / Tanggal pembayaran (jika PAID)
	BrandName string // Payer brand name (if PAID) / Nama brand pembayar (jika PAID)
	BuyerRef  string // Buyer reference (if PAID) / Referensi pembeli (jika PAID)
}

// PaymentCheckerConfig stores the configuration for payment checking.
// PaymentCheckerConfig menyimpan konfigurasi untuk pengecekan pembayaran.
type PaymentCheckerConfig struct {
	MerchantID string // Merchant ID from payment gateway / ID merchant dari payment gateway
	APIKey     string // API key for authentication / API key untuk autentikasi
	BaseURL    string // Base URL for API calls / URL dasar untuk panggilan API
}

// PaymentChecker is the main struct for payment checking operations.
// PaymentChecker adalah struct utama untuk operasi pengecekan pembayaran.
type PaymentChecker struct {
	config PaymentCheckerConfig
	client *http.Client
}

// NewPaymentChecker creates a new instance of PaymentChecker.
// NewPaymentChecker membuat instance baru dari PaymentChecker.
func NewPaymentChecker(config PaymentCheckerConfig) *PaymentChecker {
	return &PaymentChecker{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckPaymentStatus checks the payment status for a given reference and amount.
// CheckPaymentStatus mengecek status pembayaran untuk referensi dan nominal tertentu.
//
// It returns a PaymentStatus struct containing the payment information.
// Fungsi ini mengembalikan struct PaymentStatus yang berisi informasi pembayaran.
func (q *QRIS) CheckPaymentStatus(reference string, amount int64) (*PaymentStatus, error) {
	if reference == "" || amount <= 0 {
		return nil, fmt.Errorf("reference and amount must be filled correctly / reference dan amount harus diisi dengan benar")
	}

	// Create API URL
	url := "https://ftvpn.me/api/mutasi"
	log.Printf("Checking payment status for amount: %d", amount)

	// Create request body
	requestBody := map[string]string{
		"auth_token":    q.config.AuthToken,
		"auth_username": q.config.AuthUsername,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body / gagal marshal request body: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request / gagal membuat request: %v", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request / gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response / gagal membaca response: %v", err)
	}

	// Parse response
	var response struct {
		Status string `json:"status"`
		Data   []struct {
			Amount      string `json:"amount"`
			Date        string `json:"date"`
			QRIS        string `json:"qris"`
			Type        string `json:"type"`
			IssuerRef   string `json:"issuer_reff"`
			BrandName   string `json:"brand_name"`
			BuyerRef    string `json:"buyer_reff"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response / gagal parse response: %v", err)
	}

	if response.Status != "success" || len(response.Data) == 0 {
		return &PaymentStatus{
			Status:    "UNPAID",
			Amount:    amount,
			Reference: reference,
		}, nil
	}

	// Find matching transactions
	var matchingTransactions []struct {
		Amount      string `json:"amount"`
		Date        string `json:"date"`
		QRIS        string `json:"qris"`
		Type        string `json:"type"`
		IssuerRef   string `json:"issuer_reff"`
		BrandName   string `json:"brand_name"`
		BuyerRef    string `json:"buyer_reff"`
	}

	now := time.Now()
	for _, tx := range response.Data {
		txAmount, _ := strconv.Atoi(tx.Amount)
		
		// Parse transaction date
		txDate, err := time.Parse("2006-01-02 15:04:05", tx.Date)
		if err != nil {
			continue
		}

		timeDiff := now.Sub(txDate)
		
		// Check if transaction matches criteria
		if int64(txAmount) == amount &&
			tx.QRIS == "static" &&
			tx.Type == "CR" &&
			timeDiff <= 5*time.Minute {
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	if len(matchingTransactions) > 0 {
		// Get latest transaction
		latestTx := matchingTransactions[0]
		latestDate, _ := time.Parse("2006-01-02 15:04:05", latestTx.Date)
		
		for _, tx := range matchingTransactions[1:] {
			txDate, _ := time.Parse("2006-01-02 15:04:05", tx.Date)
			if txDate.After(latestDate) {
				latestTx = tx
				latestDate = txDate
			}
		}

		txAmount, _ := strconv.Atoi(latestTx.Amount)
		log.Printf("Payment found: Amount=%d, Date=%s, Brand=%s", 
			txAmount, latestTx.Date, latestTx.BrandName)

		return &PaymentStatus{
			Status:    "PAID",
			Amount:    int64(txAmount),
			Reference: latestTx.IssuerRef,
			Date:      latestTx.Date,
			BrandName: latestTx.BrandName,
			BuyerRef:  latestTx.BuyerRef,
		}, nil
	}

	log.Printf("No matching payment found for amount: %d", amount)
	return &PaymentStatus{
		Status:    "UNPAID",
		Amount:    amount,
		Reference: reference,
	}, nil
}
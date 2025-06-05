package qris

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// PaymentStatus menyimpan status pembayaran
type PaymentStatus struct {
	Status    string // Status pembayaran (PAID/UNPAID)
	Amount    int64  // Nominal pembayaran
	Reference string // Referensi pembayaran
	Date      string // Tanggal pembayaran (jika PAID)
	BrandName string // Nama brand pembayar (jika PAID)
	BuyerRef  string // Referensi pembeli (jika PAID)
}

// PaymentCheckerConfig menyimpan konfigurasi untuk pengecekan pembayaran
type PaymentCheckerConfig struct {
	MerchantID string
	APIKey     string
	BaseURL    string
}

// PaymentChecker adalah struct untuk mengecek status pembayaran
type PaymentChecker struct {
	config PaymentCheckerConfig
	client *http.Client
}

// NewPaymentChecker membuat instance PaymentChecker baru
func NewPaymentChecker(config PaymentCheckerConfig) *PaymentChecker {
	return &PaymentChecker{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckPaymentStatus mengecek status pembayaran
func (q *QRIS) CheckPaymentStatus(reference string, amount int64) (*PaymentStatus, error) {
	if reference == "" || amount <= 0 {
		return nil, fmt.Errorf("reference dan amount harus diisi dengan benar")
	}

	// Buat URL untuk request
	url := fmt.Sprintf("https://gateway.okeconnect.com/api/mutasi/qris/%s/%s", q.config.MerchantID, q.config.APIKey)
	log.Printf("Checking payment status at URL: %s", url)

	// Buat request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %v", err)
	}

	// Kirim request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request: %v", err)
	}
	defer resp.Body.Close()

	// Baca response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response: %v", err)
	}
	log.Printf("Raw response: %s", string(body))

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
		return nil, fmt.Errorf("gagal parse response: %v", err)
	}

	log.Printf("Parsed response status: %s", response.Status)
	log.Printf("Number of transactions: %d", len(response.Data))

	if response.Status != "success" || len(response.Data) == 0 {
		return &PaymentStatus{
			Status:    "UNPAID",
			Amount:    amount,
			Reference: reference,
		}, nil
	}

	// Cari transaksi yang sesuai
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
		log.Printf("Checking transaction: Amount=%s, Date=%s, QRIS=%s, Type=%s", 
			tx.Amount, tx.Date, tx.QRIS, tx.Type)

		// Coba parse tanggal dengan beberapa format
		var txDate time.Time
		var parseErr error
		dateFormats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
		}

		for _, format := range dateFormats {
			txDate, parseErr = time.Parse(format, tx.Date)
			if parseErr == nil {
				break
			}
		}

		if parseErr != nil {
			log.Printf("Warning: Could not parse date %s with any format", tx.Date)
			continue
		}

		timeDiff := now.Sub(txDate)
		log.Printf("Transaction time difference: %v", timeDiff)

		if int64(txAmount) == amount &&
			tx.QRIS == "static" &&
			tx.Type == "CR" &&
			timeDiff <= 5*time.Minute {

			log.Printf("Found matching transaction!")
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	if len(matchingTransactions) > 0 {
		// Ambil transaksi terbaru
		latestTx := matchingTransactions[0]
		latestDate, _ := time.Parse(time.RFC3339, latestTx.Date)
		
		for _, tx := range matchingTransactions[1:] {
			txDate, _ := time.Parse(time.RFC3339, tx.Date)
			if txDate.After(latestDate) {
				latestTx = tx
				latestDate = txDate
			}
		}

		txAmount, _ := strconv.Atoi(latestTx.Amount)
		return &PaymentStatus{
			Status:    "PAID",
			Amount:    int64(txAmount),
			Reference: latestTx.IssuerRef,
			Date:      latestTx.Date,
			BrandName: latestTx.BrandName,
			BuyerRef:  latestTx.BuyerRef,
		}, nil
	}

	return &PaymentStatus{
		Status:    "UNPAID",
		Amount:    amount,
		Reference: reference,
	}, nil
}
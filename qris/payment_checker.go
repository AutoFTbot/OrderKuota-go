package qris

import (
	"encoding/json"
	"fmt"
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

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("gagal parse response: %v", err)
	}

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

	for _, tx := range response.Data {
		txAmount, _ := strconv.ParseInt(tx.Amount, 10, 64)
		txDate, _ := time.Parse(time.RFC3339, tx.Date)
		timeDiff := time.Since(txDate)

		if txAmount == amount &&
			tx.QRIS == "static" &&
			tx.Type == "CR" &&
			timeDiff <= 30*time.Minute { // Ubah ke 30 menit

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

		txAmount, _ := strconv.ParseInt(latestTx.Amount, 10, 64)
		return &PaymentStatus{
			Status:    "PAID",
			Amount:    txAmount,
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
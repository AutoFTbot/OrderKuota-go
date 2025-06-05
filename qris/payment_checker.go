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
	Success bool        `json:"success"`
	Data    *StatusData `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// StatusData menyimpan detail status pembayaran
type StatusData struct {
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	Reference  string  `json:"reference"`
	Date       string  `json:"date,omitempty"`
	BrandName  string  `json:"brand_name,omitempty"`
	BuyerRef   string  `json:"buyer_reff,omitempty"`
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
func (p *PaymentChecker) CheckPaymentStatus(reference string, amount float64) (*PaymentStatus, error) {
	if reference == "" || amount <= 0 {
		return &PaymentStatus{
			Success: false,
			Error:   "Reference dan amount harus diisi dengan benar",
		}, nil
	}

	// Buat URL untuk request
	url := fmt.Sprintf("%s/api/mutasi/qris/%s/%s", p.config.BaseURL, p.config.MerchantID, p.config.APIKey)

	// Buat request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &PaymentStatus{
			Success: false,
			Error:   fmt.Sprintf("Gagal membuat request: %v", err),
		}, nil
	}

	// Kirim request
	resp, err := p.client.Do(req)
	if err != nil {
		return &PaymentStatus{
			Success: false,
			Error:   fmt.Sprintf("Gagal mengirim request: %v", err),
		}, nil
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
		return &PaymentStatus{
			Success: false,
			Error:   fmt.Sprintf("Gagal parse response: %v", err),
		}, nil
	}

	if response.Status != "success" || len(response.Data) == 0 {
		return &PaymentStatus{
			Success: true,
			Data: &StatusData{
				Status:    "UNPAID",
				Amount:    amount,
				Reference: reference,
			},
		}, nil
	}

	// Cari transaksi yang sesuai
	for _, tx := range response.Data {
		txAmount, _ := strconv.ParseFloat(tx.Amount, 64)
		txDate, _ := time.Parse(time.RFC3339, tx.Date)
		timeDiff := time.Since(txDate)

		if txAmount == amount &&
			tx.QRIS == "static" &&
			tx.Type == "CR" &&
			timeDiff <= 5*time.Minute {

			return &PaymentStatus{
				Success: true,
				Data: &StatusData{
					Status:    "PAID",
					Amount:    txAmount,
					Reference: tx.IssuerRef,
					Date:      tx.Date,
					BrandName: tx.BrandName,
					BuyerRef:  tx.BuyerRef,
				},
			}, nil
		}
	}

	return &PaymentStatus{
		Success: true,
		Data: &StatusData{
			Status:    "UNPAID",
			Amount:    amount,
			Reference: reference,
		},
	}, nil
} 
package qris

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	// BaseURL untuk payment gateway
	BaseURL = "https://gateway.okeconnect.com"
)

// QRISData menyimpan data untuk QRIS
type QRISData struct {
	Amount         float64
	TransactionID  string
}

// QRISConfig menyimpan konfigurasi untuk QRIS
type QRISConfig struct {
	MerchantID   string
	APIKey       string
	BaseQrString string
}

// QRIS adalah struct utama untuk package
type QRIS struct {
	config QRISConfig
	client *http.Client
}

// NewQRIS membuat instance QRIS baru
func NewQRIS(config QRISConfig) *QRIS {
	return &QRIS{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GenerateQRCode menghasilkan QR Code QRIS
func (q *QRIS) GenerateQRCode(data QRISData) ([]byte, error) {
	qrString := q.generateQRISString(data)
	
	// Generate QR Code
	qr, err := qrcode.Encode(qrString, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("gagal generate QR code: %v", err)
	}

	return qr, nil
}

// generateQRISString menghasilkan string QRIS sesuai standar
func (q *QRIS) generateQRISString(data QRISData) string {
	if q.config.BaseQrString != "" {
		// Jika baseQrString disediakan, gunakan itu sebagai dasar
		return q.generateQRISStringFromBase(data)
	}

	// Format: ID + Panjang + Data
	qrString := "000201010212"
	
	// ID Merchant
	qrString += "2937" + fmt.Sprintf("%02d", len(q.config.MerchantID)) + q.config.MerchantID
	
	// Amount
	if data.Amount > 0 {
		amountStr := fmt.Sprintf("%.2f", data.Amount)
		qrString += "5406" + amountStr
	}
	
	// Transaction ID
	if data.TransactionID != "" {
		qrString += "5802" + data.TransactionID
	}
	
	// Generate CRC
	crc := q.generateCRC(qrString)
	qrString += crc
	
	return qrString
}

// generateQRISStringFromBase menghasilkan string QRIS dari baseQrString
func (q *QRIS) generateQRISStringFromBase(data QRISData) string {
	if !q.config.BaseQrString.includes("5802ID") {
		return q.config.BaseQrString
	}

	qrString := q.config.BaseQrString
	if data.Amount > 0 {
		amountStr := fmt.Sprintf("%.2f", data.Amount)
		amountTag := "54" + fmt.Sprintf("%02d", len(amountStr)) + amountStr
		insertPos := qrString.indexOf("5802ID")
		qrString = qrString[:insertPos] + amountTag + qrString[insertPos:]
	}

	return qrString
}

// generateCRC menghasilkan CRC untuk QRIS string menggunakan CRC16-CCITT
func (q *QRIS) generateCRC(data string) string {
	// Implementasi CRC16-CCITT (0xFFFF)
	crc := uint16(0xFFFF)
	polynomial := uint16(0x1021)

	for _, b := range []byte(data) {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc = crc << 1
			}
		}
	}

	// Konversi ke hex string dan pastikan 4 digit
	return fmt.Sprintf("%04X", crc)
}

// ValidateQRISString memvalidasi string QRIS
func (q *QRIS) ValidateQRISString(qrString string) bool {
	// Validasi panjang minimal
	if len(qrString) < 20 {
		return false
	}

	// Validasi format dasar
	if qrString[:12] != "000201010212" {
		return false
	}

	// Validasi ID Merchant
	if !strings.Contains(qrString, "2937") {
		return false
	}

	// Validasi format amount jika ada
	if strings.Contains(qrString, "54") {
		amountIndex := strings.Index(qrString, "54")
		if amountIndex+2 >= len(qrString) {
			return false
		}
		amountLen, err := strconv.Atoi(qrString[amountIndex+2 : amountIndex+4])
		if err != nil || amountIndex+4+amountLen > len(qrString) {
			return false
		}
	}

	// Validasi CRC
	if len(qrString) < 4 {
		return false
	}
	calculatedCRC := q.generateCRC(qrString[:len(qrString)-4])
	providedCRC := qrString[len(qrString)-4:]
	if calculatedCRC != providedCRC {
		return false
	}

	return true
}

// GenerateTransactionID menghasilkan ID transaksi unik
func (q *QRIS) GenerateTransactionID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("TRX%d", timestamp)
}

// CheckPaymentStatus mengecek status pembayaran
func (q *QRIS) CheckPaymentStatus(reference string, amount float64) (*PaymentStatus, error) {
	if reference == "" || amount <= 0 {
		return &PaymentStatus{
			Success: false,
			Error:   "Reference dan amount harus diisi dengan benar",
		}, nil
	}

	// Buat URL untuk request
	url := fmt.Sprintf("%s/api/mutasi/qris/%s/%s", BaseURL, q.config.MerchantID, q.config.APIKey)

	// Buat request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &PaymentStatus{
			Success: false,
			Error:   fmt.Sprintf("Gagal membuat request: %v", err),
		}, nil
	}

	// Kirim request
	resp, err := q.client.Do(req)
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
package qris

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

// QRISConfig menyimpan konfigurasi untuk QRIS
type QRISConfig struct {
	MerchantID   string // ID merchant dari payment gateway
	APIKey       string // API key untuk autentikasi
	BaseQrString string // Base QRIS string dari merchant
}

// QRISData menyimpan data untuk generate QR code
type QRISData struct {
	Amount        int64  // Nominal pembayaran
	TransactionID string // ID transaksi unik
}

// PaymentStatus menyimpan status pembayaran
type PaymentStatus struct {
	Status    string // Status pembayaran (PAID/UNPAID)
	Amount    int64  // Nominal pembayaran
	Reference string // Referensi pembayaran
	Date      string // Tanggal pembayaran (jika PAID)
	BrandName string // Nama brand pembayar (jika PAID)
	BuyerRef  string // Referensi pembeli (jika PAID)
}

// QRIS adalah struct utama untuk operasi QRIS
type QRIS struct {
	config QRISConfig
}

// NewQRIS membuat instance baru dari QRIS
func NewQRIS(config QRISConfig) *QRIS {
	return &QRIS{
		config: config,
	}
}

// GenerateQRCode menghasilkan QR code untuk pembayaran QRIS
func (q *QRIS) GenerateQRCode(data QRISData) (*qrcode.QRCode, error) {
	if data.Amount <= 0 {
		return nil, errors.New("nominal harus lebih besar dari 0")
	}

	// Generate QRIS string
	qrString, err := q.generateQRISString(data)
	if err != nil {
		return nil, err
	}

	// Generate QR code
	qrCode, err := qrcode.New(qrString, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("gagal generate QR code: %v", err)
	}

	return qrCode, nil
}

// generateQRISString menghasilkan string QRIS sesuai format standar
func (q *QRIS) generateQRISString(data QRISData) (string, error) {
	if !strings.Contains(q.config.BaseQrString, "5802ID") {
		return "", errors.New("format QRIS tidak valid")
	}

	// Format nominal
	amountStr := fmt.Sprintf("%d", data.Amount)
	amountTag := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)

	// Insert nominal ke base string
	parts := strings.Split(q.config.BaseQrString, "5802ID")
	if len(parts) != 2 {
		return "", errors.New("format QRIS tidak valid")
	}

	qrString := parts[0] + amountTag + "5802ID" + parts[1]

	// Generate CRC
	crc := q.generateCRC(qrString)
	return qrString + crc, nil
}

// generateCRC menghasilkan CRC16-CCITT untuk string QRIS
func (q *QRIS) generateCRC(data string) string {
	var crc uint16 = 0xFFFF
	for i := 0; i < len(data); i++ {
		crc ^= uint16(data[i]) << 8
		for j := 0; j < 8; j++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	return fmt.Sprintf("%04X", crc)
}

// CheckPaymentStatus mengecek status pembayaran
func (q *QRIS) CheckPaymentStatus(reference string, amount int64) (*PaymentStatus, error) {
	// TODO: Implementasi cek status pembayaran ke payment gateway
	// Untuk sementara return status UNPAID
	return &PaymentStatus{
		Status:    "UNPAID",
		Amount:    amount,
		Reference: reference,
	}, nil
}

// ValidateQRISString memvalidasi format string QRIS
func (q *QRIS) ValidateQRISString(qrString string) bool {
	if len(qrString) < 20 {
		return false
	}

	// Validasi format dasar
	if !strings.Contains(qrString, "5802ID") {
		return false
	}

	// Validasi Merchant ID
	if !strings.Contains(qrString, q.config.MerchantID) {
		return false
	}

	// Validasi amount format
	if !strings.Contains(qrString, "54") {
		return false
	}

	// Validasi CRC
	crc := q.generateCRC(qrString[:len(qrString)-4])
	return crc == qrString[len(qrString)-4:]
} 
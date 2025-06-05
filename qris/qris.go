// Package qris provides QRIS (Quick Response Code Indonesian Standard) payment integration for Go applications.
// Package qris menyediakan integrasi pembayaran QRIS (Quick Response Code Indonesian Standard) untuk aplikasi Go.
package qris

import (
	"errors"
	"fmt"
	"strings"
	"github.com/skip2/go-qrcode"
	"image/color"
)

// QRISConfig stores the configuration for QRIS operations.
// QRISConfig menyimpan konfigurasi untuk operasi QRIS.
type QRISConfig struct {
	MerchantID   string // Merchant ID from payment gateway / ID merchant dari payment gateway
	APIKey       string // API key for authentication / API key untuk autentikasi
	BaseQrString string // Base QRIS string from merchant / Base QRIS string dari merchant
}

// QRISData stores the data needed to generate a QR code.
// QRISData menyimpan data yang diperlukan untuk generate QR code.
type QRISData struct {
	Amount        int64  // Payment amount / Nominal pembayaran
	TransactionID string // Unique transaction ID / ID transaksi unik
}

// QRIS is the main struct for QRIS operations.
// QRIS adalah struct utama untuk operasi QRIS.
type QRIS struct {
	config QRISConfig
}

// NewQRIS creates a new instance of QRIS.
// NewQRIS membuat instance baru dari QRIS.
//
// It validates the configuration and returns an error if the configuration is invalid.
// Fungsi ini memvalidasi konfigurasi dan mengembalikan error jika konfigurasi tidak valid.
func NewQRIS(config QRISConfig) (*QRIS, error) {
	if config.MerchantID == "" || config.APIKey == "" || config.BaseQrString == "" {
		return nil, errors.New("merchantID, apiKey, and baseQrString must be filled / merchantID, apiKey, dan baseQrString harus diisi")
	}

	if !strings.Contains(config.BaseQrString, "5802ID") {
		return nil, errors.New("invalid baseQrString format / format baseQrString tidak valid")
	}

	return &QRIS{
		config: config,
	}, nil
}

// GenerateQRCode generates a QR code for QRIS payment.
// GenerateQRCode menghasilkan QR code untuk pembayaran QRIS.
//
// It returns a QR code that can be saved as an image file.
// Fungsi ini mengembalikan QR code yang dapat disimpan sebagai file gambar.
func (q *QRIS) GenerateQRCode(data QRISData) (*qrcode.QRCode, error) {
	if data.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0 / nominal harus lebih besar dari 0")
	}

	if data.TransactionID == "" {
		return nil, errors.New("transactionID must be filled / transactionID harus diisi")
	}

	// Generate QRIS string
	qrString, err := q.generateQRISString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QRIS string / gagal generate QRIS string: %v", err)
	}

	// Generate QR code with high error correction level
	qrCode, err := qrcode.New(qrString, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code / gagal generate QR code: %v", err)
	}

	// Set QR code options
	qrCode.DisableBorder = false
	qrCode.ForegroundColor = color.Black
	qrCode.BackgroundColor = color.White

	return qrCode, nil
}

// generateQRISString generates a QRIS string according to the standard format.
// generateQRISString menghasilkan string QRIS sesuai format standar.
func (q *QRIS) generateQRISString(data QRISData) (string, error) {
	// Format amount
	amountStr := fmt.Sprintf("%d", data.Amount)
	amountTag := fmt.Sprintf("54%02d%s", len(amountStr), amountStr)

	// Remove existing CRC and replace 010211 with 010212
	baseString := q.config.BaseQrString[:len(q.config.BaseQrString)-4]
	baseString = strings.Replace(baseString, "010211", "010212", 1)

	// Insert amount into base string
	insertPosition := strings.Index(baseString, "5802ID")
	if insertPosition == -1 {
		return "", errors.New("invalid QRIS format: country ID not found / format QRIS tidak valid: ID negara tidak ditemukan")
	}

	qrString := baseString[:insertPosition] + amountTag + baseString[insertPosition:]

	// Generate CRC
	crc := q.generateCRC(qrString)
	return qrString + crc, nil
}

// generateCRC generates CRC16-CCITT for QRIS string.
// generateCRC menghasilkan CRC16-CCITT untuk string QRIS.
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

// ValidateQRISString validates the QRIS string format.
// ValidateQRISString memvalidasi format string QRIS.
//
// It checks the string length, country ID, merchant ID, amount format, and CRC.
// Fungsi ini memeriksa panjang string, ID negara, ID merchant, format nominal, dan CRC.
func (q *QRIS) ValidateQRISString(qrString string) error {
	if len(qrString) < 20 {
		return errors.New("QRIS string too short / string QRIS terlalu pendek")
	}

	// Basic format validation
	if !strings.Contains(qrString, "5802ID") {
		return errors.New("invalid QRIS format: country ID not found / format QRIS tidak valid: ID negara tidak ditemukan")
	}

	// Merchant ID validation
	if !strings.Contains(qrString, q.config.MerchantID) {
		return errors.New("merchant ID mismatch / merchant ID tidak sesuai")
	}

	// Amount format validation
	if !strings.Contains(qrString, "54") {
		return errors.New("invalid amount format / format nominal tidak valid")
	}

	// CRC validation
	crc := q.generateCRC(qrString[:len(qrString)-4])
	if crc != qrString[len(qrString)-4:] {
		return errors.New("invalid checksum / checksum tidak valid")
	}

	return nil
}

// GetQRISString generates a QRIS string without creating a QR code.
// GetQRISString menghasilkan string QRIS tanpa membuat QR code.
//
// It's useful when you only need the QRIS string for other purposes.
// Fungsi ini berguna ketika Anda hanya membutuhkan string QRIS untuk keperluan lain.
func (q *QRIS) GetQRISString(data QRISData) (string, error) {
	if data.Amount <= 0 {
		return "", errors.New("amount must be greater than 0 / nominal harus lebih besar dari 0")
	}

	if data.TransactionID == "" {
		return "", errors.New("transactionID must be filled / transactionID harus diisi")
	}

	return q.generateQRISString(data)
}
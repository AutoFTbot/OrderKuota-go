package main

import (
	"fmt"
	"log"

	"github.com/autoftbot/orderkuota-go/qris"
)

func main() {
	// Inisialisasi QRIS dengan konfigurasi
	config := qris.QRISConfig{
		MerchantID:   "OK2169948",
		APIKey:       "506151017388449542169948OKCTB751A34A2F8624E6A7B924038D5FE42A",
		BaseQrString: "00020101021126670016COM.NOBUBANK.WWW01189360050300000879140214158455875489000303UMI51440014ID.CO.QRIS.WWW0215ID20253762751400303UMI5204541153033605802ID5920AGIN STORE OK21699486006CIAMIS61054621162070703A0163049492",
	}

	// Buat instance QRIS
	qrisInstance := qris.NewQRIS(config)

	// Generate QR Code
	data := qris.QRISData{
		Amount:        1000,
		TransactionID: "TRX123",
	}

	qrCode, err := qrisInstance.GenerateQRCode(data)
	if err != nil {
		log.Fatalf("Error generating QR code: %v", err)
	}

	// Simpan QR code ke file
	err = qrCode.Save("qris.png")
	if err != nil {
		log.Fatalf("Error saving QR code: %v", err)
	}

	fmt.Println("QR Code berhasil dibuat dan disimpan sebagai qris.png")

	// Cek status pembayaran
	status, err := qrisInstance.CheckPaymentStatus("TRX123", 1000)
	if err != nil {
		log.Fatalf("Error checking payment status: %v", err)
	}

	fmt.Printf("Status Pembayaran: %s\n", status.Status)
	if status.Status == "PAID" {
		fmt.Printf("Amount: %d\n", status.Amount)
		fmt.Printf("Reference: %s\n", status.Reference)
		fmt.Printf("Date: %s\n", status.Date)
		fmt.Printf("Brand: %s\n", status.BrandName)
	}
} 
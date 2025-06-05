package main

import (
	"fmt"
	"log"
	"time"

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
	err = qrCode.WriteFile(256, "qris.png")
	if err != nil {
		log.Fatalf("Error saving QR code: %v", err)
	}

	fmt.Println("QR Code berhasil dibuat dan disimpan sebagai qris.png")
	fmt.Println("Silahkan scan QR code untuk melakukan pembayaran...")

	// Cek status pembayaran secara berulang
	for {
		status, err := qrisInstance.CheckPaymentStatus("TRX123", 1000)
		if err != nil {
			log.Printf("Error checking payment status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Status Pembayaran: %s\n", status.Status)
		if status.Status == "PAID" {
			fmt.Printf("Pembayaran berhasil!\n")
			fmt.Printf("Amount: %d\n", status.Amount)
			fmt.Printf("Reference: %s\n", status.Reference)
			fmt.Printf("Date: %s\n", status.Date)
			fmt.Printf("Brand: %s\n", status.BrandName)
			break
		}

		// Tunggu 5 detik sebelum cek lagi
		time.Sleep(5 * time.Second)
	}
} 
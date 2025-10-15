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
		BaseQrString: "BASE_QR_STRING",
		AuthToken:    "YOUR_AUTH_TOKEN",
		AuthUsername: "YOUR_AUTH_USERNAME",
	}

	// Buat instance QRIS
	qrisInstance, err := qris.NewQRIS(config)
	if err != nil {
		panic(err)
	}

	// Generate QR Code
	data := qris.QRISData{
		Amount:        150,
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
		fmt.Println("\nMengecek status pembayaran...")
		status, err := qrisInstance.CheckPaymentStatus("TRX123", 150)
		if err != nil {
			log.Printf("Error checking payment status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Tampilkan detail status
		fmt.Printf("Status Pembayaran: %s\n", status.Status)
		fmt.Printf("Amount yang diharapkan: %d\n", 150)
		fmt.Printf("Amount yang diterima: %d\n", status.Amount)
		fmt.Printf("Reference: %s\n", status.Reference)
		
		if status.Status == "PAID" {
			fmt.Printf("Pembayaran berhasil!\n")
			fmt.Printf("Date: %s\n", status.Date)
			fmt.Printf("Brand: %s\n", status.BrandName)
			fmt.Printf("Buyer Ref: %s\n", status.BuyerRef)
			break
		} else {
			fmt.Println("Menunggu pembayaran...")
		}

		// Tunggu 5 detik sebelum cek lagi
		time.Sleep(5 * time.Second)
	}
} 
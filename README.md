<<<<<<< HEAD
# 🚀 Paket Pembayaran QRIS untuk Go

Paket Go yang menyediakan integrasi pembayaran QRIS (Quick Response Code Indonesian Standard) untuk aplikasi Anda.

## 🎯 Fitur Utama

- Generate QRIS dengan format standar
- Pengecekan status pembayaran secara real-time
- Validasi format QRIS
- Kalkulasi checksum CRC16
- Dukungan base QRIS string
- Penanganan error yang lebih baik
- QR code dengan tingkat koreksi error tinggi

## 📦 Instalasi

```bash
go get github.com/AutoFTbot/OrderKuota-go
```

## 🚀 Penggunaan

### Inisialisasi

```go
import "github.com/AutoFTbot/OrderKuota-go/qris"

config := qris.QRISConfig{
    BaseQrString: "your-base-qr-string",
    AuthToken:    "your-auth-token",
    AuthUsername: "your-auth-username",
}

qrisInstance, err := qris.NewQRIS(config)
if err != nil {
    // handle error
}
```

### Generate QR Code

```go
data := qris.QRISData{
    Amount:        100000,
    TransactionID: "TRX123",
}

qrCode, err := qrisInstance.GenerateQRCode(data)
if err != nil {
    // handle error
}

// Simpan QR code ke file
err = qrCode.Save("qris.png")
```

### Generate QRIS String

```go
data := qris.QRISData{
    Amount:        100000,
    TransactionID: "TRX123",
}

qrString, err := qrisInstance.GetQRISString(data)
if err != nil {
    // handle error
}
```

### Cek Status Pembayaran

```go
status, err := qrisInstance.CheckPaymentStatus("TRX123", 100000)
if err != nil {
    // handle error
}

if status.Status == "PAID" {
    // Pembayaran berhasil
    fmt.Printf("Pembayaran diterima dari %s pada %s\n", 
        status.BrandName, status.Date)
}
```

### Validasi String QRIS

```go
err := qrisInstance.ValidateQRISString(qrString)
if err != nil {
    // handle error
}
```

## 📝 Dokumentasi

### QRISConfig

```go
type QRISConfig struct {
    BaseQrString string // Base QRIS string dari merchant
    AuthToken    string // Token autentikasi untuk API
    AuthUsername string // Username autentikasi untuk API
}
```

### QRISData

```go
type QRISData struct {
    Amount        int64  // Nominal pembayaran
    TransactionID string // ID transaksi unik
}
```

### PaymentStatus

```go
type PaymentStatus struct {
    Status    string // Status pembayaran (PAID/UNPAID)
    Amount    int64  // Nominal pembayaran
    Reference string // Referensi pembayaran
    Date      string // Tanggal pembayaran (jika PAID)
    BrandName string // Nama brand pembayar (jika PAID)
    BuyerRef  string // Referensi pembeli (jika PAID)
}
```

### Contoh Kode

```go
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
		BaseQrString: "your-base-qr-string",
		AuthToken:    "your-auth-token",
		AuthUsername: "your-auth-username",
	}

	// Buat instance QRIS
	qrisInstance, err := qris.NewQRIS(config)
	if err != nil {
		panic(err)
	}

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
		fmt.Println("\nMengecek status pembayaran...")
		status, err := qrisInstance.CheckPaymentStatus("TRX123", 1000)
		if err != nil {
			log.Printf("Error checking payment status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Tampilkan detail status
		fmt.Printf("Status Pembayaran: %s\n", status.Status)
		fmt.Printf("Amount yang diharapkan: %d\n", 1000)
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
```

## 🔍 Penanganan Error

Paket ini menyediakan penanganan error yang lebih baik dengan pesan error yang jelas:

- Validasi input saat inisialisasi
- Validasi format QRIS
- Validasi checksum
- Error saat generate QR code
- Error saat cek status pembayaran

## 🛠️ Praktik Terbaik

1. Selalu cek error saat inisialisasi QRIS
2. Gunakan ID transaksi unik untuk setiap transaksi
3. Validasi string QRIS sebelum digunakan
4. Gunakan penanganan error yang tepat
5. Simpan QR code dalam format PNG untuk kualitas terbaik

## 🤝 Kontribusi

Silakan kirim pull request. Untuk perubahan besar, harap buka issue terlebih dahulu untuk mendiskusikan perubahan yang diinginkan.

## 📄 Lisensi 

[MIT](https://choosealicense.com/licenses/mit/)

## Donations

Support this project by making a donation via QRIS. Your contribution helps maintain and improve this package.

### How to Donate

1. Generate a donation QR code:
```go
import (
    "github.com/AutoFTbot/OrderKuota-go/qris"
    "github.com/AutoFTbot/OrderKuota-go/ci-donation"
)

// Initialize QRIS
config := qris.QRISConfig{
    BaseQrString: "your-base-qr-string",
    AuthToken:    "your-auth-token",
    AuthUsername: "your-auth-username",
}

qr, err := qris.NewQRIS(config)
if err != nil {
    log.Fatal(err)
}

// Create donation manager
donationManager := donation.NewManager(qr, "donations.json")

// Generate QR code for donation
qrCode, err := donationManager.GenerateQR(100000, "DONATE-001")
if err != nil {
    log.Fatal(err)
}

// Save QR code to file
err = os.WriteFile("donation-qr.png", qrCode, 0644)
if err != nil {
    log.Fatal(err)
}
```

2. Scan the QR code using your mobile banking app
3. Complete the payment

### Recent Donations

Total donations: [![Total Donations](https://img.shields.io/badge/Total%20Donations-Rp%200-blue)](https://github.com/username/repo)

Latest donations:
- Rp 100.000 from John Doe - "Thank you for this amazing package!"
- Rp 50.000 from Jane Smith - "Keep up the good work!"

### Donation Records

All donations are recorded in `donations.json` and can be accessed programmatically:

```go
donations, err := donationManager.GetAll()
if err != nil {
    log.Fatal(err)
}

total, err := donationManager.GetTotal()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total donations: Rp %d\n", total)
for _, d := range donations {
    fmt.Printf("Donation: Rp %d from %s\n", d.Amount, d.DonorName)
}
``` 
=======
# 🚀 QRIS Payment Package for Go

![Go](https://github.com/AutoFTbot/OrderKuota-go/actions/workflows/go.yml/badge.svg)

A Go package that provides QRIS (Quick Response Code Indonesian Standard) payment integration for your applications.
Package Go yang menyediakan integrasi pembayaran QRIS (Quick Response Code Indonesian Standard) untuk aplikasi Anda.

## 🎯 Main Features / Fitur Utama

- Generate QRIS with standard format / Generate QRIS dengan format standar
- Real-time payment status checking / Pengecekan status pembayaran secara real-time
- QRIS format validation / Validasi format QRIS
- CRC16 checksum calculation / Kalkulasi checksum CRC16
- Base QRIS string support / Dukungan base QRIS string
- Better error handling / Penanganan error yang lebih baik
- High error correction QR code / QR code dengan tingkat koreksi error tinggi

## 📦 Installation / Instalasi

```bash
go get github.com/AutoFTbot/OrderKuota-go
```

## 🚀 Usage / Penggunaan

### Initialization / Inisialisasi

```go
import "github.com/AutoFTbot/OrderKuota-go/qris"

config := qris.QRISConfig{
    MerchantID:   "123456789",
    APIKey:       "your-api-key",
    BaseQrString: "your-base-qr-string",
}

qrisInstance, err := qris.NewQRIS(config)
if err != nil {
    // handle error
}
```

### Generate QR Code

```go
data := qris.QRISData{
    Amount:        100000,
    TransactionID: "TRX123",
}

qrCode, err := qrisInstance.GenerateQRCode(data)
if err != nil {
    // handle error
}

// Save QR code to file / Simpan QR code ke file
err = qrCode.Save("qris.png")
```

### Generate QRIS String

```go
data := qris.QRISData{
    Amount:        100000,
    TransactionID: "TRX123",
}

qrString, err := qrisInstance.GetQRISString(data)
if err != nil {
    // handle error
}
```

### Check Payment Status / Cek Status Pembayaran

```go
status, err := qrisInstance.CheckPaymentStatus("TRX123", 100000)
if err != nil {
    // handle error
}

if status.Status == "PAID" {
    // Payment successful / Pembayaran berhasil
    fmt.Printf("Payment received from %s at %s\n", 
        status.BrandName, status.Date)
}
```

### Validate QRIS String / Validasi String QRIS

```go
err := qrisInstance.ValidateQRISString(qrString)
if err != nil {
    // handle error
}
```

## 📝 Documentation / Dokumentasi

### QRISConfig

```go
type QRISConfig struct {
    MerchantID   string // Merchant ID from payment gateway / ID merchant dari payment gateway
    APIKey       string // API key for authentication / API key untuk autentikasi
    BaseQrString string // Base QRIS string from merchant / Base QRIS string dari merchant
}
```

### QRISData

```go
type QRISData struct {
    Amount        int64  // Payment amount / Nominal pembayaran
    TransactionID string // Unique transaction ID / ID transaksi unik
}
```

### PaymentStatus

```go
type PaymentStatus struct {
    Status    string // Payment status (PAID/UNPAID) / Status pembayaran (PAID/UNPAID)
    Amount    int64  // Payment amount / Nominal pembayaran
    Reference string // Payment reference / Referensi pembayaran
    Date      string // Payment date (if PAID) / Tanggal pembayaran (jika PAID)
    BrandName string // Payer brand name (if PAID) / Nama brand pembayar (jika PAID)
    BuyerRef  string // Buyer reference (if PAID) / Referensi pembeli (jika PAID)
}
```

### Examples Code

```go
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
		MerchantID:   "#",
		APIKey:       "#",
		BaseQrString: "#",
	}

	// Buat instance QRIS
	qrisInstance, err := qris.NewQRIS(config)
	if err != nil {
		panic(err)
	}

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
		fmt.Println("\nMengecek status pembayaran...")
		status, err := qrisInstance.CheckPaymentStatus("TRX123", 1000)
		if err != nil {
			log.Printf("Error checking payment status: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Tampilkan detail status
		fmt.Printf("Status Pembayaran: %s\n", status.Status)
		fmt.Printf("Amount yang diharapkan: %d\n", 1000)
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
```

## 🔍 Error Handling / Penanganan Error

The package provides better error handling with clear error messages:
Package ini menyediakan penanganan error yang lebih baik dengan pesan error yang jelas:

- Input validation during initialization / Validasi input saat inisialisasi
- QRIS format validation / Validasi format QRIS
- Checksum validation / Validasi checksum
- QR code generation errors / Error saat generate QR code
- Payment status checking errors / Error saat cek status pembayaran

## 🛠️ Best Practices / Praktik Terbaik

1. Always check for errors during QRIS initialization / Selalu cek error saat inisialisasi QRIS
2. Use unique transaction IDs for each transaction / Gunakan ID transaksi unik untuk setiap transaksi
3. Validate QRIS string before use / Validasi string QRIS sebelum digunakan
4. Use proper error handling / Gunakan penanganan error yang tepat
5. Save QR code in PNG format for best quality / Simpan QR code dalam format PNG untuk kualitas terbaik

## 🤝 Contributing / Kontribusi

Feel free to submit pull requests. For major changes, please open an issue first to discuss what you would like to change.
Silakan kirim pull request. Untuk perubahan besar, harap buka issue terlebih dahulu untuk mendiskusikan perubahan yang diinginkan.

## 📄 License / Lisensi

[AutoFtBot](https://github.com/AutoFTbot/OrderKuota-go/blob/main/LICENSE) 
>>>>>>> 35fc021c36759d7ac961797d5870716555dd47af

# 🚀 QRIS Payment Package for Go

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

[MIT](https://choosealicense.com/licenses/mit/) 
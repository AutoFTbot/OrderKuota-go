# 🚀 QRIS Payment Package for Go

Package QRIS untuk Go yang memudahkan integrasi pembayaran QRIS dalam aplikasi Anda.

## 🎯 Fitur Utama

- Generate QRIS dengan format standar
- Cek status pembayaran
- Validasi format QRIS
- Kalkulasi checksum CRC16
- Support baseQrString

## 📦 Instalasi

```bash
go get github.com/AutoFTbot/OrderKuota-go
```

## 🚀 Penggunaan

### Inisialisasi

```go
import "github.com/AutoFTbot/OrderKuota-go/qris"

config := qris.QRISConfig{
    MerchantID:   "123456789",
    APIKey:       "your-api-key",
    BaseQrString: "your-base-qr-string",
}

qrisInstance := qris.NewQRIS(config)
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

### Cek Status Pembayaran

```go
status, err := qrisInstance.CheckPaymentStatus("TRX123", 100000)
if err != nil {
    // handle error
}

if status.Status == "PAID" {
    // Pembayaran berhasil
}
```

## 📝 Dokumentasi

### QRISConfig

```go
type QRISConfig struct {
    MerchantID   string // ID merchant dari payment gateway
    APIKey       string // API key untuk autentikasi
    BaseQrString string // Base QRIS string dari merchant
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
    Status        string // Status pembayaran (PAID/UNPAID)
    Amount        int64  // Nominal pembayaran
    Reference     string // Referensi pembayaran
    Date          string // Tanggal pembayaran (jika PAID)
    BrandName     string // Nama brand pembayar (jika PAID)
    BuyerRef      string // Referensi pembeli (jika PAID)
}
```

## 🤝 Kontribusi

Silakan buat pull request untuk kontribusi. Untuk perubahan besar, harap buka issue terlebih dahulu untuk mendiskusikan perubahan yang diinginkan.

## 📄 Lisensi

[MIT](https://choosealicense.com/licenses/mit/) 
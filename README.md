# **🚀・GoKBBI - Library dan CLI KBBI Daring**

### **Pernah Bermimpi Mengakses KBBI Langsung dari Kode Go Anda? ✨**

Selamat datang di GoKBBI! Proyek keren ini adalah jembatan ajaib Anda untuk mengakses Kamus Besar Bahasa Indonesia (KBBI) Daring langsung dari aplikasi Go Anda! GoKBBI adalah re-implementasi dari [kbbi-python](https://github.com/laymonage/kbbi-python) dalam bahasa Go yang memberikan akses mudah ke KBBI melalui **library Go yang powerful** dan **command line interface yang user-friendly**. Bayangkan bisa mencari arti kata, mendapatkan etimologi, dan mengakses fitur KBBI lengkap - semua dari kode Go Anda!

---

### **📋・Daftar Isi**

- **✨・<a href="#apa-itu-gokbbi" style="text-decoration: none;">Kenali GoKBBI - Jembatan Ajaib ke KBBI!</a>**
- **🚀・<a href="#memulai" style="text-decoration: none;">Mari Memulai Petualangan Coding!</a>**
- **⚙️・<a href="#instalasi" style="text-decoration: none;">Instalasi Super Mudah!</a>**
- **🗺️・<a href="#penggunaan-cli" style="text-decoration: none;">Menguasai CLI seperti Pro!</a>**
- **💻・<a href="#penggunaan-lib" style="text-decoration: none;">Library Go yang Powerful!</a>**
- **🛠️・<a href="#opsi-command-line" style="text-decoration: none;">Panduan Lengkap Command Line!</a>**
- **🏗️・<a href="#struktur-proyek" style="text-decoration: none;">Arsitektur Proyek yang Rapi!</a>**
- **🔧・<a href="#error-handling-umum" style="text-decoration: none;">Mengatasi Error dengan Elegan!</a>**
- **⚠️・<a href="#penting-penggunaan-bertanggung-jawab" style="text-decoration: none;">Gunakan dengan Bijak dan Bertanggung Jawab!</a>**
- **💖・<a href="#contributing" style="text-decoration: none;">Bergabung dalam Komunitas!</a>**
- **📜・<a href="#license" style="text-decoration: none;">Lisensi MIT yang Bebas!</a>**
- **👋・<a href="#tentang-penulis" style="text-decoration: none;">Siapa di Balik GoKBBI?</a>**

---

### <div id="apa-itu-gokbbi">**✨・Kenali GoKBBI - Jembatan Ajaib ke KBBI!**</div>

GoKBBI (Go Kamus Besar Bahasa Indonesia) adalah library dan CLI tool super keren yang memungkinkan Anda mengakses KBBI Daring dengan mudah! 🚀 Tujuan utama kami adalah menciptakan jembatan yang seamless dan menyenangkan antara aplikasi Go Anda dengan kekayaan bahasa Indonesia yang ada di KBBI. Bersiaplah untuk melihat kosakata Indonesia beraksi dalam kode Anda!

**Fitur Utama:**
- 🔍 **Pencarian kata** dalam KBBI Daring (termasuk kata dengan spasi)
- 🔐 **Autentikasi** untuk mengakses fitur pengguna terdaftar  
- 📚 **Library Go** dengan API yang mudah digunakan
- 💻 **CLI tool** dengan berbagai opsi
- 📖 **Format output** teks dan JSON
- ⚙️ **Opsi filtering** untuk menghilangkan contoh atau kata terkait
- 🌐 **Human-like request** dengan delay dan header yang wajar
- 📱 **Cross-platform** (Windows, macOS, Linux)</div>

---

### <div id="memulai">**🚀・Mari Memulai Petualangan Coding!**</div>

Siap untuk menghidupkan KBBI di aplikasi Go Anda? Berikut cara mendapatkan GoKBBI up and running dalam sekejap:

1.  **Clone keajaibannya!** ✨
    
    ```bash
    git clone https://github.com/ZulfaNurhuda/GoKBBI.project.git
    cd GoKBBI.project
    ```

2.  **Install Go Toolchain** ⚙️

    Untuk menjalankan GoKBBI, Anda memerlukan Go 1.21 atau yang lebih baru. Berikut instruksi setup untuk platform yang berbeda.

    #### **Windows**

    1. **Download Go:** Kunjungi [https://golang.org/dl/](https://golang.org/dl/) dan download installer Windows.
    2. **Install:** Jalankan installer dan ikuti petunjuk instalasi.
    3. **Verify:** Buka Command Prompt dan ketik `go version` untuk memastikan instalasi berhasil.

    #### **Linux**

    - **Debian/Ubuntu-based:**
    
       ```bash
       sudo apt update
       sudo apt install golang-go
       ```
 
    - **Fedora/RHEL-based:**
    
      ```bash
      sudo dnf install golang
      ```

    - **Arch Linux:**
    
      ```bash
      sudo pacman -S go
      ```

    #### **macOS**

    ```bash
    brew install go
    ```

---

### <div id="instalasi">**⚙️・Instalasi Super Mudah!**</div>

GoKBBI mendukung dua cara instalasi: **Build dari source** (untuk development) dan **sebagai library** (untuk proyek Anda).

#### **Opsi 1: Build dari Source (Cara Classic 📜)**

Perfect untuk development atau jika Anda ingin CLI tool:

```bash
# Clone dan setup dependencies
git clone https://github.com/ZulfaNurhuda/GoKBBI.project.git
cd GoKBBI.project
go mod tidy

# Build CLI tool
go build -o bin/kbbi cmd/kbbi/main.go
```

Setelah build berhasil, executable `kbbi` akan menunggu Anda di direktori `bin/`! 🎉

#### **Opsi 2: Sebagai Library Go (Cara Modern ✨)**

Untuk menggunakan GoKBBI dalam proyek Go Anda:

```bash
go get github.com/ZulfaNurhuda/GoKBBI.project
```

---

### <div id="penggunaan-cli">**🗺️・Menguasai CLI seperti Pro!**</div>

Setelah GoKBBI di-build, Anda hanya tinggal satu perintah saja untuk melihat keajaiban pencarian KBBI! Seperti magic, tapi dengan code:

#### **Pencarian Dasar CLI**

```bash
./bin/kbbi cinta
./bin/kbbi --kata rumah
```

#### **Autentikasi CLI**

```bash
# Login dan simpan kuki
./bin/kbbi --email your@email.com --sandi yourpassword --autentikasi

# Setelah login, kuki akan digunakan otomatis
./bin/kbbi cinta
```

#### **Format Output CLI**

```bash
# Output default (teks)
./bin/kbbi cinta

# Output JSON
./bin/kbbi --kata cinta --json

# Output JSON dengan indentasi
./bin/kbbi --kata cinta --json --indent
```

#### **Opsi Filtering CLI**

```bash
# Tanpa contoh penggunaan
./bin/kbbi --kata cinta --tanpa-contoh

# Tanpa kata terkait
./bin/kbbi --kata cinta --tanpa-terkait

# Mode nonpengguna (menonaktifkan fitur khusus)
./bin/kbbi --kata cinta --nonpengguna
```

#### **Manajemen Kuki CLI**

```bash
# Hapus kuki tersimpan
./bin/kbbi --bersihkan-kuki

# Gunakan lokasi kuki custom
./bin/kbbi --lokasi-kuki /path/to/cookies.json
```

---

### <div id="penggunaan-lib">**💻・Library Go yang Powerful!**</div>

#### **Pencarian Dasar**

```go
package main

import (
    "fmt"
    "log"
    gokbbi "github.com/ZulfaNurhuda/GoKBBI.project"
)

func main() {
    // Pencarian tanpa autentikasi
    definisi, err := gokbbi.Cari("rumah")
    if err != nil {
        if err == gokbbi.ErrModaTerbatas {
            log.Fatal("Perlu autentikasi: KBBI dalam moda terbatas")
        }
        log.Fatal(err)
    }
    
    // Tampilkan hasil
    fmt.Println(definisi.String())
    
    // Akses data terstruktur
    for _, entri := range definisi.Entri {
        fmt.Printf("Entri: %s\n", entri.Nama)
        for i, makna := range entri.Makna {
            fmt.Printf("%d. %s\n", i+1, makna.String())
        }
    }
}
```

#### **Dengan Autentikasi**

```go
// Autentikasi baru
auth, err := gokbbi.NewAuth("email@example.com", "password", "")
if err != nil {
    log.Fatal(err)
}

// Simpan kuki untuk penggunaan berikutnya  
err = auth.SimpanKuki()
if err != nil {
    log.Printf("Warning: %v", err)
}

// Pencarian dengan autentikasi (akses fitur lengkap)
definisi, err := gokbbi.CariDenganAuth("cinta", auth)
if err != nil {
    log.Fatal(err)
}

// Dengan autentikasi, dapat akses etimologi dan kata terkait
for _, entri := range definisi.Entri {
    if entri.Etimologi != nil {
        fmt.Printf("Etimologi: %s\n", entri.Etimologi.String())
    }
    fmt.Printf("Kata Turunan: %s\n", strings.Join(entri.KataTurunan, ", "))
}
```

#### **Load Kuki Tersimpan**

```go
// Muat kuki yang sudah tersimpan
auth, err := gokbbi.LoadAuth("")
if err != nil {
    // Kuki tidak ada, perlu login ulang
    auth, err = gokbbi.NewAuth("email@example.com", "password", "")
    if err != nil {
        log.Fatal(err)
    }
    auth.SimpanKuki()
}

definisi, err := gokbbi.CariDenganAuth("kata", auth)
```

#### **Export ke JSON**

```go
definisi, err := gokbbi.Cari("buku")
if err != nil {
    log.Fatal(err)
}

// Konversi ke JSON
jsonStr, err := definisi.ToJSON(true) // true = dengan indentasi
if err != nil {
    log.Fatal(err)
}

fmt.Println(jsonStr)
```

#### **Error Handling**

```go
definisi, err := gokbbi.Cari("katayangtidakada")
if err != nil {
    switch err {
    case gokbbi.ErrTidakDitemukan:
        fmt.Println("Kata tidak ditemukan")
        // Cek apakah ada saran
        if definisi != nil && len(definisi.SaranEntri) > 0 {
            fmt.Printf("Saran: %s\n", strings.Join(definisi.SaranEntri, ", "))
        }
    case gokbbi.ErrBatasSehari:
        fmt.Println("Batas pencarian harian tercapai")
    case gokbbi.ErrModaTerbatas:
        fmt.Println("KBBI dalam moda terbatas, perlu autentikasi")
    case gokbbi.ErrAkunDibekukan:
        fmt.Println("Akun dibekukan")
    default:
        fmt.Printf("Error lain: %v\n", err)
    }
}
```

---

### <div id="opsi-command-line">**🛠️・Panduan Lengkap Command Line!**</div>

#### **Pencarian**
- `--kata <kata>` - Kata yang ingin dicari
- `--json` - Output dalam format JSON
- `--indent` - Gunakan indentasi untuk JSON
- `--tanpa-contoh` - Jangan tampilkan contoh penggunaan
- `--tanpa-terkait` - Jangan tampilkan kata terkait
- `--nonpengguna` - Nonaktifkan fitur khusus pengguna

#### **Autentikasi**
- `--email <email>` - Alamat email akun KBBI
- `--sandi <password>` - Kata sandi akun KBBI
- `--autentikasi` - Lakukan proses autentikasi
- `--lokasi-kuki <path>` - Lokasi file kuki
- `--bersihkan-kuki` - Hapus kuki tersimpan

#### **Lainnya**

- `--bantuan, --help` - Tampilkan bantuan
- `--version` - Tampilkan versi aplikasi

---

### <div id="struktur-proyek">**🏗️・Arsitektur Proyek yang Rapi!**</div>

```bash
GoKBBI.project/
├── cmd/kbbi/          # Main CLI application
├── internal/
│   ├── auth/          # Autentikasi KBBI
│   ├── fetcher/       # HTTP client untuk mengambil halaman
│   ├── model/         # Data structures
│   └── parser/        # HTML parser
├── go.mod
├── go.sum
└── README.md
```

---

### <div id="error-handling-umum">**🔧・Mengatasi Error dengan Elegan!**</div>

Aplikasi menangani berbagai kondisi error:

- **Entri tidak ditemukan**: Menampilkan saran kata alternatif
- **Batas harian tercapai**: Error dengan pesan informatif
- **Akun dibekukan**: Error dengan pesan peringatan
- **Moda terbatas**: Error ketika KBBI dalam mode terbatas
- **Koneksi gagal**: Error jaringan dengan retry mechanism

---

### <div id="penting-penggunaan-bertanggung-jawab">**⚠️・Gunakan dengan Bijak dan Bertanggung Jawab!**</div>

- Gunakan delay yang wajar antara request untuk menghindari pemblokiran
- Jangan melakukan request berlebihan dalam waktu singkat
- Hormati _terms of service_ KBBI Daring
- Gunakan fitur retry dengan bijak

---

### <div id="contributing">**💖・Bergabung dalam Komunitas!**</div>

Punya ide keren? Menemukan bug yang menyebalkan? Ingin menambah fitur awesome? Kami akan sangat senang dengan bantuan Anda! GoKBBI adalah effort komunitas, dan setiap kontribusi, besar atau kecil, membuat perbedaan. Cek `CONTRIBUTING.md` kami untuk detail lengkap tentang bagaimana Anda bisa ikut serta dan membuat GoKBBI lebih keren lagi. Mari bangun sesuatu yang amazing bersama-sama!

---

### <div id="license">**📜・Lisensi MIT yang Bebas!**</div>

Proyek ini open-source dan dengan bangga didistribusikan di bawah MIT License. Ini artinya Anda bebas untuk explore, use, modify, dan share GoKBBI! Anda bisa menemukan semua detail lengkapnya di file `LICENSE`. Happy coding!

---

### **⚠️・Penting: Penggunaan Bertanggung Jawab**

- 🕒 Gunakan delay yang wajar antara request untuk menghindari pemblokiran
- 🚫 Jangan melakukan request berlebihan dalam waktu singkat  
- 📋 Hormati _terms of service_ KBBI Daring
- 🔄 Gunakan fitur retry dengan bijak

---

### <div id="tentang-penulis">**👋・Siapa di Balik GoKBBI?**</div>

**Muhammad Zulfa Fauzan Nurhuda** (18224064)

Seorang mahasiswa STI ITB yang sangat tertarik dalam dunia programming! 😄 Selalu antusias belajar dan membangun project keren seperti GoKBBI ini! 🚀

<img src="https://i.imgur.com/Zp8msEG.png" alt="Logo ITB" height="90" style="border-radius: 10px">
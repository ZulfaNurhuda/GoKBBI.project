// Package gokbbi menyediakan akses terprogram ke KBBI Daring
//
// GoKBBI adalah library Go untuk mengakses Kamus Besar Bahasa Indonesia (KBBI) Daring.
// Library ini menyediakan API yang mudah digunakan untuk mencari kata, mengambil definisi,
// dan mengakses fitur-fitur KBBI lainnya.
//
// Contoh penggunaan dasar:
//
//	// Pencarian tanpa autentikasi
//	definisi, err := gokbbi.Cari("rumah")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(definisi.String())
//
//	// Pencarian dengan autentikasi
//	auth, err := gokbbi.NewAuth("email@example.com", "password", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	definisi, err := gokbbi.CariDenganAuth("cinta", auth)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(definisi.String())
//
// Untuk fitur lengkap dan dokumentasi detail, lihat README.md
package gokbbi

import (
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/auth"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/fetcher"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/model"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/parser"
)

// Definisi adalah struktur data hasil pencarian KBBI
type Definisi = model.Definisi

// Entri adalah struktur data entri dalam KBBI
type Entri = model.Entri

// Makna adalah struktur data makna dari entri
type Makna = model.Makna

// Etimologi adalah struktur data etimologi kata
type Etimologi = model.Etimologi

// KelasKata adalah struktur data kelas kata
type KelasKata = model.KelasKata

// Auth adalah struktur untuk autentikasi KBBI
type Auth = auth.AutentikasiKBBI

// Error types yang bisa dikembalikan oleh library
var (
	ErrTidakDitemukan   = fetcher.ErrTidakDitemukan
	ErrBatasSehari      = fetcher.ErrBatasSehari
	ErrModaTerbatas     = fetcher.ErrModaTerbatas
	ErrTerjadiKesalahan = fetcher.ErrTerjadiKesalahan
	ErrAkunDibekukan    = fetcher.ErrAkunDibekukan
)

// Cari mencari kata dalam KBBI tanpa autentikasi
//
// Parameter:
//   - kata: kata atau frasa yang ingin dicari
//
// Return:
//   - *Definisi: hasil pencarian berisi entri, makna, dll
//   - error: error jika terjadi masalah dalam pencarian
//
// Contoh:
//
//	definisi, err := gokbbi.Cari("rumah")
//	if err != nil {
//		// Handle error
//		if err == gokbbi.ErrTidakDitemukan {
//			fmt.Println("Kata tidak ditemukan")
//		}
//		return err
//	}
//
//	// Tampilkan hasil
//	fmt.Println(definisi.String())
//
//	// Atau akses data terstruktur
//	for _, entri := range definisi.Entri {
//		fmt.Printf("Entri: %s\n", entri.Nama)
//		for i, makna := range entri.Makna {
//			fmt.Printf("%d. %s\n", i+1, makna.String())
//		}
//	}
func Cari(kata string) (*Definisi, error) {
	return CariDenganAuth(kata, nil)
}

// CariDenganAuth mencari kata dalam KBBI dengan autentikasi
//
// Parameter:
//   - kata: kata atau frasa yang ingin dicari
//   - auth: objek autentikasi, bisa nil untuk pencarian tanpa auth
//
// Return:
//   - *Definisi: hasil pencarian berisi entri, makna, dll
//   - error: error jika terjadi masalah dalam pencarian
//
// Contoh dengan autentikasi:
//
//	auth, err := gokbbi.NewAuth("email@example.com", "password", "")
//	if err != nil {
//		return err
//	}
//
//	definisi, err := gokbbi.CariDenganAuth("cinta", auth)
//	if err != nil {
//		return err
//	}
//
//	// Dengan autentikasi, bisa akses etimologi dan kata terkait
//	for _, entri := range definisi.Entri {
//		if entri.Etimologi != nil {
//			fmt.Printf("Etimologi: %s\n", entri.Etimologi.String())
//		}
//		fmt.Printf("Kata Turunan: %s\n", strings.Join(entri.KataTurunan, ", "))
//	}
func CariDenganAuth(kata string, autentikasi *Auth) (*Definisi, error) {
	// Ambil halaman HTML
	html, err := fetcher.AmbilHalamanDenganRetry(kata, autentikasi, 3)
	if err != nil {
		// Jika error adalah TidakDitemukan dan ada HTML, parse untuk saran
		if kesalahanKBBI, ok := err.(*fetcher.KesalahanKBBI); ok && kesalahanKBBI.Jenis == "TidakDitemukan" && html != "" {
			// Parse saran entri
			terautentikasi := autentikasi != nil && autentikasi.Terautentikasi
			definisi, parseErr := parser.ParseDefinisi(html, terautentikasi)
			if parseErr != nil {
				return nil, parseErr
			}
			parser.SetPranala(definisi, kata)
			return definisi, err // Kembalikan definisi dengan saran DAN error TidakDitemukan
		}
		return nil, err
	}

	// Parse HTML menjadi definisi
	terautentikasi := autentikasi != nil && autentikasi.Terautentikasi
	definisi, err := parser.ParseDefinisi(html, terautentikasi)
	if err != nil {
		return nil, err
	}

	// Set pranala
	parser.SetPranala(definisi, kata)

	return definisi, nil
}

// NewAuth membuat objek autentikasi baru
//
// Parameter:
//   - email: alamat email akun KBBI Daring
//   - sandi: kata sandi akun KBBI Daring
//   - lokasiKuki: lokasi file kuki, kosong untuk default (~/.kbbi/kuki.json)
//
// Return:
//   - *Auth: objek autentikasi yang bisa digunakan untuk pencarian
//   - error: error jika autentikasi gagal
//
// Contoh:
//
//	// Autentikasi dengan email dan password
//	auth, err := gokbbi.NewAuth("email@example.com", "password", "")
//	if err != nil {
//		return err
//	}
//
//	// Simpan kuki untuk penggunaan berikutnya
//	err = auth.SimpanKuki()
//	if err != nil {
//		return err
//	}
//
//	// Gunakan untuk pencarian
//	definisi, err := gokbbi.CariDenganAuth("kata", auth)
func NewAuth(email, sandi, lokasiKuki string) (*Auth, error) {
	return auth.BaruAuth(email, sandi, lokasiKuki)
}

// LoadAuth memuat autentikasi dari kuki yang tersimpan
//
// Parameter:
//   - lokasiKuki: lokasi file kuki, kosong untuk default (~/.kbbi/kuki.json)
//
// Return:
//   - *Auth: objek autentikasi yang dimuat dari kuki
//   - error: error jika kuki tidak ditemukan atau tidak valid
//
// Contoh:
//
//	// Muat kuki yang sudah tersimpan
//	auth, err := gokbbi.LoadAuth("")
//	if err != nil {
//		// Kuki tidak ada, perlu login ulang
//		auth, err = gokbbi.NewAuth("email@example.com", "password", "")
//		if err != nil {
//			return err
//		}
//		auth.SimpanKuki()
//	}
//
//	// Gunakan untuk pencarian
//	definisi, err := gokbbi.CariDenganAuth("kata", auth)
func LoadAuth(lokasiKuki string) (*Auth, error) {
	return auth.BaruAuth("", "", lokasiKuki)
}

// CekKoneksi memeriksa koneksi ke KBBI Daring
//
// Return:
//   - error: nil jika koneksi berhasil, error jika gagal
//
// Contoh:
//
//	err := gokbbi.CekKoneksi()
//	if err != nil {
//		fmt.Println("Tidak dapat terhubung ke KBBI:", err)
//		return
//	}
//	fmt.Println("Koneksi ke KBBI berhasil")
func CekKoneksi() error {
	return fetcher.CekKoneksi()
}


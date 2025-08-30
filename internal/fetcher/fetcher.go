// Package fetcher menyediakan fungsi untuk mengambil halaman dari KBBI Daring
package fetcher

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/auth"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/cache"
)

const (
	HostKBBI = "https://kbbi.kemdikbud.go.id"
)

// KesalahanKBBI merepresentasikan berbagai kesalahan dari KBBI
type KesalahanKBBI struct {
	Jenis  string
	Pesan  string
}

func (e *KesalahanKBBI) Error() string {
	return e.Pesan
}

var (
	ErrTidakDitemukan = &KesalahanKBBI{
		Jenis: "TidakDitemukan",
		Pesan: "Entri tidak ditemukan dalam KBBI",
	}
	ErrBatasSehari = &KesalahanKBBI{
		Jenis: "BatasSehari", 
		Pesan: "Pencarian Anda telah mencapai batas maksimum dalam sehari",
	}
	ErrModaTerbatas = &KesalahanKBBI{
		Jenis: "ModaTerbatas",
		Pesan: "KBBI Daring sedang dalam moda terbatas. Fitur pencarian dibatasi untuk pengguna umum",
	}
	ErrTerjadiKesalahan = &KesalahanKBBI{
		Jenis: "TerjadiKesalahan",
		Pesan: "Terjadi kesalahan saat memproses permintaan Anda",
	}
	ErrAkunDibekukan = &KesalahanKBBI{
		Jenis: "AkunDibekukan", 
		Pesan: "Akun ini sedang dibekukan, tidak dapat digunakan",
	}
)

// AmbilHalaman mengambil halaman dari KBBI berdasarkan kata pencarian
func AmbilHalaman(kata string, autentikasi *auth.AutentikasiKBBI) (string, error) {
	return AmbilHalamanDenganCache(kata, autentikasi, "", false)
}

// AmbilHalamanDenganCache mengambil halaman dari KBBI dengan dukungan cache
func AmbilHalamanDenganCache(kata string, autentikasi *auth.AutentikasiKBBI, lokasiKuki string, tanpaCache bool) (string, error) {
	var managerCache *cache.ManagerCache
	var err error

	// Inisialisasi cache manager jika cache digunakan
	if !tanpaCache {
		managerCache, err = cache.BaruManagerCache(lokasiKuki)
		if err != nil {
			// Jika gagal membuat cache manager, lanjutkan tanpa cache
			managerCache = nil
		}
	}

	// Coba ambil dari cache terlebih dahulu jika cache aktif
	if managerCache != nil {
		if htmlCache, found := managerCache.AmbilCache(kata); found {
			return htmlCache, nil
		}
	}

	// Jika tidak ada di cache atau cache dinonaktifkan, ambil dari KBBI
	html, err := ambilHalamanLangsung(kata, autentikasi)
	
	// Simpan ke cache hanya jika berhasil (tidak ada error) dan cache aktif
	if managerCache != nil && err == nil {
		// Simpan ke cache, abaikan error penyimpanan
		managerCache.SimpanCache(kata, html)
	}

	return html, err
}

// ambilHalamanLangsung mengambil halaman langsung dari KBBI tanpa cache
func ambilHalamanLangsung(kata string, autentikasi *auth.AutentikasiKBBI) (string, error) {
	var client *http.Client
	
	if autentikasi != nil {
		client = autentikasi.GetClient()
	} else {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// Tentukan URL berdasarkan kata pencarian
	lokasi := tentukanLokasi(kata)
	urlLengkap := fmt.Sprintf("%s/%s", HostKBBI, lokasi)

	// Buat request dengan header yang wajar
	req, err := http.NewRequest("GET", urlLengkap, nil)
	if err != nil {
		return "", fmt.Errorf("gagal membuat request: %w", err)
	}

	// Set header untuk terlihat seperti browser biasa
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Tambahkan delay kecil untuk menghindari rate limiting
	time.Sleep(500 * time.Millisecond)

	// Kirim request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengambil halaman: %w", err)
	}
	defer resp.Body.Close()

	// Periksa status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server mengembalikan status code: %d", resp.StatusCode)
	}

	// Baca response body dengan dekompres jika perlu
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("gagal membuat gzip reader: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("gagal membaca response body: %w", err)
	}

	htmlContent := string(body)
	
	// Update status autentikasi jika ada objek auth
	if autentikasi != nil {
		autentikasi.CekAutentikasi(htmlContent)
	}

	// Periksa kesalahan berdasarkan URL redirect atau konten
	if err := cekKesalahan(resp.Request.URL.String(), htmlContent); err != nil {
		return htmlContent, err // Kembalikan HTML untuk saran entri
	}

	return htmlContent, nil
}

// tentukanLokasi menentukan path URL berdasarkan kata pencarian
func tentukanLokasi(kata string) string {
	// Kasus khusus yang memerlukan pencarian via Cari/Hasil
	kasusKhusus := []bool{
		strings.Contains(kata, "."),
		strings.Contains(kata, "?"),
		strings.ToLower(kata) == "nul",
		strings.ToLower(kata) == "bin",
	}

	for _, kondisi := range kasusKhusus {
		if kondisi {
			return fmt.Sprintf("Cari/Hasil?frasa=%s", url.QueryEscape(kata))
		}
	}

	// Kasus normal - akses langsung ke entri (termasuk kata dengan spasi)
	// Gunakan PathEscape untuk URL path, bukan QueryEscape
	return fmt.Sprintf("entri/%s", url.PathEscape(kata))
}

// cekKesalahan memeriksa apakah ada kesalahan dalam response
func cekKesalahan(urlResponse, htmlContent string) error {
	// Periksa URL redirect
	if strings.Contains(urlResponse, "Beranda/Error") {
		return ErrTerjadiKesalahan
	}
	
	if strings.Contains(urlResponse, "Beranda/BatasSehari") {
		return ErrBatasSehari
	}
	
	if strings.Contains(urlResponse, "Beranda/ModaTerbatas") {
		return ErrModaTerbatas
	}
	
	if strings.Contains(urlResponse, "Account/Banned") {
		return ErrAkunDibekukan
	}

	// Periksa konten HTML untuk berbagai error
	if strings.Contains(htmlContent, "Entri tidak ditemukan.") {
		return ErrTidakDitemukan
	}
	
	// Periksa konten HTML untuk moda terbatas
	if strings.Contains(htmlContent, "Moda terbatas sedang diaktifkan") || 
	   strings.Contains(htmlContent, "pengguna tidak terdaftar tidak dapat dilayani") ||
	   strings.Contains(htmlContent, "moda terbatas") {
		return ErrModaTerbatas
	}

	return nil
}

// AmbilHalamanDenganRetry mengambil halaman dengan retry mechanism
func AmbilHalamanDenganRetry(kata string, autentikasi *auth.AutentikasiKBBI, maxRetry int) (string, error) {
	return AmbilHalamanDenganRetrydanCache(kata, autentikasi, maxRetry, "", false)
}

// AmbilHalamanDenganRetrydanCache mengambil halaman dengan retry mechanism dan dukungan cache
func AmbilHalamanDenganRetrydanCache(kata string, autentikasi *auth.AutentikasiKBBI, maxRetry int, lokasiKuki string, tanpaCache bool) (string, error) {
	var lastErr error
	
	for i := 0; i < maxRetry; i++ {
		html, err := AmbilHalamanDenganCache(kata, autentikasi, lokasiKuki, tanpaCache)
		if err != nil {
			// Jika error adalah kesalahan KBBI tertentu, jangan retry
			if kesalahanKBBI, ok := err.(*KesalahanKBBI); ok {
				switch kesalahanKBBI.Jenis {
				case "TidakDitemukan", "BatasSehari", "ModaTerbatas", "AkunDibekukan":
					return html, err
				}
			}
			
			lastErr = err
			// Tambahkan delay yang semakin lama untuk retry
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		
		return html, nil
	}
	
	return "", fmt.Errorf("gagal mengambil halaman setelah %d percobaan: %w", maxRetry, lastErr)
}

// CekKoneksi memeriksa koneksi ke KBBI
func CekKoneksi() error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(HostKBBI)
	if err != nil {
		return fmt.Errorf("tidak dapat terhubung ke KBBI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("KBBI mengembalikan status code: %d", resp.StatusCode)
	}

	return nil
}
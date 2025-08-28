// Package auth menyediakan fungsi autentikasi untuk KBBI Daring
package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	HostKBBI      = "https://kbbi.kemdikbud.go.id"
	LokasiLogin   = "Account/Login"
	NamaKukiUtama = ".AspNet.ApplicationCookie"
)

// AutentikasiKBBI mengelola autentikasi dengan KBBI Daring
type AutentikasiKBBI struct {
	Email       string
	Sandi       string
	LokasiKuki  string
	client      *http.Client
	Terautentikasi bool
}

// BaruAuth membuat objek AutentikasiKBBI baru
func BaruAuth(email, sandi, lokasiKuki string) (*AutentikasiKBBI, error) {
	// Buat cookie jar untuk mengelola session
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	// Tentukan lokasi kuki default jika tidak disediakan
	if lokasiKuki == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan home directory: %w", err)
		}
		lokasiKuki = filepath.Join(homeDir, ".kbbi", "kuki.json")
	}

	auth := &AutentikasiKBBI{
		Email:      email,
		Sandi:      sandi,
		LokasiKuki: lokasiKuki,
		client:     client,
	}

	// Jika email dan sandi kosong, coba muat kuki
	if email == "" && sandi == "" {
		err = auth.MuatKuki()
		if err != nil {
			return nil, fmt.Errorf("tidak dapat memuat kuki: %w", err)
		}
	} else {
		// Lakukan login dengan email dan sandi
		err = auth.Login()
		if err != nil {
			return nil, err
		}
	}

	return auth, nil
}

// Login melakukan autentikasi ke KBBI Daring
func (a *AutentikasiKBBI) Login() error {
	// Ambil token CSRF
	token, err := a.ambilToken()
	if err != nil {
		return fmt.Errorf("gagal mengambil token: %w", err)
	}

	// Siapkan data form untuk login
	data := url.Values{}
	data.Set("__RequestVerificationToken", token)
	data.Set("Posel", a.Email)
	data.Set("KataSandi", a.Sandi)
	data.Set("IngatSaya", "true")

	// Kirim permintaan login
	resp, err := a.client.PostForm(fmt.Sprintf("%s/%s", HostKBBI, LokasiLogin), data)
	if err != nil {
		return fmt.Errorf("gagal melakukan login: %w", err)
	}
	defer resp.Body.Close()

	// Periksa hasil login
	if strings.Contains(resp.Request.URL.String(), "Beranda/Error") {
		return fmt.Errorf("terjadi kesalahan saat memproses permintaan login")
	}
	
	if strings.Contains(resp.Request.URL.String(), "Account/Login") {
		return fmt.Errorf("gagal melakukan autentikasi dengan alamat posel dan sandi yang diberikan")
	}

	a.Terautentikasi = true
	return nil
}

// SimpanKuki menyimpan kuki autentikasi ke file
func (a *AutentikasiKBBI) SimpanKuki() error {
	// Buat direktori jika belum ada
	dir := filepath.Dir(a.LokasiKuki)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori: %w", err)
	}

	// Ambil kuki dari client
	u, _ := url.Parse(HostKBBI)
	cookies := a.client.Jar.Cookies(u)
	
	kukiData := make(map[string]string)
	for _, cookie := range cookies {
		if cookie.Name == NamaKukiUtama {
			kukiData[cookie.Name] = cookie.Value
			break
		}
	}

	// Simpan ke file JSON
	data, err := json.Marshal(kukiData)
	if err != nil {
		return fmt.Errorf("gagal mengenkode kuki: %w", err)
	}

	if err := os.WriteFile(a.LokasiKuki, data, 0600); err != nil {
		return fmt.Errorf("gagal menyimpan kuki: %w", err)
	}

	return nil
}

// MuatKuki memuat kuki autentikasi dari file
func (a *AutentikasiKBBI) MuatKuki() error {
	data, err := os.ReadFile(a.LokasiKuki)
	if err != nil {
		return fmt.Errorf("kuki tidak ditemukan pada %s", a.LokasiKuki)
	}

	var kukiData map[string]string
	if err := json.Unmarshal(data, &kukiData); err != nil {
		return fmt.Errorf("gagal membaca kuki: %w", err)
	}

	// Set kuki ke client
	u, _ := url.Parse(HostKBBI)
	for name, value := range kukiData {
		if name == NamaKukiUtama {
			cookie := &http.Cookie{
				Name:  name,
				Value: value,
			}
			a.client.Jar.SetCookies(u, []*http.Cookie{cookie})
			a.Terautentikasi = true
			break
		}
	}

	return nil
}

// GetClient mengembalikan http.Client yang sudah terautentikasi
func (a *AutentikasiKBBI) GetClient() *http.Client {
	return a.client
}

// CekAutentikasi memeriksa apakah sesi masih terautentikasi
func (a *AutentikasiKBBI) CekAutentikasi(htmlContent string) bool {
	a.Terautentikasi = !strings.Contains(htmlContent, "loginLink")
	return a.Terautentikasi
}

// ambilToken mengambil token CSRF dari halaman login
func (a *AutentikasiKBBI) ambilToken() (string, error) {
	resp, err := a.client.Get(fmt.Sprintf("%s/%s", HostKBBI, LokasiLogin))
	if err != nil {
		return "", fmt.Errorf("gagal mengakses halaman login: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca halaman login: %w", err)
	}

	// Cari token CSRF menggunakan regex
	re := regexp.MustCompile(`<input name="__RequestVerificationToken".*value="([^"]*)"`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("token CSRF tidak ditemukan")
	}

	return string(matches[1]), nil
}

// HapusKuki menghapus file kuki yang tersimpan
func (a *AutentikasiKBBI) HapusKuki() error {
	if err := os.Remove(a.LokasiKuki); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("kuki tidak ditemukan pada %s", a.LokasiKuki)
		}
		return fmt.Errorf("gagal menghapus kuki: %w", err)
	}
	return nil
}
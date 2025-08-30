// Package cache menyediakan fungsi cache untuk mengurangi request ke KBBI
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// Direktori cache relatif terhadap direktori kuki
	DirCache = "cache"
	
	// Durasi cache expired (30 hari)
	DurasiCache = 30 * 24 * time.Hour
)

// EntriCache merepresentasikan satu entri cache
type EntriCache struct {
	Kata      string    `json:"kata"`
	HTML      string    `json:"html"`
	Timestamp time.Time `json:"timestamp"`
	Expired   time.Time `json:"expired"`
}

// ManagerCache mengelola operasi cache
type ManagerCache struct {
	DirektorCache string
}

// BaruManagerCache membuat manager cache baru
func BaruManagerCache(lokasiKuki string) (*ManagerCache, error) {
	// Tentukan direktori cache berdasarkan lokasi kuki
	var dirCache string
	
	if lokasiKuki == "" {
		// Gunakan default path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan home directory: %w", err)
		}
		dirCache = filepath.Join(homeDir, ".kbbi", DirCache)
	} else {
		// Ambil directory dari lokasi kuki dan tambahkan subdirektori cache
		dirKuki := filepath.Dir(lokasiKuki)
		dirCache = filepath.Join(dirKuki, DirCache)
	}

	// Buat direktori cache jika belum ada
	if err := os.MkdirAll(dirCache, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori cache: %w", err)
	}

	return &ManagerCache{
		DirektorCache: dirCache,
	}, nil
}

// buatKey membuat key cache dari kata pencarian
func (m *ManagerCache) buatKey(kata string) string {
	hash := sha256.Sum256([]byte(kata))
	return hex.EncodeToString(hash[:])
}

// namaFileCache mengembalikan nama file cache untuk kata tertentu
func (m *ManagerCache) namaFileCache(kata string) string {
	key := m.buatKey(kata)
	return filepath.Join(m.DirektorCache, key+".json")
}

// AmbilCache mengambil data cache untuk kata tertentu
func (m *ManagerCache) AmbilCache(kata string) (string, bool) {
	namaFile := m.namaFileCache(kata)
	
	// Cek apakah file cache ada
	if _, err := os.Stat(namaFile); os.IsNotExist(err) {
		return "", false
	}

	// Baca file cache
	data, err := os.ReadFile(namaFile)
	if err != nil {
		return "", false
	}

	// Parse JSON
	var entri EntriCache
	if err := json.Unmarshal(data, &entri); err != nil {
		return "", false
	}

	// Cek apakah cache expired
	if time.Now().After(entri.Expired) {
		// Hapus file cache yang expired
		os.Remove(namaFile)
		return "", false
	}

	return entri.HTML, true
}

// SimpanCache menyimpan data HTML ke cache
func (m *ManagerCache) SimpanCache(kata, html string) error {
	namaFile := m.namaFileCache(kata)
	
	// Buat entri cache
	now := time.Now()
	entri := EntriCache{
		Kata:      kata,
		HTML:      html,
		Timestamp: now,
		Expired:   now.Add(DurasiCache),
	}

	// Konversi ke JSON
	data, err := json.MarshalIndent(entri, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal mengenkode cache: %w", err)
	}

	// Simpan ke file
	if err := os.WriteFile(namaFile, data, 0644); err != nil {
		return fmt.Errorf("gagal menyimpan cache: %w", err)
	}

	return nil
}

// HapusCache menghapus cache untuk kata tertentu
func (m *ManagerCache) HapusCache(kata string) error {
	namaFile := m.namaFileCache(kata)
	
	if err := os.Remove(namaFile); err != nil {
		if os.IsNotExist(err) {
			return nil // Tidak ada yang perlu dihapus
		}
		return fmt.Errorf("gagal menghapus cache: %w", err)
	}

	return nil
}

// BersihkanCacheExpired menghapus semua cache yang sudah expired
func (m *ManagerCache) BersihkanCacheExpired() error {
	files, err := filepath.Glob(filepath.Join(m.DirektorCache, "*.json"))
	if err != nil {
		return fmt.Errorf("gagal membaca direktori cache: %w", err)
	}

	for _, file := range files {
		// Baca file
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip file yang error
		}

		// Parse JSON
		var entri EntriCache
		if err := json.Unmarshal(data, &entri); err != nil {
			continue // Skip file yang rusak
		}

		// Hapus jika expired
		if time.Now().After(entri.Expired) {
			os.Remove(file)
		}
	}

	return nil
}

// HitungUkuranCache menghitung jumlah file cache dan total ukuran
func (m *ManagerCache) HitungUkuranCache() (int, int64, error) {
	files, err := filepath.Glob(filepath.Join(m.DirektorCache, "*.json"))
	if err != nil {
		return 0, 0, fmt.Errorf("gagal membaca direktori cache: %w", err)
	}

	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	return len(files), totalSize, nil
}

// HapusSemuaCache menghapus semua file cache
func (m *ManagerCache) HapusSemuaCache() error {
	files, err := filepath.Glob(filepath.Join(m.DirektorCache, "*.json"))
	if err != nil {
		return fmt.Errorf("gagal membaca direktori cache: %w", err)
	}

	for _, file := range files {
		os.Remove(file)
	}

	return nil
}
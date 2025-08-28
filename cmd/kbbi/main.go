// Package main menyediakan CLI untuk KBBI Go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZulfaNurhuda/GoKBBI.project/internal/auth"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/fetcher"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/model"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/parser"
)

var (
	// Flag untuk pencarian
	kata          = flag.String("kata", "", "kata yang ingin dicari dalam KBBI")
	outputJSON    = flag.Bool("json", false, "tampilkan hasil dalam format JSON")
	indentJSON    = flag.Bool("indent", false, "gunakan indentasi untuk output JSON")
	tanpaContoh   = flag.Bool("tanpa-contoh", false, "jangan tampilkan contoh penggunaan")
	tanpaTerkait  = flag.Bool("tanpa-terkait", false, "jangan tampilkan kata terkait")
	nonpengguna   = flag.Bool("nonpengguna", false, "nonaktifkan fitur khusus pengguna")

	// Flag untuk autentikasi
	email        = flag.String("email", "", "alamat email untuk autentikasi KBBI")
	sandi        = flag.String("sandi", "", "kata sandi untuk autentikasi KBBI")
	lokasiKuki   = flag.String("lokasi-kuki", "", "lokasi file kuki untuk autentikasi")
	autentikasi  = flag.Bool("autentikasi", false, "lakukan autentikasi dengan email dan sandi")
	hapusKuki = flag.Bool("bersihkan-kuki", false, "hapus kuki yang tersimpan")

	// Flag bantuan
	bantuan = flag.Bool("bantuan", false, "tampilkan bantuan penggunaan")
	help    = flag.Bool("help", false, "tampilkan bantuan penggunaan")
	version = flag.Bool("version", false, "tampilkan versi aplikasi")
)

const (
	AppName    = "GoKBBI"
	AppVersion = "1.0.0"
	AppDesc    = "Aplikasi CLI untuk mengakses KBBI Daring menggunakan Go"
)

func main() {
	flag.Parse()

	// Tampilkan bantuan jika diminta atau tidak ada kata
	if *bantuan || *help {
		tampilkanBantuan()
		return
	}

	// Tampilkan versi jika diminta
	if *version {
		fmt.Printf("%s v%s\n", AppName, AppVersion)
		return
	}

	// Handle perintah autentikasi
	if *autentikasi && (*email != "" || *sandi != "") {
		if err := lakukanAutentikasi(); err != nil {
			fmt.Fprintf(os.Stderr, "Error autentikasi: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle pembersihan kuki
	if *hapusKuki {
		if err := bersihkanKuki(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Jika tidak ada kata yang diberikan
	if *kata == "" {
		if len(flag.Args()) > 0 {
			*kata = flag.Args()[0]
		} else {
			fmt.Fprintf(os.Stderr, "Error: Tidak ada kata yang diberikan\n")
			fmt.Fprintf(os.Stderr, "Gunakan --bantuan untuk melihat panduan penggunaan\n")
			os.Exit(1)
		}
	}

	// Lakukan pencarian
	if err := lakukanPencarian(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// tampilkanBantuan menampilkan panduan penggunaan
func tampilkanBantuan() {
	fmt.Printf("%s v%s - %s\n\n", AppName, AppVersion, AppDesc)
	
	fmt.Println("PENGGUNAAN:")
	fmt.Printf("  %s [OPTIONS] <kata>\n", os.Args[0])
	fmt.Printf("  %s --kata <kata> [OPTIONS]\n\n", os.Args[0])
	
	fmt.Println("CONTOH:")
	fmt.Printf("  %s cinta\n", os.Args[0])
	fmt.Printf("  %s --kata rumah --json\n", os.Args[0])
	fmt.Printf("  %s --email user@email.com --sandi password --autentikasi\n\n", os.Args[0])
	
	fmt.Println("OPTIONS:")
	fmt.Println("  Pencarian:")
	fmt.Println("    --kata <kata>           Kata yang ingin dicari")
	fmt.Println("    --json                  Tampilkan hasil dalam format JSON")
	fmt.Println("    --indent                Gunakan indentasi untuk JSON (hanya dengan --json)")
	fmt.Println("    --tanpa-contoh          Jangan tampilkan contoh penggunaan")
	fmt.Println("    --tanpa-terkait         Jangan tampilkan kata terkait")
	fmt.Println("    --nonpengguna           Nonaktifkan fitur khusus pengguna")
	
	fmt.Println("\n  Autentikasi:")
	fmt.Println("    --email <email>         Alamat email akun KBBI")
	fmt.Println("    --sandi <password>      Kata sandi akun KBBI")
	fmt.Println("    --autentikasi           Lakukan proses autentikasi")
	fmt.Println("    --lokasi-kuki <path>    Lokasi file kuki (default: ~/.kbbi/kuki.json)")
	fmt.Println("    --bersihkan-kuki        Hapus kuki yang tersimpan")
	
	fmt.Println("\n  Lainnya:")
	fmt.Println("    --bantuan, --help       Tampilkan bantuan ini")
	fmt.Println("    --version               Tampilkan versi aplikasi")
	
	fmt.Println("\nCATATAN:")
	fmt.Println("  - Fitur pengguna terdaftar memerlukan autentikasi dengan akun KBBI")
	fmt.Println("  - Setelah autentikasi berhasil, kuki akan disimpan otomatis")
	fmt.Println("  - Gunakan pencarian secara wajar untuk menghindari pemblokiran akun")
}

// lakukanAutentikasi menangani proses autentikasi
func lakukanAutentikasi() error {
	if *email == "" {
		fmt.Print("Masukkan email: ")
		fmt.Scanln(email)
	}
	
	if *sandi == "" {
		fmt.Print("Masukkan kata sandi: ")
		fmt.Scanln(sandi)
	}

	// Buat objek autentikasi
	autentikasiObj, err := auth.BaruAuth(*email, *sandi, *lokasiKuki)
	if err != nil {
		return fmt.Errorf("gagal melakukan autentikasi: %w", err)
	}

	// Simpan kuki
	if err := autentikasiObj.SimpanKuki(); err != nil {
		return fmt.Errorf("gagal menyimpan kuki: %w", err)
	}

	fmt.Println("Autentikasi berhasil!")
	fmt.Printf("Kuki telah disimpan di: %s\n", autentikasiObj.LokasiKuki)
	fmt.Println("Kuki akan otomatis digunakan pada pencarian berikutnya.")
	
	return nil
}

// bersihkanKuki menghapus file kuki
func bersihkanKuki() error {
	lokasiFile := *lokasiKuki
	if lokasiFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("gagal mendapatkan home directory: %w", err)
		}
		lokasiFile = filepath.Join(homeDir, ".kbbi", "kuki.json")
	}

	if err := os.Remove(lokasiFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Kuki tidak ditemukan di: %s\n", lokasiFile)
			return nil
		}
		return fmt.Errorf("gagal menghapus kuki: %w", err)
	}

	fmt.Printf("Kuki berhasil dihapus dari: %s\n", lokasiFile)
	return nil
}

// lakukanPencarian menangani proses pencarian kata
func lakukanPencarian() error {
	var autentikasiObj *auth.AutentikasiKBBI
	var err error

	// Cek apakah ada kuki yang tersimpan (kecuali jika nonpengguna)
	if !*nonpengguna {
		// Coba buat autentikasi dengan kuki yang ada
		autentikasiObj, err = auth.BaruAuth("", "", *lokasiKuki)
		if err != nil {
			// Jika gagal, lanjutkan tanpa autentikasi
			autentikasiObj = nil
		}
	}

	// Ambil halaman KBBI
	html, err := fetcher.AmbilHalamanDenganRetry(*kata, autentikasiObj, 3)
	if err != nil {
		// Jika error adalah TidakDitemukan dan ada HTML, parse untuk saran
		if kesalahanKBBI, ok := err.(*fetcher.KesalahanKBBI); ok && kesalahanKBBI.Jenis == "TidakDitemukan" && html != "" {
			return tampilkanSaran(html, autentikasiObj)
		}
		return err
	}

	// Parse HTML menjadi definisi
	terautentikasi := autentikasiObj != nil && autentikasiObj.Terautentikasi
	definisi, err := parser.ParseDefinisi(html, terautentikasi)
	if err != nil {
		return fmt.Errorf("gagal parsing HTML: %w", err)
	}

	// Set pranala
	parser.SetPranala(definisi, *kata)

	// Tampilkan hasil
	return tampilkanHasil(definisi)
}

// tampilkanSaran menangani kasus entri tidak ditemukan dengan saran
func tampilkanSaran(html string, autentikasiObj *auth.AutentikasiKBBI) error {
	terautentikasi := autentikasiObj != nil && autentikasiObj.Terautentikasi
	definisi, err := parser.ParseDefinisi(html, terautentikasi)
	if err != nil {
		return fmt.Errorf("gagal parsing saran: %w", err)
	}

	if !*outputJSON {
		fmt.Printf("%s tidak ditemukan dalam KBBI.\n", *kata)
	}

	if len(definisi.SaranEntri) > 0 && (terautentikasi || *outputJSON) {
		return tampilkanHasil(definisi)
	}

	return nil
}

// tampilkanHasil menampilkan hasil pencarian
func tampilkanHasil(definisi *model.Definisi) error {
	if *outputJSON {
		jsonStr, err := definisi.ToJSON(*indentJSON)
		if err != nil {
			return fmt.Errorf("gagal mengkonversi ke JSON: %w", err)
		}
		fmt.Println(jsonStr)
	} else {
		// Filter output berdasarkan flag
		output := definisi.String()
		
		// Filter contoh jika diminta
		if *tanpaContoh {
			output = hapusContoh(output)
		}
		
		// Filter terkait jika diminta  
		if *tanpaTerkait {
			output = hapusTerkait(output)
		}
		
		fmt.Println(output)
	}
	
	return nil
}

// hapusContoh menghapus bagian contoh dari output
func hapusContoh(text string) string {
	lines := strings.Split(text, "\n")
	var hasil []string
	
	for _, line := range lines {
		// Hapus bagian setelah ": " yang merupakan contoh
		if idx := strings.Index(line, ": "); idx != -1 {
			// Kecuali jika baris dimulai dengan angka (penomoran makna)
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 0 && (trimmed[0] >= '1' && trimmed[0] <= '9') {
				// Ini adalah makna bernomor, hapus bagian contoh
				hasil = append(hasil, line[:idx])
			} else {
				hasil = append(hasil, line)
			}
		} else {
			hasil = append(hasil, line)
		}
	}
	
	return strings.Join(hasil, "\n")
}

// hapusTerkait menghapus bagian kata terkait dari output
func hapusTerkait(text string) string {
	lines := strings.Split(text, "\n")
	var hasil []string
	skipMode := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip baris yang dimulai dengan angka dan → (rujukan internal dalam makna)
		if skipMode || strings.HasPrefix(trimmed, "Kata Turunan") ||
		   strings.HasPrefix(trimmed, "Gabungan Kata") ||
		   strings.HasPrefix(trimmed, "Peribahasa") ||
		   strings.HasPrefix(trimmed, "Idiom") {
			skipMode = true
			continue
		}
		
		// Skip baris yang berisi rujukan internal
		if strings.Contains(trimmed, "→") {
			continue
		}
		
		// Reset skip mode jika menemukan baris kosong dan bukan dalam daftar terkait
		if skipMode && trimmed == "" {
			skipMode = false
			continue
		}
		
		// Skip jika dalam mode skip dan berisi ; (indikator daftar)
		if skipMode && strings.Contains(trimmed, ";") {
			continue
		}
		
		// Reset skip mode jika bukan bagian dari daftar terkait
		if skipMode && !strings.Contains(trimmed, ";") && trimmed != "" {
			skipMode = false
		}
		
		if !skipMode {
			hasil = append(hasil, line)
		}
	}
	
	return strings.Join(hasil, "\n")
}
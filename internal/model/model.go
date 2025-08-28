// Package model menyediakan struktur data untuk entri KBBI
package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Definisi merepresentasikan hasil pencarian dalam KBBI
type Definisi struct {
	Pranala    string   `json:"pranala"`
	Entri      []Entri  `json:"entri"`
	SaranEntri []string `json:"saran_entri,omitempty"`
}

// Entri merepresentasikan satu entri dalam KBBI
type Entri struct {
	Nama             string      `json:"nama"`
	Nomor            string      `json:"nomor"`
	KataDasar        []string    `json:"kata_dasar"`
	Varian           []string    `json:"varian"`
	BentukTidakBaku  []string    `json:"bentuk_tidak_baku,omitempty"`
	Pelafalan        string      `json:"pelafalan"`
	Makna            []Makna     `json:"makna"`
	Etimologi        *Etimologi  `json:"etimologi,omitempty"`
	KataTurunan      []string    `json:"kata_turunan,omitempty"`
	GabunganKata     []string    `json:"gabungan_kata,omitempty"`
	Peribahasa       []string    `json:"peribahasa,omitempty"`
	Idiom            []string    `json:"idiom,omitempty"`
}

// Makna merepresentasikan makna dari sebuah entri
type Makna struct {
	Kelas    []KelasKata `json:"kelas"`
	Submakna []string    `json:"submakna"`
	Info     string      `json:"info"`
	Contoh   []string    `json:"contoh"`
}

// KelasKata merepresentasikan kelas kata (noun, verb, dll)
type KelasKata struct {
	Kode      string `json:"kode"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
}

// Etimologi merepresentasikan asal usul kata
type Etimologi struct {
	Kelas     []string `json:"kelas"`
	Bahasa    string   `json:"bahasa"`
	AsalKata  string   `json:"asal_kata"`
	Pelafalan string   `json:"pelafalan"`
	Arti      []string `json:"arti"`
}

// String mengembalikan representasi string dari Definisi
func (d *Definisi) String() string {
	if len(d.SaranEntri) > 0 && len(d.Entri) == 0 {
		return fmt.Sprintf("Berikut beberapa saran entri lain yang mirip.\n%s", 
			strings.Join(d.SaranEntri, ", "))
	}
	
	var hasil []string
	for _, entri := range d.Entri {
		hasil = append(hasil, entri.String())
	}
	return strings.Join(hasil, "\n\n")
}

// String mengembalikan representasi string dari Entri
func (e *Entri) String() string {
	var hasil []string
	
	// Nama entri dengan kata dasar jika ada
	nama := e.Nama
	if e.Nomor != "" {
		nama += fmt.Sprintf(" (%s)", e.Nomor)
	}
	if len(e.KataDasar) > 0 {
		nama = fmt.Sprintf("%s » %s", strings.Join(e.KataDasar, " » "), nama)
	}
	
	// Tambahkan pelafalan jika ada
	if e.Pelafalan != "" {
		nama += fmt.Sprintf("  %s", e.Pelafalan)
	}
	hasil = append(hasil, nama)
	
	// Varian atau bentuk tidak baku
	if len(e.BentukTidakBaku) > 0 {
		hasil = append(hasil, fmt.Sprintf("bentuk tidak baku: %s", 
			strings.Join(e.BentukTidakBaku, ", ")))
	} else if len(e.Varian) > 0 {
		hasil = append(hasil, fmt.Sprintf("varian: %s", 
			strings.Join(e.Varian, ", ")))
	}
	
	// Etimologi
	if e.Etimologi != nil {
		hasil = append(hasil, fmt.Sprintf("Etimologi: %s", e.Etimologi.String()))
	}
	
	// Makna
	if len(e.Makna) > 0 {
		if len(e.Makna) > 1 {
			for i, makna := range e.Makna {
				hasil = append(hasil, fmt.Sprintf("%d. %s", i+1, makna.String()))
			}
		} else {
			hasil = append(hasil, e.Makna[0].String())
		}
	}
	
	// Kata terkait
	namaMurni := strings.ReplaceAll(e.Nama, ".", "")
	if len(e.KataTurunan) > 0 {
		hasil = append(hasil, fmt.Sprintf("\nKata Turunan\n%s", 
			strings.Join(e.KataTurunan, "; ")))
	}
	if len(e.GabunganKata) > 0 {
		hasil = append(hasil, fmt.Sprintf("\nGabungan Kata\n%s", 
			strings.Join(e.GabunganKata, "; ")))
	}
	if len(e.Peribahasa) > 0 {
		hasil = append(hasil, fmt.Sprintf("\nPeribahasa (mengandung [%s])\n%s", 
			namaMurni, strings.Join(e.Peribahasa, "; ")))
	}
	if len(e.Idiom) > 0 {
		hasil = append(hasil, fmt.Sprintf("\nIdiom (mengandung [%s])\n%s", 
			namaMurni, strings.Join(e.Idiom, "; ")))
	}
	
	return strings.Join(hasil, "\n")
}

// String mengembalikan representasi string dari Makna
func (m *Makna) String() string {
	var hasil []string
	
	// Kelas kata
	if len(m.Kelas) > 0 {
		var kelas []string
		for _, k := range m.Kelas {
			kelas = append(kelas, fmt.Sprintf("(%s)", k.Kode))
		}
		hasil = append(hasil, strings.Join(kelas, " "))
	}
	
	// Submakna
	hasil = append(hasil, strings.Join(m.Submakna, "; "))
	
	// Info tambahan
	if m.Info != "" {
		hasil = append(hasil, m.Info)
	}
	
	// Contoh
	if len(m.Contoh) > 0 {
		return fmt.Sprintf("%s: %s", strings.Join(hasil, "  "), 
			strings.Join(m.Contoh, "; "))
	}
	
	return strings.Join(hasil, "  ")
}

// String mengembalikan representasi string dari Etimologi
func (e *Etimologi) String() string {
	var hasil []string
	
	// Bahasa asal
	if e.Bahasa != "" {
		hasil = append(hasil, fmt.Sprintf("[%s]", e.Bahasa))
	}
	
	// Kelas kata
	if len(e.Kelas) > 0 {
		var kelas []string
		for _, k := range e.Kelas {
			kelas = append(kelas, fmt.Sprintf("(%s)", k))
		}
		hasil = append(hasil, strings.Join(kelas, " "))
	}
	
	// Kata asal dan pelafalan
	asalKata := e.AsalKata
	if e.Pelafalan != "" {
		asalKata += fmt.Sprintf(" %s", e.Pelafalan)
	}
	hasil = append(hasil, asalKata)
	
	// Arti
	if len(e.Arti) > 0 {
		return fmt.Sprintf("%s: %s", strings.Join(hasil, " "), 
			strings.Join(e.Arti, "; "))
	}
	
	return strings.Join(hasil, " ")
}

// ToJSON mengkonversi Definisi ke JSON string
func (d *Definisi) ToJSON(indent bool) (string, error) {
	var data []byte
	var err error
	
	if indent {
		data, err = json.MarshalIndent(d, "", "  ")
	} else {
		data, err = json.Marshal(d)
	}
	
	if err != nil {
		return "", fmt.Errorf("gagal mengkonversi ke JSON: %w", err)
	}
	
	return string(data), nil
}
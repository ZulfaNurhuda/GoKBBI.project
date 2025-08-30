// Package parser menyediakan fungsi untuk parsing HTML KBBI
package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ZulfaNurhuda/GoKBBI.project/internal/model"
)

// ParseDefinisi mengurai HTML menjadi struktur Definisi
func ParseDefinisi(html string, terautentikasi bool) (*model.Definisi, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("gagal parsing HTML: %w", err)
	}

	definisi := &model.Definisi{
		Entri:      []model.Entri{},
		Peribahasa: []string{},
		Idiom:      []string{},
		SaranEntri: []string{},
	}

	// Cek apakah ada saran entri (ketika entri tidak ditemukan)
	if strings.Contains(html, "Berikut beberapa saran entri lain yang mirip.") {
		definisi.SaranEntri = parseSaranEntri(doc)
		return definisi, nil
	}

	// Parse entri normal
	definisi.Entri = parseEntriList(doc, terautentikasi)
	
	// Parse Peribahasa dan Idiom di level definisi
	parsePeribahawanIdiom(doc, definisi)
	
	return definisi, nil
}

// parseSaranEntri mengurai saran entri dari HTML
func parseSaranEntri(doc *goquery.Document) []string {
	var saranEntri []string
	
	doc.Find(".col-md-3").Each(func(i int, s *goquery.Selection) {
		saran := strings.TrimSpace(s.Text())
		if saran != "" {
			saranEntri = append(saranEntri, saran)
		}
	})

	return saranEntri
}

// parseEntriList mengurai daftar entri dari HTML
func parseEntriList(doc *goquery.Document, terautentikasi bool) []model.Entri {
	var entris []model.Entri
	var currentEntri strings.Builder
	var finished bool
	
	// Cari elemen hr pertama sebagai penanda awal
	doc.Find("hr").First().NextAll().Each(func(i int, s *goquery.Selection) {
		// Jika menemukan hr tanpa style, itu penanda akhir
		if goquery.NodeName(s) == "hr" && s.AttrOr("style", "") == "" {
			if currentEntri.Len() > 0 {
				entri := parseEntri(currentEntri.String(), terautentikasi)
				if entri.Nama != "" {
					entris = append(entris, entri)
				}
			}
			finished = true
			return
		}

		// Jika menemukan h2, itu awal entri baru
		if goquery.NodeName(s) == "h2" {
			// Simpan entri sebelumnya jika ada
			if currentEntri.Len() > 0 {
				entri := parseEntri(currentEntri.String(), terautentikasi)
				if entri.Nama != "" {
					entris = append(entris, entri)
				}
				currentEntri.Reset()
			}
			
			// Skip jika ini adalah lampiran (style="color:gray")
			if s.AttrOr("style", "") == "color:gray" {
				return
			}
		}

		// Tambahkan HTML ke current entri
		html, _ := s.Html()
		currentEntri.WriteString(fmt.Sprintf("<%s>%s</%s>", goquery.NodeName(s), html, goquery.NodeName(s)))
	})

	// Proses entri terakhir hanya jika belum diproses
	if !finished && currentEntri.Len() > 0 {
		entri := parseEntri(currentEntri.String(), terautentikasi)
		if entri.Nama != "" {
			entris = append(entris, entri)
		}
	}

	return entris
}

// parseEntri mengurai satu entri dari HTML
func parseEntri(htmlEntri string, terautentikasi bool) model.Entri {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlEntri))
	if err != nil {
		return model.Entri{}
	}

	entri := model.Entri{
		KataDasar:       []string{},
		Varian:          []string{},
		BentukTidakBaku: []string{},
		Makna:           []model.Makna{},
	}

	// Parse header entri (h2)
	judul := doc.Find("h2").First()
	if judul.Length() == 0 {
		return entri
	}

	parseNamaEntri(judul, &entri)
	parseNomorEntri(judul, &entri)
	parseKataDasar(judul, &entri)
	parsePelafalan(judul, &entri)
	parseVarian(judul, &entri, terautentikasi)

	// Parse etimologi jika terautentikasi
	if terautentikasi {
		parseEtimologi(doc, &entri)
		parseTerkait(doc, &entri)
	}

	// Parse makna
	parseMakna(doc, &entri, terautentikasi)

	return entri
}

// parseNamaEntri mengurai nama entri
func parseNamaEntri(judul *goquery.Selection, entri *model.Entri) {
	// Ambil text dari elemen i jika ada, jika tidak ambil text biasa
	italic := judul.Find("i")
	if italic.Length() > 0 {
		entri.Nama = strings.TrimSpace(italic.Text())
	} else {
		// Ambil text direct children saja
		var textParts []string
		judul.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					textParts = append(textParts, text)
				}
			}
		})
		entri.Nama = strings.Join(textParts, " ")
	}
}

// parseNomorEntri mengurai nomor entri
func parseNomorEntri(judul *goquery.Selection, entri *model.Entri) {
	sup := judul.Find("sup").First()
	if sup.Length() > 0 {
		entri.Nomor = strings.TrimSpace(sup.Text())
	}
}

// parseKataDasar mengurai kata dasar
func parseKataDasar(judul *goquery.Selection, entri *model.Entri) {
	judul.Find(".rootword").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		if link.Length() > 0 {
			kata := ambilTeksDalamLabel(link)
			nomor := link.Find("sup")
			if nomor.Length() > 0 {
				kata = fmt.Sprintf("%s (%s)", kata, strings.TrimSpace(nomor.Text()))
			}
			entri.KataDasar = append(entri.KataDasar, kata)
		}
	})
}

// parsePelafalan mengurai pelafalan
func parsePelafalan(judul *goquery.Selection, entri *model.Entri) {
	lafal := judul.Find(".syllable")
	if lafal.Length() > 0 {
		entri.Pelafalan = strings.TrimSpace(lafal.Text())
	}
}

// parseVarian mengurai varian kata
func parseVarian(judul *goquery.Selection, entri *model.Entri, terautentikasi bool) {
	var varian *goquery.Selection

	if terautentikasi {
		// Untuk pengguna terautentikasi, cari small yang bukan entrisButton
		judul.Find("small").Each(func(i int, s *goquery.Selection) {
			span := s.Find("span")
			if span.Length() > 0 {
				class := span.AttrOr("class", "")
				if strings.Contains(class, "entrisButton") {
					return // skip ini
				}
			}
			varian = s
		})
	} else {
		varian = judul.Find("small").First()
	}

	if varian == nil || varian.Length() == 0 {
		return
	}

	// Cek apakah ada bentuk tidak baku (elemen b)
	bentukTidakBaku := varian.Find("b")
	if bentukTidakBaku.Length() > 0 {
		bentukTidakBaku.Each(func(i int, s *goquery.Selection) {
			nama := strings.TrimSpace(s.Text())
			nama = strings.TrimLeft(nama, ", ")
			
			nomor := s.Find("sup")
			if nomor.Length() > 0 {
				nama = fmt.Sprintf("%s (%s)", nama, strings.TrimSpace(nomor.Text()))
			}
			entri.BentukTidakBaku = append(entri.BentukTidakBaku, nama)
		})
	} else {
		// Varian biasa
		text := strings.TrimSpace(varian.Text())
		if strings.HasPrefix(text, "varian: ") {
			varianText := strings.TrimPrefix(text, "varian: ")
			entri.Varian = strings.Split(varianText, ", ")
			// Trim setiap varian
			for i, v := range entri.Varian {
				entri.Varian[i] = strings.TrimSpace(v)
			}
		}
	}
}

// parseEtimologi mengurai etimologi
func parseEtimologi(doc *goquery.Document, entri *model.Entri) {
	// Cari text "Etimologi:"
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Etimologi:") {
			// Ambil sibling berikutnya
			next := s.Next()
			if next.Length() > 0 {
				etimologiHTML, _ := next.Html()
				etimologi := parseEtimologiHTML(etimologiHTML)
				if etimologi != nil {
					entri.Etimologi = etimologi
				}
			}
		}
	})
}

// parseEtimologiHTML mengurai HTML etimologi
func parseEtimologiHTML(html string) *model.Etimologi {
	// Hapus kurung siku di awal dan akhir
	html = strings.TrimPrefix(html, "[")
	html = strings.TrimSuffix(html, "]")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
	}

	etimologi := &model.Etimologi{
		Kelas: []string{},
		Arti:  []string{},
	}

	// Parse bahasa (i dengan style color:darkred)
	doc.Find(`i[style*="color:darkred"]`).Each(func(i int, s *goquery.Selection) {
		etimologi.Bahasa = strings.TrimSpace(s.Text())
	})

	// Parse kelas (span dengan style color:red)
	doc.Find(`span[style*="color:red"]`).Each(func(i int, s *goquery.Selection) {
		kelas := strings.TrimSpace(s.Text())
		if kelas != "" {
			etimologi.Kelas = append(etimologi.Kelas, kelas)
		}
	})

	// Parse asal kata (b)
	doc.Find("b").Each(func(i int, s *goquery.Selection) {
		etimologi.AsalKata = strings.TrimSpace(s.Text())
	})

	// Parse pelafalan (span dengan style color:darkgreen)
	doc.Find(`span[style*="color:darkgreen"]`).Each(func(i int, s *goquery.Selection) {
		etimologi.Pelafalan = strings.TrimSpace(s.Text())
	})

	// Parse arti (ambil text keseluruhan dan bersihkan)
	text := doc.Text()
	text = strings.Trim(text, `'"`)
	if strings.Contains(text, "; ") {
		etimologi.Arti = strings.Split(text, "; ")
		// Trim setiap arti
		for i, arti := range etimologi.Arti {
			etimologi.Arti[i] = strings.TrimSpace(arti)
		}
	} else if text != "" {
		etimologi.Arti = []string{strings.TrimSpace(text)}
	}

	return etimologi
}

// parseTerkait mengurai kata terkait
func parseTerkait(doc *goquery.Document, entri *model.Entri) {
	// Mapping header ke field
	headerMap := map[string]*[]string{
		"Kata Turunan":   &entri.KataTurunan,
		"Gabungan Kata":  &entri.GabunganKata,
	}

	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
		headerText := strings.TrimSpace(s.Text())
		
		for header, field := range headerMap {
			if strings.Contains(headerText, header) {
				// Ambil link-link di sibling berikutnya
				next := s.Next()
				if next.Length() > 0 {
					next.Find("a").Each(func(j int, link *goquery.Selection) {
						text := strings.TrimSpace(link.Text())
						if text != "" {
							*field = append(*field, text)
						}
					})
				}
				break
			}
		}
	})
}

// parseMakna mengurai makna-makna entri
func parseMakna(doc *goquery.Document, entri *model.Entri, terautentikasi bool) {
	// Cari makna prakategorial (dengan color="darkgreen")
	prakategorial := doc.Find(`[color="darkgreen"]`)
	if prakategorial.Length() > 0 {
		makna := parseMaknaSingle(prakategorial.First())
		entri.Makna = append(entri.Makna, makna)
		return
	}

	// Parse makna dari li
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		// Skip jika terautentikasi dan mengandung "Usulkan makna baru"
		if terautentikasi && strings.Contains(s.Text(), "Usulkan makna baru") {
			return
		}

		// Skip jika li hanya berisi rujukan internal (dimulai dengan →)
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "→") {
			return
		}

		makna := parseMaknaSingle(s)
		
		// Filter tambahan: skip jika makna hanya berisi submakna dengan rujukan
		if len(makna.Submakna) > 0 {
			isAllRujukan := true
			for _, submakna := range makna.Submakna {
				if !strings.HasPrefix(strings.TrimSpace(submakna), "→") {
					isAllRujukan = false
					break
				}
			}
			// Skip jika semua submakna adalah rujukan dan tidak ada kelas kata
			if isAllRujukan && len(makna.Kelas) == 0 {
				return
			}
			
			entri.Makna = append(entri.Makna, makna)
		}
	})
}

// parseMaknaSingle mengurai satu makna
func parseMaknaSingle(s *goquery.Selection) model.Makna {
	makna := model.Makna{
		Kelas:    []model.KelasKata{},
		Submakna: []string{},
		Contoh:   []string{},
	}

	// Hapus entrisButton jika ada
	s.Find("span.entrisButton").Remove()

	// Parse kelas kata (span dengan color="red")
	s.Find(`[color="red"]`).Each(func(i int, kelasElement *goquery.Selection) {
		kelasElement.Find("span").Each(func(j int, span *goquery.Selection) {
			kode := strings.TrimSpace(span.Text())
			title := span.AttrOr("title", "")
			
			parts := strings.Split(title, ": ")
			nama := ""
			deskripsi := ""
			
			if len(parts) > 0 {
				nama = strings.TrimSpace(parts[0])
			}
			if len(parts) > 1 {
				deskripsi = strings.TrimSpace(parts[1])
			}

			if kode != "" {
				makna.Kelas = append(makna.Kelas, model.KelasKata{
					Kode:      kode,
					Nama:      nama,
					Deskripsi: deskripsi,
				})
			}
		})
	})

	// Parse kelas prakategorial
	if s.AttrOr("color", "") == "darkgreen" {
		kode := strings.TrimSpace(s.Text())
		title := s.AttrOr("title", "")
		
		parts := strings.Split(title, ": ")
		nama := ""
		deskripsi := ""
		
		if len(parts) > 0 {
			nama = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			deskripsi = strings.TrimSpace(parts[1])
		}

		makna.Kelas = append(makna.Kelas, model.KelasKata{
			Kode:      kode,
			Nama:      nama,
			Deskripsi: deskripsi,
		})
	}

	// Parse info tambahan (color="green")
	info := s.Find(`[color="green"]`)
	if info.Length() > 0 {
		infoText := strings.TrimSpace(info.Text())
		// Pastikan info tidak duplikat dengan kelas
		isDuplicate := false
		for _, kelas := range makna.Kelas {
			if infoText == kelas.Kode {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			makna.Info = infoText
		}
	}

	// Parse rujukan (link a tanpa span style color:red)
	rujukan := s.Find("a")
	if rujukan.Length() > 0 && rujukan.Find(`span[style*="color:red"]`).Length() == 0 {
		rujukanText := ambilTeksDalamLabel(rujukan)
		nomor := rujukan.Find("sup")
		if nomor.Length() > 0 {
			rujukanText = fmt.Sprintf("%s (%s)", rujukanText, strings.TrimSpace(nomor.Text()))
		}
		makna.Submakna = append(makna.Submakna, fmt.Sprintf("→ %s", rujukanText))
	} else if s.AttrOr("color", "") == "darkgreen" {
		// Prakategorial
		next := s.Get(0).NextSibling
		if next != nil {
			text := strings.TrimSpace(next.Data)
			if text != "" {
				makna.Submakna = append(makna.Submakna, text)
			}
		}
	} else {
		// Parse submakna normal
		text := ""
		s.Contents().Each(func(i int, content *goquery.Selection) {
			if goquery.NodeName(content) == "#text" {
				text += content.Text()
			} else if goquery.NodeName(content) != "font" {
				// Ignore font elements untuk warna
				text += content.Text()
			}
		})
		
		text = strings.TrimSpace(text)
		text = strings.TrimSuffix(text, ":")
		
		if text != "" {
			if strings.Contains(text, "; ") {
				makna.Submakna = strings.Split(text, "; ")
			} else {
				makna.Submakna = []string{text}
			}
		}
	}

	// Parse contoh (setelah ": ")
	fullText := s.Text()
	if idx := strings.Index(fullText, ": "); idx != -1 {
		contohText := strings.TrimSpace(fullText[idx+2:])
		if contohText != "" {
			makna.Contoh = strings.Split(contohText, "; ")
			// Trim setiap contoh
			for i, contoh := range makna.Contoh {
				makna.Contoh[i] = strings.TrimSpace(contoh)
			}
		}
	}

	// Trim semua submakna
	for i, sub := range makna.Submakna {
		makna.Submakna[i] = strings.TrimSpace(sub)
	}

	return makna
}

// ambilTeksDalamLabel mengambil text direct children dari element
func ambilTeksDalamLabel(s *goquery.Selection) string {
	var textParts []string
	
	s.Contents().Each(func(i int, content *goquery.Selection) {
		if goquery.NodeName(content) == "#text" {
			text := strings.TrimSpace(content.Text())
			if text != "" {
				textParts = append(textParts, text)
			}
		}
	})
	
	return strings.Join(textParts, " ")
}

// parsePeribahawanIdiom mengurai Peribahasa dan Idiom di level definisi
func parsePeribahawanIdiom(doc *goquery.Document, definisi *model.Definisi) {
	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
		headerText := strings.TrimSpace(s.Text())
		
		if strings.Contains(headerText, "Peribahasa") {
			// Ambil link-link di sibling berikutnya
			next := s.Next()
			if next.Length() > 0 {
				next.Find("a").Each(func(j int, link *goquery.Selection) {
					text := strings.TrimSpace(link.Text())
					if text != "" {
						definisi.Peribahasa = append(definisi.Peribahasa, text)
					}
				})
			}
		} else if strings.Contains(headerText, "Idiom") {
			// Ambil link-link di sibling berikutnya
			next := s.Next()
			if next.Length() > 0 {
				next.Find("a").Each(func(j int, link *goquery.Selection) {
					text := strings.TrimSpace(link.Text())
					if text != "" {
						definisi.Idiom = append(definisi.Idiom, text)
					}
				})
			}
		}
	})
}

// SetPranala mengatur pranala dalam definisi
func SetPranala(d *model.Definisi, kata string) {
	// Buat pranala berdasarkan kata
	d.Pranala = fmt.Sprintf("https://kbbi.kemdikbud.go.id/entri/%s", kata)
}
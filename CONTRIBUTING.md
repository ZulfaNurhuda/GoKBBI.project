# Contributing ke GoKBBI

Kami sangat menyambut kontribusi dari komunitas! Baik itu laporan bug, fitur baru, atau perbaikan dokumentasi, bantuan Anda selalu diterima dengan tangan terbuka. Silakan luangkan waktu sejenak untuk membaca dokumen ini agar proses kontribusi berjalan dengan lancar.

## Bagaimana Cara Berkontribusi?

### Melaporkan Bug

Jika Anda menemukan bug, silakan buka issue di repository GitHub kami. Saat melaporkan bug, mohon sertakan:

*   Deskripsi yang jelas dan ringkas tentang bug tersebut.
*   Langkah-langkah untuk mereproduksi masalah.
*   Perilaku yang diharapkan.
*   Perilaku yang sebenarnya terjadi.
*   Screenshot atau pesan error (jika ada).
*   Sistem operasi dan versi Go yang Anda gunakan.

### Menyarankan Peningkatan

Kami menyambut saran untuk fitur baru atau perbaikan. Silakan buka issue di GitHub dan jelaskan ide Anda secara detail. Jelaskan mengapa Anda pikir itu akan menjadi tambahan yang berharga untuk GoKBBI.

### Kontribusi Kode

Jika Anda ingin berkontribusi kode, silakan ikuti langkah-langkah berikut:

1.  **Fork repository** di GitHub.
2.  **Clone repository yang sudah di-fork** ke mesin lokal Anda.
    ```bash
    git clone https://github.com/username-anda/GoKBBI.project.git
    cd GoKBBI.project
    ```
3.  **Buat branch baru** untuk fitur atau perbaikan bug Anda.
    ```bash
    git checkout -b feature/nama-fitur-anda
    # atau
    git checkout -b bugfix/deskripsi-perbaikan
    ```
4.  **Lakukan perubahan Anda.** Pastikan kode Anda mengikuti gaya coding dan konvensi yang ada.
5.  **Tulis test** untuk perubahan Anda. Semua fitur baru dan perbaikan bug harus memiliki unit test atau integration test yang sesuai.
6.  **Jalankan test** untuk memastikan semuanya berjalan dengan benar dan tidak merusak fungsionalitas yang ada.
    ```bash
    go test ./...
    go test -v ./internal/...
    ```
7.  **Build proyek** untuk memastikan tidak ada compilation error.
    ```bash
    go build -o bin/kbbi cmd/kbbi/main.go
    go build .
    ```
8.  **Commit perubahan Anda** dengan pesan commit yang jelas dan ringkas. Ikuti spesifikasi [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) (contoh: `feat: tambah fitur baru`, `fix: perbaiki bug`).
    ```bash
    git add .
    git commit -m "feat: tambah fitur pencarian dengan filter etimologi"
    ```
9.  **Push branch Anda** ke repository yang sudah di-fork.
    ```bash
    git push origin feature/nama-fitur-anda
    ```
10. **Buka Pull Request (PR)** di repository GoKBBI utama. Berikan deskripsi detail tentang perubahan Anda dan referensikan issue terkait jika ada.

## Gaya Kode dan Konvensi

*   **Format:** Gunakan `gofmt` untuk memformat kode Go Anda secara konsisten.
*   **Naming:** Ikuti konvensi Go: PascalCase untuk exported identifiers, camelCase untuk unexported identifiers.
*   **Comments:** Tambahkan komentar GoDoc untuk fungsi dan tipe yang di-export. Gunakan bahasa Indonesia untuk konsistensi.
*   **Error Handling:** Gunakan error handling yang proper sesuai idiom Go. Jangan ignore error tanpa alasan yang jelas.
*   **Package Structure:** Ikuti struktur package yang ada dengan `internal/` untuk kode internal dan public API di root.
*   **Import Grouping:** Kelompokkan import: standard library, third-party packages, kemudian internal packages.

## Pedoman Khusus GoKBBI

*   **Parsing HTML:** Gunakan goquery untuk parsing HTML KBBI dengan hati-hati agar tidak merusak struktur data.
*   **HTTP Requests:** Implementasikan retry mechanism dan rate limiting untuk menghormati server KBBI.
*   **Authentication:** Jaga keamanan kredensial pengguna dan implementasikan cookie management yang aman.
*   **Error Messages:** Gunakan bahasa Indonesia untuk pesan error yang user-friendly.
*   **Testing:** Test dengan data real dari KBBI, tapi gunakan mock untuk CI/CD pipeline.

## Contoh Kontribusi yang Dibutuhkan

*   🐛 **Bug Fixes:** Perbaikan parsing untuk kasus edge case tertentu
*   ✨ **Features:** Fitur pencarian dengan filter (kelas kata, etimologi, dll)
*   📚 **Documentation:** Penambahan contoh penggunaan atau perbaikan dokumentasi
*   🔧 **Performance:** Optimisasi parsing atau HTTP requests
*   🧪 **Testing:** Penambahan test coverage untuk fitur yang belum di-test

## License

Dengan berkontribusi ke GoKBBI, Anda setuju bahwa kontribusi Anda akan dilisensikan di bawah MIT License. Lihat file `LICENSE` untuk detail lebih lanjut.

Terima kasih telah berkontribusi ke GoKBBI! ✨

---

**Catatan Penting:** Pastikan untuk menghormati Terms of Service KBBI Daring dalam semua kontribusi Anda. Jangan membuat fitur yang dapat membebani server KBBI secara berlebihan.
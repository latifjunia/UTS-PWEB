package main

import (
    "html/template"
    "log"
    "net/http"
    "strconv"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var db *gorm.DB

func connectDB() {
    var err error
    dsn := "root:@tcp(localhost:3306)/koleksibunga?charset=utf8mb4&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Gagal koneksi ke database:", err)
    }
}

func migrasiDB() {
    db.AutoMigrate(&JenisBunga{}, &Bunga{})
    log.Println("Migrasi berhasil.")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    var jenis []JenisBunga
    var bunga []Bunga

    db.Preload("Bunga").Find(&jenis)
    db.Find(&bunga)

    tmpl := template.Must(template.ParseFiles("template/index.html"))
    tmpl.Execute(w, map[string]interface{}{
        "Jenis": jenis,
        "Bunga": bunga,
    })
}

func tambahJenisHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        nama := r.FormValue("nama_jenis")
        db.Create(&JenisBunga{Nama: nama})
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func tambahBungaHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        nama := r.FormValue("nama_bunga")
        warna := r.FormValue("warna")
        jenisID, _ := strconv.Atoi(r.FormValue("jenis_id"))

        db.Create(&Bunga{
            Nama:    nama,
            Warna:   warna,
            JenisID: uint(jenisID),
        })
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hapusBungaHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    bid, _ := strconv.Atoi(id)
    db.Delete(&Bunga{}, bid)

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
    connectDB()
    migrasiDB()

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/tambahjenis", tambahJenisHandler)
    http.HandleFunc("/tambahbunga", tambahBungaHandler)
    http.HandleFunc("/hapusbunga", hapusBungaHandler)

    log.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
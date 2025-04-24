package main

type JenisBunga struct {
	ID    uint
	Nama  string
	Bunga []Bunga `gorm:"foreignKey:JenisID"`
}

type Bunga struct {
	ID      uint
	Nama    string
	Warna   string
	JenisID uint
}

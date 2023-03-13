package penyusutan

import (
	"pvg/simada/service-golang/domains"
	"time"
)

type Penyusutan interface{}

type BasePenyusutan struct {
	tableName                      struct{} `pg:"_"`
	Id                             int      `json:"id"`
	InventarisID                   int      `json:"inventaris_id"`
	TotalAtribusi                  float64  `json:"total_atribusi"`
	BebanPenyusutanSebelumAtribusi float64  `json:"beban_penyusutan_sebelum_atribusi"`
	BebanPenyusutanSetelahAtribusi float64  `json:"beban_penyusutan_setelah_atribusi"`
	BebanPenyusutanTahunBerkenaan  float64  `json:"beban_penyusutan_tahun_berkenaan"`

	NilaiBukuSebelumBulanAtribusi       float64 `json:"nilai_buku_sebelum_bulan_atribusi"`
	NilaiBukuSetelahBulanAtribusi       float64 `json:"nilai_buku_setelah_bulan_atribusi"`
	NilaiBukuTahunBerkenaan             float64 `json:"nilai_buku_tahun_berkenaan"`
	AkumulasiPenyusutanSDTahunBerkenaan float64 `json:"akumulasi_penyusutan_sd_tahun_berkenaan"`

	BulanAtribusi                     float64 `json:"bulan_atribusi"`
	PemakaianSDTahunBerkenaan         int     `json:"pemakaian_sd_tahun_berkenaan"`
	SisaUmurEkonomisSDTahunSebelumnya int     `json:"sisa_umur_ekonomis_sd_tahun_sebelumnya"`
	SisaUmurEkonomisSDBulanBerkenaan  int     `json:"sisa_umur_ekonomis_sd_bulan_berkenaan"`
	SisaUmurEkonomisSetelahAtribusi   int     `json:"sisa_umur_ekonomis_setelah_atribusi"`

	PenambahanAtribusi     float64    `json:"penambahan_atribusi"`
	PenambahanUmurEkonomis int        `json:"penambahan_umur_ekonomis"`
	RunningPenyusutan      *time.Time `json:"running_penyusutan"`
	CreatedDate            *time.Time `json:"created_date"`
	UpdatedDate            *time.Time `json:"updated_date"`
	domains.GenericModel
}

func (p *BasePenyusutan) Table() string {
	return ""
}

type PenyusutanTahun2014 struct {
	tableName struct{} `pg:"report_penyusutan_2014"`
	*BasePenyusutan
}

type PenyusutanTahun2015 struct {
	tableName struct{} `pg:"report_penyusutan_2015"`
	*BasePenyusutan
}

type PenyusutanTahun2016 struct {
	tableName struct{} `pg:"report_penyusutan_2016"`
	*BasePenyusutan
}

type PenyusutanTahun2017 struct {
	tableName struct{} `pg:"report_penyusutan_2017"`
	*BasePenyusutan
}

type PenyusutanTahun2018 struct {
	tableName struct{} `pg:"report_penyusutan_2018"`
	*BasePenyusutan
}

type PenyusutanTahun2019 struct {
	tableName struct{} `pg:"report_penyusutan_2019"`
	*BasePenyusutan
}

type PenyusutanTahun2020 struct {
	tableName struct{} `pg:"report_penyusutan_2020"`
	*BasePenyusutan
}

type PenyusutanTahun2021 struct {
	tableName struct{} `pg:"report_penyusutan_2021"`
	*BasePenyusutan
}

type PenyusutanTahun2022 struct {
	tableName struct{} `pg:"report_penyusutan_2022"`
	*BasePenyusutan
	Penyusutan
}

type PenyusutanTahun2023 struct {
	tableName struct{} `pg:"report_penyusutan_2023"`
	*BasePenyusutan
}

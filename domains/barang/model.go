package barang

import (
	"pvg/simada/service-golang/domains"
	"time"
)

type Model struct {
	tableName              struct{} `pg:"m_barang,discard_unknown_columns,alias:barang"`
	ID                     int      `json:"id"`
	NamaRekAset            string   `json:"nama_rek_aset"`
	KodeAkun               string   `json:"kode_akun"`
	KodeKelompok           string   `json:"kode_kelompok"`
	KodeJenis              string   `json:"kode_jenis"`
	KodeObjek              string   `json:"kode_objek"`
	KodeRincianObjek       string   `json:"kode_rincian_objek"`
	KodeSubRincianObjek    string   `json:"kode_sub_rincian_object"`
	KodeSubSubRincianObjek string   `json:"kode_sub_sub_rincian_object"`
	UmurEkonomis           int      `json:"umur_ekonomis"`
	domains.GenericModel   `swaggerignore:"true" json:"-"`
}

type MesinModel struct {
	// Merk   string
	tableName            struct{}         `pg:"detil_mesin,discard_unknown_columns,alias:mesin"`
	ID                   int              `json:"id"`
	Pidinventaris        int              `json:"pidinventaris"`
	Ukuran               string           `json:"ukuran"`
	Bahan                string           `json:"bahan"`
	Norangka             string           `json:"norangka"`
	Nomesin              string           `json:"nomesin"`
	Nopol                string           `json:"nopol"`
	BPKB                 string           `json:"bpkb"`
	Keterangan           string           `json:"keterangan"`
	Inventaris           *InventarisModel `pg:",fk:pidinventaris"`
	domains.GenericModel `swaggerignore:"true" json:"-"`
}

type InventarisModel struct {
	tableName            struct{}   `pg:"inventaris,discard_unknown_columns,alias:inventaris"`
	ID                   int        `json:"id"`
	PIDBarang            int        `pg:"pidbarang" json:"pidbarang"`
	TahunPerolehan       int        `pg:"tahun_perolehan" json:"tahun_perolehan"`
	TglPerolehan         *time.Time `pg:"tgl_perolehan" json:"tgl_perolehan"`
	HargaSatuan          float64    `pg:"harga_satuan" json:"harga_satuan"`
	KodeBarang           string     `json:"kode_barang"`
	Barang               *Model     `pg:",fk:pidbarang"`
	domains.GenericModel `swaggerignore:"true" json:"-"`
}

type ParamPagination struct {
	take   int    `example:"10"`
	page   int    `example:"1"`
	search string `example:"smk"`
}

type ParamCheckNumberPlate struct {
	NumberPlate string `example:"W 600 X"`
	Pidopd      int    `example:"1"`
	SubPidopd   int    `example:"2"`
	Pidupt      int    `example:"3"`
}

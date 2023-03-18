package inventaris

import (
	"pvg/simada/service-golang/domains"
	"time"
)

type Model struct {
	tableName      struct{}     `pg:"inventaris,discard_unknown_columns,alias:inventaris"`
	ID             int          `json:"id"`
	PIDBarang      int          `pg:"pidbarang" json:"pidbarang"`
	TahunPerolehan int          `pg:"tahun_perolehan" json:"tahun_perolehan"`
	TglPerolehan   *time.Time   `pg:"tgl_perolehan" json:"tgl_perolehan"`
	HargaSatuan    float64      `pg:"harga_satuan" json:"harga_satuan"`
	KodeBarang     string       `json:"kode_barang"`
	Barang         *BarangModel `pg:",fk:pidbarang"`
	domains.GenericModel
}

type BarangModel struct {
	tableName    struct{} `pg:"m_barang,discard_unknown_columns,alias:barang"`
	ID           int      `json:"id"`
	UmurEkonomis int      `json:"umur_ekonomis"`
	domains.GenericModel
}

func (i *Model) Table() string {
	return "inventaris"
}

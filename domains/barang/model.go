package barang

import "pvg/simada/service-golang/domains"

type Model struct {
	tableName    struct{} `pg:"m_barang,discard_unknown_columns,alias:barang"`
	ID           int      `json:"id"`
	UmurEkonomis int      `json:"umur_ekonomis"`
	domains.GenericModel
}

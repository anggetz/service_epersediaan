package kriteria_masa_manfaat

import "pvg/simada/service-golang/domains"

type Model struct {
	tableName     struct{} `pg:"m_kriteria_barang_penyusutan,discard_unknown_columns"`
	ID            int      `json:"id"`
	TahunTambahan int      `json:"tahun_tambahan"`
	domains.GenericModel
}

func (i *Model) Table() string {
	return "m_kriteria_barang_penyusutan"
}

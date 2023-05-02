package organisasi

type OrganisasiModel struct {
	tableName  struct{}         `pg:"m_organisasi,alias:mo,discard_unknown_columns"`
	ID         int              `json:"id"`
	PID        int              `json:"pid"`
	NAMA       string           `json:"nama"`
	ALAMAT     string           `json:"alamat"`
	AKTIF      int              `json:"aktif"`
	KODE       string           `json:"kode"`
	LEVEL      string           `json:"level"`
	SETTING    int              `json:"setting"`
	ORGANISASI *OrganisasiModel `json:"organisasi" sql:"-" pg:"rel:has-one,fk:pid"`
}

type Level struct {
	ID   int    `json:"id"`
	NAMA string `json:"nama"`
}

type ParamPagination struct {
	take   int    `example:"10"`
	page   int    `example:"1"`
	search string `example:"smk"`
	level  *int   `example:"2,1,0,0,-1,-2"`
}

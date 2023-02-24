package user

import organisasiDomain "pvg/simada/service-epersediaan/domains/organisasi"

type UserModel struct {
	tableName     struct{} `pg:"users"`
	ID            int      `json:"id"`
	USERNAME      string   `json:"username"`
	EMAIL         string   `json:"email"`
	PASSWORD      string   `json:"password"`
	NIP           string   `json:"nip"`
	NO_HP         string   `json:"no_hp"`
	TGL_LAHIR     string   `json:"tgl_lahir"`
	ROLE          int      `json:"rike"`
	AKTIF         string   `json:"aktif"`
	JABATAN       int      `json:"jabatan"`
	PidOrganisasi int      `json:"pid_organisasi" pg:"pid_organisasi"`
}

type ResponseIAM struct {
	Username   string
	Organisasi organisasiDomain.OrganisasiModel `pg:"fk:organisasi_id"`
}

type RequestToken struct {
	Username string `example:"NOVIRAHYANTI" json:"username"`
	Password string `example:"bdg230683" json:"password"`
}

type ResponseToken struct {
	Token string `json:"token"`
}

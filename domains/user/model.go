package user

type UserModel struct {
	tableName struct{} `pg:"users"`
	ID        int      `json:"id"`
	USERNAME  string   `json:"username"`
	EMAIL     string   `json:"email"`
	PASSWORD  string   `json:"password"`
	NIP       string   `json:"nip"`
	NO_HP     string   `json:"no_hp"`
	TGL_LAHIR string   `json:"tgl_lahir"`
	ROLE      int      `json:"rike"`
	AKTIF     string   `json:"aktif"`
	JABATAN   int      `json:"jabatan"`
}

type ResponseUser struct {
	Username string
}

type RequestToken struct {
	Username string `example:"NOVIRAHYANTI"`
	Password string `example:"bdg230683"`
}

type ResponseToken struct {
	Token string
}

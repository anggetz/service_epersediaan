package user

type UserModel struct {
	tableName struct{} `pg:"users"`
	Username  string
	Password  string
}

type RequestToken struct {
	Username string `example:"NOVIRAHYANTI"`
	Password string `example:"bdg230683"`
}

type ResponseToken struct {
	Token string
}

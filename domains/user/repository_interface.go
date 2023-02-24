package user

type Repository interface {
	GetByUsername(string) (*UserModel, error)
	GetAllUser() ([]UserModel, error)
}

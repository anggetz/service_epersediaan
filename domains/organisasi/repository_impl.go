package organisasi

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

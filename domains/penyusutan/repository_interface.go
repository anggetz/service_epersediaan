package penyusutan

type Repository[T Penyusutan] interface {
	GetAll(tableName string, page int, take int, offset int, search string) ([]T, error)
}

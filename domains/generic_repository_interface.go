package domains

import "github.com/go-pg/pg/v10/orm"

type GenericRepository[T GenericModel] interface {
	One(int) (T, error)
	Exists(func(*orm.Query)) (bool, error)
	All(int, int, func(*orm.Query)) ([]T, error)
	Count(int, int, func(*orm.Query)) (int, error)
	Insert(T) error
	Migrate(bool) GenericRepository[T]
	SetTableName(string) GenericRepository[T]
	GetTableName() string
}

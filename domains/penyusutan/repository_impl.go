package penyusutan

import (
	"pvg/simada/service-golang/util"
)

type RepositoryImpl[T Penyusutan] struct{}

func NewRepository[T Penyusutan]() Repository[T] {
	return &RepositoryImpl[T]{}
}

func (r *RepositoryImpl[T]) GetAll(tableName string, page int, take int, offset int, search string) ([]T, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	dataPenyusutan := []T{}

	err := dbUtil.GetDB().Model().
		TableExpr(tableName).
		Column("*").
		Limit(take).
		Offset(offset).
		Select(&dataPenyusutan)

	if err != nil {
		return nil, err
	}

	return dataPenyusutan, nil
}

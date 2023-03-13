package domains

import (
	"log"
	"pvg/simada/service-golang/util"

	"github.com/go-pg/pg/v10/orm"
)

type GenericRepositoryImpl[T GenericModel] struct {
	tableName string
}

func NewGenericRepository[T GenericModel]() GenericRepository[T] {
	return &GenericRepositoryImpl[T]{
		tableName: "",
	}
}

func (g *GenericRepositoryImpl[T]) Insert(data T) error {

	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	_, err := dbUtil.GetDB().Model(data).
		Table(g.GetTableName()).
		Insert()

	if err != nil {
		return err
	}

	return nil
}

func (g *GenericRepositoryImpl[T]) SetTableName(tableName string) GenericRepository[T] {
	g.tableName = tableName

	return g

}

func (g *GenericRepositoryImpl[T]) GetTableName() string {

	var data T

	retTableName := data.Table()

	if g.tableName != "" {
		retTableName = g.tableName
	}

	return retTableName
}

func (g *GenericRepositoryImpl[T]) Migrate(needMigrate bool) GenericRepository[T] {
	var dbModel T

	if needMigrate == false {
		return g
	}

	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	log.Println("migrating", g.GetTableName())

	qOrm := dbUtil.GetDB().Model(dbModel)
	qOrm.TableExpr(g.GetTableName())

	err := qOrm.CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})

	if err != nil {
		panic(err)
	}

	return g
}

func (g *GenericRepositoryImpl[T]) Exists(qFunc func(*orm.Query)) (bool, error) {

	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	qOrm := dbUtil.GetDB().Model().
		TableExpr(g.GetTableName())

	qFunc(qOrm)

	return qOrm.Exists()
}

func (g *GenericRepositoryImpl[T]) All(take int, offset int, qFunc func(*orm.Query)) ([]T, error) {
	data := []T{}

	var dModel T

	var qOrm *orm.Query

	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	if g.tableName != "" {
		qOrm = dbUtil.GetDB().Model().
			TableExpr(g.GetTableName())
	} else {
		qOrm = dbUtil.GetDB().Model(dModel)
	}

	if take > 0 {
		qOrm = qOrm.
			Limit(take).
			Offset(offset)
	}

	qFunc(qOrm)

	if g.tableName != "" {
		err := qOrm.Select(&data)
		if err != nil {
			return data, err
		}
	} else {
		err := qOrm.Select(&data)
		if err != nil {
			return data, err
		}
	}

	return data, nil
}

func (g *GenericRepositoryImpl[T]) One(id int) (T, error) {
	data := []T{}
	var dModel T

	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	err := dbUtil.GetDB().Model().
		TableExpr(dModel.Table()).
		Column("*").
		Where("id = ?", id).
		Limit(1).
		Select(&data)

	if err != nil {
		return dModel, err
	}

	return data[0], nil
}

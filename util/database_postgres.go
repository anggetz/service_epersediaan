package util

import (
	"fmt"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
)

type DatabasePostgresImpl struct {
	DB *pg.DB
}

type DatabasePostgres interface {
	Close() error
	Connect() error
	GetDB() *pg.DB
}

func NewDatabasePostgres() DatabasePostgres {

	return &DatabasePostgresImpl{}
}

func (d *DatabasePostgresImpl) Connect() error {
	db := pg.Connect(&pg.Options{
		User:        os.Getenv("DB_USERNAME"),
		Password:    os.Getenv("DB_PASSWORD"),
		Database:    os.Getenv("DB_DATABASE"),
		PoolSize:    50,
		PoolTimeout: time.Second * 30,
		Addr:        os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
	})

	d.DB = db

	return nil
}

func (d *DatabasePostgresImpl) Close() error {
	err := d.DB.Close()
	if err != nil {
		return fmt.Errorf("Could not close database connection: %v", err.Error())
	}

	return nil
}

func (d *DatabasePostgresImpl) GetDB() *pg.DB {
	return d.DB
}

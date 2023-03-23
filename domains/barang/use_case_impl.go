package barang

import (
	"fmt"
	"pvg/simada/service-golang/domains"

	"github.com/go-pg/pg/v10/orm"
)

type UseCaseImpl struct{}

func NewUseCase() UseCase {
	return &UseCaseImpl{}
}

func (u *UseCaseImpl) GetApelMaster(page int, offset int, search string) ([]*Model, int, error) {
	funcClause := func(isCountClause bool) func(*orm.Query) {
		return func(q *orm.Query) {
			q.Where("barang.kode_akun = ?", "1")
			q.Where("barang.kode_kelompok = ?", "3")
			q.Where("barang.kode_jenis = ?", "2")
			q.Where("barang.kode_objek = ?", "02")

			if search != "" && !isCountClause {
				q.Where("barang.nama_rek_aset LIKE '?'", search)
			}
		}
	}

	data, err := domains.NewGenericRepository[*Model]().All(page, offset, funcClause(false))

	if err != nil {
		return nil, 0, fmt.Errorf("could not get data barang :%v", err.Error())
	}

	totalData, err := domains.NewGenericRepository[*Model]().Count(page, offset, funcClause(true))

	if err != nil {
		return nil, 0, fmt.Errorf("could not total data barang :%v", err.Error())
	}

	return data, totalData, nil

}

func (u *UseCaseImpl) CheckPlatNumber(platNumber string, opdid int, opdid_cabang int, uptid int) (*MesinModel, error) {
	inventarises, err := domains.NewGenericRepository[*MesinModel]().All(-1, -1, func(q *orm.Query) {
		// q.Relation("Barang", func(q *orm.Query) (*orm.Query, error) {
		// 	// q.Where("barang.kode_akun = ?", "1")
		// 	// q.Where("barang.kode_kelompok = ?", "3")
		// 	// q.Where("barang.kode_jenis = ?", "2")
		// 	// q.Where("barang.kode_objek = ?", "02")

		// 	return q, nil
		// })

		q.Where("mesin.nopol = ?", platNumber)

		q.Relation("Inventaris", func(q *orm.Query) (*orm.Query, error) {
			if opdid != 0 {
				q.Where("inventaris.pidopd = ?", opdid)
			}

			if opdid_cabang != 0 {
				q.Where("inventaris.pidopd_cabang = ?", opdid_cabang)
			}

			if uptid != 0 {
				q.Where("inventaris.pidupt = ?", uptid)
			}

			return q, nil
		})

		q.Relation("Inventaris.Barang", func(q *orm.Query) (*orm.Query, error) {

			return q, nil
		})

	})

	if err != nil {
		return nil, fmt.Errorf("could not get data: %v", err.Error())
	}

	if len(inventarises) < 1 {
		return nil, fmt.Errorf("data not found")
	} else {
		return inventarises[0], nil
	}
}

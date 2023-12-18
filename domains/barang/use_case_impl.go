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

func (u *UseCaseImpl) GetRegisteredDataTransportation(page int, offset int, pidopd int, search string) ([]*MesinModel, int, error) {
	funcClause := func(isCountClause bool) func(*orm.Query) {
		return func(q *orm.Query) {
			q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				q.Where("mesin.nopol != ?", "-")
				q.Where("mesin.nopol != ?", "")
				return q, nil
			})

			if search != "" {
				q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
					q.WhereOr("mesin.nopol LIKE ?", "%"+search+"%")
					q.WhereOr("mesin.norangka LIKE ?", "%"+search+"%")
					q.WhereOr("mesin.nomesin LIKE ?", "%"+search+"%")

					return q, nil
				})

			}

			q.Relation("Inventaris", func(q *orm.Query) (*orm.Query, error) {
				if pidopd != 0 {
					q.Where("inventaris.pidopd = ?", pidopd)
				}

				return q, nil
			})

		}
	}

	data, err := domains.NewGenericRepository[*MesinModel]().All(page, offset, funcClause(false))
	if err != nil {
		return nil, 0, fmt.Errorf("could not get data barang :%v", err.Error())
	}

	totalData, err := domains.NewGenericRepository[*MesinModel]().Count(page, offset, funcClause(true))

	if err != nil {
		return nil, 0, fmt.Errorf("could not total data barang :%v", err.Error())
	}

	return data, totalData, nil
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

func (u *UseCaseImpl) GetApelMasterNonKendaraan(page int, offset int, opd_id int, search string) ([]*ResponseNonALatAngkut, int, error) {
	funcClause := func(isCountClause bool) func(*orm.Query) {

		return func(q *orm.Query) {
			q.ColumnExpr(
				`
				inventaris.kode_barang,
				m_barang.nama_rek_aset as nama_barang,
				m_merk_barang.nama as merk_type,
				inventaris.tahun_perolehan,
				inventaris.harga_satuan as nilai_perolehan,
				org1.nama as pengguna_barang,
				org2.nama as kuasa_pengguna_barang`,
			).
				Join(" INNER JOIN inventaris ON inventaris.id = mesin.pidinventaris").
				Join(" INNER JOIN m_barang ON m_barang.id = inventaris.pidbarang").
				Join(" LEFT JOIN m_merk_barang ON m_merk_barang.id = mesin.merk").
				Join(" INNER JOIN m_organisasi as org1 ON org1.id = inventaris.pidopd").
				Join(" INNER JOIN m_organisasi as org2 ON org2.id = inventaris.pidopd_cabang")

			q.Where("m_barang.kode_akun = ?", "1")
			q.Where("m_barang.kode_kelompok = ?", "3")
			q.Where("m_barang.kode_jenis = ?", "2")
			q.Where("m_barang.kode_objek != ?", "02")

			if search != "" && !isCountClause {
				q.Where("m_barang.nama_rek_aset LIKE '?'", search)
			}

			if opd_id != 0 {
				q.Where("org1.id = ?", opd_id)
			}

			value, _ := q.AppendQuery(orm.NewFormatter(), nil)
			fmt.Println(string(value))

		}
	}

	data, err := domains.NewGenericRepository[*ResponseNonALatAngkut]().All(page, offset, funcClause(false))

	if err != nil {
		return nil, 0, fmt.Errorf("could not get data barang :%v", err.Error())
	}

	totalData, err := domains.NewGenericRepository[*ResponseNonALatAngkut]().Count(page, offset, funcClause(true))

	if err != nil {
		return nil, 0, fmt.Errorf("could not total data barang :%v", err.Error())
	}

	return data, totalData, nil

}

func (u *UseCaseImpl) CheckPlatNumberChassisNumberAndMachineNumber(platNumber string, chassisNumber string, machineNumber string, opdid int, opdid_cabang int, uptid int) (*MesinModel, error) {
	if platNumber == "" && chassisNumber == "" && machineNumber == "" {
		return nil, fmt.Errorf("nopol, nomor rangka, atau nomor mesin salah satu harus terisi!")
	}

	inventarises, err := domains.NewGenericRepository[*MesinModel]().All(-1, -1, func(q *orm.Query) {
		// q.Relation("Barang", func(q *orm.Query) (*orm.Query, error) {
		// 	// q.Where("barang.kode_akun = ?", "1")
		// 	// q.Where("barang.kode_kelompok = ?", "3")
		// 	// q.Where("barang.kode_jenis = ?", "2")
		// 	// q.Where("barang.kode_objek = ?", "02")

		// 	return q, nil
		// })

		if platNumber != "" {
			q.Where("mesin.nopol = ?", platNumber)
		}

		if chassisNumber != "" {
			q.Where("mesin.norangka = ?", chassisNumber)
		}

		if machineNumber != "" {
			q.Where("mesin.nomesin = ?", machineNumber)
		}

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

		q.Relation("Inventaris.PenggunaBarang", func(q *orm.Query) (*orm.Query, error) {

			return q, nil
		})

		q.Relation("Inventaris.KuasaPenggunaBarang", func(q *orm.Query) (*orm.Query, error) {

			return q, nil
		})

		q.Relation("Inventaris.SubKuasaPenggunaBarang", func(q *orm.Query) (*orm.Query, error) {

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

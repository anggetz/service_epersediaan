package organisasi

import (
	"pvg/simada/service-golang/util"
)

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) GetOneOrganisasiById(id int) (*OrganisasiModel, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	organisasi := OrganisasiModel{}

	err := dbUtil.GetDB().Model(&organisasi).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}

	return &organisasi, nil
}

func (r *RepositoryImpl) GetAllOrganisasi(page int, take int, offset int, search string, level *int) ([]OrganisasiModel, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	OrganisasiModel := []OrganisasiModel{}

	q := dbUtil.GetDB().Model(&OrganisasiModel).
		Column("mo.*").
		Relation("ORGANISASI").
		Limit(take).
		Offset(offset).
		Where("mo.nama like ?", "%"+search+"%")

	if level != nil {
		q.Where("mo.level = ?", level)
	}

	err := q.Select()

	if err != nil {
		return nil, err
	}

	return OrganisasiModel, nil
}

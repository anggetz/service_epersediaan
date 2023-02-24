package organisasi

import (
	"pvg/simada/service-epersediaan/util"
)

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) GetAllOrganisasi(page int, take int, offset int, search string) ([]OrganisasiModel, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	OrganisasiModel := []OrganisasiModel{}

	err := dbUtil.GetDB().Model(&OrganisasiModel).
		Column("mo.*").
		Relation("ORGANISASI").
		Limit(take).
		Offset(offset).
		Where("mo.nama like ?", "%"+search+"%").
		Select()

	if err != nil {
		return nil, err
	}

	return OrganisasiModel, nil
}

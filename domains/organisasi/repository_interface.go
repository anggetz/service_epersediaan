package organisasi

type Repository interface {
	GetAllOrganisasi(page int, take int, offset int, search string, level *int) ([]OrganisasiModel, error)
	GetOneOrganisasiById(int) (*OrganisasiModel, error)
}

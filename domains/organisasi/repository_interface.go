package organisasi

type Repository interface {
	GetAllOrganisasi(page int, take int, offset int, search string) ([]OrganisasiModel, error)
	GetOneOrganisasiById(int) (*OrganisasiModel, error)
}

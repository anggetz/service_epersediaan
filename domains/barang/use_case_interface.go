package barang

type UseCase interface {
	GetApelMaster(page int, offset int, search string) ([]*Model, int, error)
	CheckPlatNumber(platNumber string, opdid int, opdid_cabang int, uptid int) (*MesinModel, error)
}

package barang

type UseCase interface {
	GetApelMaster(page int, offset int, search string) ([]*Model, int, error)
	GetRegisteredDataTransportation(page int, offset int, pidopd int, search string) ([]*MesinModel, int, error)
	CheckPlatNumber(platNumber string, opdid int, opdid_cabang int, uptid int) (*MesinModel, error)
	GetApelMasterNonKendaraan(page int, offset int, opd_id int, search string) ([]*ResponseNonALatAngkut, int, error)
	CheckPlatNumberChassisNumberAndMachineNumber(platNumber string, chassisNumber string, machineNumber string, opdid int, opdid_cabang int, uptid int) (*MesinModel, error)
}

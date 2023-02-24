package user

import (
	"pvg/simada/service-epersediaan/util"
)

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) GetByUsername(username string) (*UserModel, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	userModel := UserModel{}

	err := dbUtil.GetDB().Model(&userModel).Where("username = ?", username).Select()

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (r *RepositoryImpl) GetAllUser() ([]UserModel, error) {
	dbUtil := util.NewDatabasePostgres()
	dbUtil.Connect()

	defer dbUtil.Close()

	userModel := []UserModel{}

	err := dbUtil.GetDB().Model(&userModel).Select()

	if err != nil {
		return nil, err
	}

	return userModel, nil
}

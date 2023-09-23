package appInfoUseCases

import (
	"github.com/pandakn/cafe-beans/modules/appInfo"
	"github.com/pandakn/cafe-beans/modules/appInfo/appInfoRepositories"
)

type IAppInfoUseCase interface {
	FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error)
	InsertCategory(req []*appInfo.Category) error
	DeleteCategory(categoryId int) error
}

type appInfoUseCase struct {
	appInfoRepository appInfoRepositories.IAppInfoRepository
}

func AppInfoUseCase(appInfoRepository appInfoRepositories.IAppInfoRepository) IAppInfoUseCase {
	return &appInfoUseCase{
		appInfoRepository: appInfoRepository,
	}
}

func (u *appInfoUseCase) FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error) {
	categories, err := u.appInfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (u *appInfoUseCase) InsertCategory(req []*appInfo.Category) error {
	if err := u.appInfoRepository.InsertCategory(req); err != nil {
		return err
	}

	return nil
}

func (u *appInfoUseCase) DeleteCategory(categoryId int) error {
	if err := u.appInfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}

	return nil
}

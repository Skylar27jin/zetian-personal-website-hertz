package school_repo

import (
	"context"
	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// CreateSchool 创建学校
func CreateSchool(ctx context.Context, school *domain.School) error {
	return DB.DB.WithContext(ctx).Create(school).Error
}

// GetSchoolByID 按ID查询学校
func GetSchoolByID(ctx context.Context, id int64) (*domain.School, error) {
	var school domain.School
	err := DB.DB.WithContext(ctx).Where("id = ?", id).First(&school).Error
	if err != nil {
		return nil, err
	}
	return &school, nil
}

// GetSchoolByShortName 按简称（如BU）查询
func GetSchoolByShortName(ctx context.Context, shortName string) (*domain.School, error) {
	var school domain.School
	err := DB.DB.WithContext(ctx).Where("short_name = ?", shortName).First(&school).Error
	if err != nil {
		return nil, err
	}
	return &school, nil
}

// ListAllSchools 获取所有学校
func ListAllSchools(ctx context.Context) ([]domain.School, error) {
	var schools []domain.School
	err := DB.DB.WithContext(ctx).Find(&schools).Error
	return schools, err
}

func GetTableName() string {
	return "schools"
}

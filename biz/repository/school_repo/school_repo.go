package school_repo

import (
	"context"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// 按 id 查一条学校记录
func GetSchoolByID(ctx context.Context, id int64) (*domain.School, error) {
	var s domain.School
	if err := DB.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// 按 id 只取 name（如果你只想拿全名）
func GetSchoolNameByID(ctx context.Context, id int64) (string, error) {
	var s domain.School
	if err := DB.DB.WithContext(ctx).
		Select("id, name").
		Where("id = ?", id).
		First(&s).Error; err != nil {
		return "", err
	}
	return s.Name, nil
}

// 批量按 id 查学校，用于列表接口
func GetSchoolsByIDs(ctx context.Context, ids []int64) (map[int64]*domain.School, error) {
	if len(ids) == 0 {
		return map[int64]*domain.School{}, nil
	}

	var schools []domain.School
	if err := DB.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&schools).Error; err != nil {
		return nil, err
	}

	res := make(map[int64]*domain.School, len(schools))
	for i := range schools {
		s := schools[i]
		res[s.ID] = &s
	}
	return res, nil
}

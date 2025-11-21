package school_repo

import (
	"context"
	"fmt"
	"strings"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

var (
	//cache maps, to speed up school lookups and reduce DB load
	schoolIDToSchoolCache = make(map[int64]*domain.School)
	schoolNameToSchoolCache = make(map[string]*domain.School)
)

//every time the server starts, load all schools into memory
func InitSchoolCache() error {
	var schools []domain.School
	if err := DB.DB.Find(&schools).Error; err != nil {
		return err
	}

	for i := range schools {
		s := schools[i]
		schoolIDToSchoolCache[s.ID] = &s
		schoolNameToSchoolCache[s.Name] = &s
	}
	fmt.Println("len(schoolIDToSchoolCache): ", len(schoolIDToSchoolCache))
	fmt.Println("len(schoolNameToSchoolCache): ", len(schoolNameToSchoolCache))
	return nil
}

func GetSchoolByIDInCache(id int64) (*domain.School, error) {
	s, ok := schoolIDToSchoolCache[id]
	if !ok {
		return nil, fmt.Errorf("school not found in cache")
	}
	return s, nil
}

func GetSchoolByNameInCache(name string) (*domain.School, error) {
	s, ok := schoolNameToSchoolCache[name]
	if !ok {
		return nil, fmt.Errorf("school not found in cache")
	}
	return s, nil
}

func GetSchoolLikeNameInCache(name string, limit int) ([]domain.School, error) {
	similarSchools := make([]domain.School, 0, limit)
	for _, s := range schoolNameToSchoolCache {
		if s.Name == name {
			return []domain.School{*s}, nil
		}
		if likes(name, s.Name) {
			similarSchools = append(similarSchools, *s)
			if len(similarSchools) >= limit {
				break
			}
		}
	}
	return similarSchools, nil
}

// a simple substring match
func likes(name, s string) bool {
    if name == "" || s == "" {
        return false
    }
    name = strings.ToLower(name)
    s = strings.ToLower(s)
    return strings.Contains(s, name)
}

func GetAllSchoolsInCache() ([]*domain.School, error) {
	schools := make([]*domain.School, 0, len(schoolIDToSchoolCache))
	for _, s := range schoolIDToSchoolCache {
		schools = append(schools, s)
	}
	return schools, nil
}


func GetSchoolsByIDsInCache(ids []int64) (map[int64]*domain.School, error) {
	res := make(map[int64]*domain.School)
	for _, id := range ids {
		s, ok := schoolIDToSchoolCache[id]
		if ok {
			res[id] = s
		} else {
			res[id] = nil
		}
	}
	return res, nil
}


// get a school by its ID from DB
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

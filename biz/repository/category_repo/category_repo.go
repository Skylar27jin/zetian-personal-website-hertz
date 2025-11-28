package category_repo

import (
	"context"
	"fmt"
	"strings"

	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

var (
	// cache maps：加速板块查找
	categoryIDToCategoryCache  = make(map[int64]*domain.Category)
	categoryKeyToCategoryCache = make(map[string]*domain.Category)
)

// InitCategoryCache 在服务启动时加载所有 category 到内存
func InitCategoryCache() error {
	var categories []domain.Category
	if err := DB.DB.Find(&categories).Error; err != nil {
		return err
	}

	for i := range categories {
		c := categories[i]
		categoryIDToCategoryCache[c.ID] = &c
		categoryKeyToCategoryCache[strings.ToLower(c.Key)] = &c
	}
	fmt.Println("len(categoryIDToCategoryCache): ", len(categoryIDToCategoryCache))
	fmt.Println("len(categoryKeyToCategoryCache): ", len(categoryKeyToCategoryCache))
	return nil
}

func GetCategoryByIDInCache(id int64) (*domain.Category, error) {
	c, ok := categoryIDToCategoryCache[id]
	if !ok {
		return nil, fmt.Errorf("category not found in cache")
	}
	return c, nil
}

func GetCategoryByKeyInCache(key string) (*domain.Category, error) {
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}
	c, ok := categoryKeyToCategoryCache[strings.ToLower(key)]
	if !ok {
		return nil, fmt.Errorf("category not found in cache")
	}
	return c, nil
}

// 简单 substring 匹配，用于搜索
func likes(name, s string) bool {
	if name == "" || s == "" {
		return false
	}
	name = strings.ToLower(name)
	s = strings.ToLower(s)
	return strings.Contains(s, name)
}

// 按 name 模糊匹配（可选）
func GetCategoryLikeNameInCache(name string, limit int) ([]domain.Category, error) {
	name = strings.ToLower(name)
	similar := make([]domain.Category, 0, limit)
	for _, c := range categoryIDToCategoryCache {
		if strings.ToLower(c.Name) == name {
			return []domain.Category{*c}, nil
		}
		if likes(name, c.Name) {
			similar = append(similar, *c)
			if len(similar) >= limit {
				break
			}
		}
	}
	return similar, nil
}

func GetAllCategoriesInCache() ([]*domain.Category, error) {
	categories := make([]*domain.Category, 0, len(categoryIDToCategoryCache))
	for _, c := range categoryIDToCategoryCache {
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategoriesByIDsInCache(ids []int64) (map[int64]*domain.Category, error) {
	res := make(map[int64]*domain.Category)
	for _, id := range ids {
		c, ok := categoryIDToCategoryCache[id]
		if ok {
			res[id] = c
		} else {
			res[id] = nil
		}
	}
	return res, nil
}


func GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	var c domain.Category
	if err := DB.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

package service

import (
	"context"
	"goblog/internal/domain"
)

// CategoryService 分类服务实现
type CategoryService struct {
	categoryRepo domain.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo domain.CategoryRepository) domain.CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// Create 创建分类
func (s *CategoryService) Create(ctx context.Context, req *domain.CategoryCreateRequest) (*domain.Category, error) {
	// 检查分类名称是否已存在
	_, err := s.categoryRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrDuplicateResource
	}
	if err != domain.ErrNotFound {
		return nil, err
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.categoryRepo.Create(ctx, category)
}

// GetByID 根据ID获取分类
func (s *CategoryService) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// Update 更新分类
func (s *CategoryService) Update(ctx context.Context, id int, req *domain.CategoryUpdateRequest) (*domain.Category, error) {
	// 检查分类是否存在
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查新名称是否与其他分类重复
	existingCategory, err := s.categoryRepo.GetByName(ctx, req.Name)
	if err == nil && existingCategory.ID != id {
		return nil, domain.ErrDuplicateResource
	}
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.categoryRepo.Update(ctx, id, category)
}

// Delete 删除分类
func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.categoryRepo.Delete(ctx, id)
}

// List 获取分类列表
func (s *CategoryService) List(ctx context.Context) ([]*domain.Category, error) {
	return s.categoryRepo.List(ctx)
}

package service

import (
	"context"
	"goblog/internal/domain"
)

// TagService 标签服务实现
type TagService struct {
	tagRepo domain.TagRepository
}

// NewTagService 创建标签服务
func NewTagService(tagRepo domain.TagRepository) domain.TagService {
	return &TagService{tagRepo: tagRepo}
}

// Create 创建标签
func (s *TagService) Create(ctx context.Context, req *domain.TagCreateRequest) (*domain.Tag, error) {
	// 检查标签名称是否已存在
	_, err := s.tagRepo.GetByName(ctx, req.Name)
	if err == nil {
		return nil, domain.ErrDuplicateResource
	}
	if err != domain.ErrNotFound {
		return nil, err
	}

	tag := &domain.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	// 如果没有提供颜色，使用默认颜色
	if tag.Color == "" {
		tag.Color = "#007bff"
	}

	return s.tagRepo.Create(ctx, tag)
}

// GetByID 根据ID获取标签
func (s *TagService) GetByID(ctx context.Context, id int) (*domain.Tag, error) {
	return s.tagRepo.GetByID(ctx, id)
}

// Update 更新标签
func (s *TagService) Update(ctx context.Context, id int, req *domain.TagUpdateRequest) (*domain.Tag, error) {
	// 检查标签是否存在
	_, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查新名称是否与其他标签重复
	existingTag, err := s.tagRepo.GetByName(ctx, req.Name)
	if err == nil && existingTag.ID != id {
		return nil, domain.ErrDuplicateResource
	}
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}

	tag := &domain.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	// 如果没有提供颜色，使用默认颜色
	if tag.Color == "" {
		tag.Color = "#007bff"
	}

	return s.tagRepo.Update(ctx, id, tag)
}

// Delete 删除标签
func (s *TagService) Delete(ctx context.Context, id int) error {
	return s.tagRepo.Delete(ctx, id)
}

// List 获取标签列表
func (s *TagService) List(ctx context.Context) ([]*domain.Tag, error) {
	return s.tagRepo.List(ctx)
}

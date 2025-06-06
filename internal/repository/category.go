package repository

import (
	"context"
	"goblog/ent"
	"goblog/ent/category"
	"goblog/internal/domain"
)

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	client *ent.Client
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(client *ent.Client) domain.CategoryRepository {
	return &CategoryRepository{client: client}
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	create := r.client.Category.Create().
		SetName(cat.Name)

	if cat.Description != "" {
		create = create.SetDescription(cat.Description)
	}

	entCategory, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.entToDomain(entCategory), nil
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	entCategory, err := r.client.Category.Query().
		Where(category.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entCategory), nil
}

// Update 更新分类
func (r *CategoryRepository) Update(ctx context.Context, id int, cat *domain.Category) (*domain.Category, error) {
	update := r.client.Category.UpdateOneID(id).
		SetName(cat.Name)

	if cat.Description != "" {
		update = update.SetDescription(cat.Description)
	}

	entCategory, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entCategory), nil
}

// Delete 删除分类
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	err := r.client.Category.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

// List 获取分类列表
func (r *CategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	entCategories, err := r.client.Category.Query().
		Order(ent.Desc(category.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(entCategories))
	for i, entCategory := range entCategories {
		categories[i] = r.entToDomain(entCategory)
	}

	return categories, nil
}

// GetByName 根据名称获取分类
func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	entCategory, err := r.client.Category.Query().
		Where(category.Name(name)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entCategory), nil
}

// entToDomain 将ent实体转换为领域模型
func (r *CategoryRepository) entToDomain(entCategory *ent.Category) *domain.Category {
	return &domain.Category{
		ID:          entCategory.ID,
		Name:        entCategory.Name,
		Description: entCategory.Description,
		CreatedAt:   entCategory.CreatedAt,
		UpdatedAt:   entCategory.UpdatedAt,
	}
}

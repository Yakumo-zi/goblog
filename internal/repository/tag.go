package repository

import (
	"context"
	"goblog/ent"
	"goblog/ent/tag"
	"goblog/internal/domain"
)

// TagRepository 标签仓储实现
type TagRepository struct {
	client *ent.Client
}

// NewTagRepository 创建标签仓储
func NewTagRepository(client *ent.Client) domain.TagRepository {
	return &TagRepository{client: client}
}

// Create 创建标签
func (r *TagRepository) Create(ctx context.Context, t *domain.Tag) (*domain.Tag, error) {
	create := r.client.Tag.Create().
		SetName(t.Name)

	if t.Color != "" {
		create = create.SetColor(t.Color)
	}

	entTag, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.entToDomain(entTag), nil
}

// GetByID 根据ID获取标签
func (r *TagRepository) GetByID(ctx context.Context, id int) (*domain.Tag, error) {
	entTag, err := r.client.Tag.Query().
		Where(tag.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entTag), nil
}

// Update 更新标签
func (r *TagRepository) Update(ctx context.Context, id int, t *domain.Tag) (*domain.Tag, error) {
	update := r.client.Tag.UpdateOneID(id).
		SetName(t.Name)

	if t.Color != "" {
		update = update.SetColor(t.Color)
	}

	entTag, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entTag), nil
}

// Delete 删除标签
func (r *TagRepository) Delete(ctx context.Context, id int) error {
	err := r.client.Tag.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

// List 获取标签列表
func (r *TagRepository) List(ctx context.Context) ([]*domain.Tag, error) {
	entTags, err := r.client.Tag.Query().
		Order(ent.Desc(tag.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]*domain.Tag, len(entTags))
	for i, entTag := range entTags {
		tags[i] = r.entToDomain(entTag)
	}

	return tags, nil
}

// GetByName 根据名称获取标签
func (r *TagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	entTag, err := r.client.Tag.Query().
		Where(tag.Name(name)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entTag), nil
}

// GetByIDs 根据ID列表获取标签
func (r *TagRepository) GetByIDs(ctx context.Context, ids []int) ([]*domain.Tag, error) {
	entTags, err := r.client.Tag.Query().
		Where(tag.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]*domain.Tag, len(entTags))
	for i, entTag := range entTags {
		tags[i] = r.entToDomain(entTag)
	}

	return tags, nil
}

// entToDomain 将ent实体转换为领域模型
func (r *TagRepository) entToDomain(entTag *ent.Tag) *domain.Tag {
	return &domain.Tag{
		ID:        entTag.ID,
		Name:      entTag.Name,
		Color:     entTag.Color,
		CreatedAt: entTag.CreatedAt,
		UpdatedAt: entTag.UpdatedAt,
	}
}

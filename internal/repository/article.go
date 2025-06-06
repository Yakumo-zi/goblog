package repository

import (
	"context"
	"goblog/ent"
	"goblog/ent/article"
	"goblog/ent/category"
	"goblog/ent/tag"
	"goblog/internal/domain"
)

// ArticleRepository 文章仓储实现
type ArticleRepository struct {
	client *ent.Client
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(client *ent.Client) domain.ArticleRepository {
	return &ArticleRepository{client: client}
}

// Create 创建文章
func (r *ArticleRepository) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	create := r.client.Article.Create().
		SetTitle(article.Title).
		SetContent(article.Content).
		SetPublished(article.Published)

	if article.Summary != "" {
		create = create.SetSummary(article.Summary)
	}

	if article.Category != nil {
		create = create.SetCategoryID(article.Category.ID)
	}

	entArticle, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 添加标签关联
	if len(article.Tags) > 0 {
		tagIDs := make([]int, len(article.Tags))
		for i, tag := range article.Tags {
			tagIDs[i] = tag.ID
		}
		_, err = r.client.Article.UpdateOneID(entArticle.ID).AddTagIDs(tagIDs...).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, entArticle.ID)
}

// GetByID 根据ID获取文章
func (r *ArticleRepository) GetByID(ctx context.Context, id int) (*domain.Article, error) {
	entArticle, err := r.client.Article.Query().
		Where(article.ID(id)).
		WithCategory().
		WithTags().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.entToDomain(entArticle), nil
}

// Update 更新文章
func (r *ArticleRepository) Update(ctx context.Context, id int, article *domain.Article) (*domain.Article, error) {
	update := r.client.Article.UpdateOneID(id).
		SetTitle(article.Title).
		SetContent(article.Content).
		SetPublished(article.Published)

	if article.Summary != "" {
		update = update.SetSummary(article.Summary)
	}

	if article.Category != nil {
		update = update.SetCategoryID(article.Category.ID)
	} else {
		update = update.ClearCategory()
	}

	// 清除现有标签关联
	update = update.ClearTags()

	entArticle, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// 添加新的标签关联
	if len(article.Tags) > 0 {
		tagIDs := make([]int, len(article.Tags))
		for i, tag := range article.Tags {
			tagIDs[i] = tag.ID
		}
		_, err = r.client.Article.UpdateOneID(entArticle.ID).AddTagIDs(tagIDs...).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, entArticle.ID)
}

// Delete 删除文章
func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
	err := r.client.Article.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

// List 获取文章列表
func (r *ArticleRepository) List(ctx context.Context, params domain.QueryParams) ([]*domain.Article, int64, error) {
	query := r.client.Article.Query().
		WithCategory().
		WithTags().
		Order(ent.Desc(article.FieldCreatedAt))

	// 添加过滤条件
	if params.Published != nil {
		query = query.Where(article.Published(*params.Published))
	}

	if params.Search != "" {
		query = query.Where(article.Or(
			article.TitleContains(params.Search),
			article.ContentContains(params.Search),
		))
	}

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页
	if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	entArticles, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	articles := make([]*domain.Article, len(entArticles))
	for i, entArticle := range entArticles {
		articles[i] = r.entToDomain(entArticle)
	}

	return articles, int64(total), nil
}

// ListByCategory 按分类获取文章
func (r *ArticleRepository) ListByCategory(ctx context.Context, categoryID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	query := r.client.Article.Query().
		Where(article.HasCategoryWith(category.ID(categoryID))).
		WithCategory().
		WithTags().
		Order(ent.Desc(article.FieldCreatedAt))

	// 添加过滤条件
	if params.Published != nil {
		query = query.Where(article.Published(*params.Published))
	}

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页
	if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	entArticles, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	articles := make([]*domain.Article, len(entArticles))
	for i, entArticle := range entArticles {
		articles[i] = r.entToDomain(entArticle)
	}

	return articles, int64(total), nil
}

// ListByTag 按标签获取文章
func (r *ArticleRepository) ListByTag(ctx context.Context, tagID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	query := r.client.Article.Query().
		Where(article.HasTagsWith(tag.ID(tagID))).
		WithCategory().
		WithTags().
		Order(ent.Desc(article.FieldCreatedAt))

	// 添加过滤条件
	if params.Published != nil {
		query = query.Where(article.Published(*params.Published))
	}

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 分页
	if params.Page > 0 && params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	entArticles, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	articles := make([]*domain.Article, len(entArticles))
	for i, entArticle := range entArticles {
		articles[i] = r.entToDomain(entArticle)
	}

	return articles, int64(total), nil
}

// entToDomain 将ent实体转换为领域模型
func (r *ArticleRepository) entToDomain(entArticle *ent.Article) *domain.Article {
	article := &domain.Article{
		ID:        entArticle.ID,
		Title:     entArticle.Title,
		Content:   entArticle.Content,
		Summary:   entArticle.Summary,
		Published: entArticle.Published,
		CreatedAt: entArticle.CreatedAt,
		UpdatedAt: entArticle.UpdatedAt,
	}

	// 转换分类
	if cat := entArticle.Edges.Category; cat != nil {
		article.Category = &domain.Category{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	// 转换标签
	for _, entTag := range entArticle.Edges.Tags {
		article.Tags = append(article.Tags, domain.Tag{
			ID:        entTag.ID,
			Name:      entTag.Name,
			Color:     entTag.Color,
			CreatedAt: entTag.CreatedAt,
			UpdatedAt: entTag.UpdatedAt,
		})
	}

	return article
}

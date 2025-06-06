package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"goblog/internal/domain"
)

// ArticleService 文章服务实现
type ArticleService struct {
	articleRepo  domain.ArticleRepository
	categoryRepo domain.CategoryRepository
	tagRepo      domain.TagRepository
}

// NewArticleService 创建文章服务
func NewArticleService(
	articleRepo domain.ArticleRepository,
	categoryRepo domain.CategoryRepository,
	tagRepo domain.TagRepository,
) domain.ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, req *domain.ArticleCreateRequest) (*domain.Article, error) {
	article := &domain.Article{
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Published: req.Published,
	}

	// 验证分类是否存在
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		article.Category = category
	}

	// 验证标签是否存在
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}
		if len(tags) != len(req.TagIDs) {
			return nil, domain.ErrInvalidInput
		}
		article.Tags = make([]domain.Tag, len(tags))
		for i, tag := range tags {
			article.Tags[i] = *tag
		}
	}

	return s.articleRepo.Create(ctx, article)
}

// GetByID 根据ID获取文章
func (s *ArticleService) GetByID(ctx context.Context, id int) (*domain.Article, error) {
	return s.articleRepo.GetByID(ctx, id)
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, id int, req *domain.ArticleUpdateRequest) (*domain.Article, error) {
	// 检查文章是否存在
	_, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	article := &domain.Article{
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Published: req.Published,
	}

	// 验证分类是否存在
	if req.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		article.Category = category
	}

	// 验证标签是否存在
	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.GetByIDs(ctx, req.TagIDs)
		if err != nil {
			return nil, err
		}
		if len(tags) != len(req.TagIDs) {
			return nil, domain.ErrInvalidInput
		}
		article.Tags = make([]domain.Tag, len(tags))
		for i, tag := range tags {
			article.Tags[i] = *tag
		}
	}

	return s.articleRepo.Update(ctx, id, article)
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, id int) error {
	return s.articleRepo.Delete(ctx, id)
}

// List 获取文章列表
func (s *ArticleService) List(ctx context.Context, params domain.QueryParams) ([]*domain.Article, int64, error) {
	return s.articleRepo.List(ctx, params)
}

// ListByCategory 按分类获取文章
func (s *ArticleService) ListByCategory(ctx context.Context, categoryID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	// 验证分类是否存在
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	return s.articleRepo.ListByCategory(ctx, categoryID, params)
}

// ListByTag 按标签获取文章
func (s *ArticleService) ListByTag(ctx context.Context, tagID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	// 验证标签是否存在
	_, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, 0, err
	}

	return s.articleRepo.ListByTag(ctx, tagID, params)
}

// BackupAll 备份所有文章为ZIP压缩包
func (s *ArticleService) BackupAll(ctx context.Context) ([]byte, error) {
	// 获取所有文章（不分页）
	params := domain.QueryParams{
		Page:  1,
		Limit: 10000, // 设置一个足够大的值来获取所有文章
	}

	articles, _, err := s.articleRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("获取文章列表失败: %w", err)
	}

	// 创建ZIP缓冲区
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	// 准备备份数据结构
	type BackupData struct {
		BackupTime   time.Time         `json:"backup_time"`
		ArticleCount int               `json:"article_count"`
		Articles     []*domain.Article `json:"articles"`
	}

	backupData := BackupData{
		BackupTime:   time.Now(),
		ArticleCount: len(articles),
		Articles:     articles,
	}

	// 将数据序列化为JSON
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化文章数据失败: %w", err)
	}

	// 添加主备份文件到ZIP
	mainFile, err := zipWriter.Create("articles_backup.json")
	if err != nil {
		return nil, fmt.Errorf("创建备份文件失败: %w", err)
	}

	if _, err := mainFile.Write(jsonData); err != nil {
		return nil, fmt.Errorf("写入备份文件失败: %w", err)
	}

	// 为每个文章创建单独的文件（便于单独恢复）
	for _, article := range articles {
		filename := fmt.Sprintf("articles/%d_%s.json", article.ID, sanitizeFilename(article.Title))

		articleFile, err := zipWriter.Create(filename)
		if err != nil {
			continue // 跳过有问题的文件，不中断整个备份过程
		}

		articleData, err := json.MarshalIndent(article, "", "  ")
		if err != nil {
			continue
		}

		articleFile.Write(articleData)
	}

	// 添加备份信息文件
	infoFile, err := zipWriter.Create("backup_info.txt")
	if err == nil {
		info := fmt.Sprintf(`博客文章备份
备份时间: %s
文章总数: %d
备份格式: JSON
备份工具: goblog backend

文件说明:
- articles_backup.json: 完整的文章备份数据
- articles/: 每个文章的单独文件
- backup_info.txt: 备份信息（本文件）
`, backupData.BackupTime.Format("2006-01-02 15:04:05"), backupData.ArticleCount)

		infoFile.Write([]byte(info))
	}

	// 完成ZIP写入
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("完成ZIP文件写入失败: %w", err)
	}

	return buf.Bytes(), nil
}

// sanitizeFilename 清理文件名，移除不安全的字符
func sanitizeFilename(name string) string {
	// 简单的文件名清理，移除常见的不安全字符
	safeName := ""
	for _, r := range name {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			safeName += "_"
		default:
			if r < 32 || r > 126 { // 非打印字符
				safeName += "_"
			} else {
				safeName += string(r)
			}
		}
	}

	// 限制长度
	if len(safeName) > 50 {
		safeName = safeName[:50]
	}

	return safeName
}

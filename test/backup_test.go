package test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"goblog/internal/domain"
	"goblog/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestArticleService_BackupAll 测试文章备份功能
func TestArticleService_BackupAll(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 测试数据
	testArticles := []*domain.Article{
		{
			ID:        1,
			Title:     "Go语言入门",
			Content:   "这是一篇关于Go语言的入门文章",
			Summary:   "Go语言基础教程",
			Published: true,
			Category: &domain.Category{
				ID:   1,
				Name: "技术",
			},
			Tags: []domain.Tag{
				{ID: 1, Name: "Go语言", Color: "#00ADD8"},
				{ID: 2, Name: "编程", Color: "#007bff"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Docker容器化部署",
			Content:   "这是一篇关于Docker的文章",
			Summary:   "Docker部署指南",
			Published: true,
			Category: &domain.Category{
				ID:   2,
				Name: "运维",
			},
			Tags: []domain.Tag{
				{ID: 3, Name: "Docker", Color: "#2496ED"},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 设置Mock期望
	mockArticleRepo.On("List", mock.Anything, mock.AnythingOfType("domain.QueryParams")).
		Return(testArticles, int64(len(testArticles)), nil)

	// 执行测试
	ctx := context.Background()
	backupData, err := articleService.BackupAll(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, backupData)
	assert.Greater(t, len(backupData), 0, "备份数据不应为空")

	// 验证ZIP文件结构
	zipReader, err := zip.NewReader(bytes.NewReader(backupData), int64(len(backupData)))
	assert.NoError(t, err)

	// 检查ZIP文件中的文件
	expectedFiles := map[string]bool{
		"articles_backup.json": false,
		"backup_info.txt":      false,
	}

	articleFiles := 0
	for _, file := range zipReader.File {
		if file.Name == "articles_backup.json" {
			expectedFiles["articles_backup.json"] = true

			// 验证主备份文件内容
			rc, err := file.Open()
			assert.NoError(t, err)
			defer rc.Close()

			var backupContent struct {
				BackupTime   time.Time         `json:"backup_time"`
				ArticleCount int               `json:"article_count"`
				Articles     []*domain.Article `json:"articles"`
			}

			err = json.NewDecoder(rc).Decode(&backupContent)
			assert.NoError(t, err)
			assert.Equal(t, 2, backupContent.ArticleCount)
			assert.Len(t, backupContent.Articles, 2)
			assert.Equal(t, "Go语言入门", backupContent.Articles[0].Title)
			assert.Equal(t, "Docker容器化部署", backupContent.Articles[1].Title)

		} else if file.Name == "backup_info.txt" {
			expectedFiles["backup_info.txt"] = true

			// 验证备份信息文件
			rc, err := file.Open()
			assert.NoError(t, err)
			defer rc.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)
			content := buf.String()
			assert.Contains(t, content, "博客文章备份")
			assert.Contains(t, content, "文章总数: 2")

		} else if len(file.Name) > 9 && file.Name[:9] == "articles/" {
			articleFiles++

			// 验证单个文章文件
			rc, err := file.Open()
			assert.NoError(t, err)
			defer rc.Close()

			var article domain.Article
			err = json.NewDecoder(rc).Decode(&article)
			assert.NoError(t, err)
			assert.NotEmpty(t, article.Title)
		}
	}

	// 验证所有期望的文件都存在
	for fileName, found := range expectedFiles {
		assert.True(t, found, "文件 %s 应该存在于备份中", fileName)
	}

	// 验证单个文章文件数量
	assert.Equal(t, 2, articleFiles, "应该有2个单独的文章文件")

	// 验证Mock调用
	mockArticleRepo.AssertExpectations(t)
}

// TestArticleService_BackupAll_EmptyArticles 测试空文章列表的备份
func TestArticleService_BackupAll_EmptyArticles(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 设置Mock期望 - 返回空文章列表
	mockArticleRepo.On("List", mock.Anything, mock.AnythingOfType("domain.QueryParams")).
		Return([]*domain.Article{}, int64(0), nil)

	// 执行测试
	ctx := context.Background()
	backupData, err := articleService.BackupAll(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, backupData)
	assert.Greater(t, len(backupData), 0, "即使没有文章，备份数据也不应为空")

	// 验证ZIP文件结构
	zipReader, err := zip.NewReader(bytes.NewReader(backupData), int64(len(backupData)))
	assert.NoError(t, err)

	// 应该至少包含主备份文件和信息文件
	assert.GreaterOrEqual(t, len(zipReader.File), 2)

	// 验证主备份文件内容
	for _, file := range zipReader.File {
		if file.Name == "articles_backup.json" {
			rc, err := file.Open()
			assert.NoError(t, err)
			defer rc.Close()

			var backupContent struct {
				ArticleCount int `json:"article_count"`
			}

			err = json.NewDecoder(rc).Decode(&backupContent)
			assert.NoError(t, err)
			assert.Equal(t, 0, backupContent.ArticleCount)
			break
		}
	}

	// 验证Mock调用
	mockArticleRepo.AssertExpectations(t)
}

// TestArticleService_BackupAll_RepositoryError 测试仓储层错误
func TestArticleService_BackupAll_RepositoryError(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 设置Mock期望 - 返回错误
	mockArticleRepo.On("List", mock.Anything, mock.AnythingOfType("domain.QueryParams")).
		Return([]*domain.Article(nil), int64(0), domain.ErrNotFound)

	// 执行测试
	ctx := context.Background()
	backupData, err := articleService.BackupAll(ctx)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, backupData)
	assert.Contains(t, err.Error(), "获取文章列表失败")

	// 验证Mock调用
	mockArticleRepo.AssertExpectations(t)
}

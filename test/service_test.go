package test

import (
	"context"
	"testing"
	"time"

	"goblog/internal/domain"
	"goblog/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockArticleRepository 文章仓储Mock
type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	args := m.Called(ctx, article)
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id int) (*domain.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *MockArticleRepository) Update(ctx context.Context, id int, article *domain.Article) (*domain.Article, error) {
	args := m.Called(ctx, id, article)
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleRepository) List(ctx context.Context, params domain.QueryParams) ([]*domain.Article, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*domain.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) ListByCategory(ctx context.Context, categoryID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	args := m.Called(ctx, categoryID, params)
	return args.Get(0).([]*domain.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) ListByTag(ctx context.Context, tagID int, params domain.QueryParams) ([]*domain.Article, int64, error) {
	args := m.Called(ctx, tagID, params)
	return args.Get(0).([]*domain.Article), args.Get(1).(int64), args.Error(2)
}

// MockCategoryRepository 分类仓储Mock
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, id int, category *domain.Category) (*domain.Category, error) {
	args := m.Called(ctx, id, category)
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

// MockTagRepository 标签仓储Mock
type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Create(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	args := m.Called(ctx, tag)
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByID(ctx context.Context, id int) (*domain.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) Update(ctx context.Context, id int, tag *domain.Tag) (*domain.Tag, error) {
	args := m.Called(ctx, id, tag)
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTagRepository) List(ctx context.Context) ([]*domain.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByIDs(ctx context.Context, ids []int) ([]*domain.Tag, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Tag), args.Error(1)
}

// TestArticleService_Create 测试创建文章
func TestArticleService_Create(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 测试数据
	req := &domain.ArticleCreateRequest{
		Title:      "测试文章",
		Content:    "这是测试内容",
		Summary:    "测试摘要",
		Published:  true,
		CategoryID: &[]int{1}[0],
		TagIDs:     []int{1, 2},
	}

	expectedCategory := &domain.Category{
		ID:          1,
		Name:        "技术",
		Description: "技术分类",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedTags := []*domain.Tag{
		{
			ID:        1,
			Name:      "Go",
			Color:     "#00ADD8",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "测试",
			Color:     "#007bff",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	expectedArticle := &domain.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		Summary:   req.Summary,
		Published: req.Published,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Category:  expectedCategory,
		Tags:      []domain.Tag{*expectedTags[0], *expectedTags[1]},
	}

	// 设置Mock期望
	mockCategoryRepo.On("GetByID", mock.Anything, 1).Return(expectedCategory, nil)
	mockTagRepo.On("GetByIDs", mock.Anything, []int{1, 2}).Return(expectedTags, nil)
	mockArticleRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(expectedArticle, nil)

	// 执行测试
	ctx := context.Background()
	result, err := articleService.Create(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedArticle.Title, result.Title)
	assert.Equal(t, expectedArticle.Content, result.Content)
	assert.Equal(t, expectedArticle.Category.ID, result.Category.ID)
	assert.Len(t, result.Tags, 2)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
	mockTagRepo.AssertExpectations(t)
	mockArticleRepo.AssertExpectations(t)
}

// TestArticleService_GetByID 测试获取文章
func TestArticleService_GetByID(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 测试数据
	expectedArticle := &domain.Article{
		ID:        1,
		Title:     "测试文章",
		Content:   "测试内容",
		Published: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望
	mockArticleRepo.On("GetByID", mock.Anything, 1).Return(expectedArticle, nil)

	// 执行测试
	ctx := context.Background()
	result, err := articleService.GetByID(ctx, 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedArticle.ID, result.ID)
	assert.Equal(t, expectedArticle.Title, result.Title)

	// 验证Mock调用
	mockArticleRepo.AssertExpectations(t)
}

// TestArticleService_GetByID_NotFound 测试文章不存在的情况
func TestArticleService_GetByID_NotFound(t *testing.T) {
	// 准备Mock
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	articleService := service.NewArticleService(mockArticleRepo, mockCategoryRepo, mockTagRepo)

	// 设置Mock期望
	mockArticleRepo.On("GetByID", mock.Anything, 999).Return(nil, domain.ErrNotFound)

	// 执行测试
	ctx := context.Background()
	result, err := articleService.GetByID(ctx, 999)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, result)

	// 验证Mock调用
	mockArticleRepo.AssertExpectations(t)
}

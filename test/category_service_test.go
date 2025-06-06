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

// TestCategoryService_Create 测试创建分类
func TestCategoryService_Create(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 测试数据
	req := &domain.CategoryCreateRequest{
		Name:        "技术",
		Description: "技术相关分类",
	}

	expectedCategory := &domain.Category{
		ID:          1,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 设置Mock期望 - 检查名称不存在
	mockCategoryRepo.On("GetByName", mock.Anything, req.Name).Return(nil, domain.ErrNotFound)
	// 创建分类
	mockCategoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(expectedCategory, nil)

	// 执行测试
	ctx := context.Background()
	result, err := categoryService.Create(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCategory.Name, result.Name)
	assert.Equal(t, expectedCategory.Description, result.Description)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Create_Duplicate 测试创建重复分类名称
func TestCategoryService_Create_Duplicate(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 测试数据
	req := &domain.CategoryCreateRequest{
		Name:        "技术",
		Description: "技术相关分类",
	}

	existingCategory := &domain.Category{
		ID:          1,
		Name:        req.Name,
		Description: "已存在的分类",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 设置Mock期望 - 名称已存在
	mockCategoryRepo.On("GetByName", mock.Anything, req.Name).Return(existingCategory, nil)

	// 执行测试
	ctx := context.Background()
	result, err := categoryService.Create(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateResource, err)
	assert.Nil(t, result)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_GetByID 测试获取分类
func TestCategoryService_GetByID(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 测试数据
	expectedCategory := &domain.Category{
		ID:          1,
		Name:        "技术",
		Description: "技术相关分类",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 设置Mock期望
	mockCategoryRepo.On("GetByID", mock.Anything, 1).Return(expectedCategory, nil)

	// 执行测试
	ctx := context.Background()
	result, err := categoryService.GetByID(ctx, 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCategory.ID, result.ID)
	assert.Equal(t, expectedCategory.Name, result.Name)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Update 测试更新分类
func TestCategoryService_Update(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 测试数据
	req := &domain.CategoryUpdateRequest{
		Name:        "技术更新",
		Description: "更新后的技术分类",
	}

	existingCategory := &domain.Category{
		ID:          1,
		Name:        "技术",
		Description: "技术相关分类",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedCategory := &domain.Category{
		ID:          1,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   existingCategory.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// 设置Mock期望
	mockCategoryRepo.On("GetByID", mock.Anything, 1).Return(existingCategory, nil)
	mockCategoryRepo.On("GetByName", mock.Anything, req.Name).Return(nil, domain.ErrNotFound)
	mockCategoryRepo.On("Update", mock.Anything, 1, mock.AnythingOfType("*domain.Category")).Return(updatedCategory, nil)

	// 执行测试
	ctx := context.Background()
	result, err := categoryService.Update(ctx, 1, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedCategory.Name, result.Name)
	assert.Equal(t, updatedCategory.Description, result.Description)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_Delete 测试删除分类
func TestCategoryService_Delete(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 设置Mock期望
	mockCategoryRepo.On("Delete", mock.Anything, 1).Return(nil)

	// 执行测试
	ctx := context.Background()
	err := categoryService.Delete(ctx, 1)

	// 验证结果
	assert.NoError(t, err)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

// TestCategoryService_List 测试获取分类列表
func TestCategoryService_List(t *testing.T) {
	// 准备Mock
	mockCategoryRepo := new(MockCategoryRepository)

	// 创建服务
	categoryService := service.NewCategoryService(mockCategoryRepo)

	// 测试数据
	expectedCategories := []*domain.Category{
		{
			ID:          1,
			Name:        "技术",
			Description: "技术相关分类",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "生活",
			Description: "生活相关分类",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// 设置Mock期望
	mockCategoryRepo.On("List", mock.Anything).Return(expectedCategories, nil)

	// 执行测试
	ctx := context.Background()
	result, err := categoryService.List(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedCategories[0].Name, result[0].Name)
	assert.Equal(t, expectedCategories[1].Name, result[1].Name)

	// 验证Mock调用
	mockCategoryRepo.AssertExpectations(t)
}

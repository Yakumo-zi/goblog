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

// TestTagService_Create 测试创建标签
func TestTagService_Create(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	req := &domain.TagCreateRequest{
		Name:  "Go语言",
		Color: "#00ADD8",
	}

	expectedTag := &domain.Tag{
		ID:        1,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望 - 检查名称不存在
	mockTagRepo.On("GetByName", mock.Anything, req.Name).Return(nil, domain.ErrNotFound)
	// 创建标签
	mockTagRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Tag")).Return(expectedTag, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.Create(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTag.Name, result.Name)
	assert.Equal(t, expectedTag.Color, result.Color)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_Create_WithoutColor 测试创建标签（不提供颜色）
func TestTagService_Create_WithoutColor(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	req := &domain.TagCreateRequest{
		Name:  "Go语言",
		Color: "", // 不提供颜色
	}

	expectedTag := &domain.Tag{
		ID:        1,
		Name:      req.Name,
		Color:     "#007bff", // 应该使用默认颜色
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望 - 检查名称不存在
	mockTagRepo.On("GetByName", mock.Anything, req.Name).Return(nil, domain.ErrNotFound)
	// 创建标签
	mockTagRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Tag")).Return(expectedTag, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.Create(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTag.Name, result.Name)
	assert.Equal(t, "#007bff", result.Color) // 验证默认颜色

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_Create_Duplicate 测试创建重复标签名称
func TestTagService_Create_Duplicate(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	req := &domain.TagCreateRequest{
		Name:  "Go语言",
		Color: "#00ADD8",
	}

	existingTag := &domain.Tag{
		ID:        1,
		Name:      req.Name,
		Color:     "#FF0000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望 - 名称已存在
	mockTagRepo.On("GetByName", mock.Anything, req.Name).Return(existingTag, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.Create(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateResource, err)
	assert.Nil(t, result)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_GetByID 测试获取标签
func TestTagService_GetByID(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	expectedTag := &domain.Tag{
		ID:        1,
		Name:      "Go语言",
		Color:     "#00ADD8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望
	mockTagRepo.On("GetByID", mock.Anything, 1).Return(expectedTag, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.GetByID(ctx, 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTag.ID, result.ID)
	assert.Equal(t, expectedTag.Name, result.Name)
	assert.Equal(t, expectedTag.Color, result.Color)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_Update 测试更新标签
func TestTagService_Update(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	req := &domain.TagUpdateRequest{
		Name:  "Go语言更新",
		Color: "#FF0000",
	}

	existingTag := &domain.Tag{
		ID:        1,
		Name:      "Go语言",
		Color:     "#00ADD8",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedTag := &domain.Tag{
		ID:        1,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: existingTag.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// 设置Mock期望
	mockTagRepo.On("GetByID", mock.Anything, 1).Return(existingTag, nil)
	mockTagRepo.On("GetByName", mock.Anything, req.Name).Return(nil, domain.ErrNotFound)
	mockTagRepo.On("Update", mock.Anything, 1, mock.AnythingOfType("*domain.Tag")).Return(updatedTag, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.Update(ctx, 1, req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedTag.Name, result.Name)
	assert.Equal(t, updatedTag.Color, result.Color)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_Delete 测试删除标签
func TestTagService_Delete(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 设置Mock期望
	mockTagRepo.On("Delete", mock.Anything, 1).Return(nil)

	// 执行测试
	ctx := context.Background()
	err := tagService.Delete(ctx, 1)

	// 验证结果
	assert.NoError(t, err)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

// TestTagService_List 测试获取标签列表
func TestTagService_List(t *testing.T) {
	// 准备Mock
	mockTagRepo := new(MockTagRepository)

	// 创建服务
	tagService := service.NewTagService(mockTagRepo)

	// 测试数据
	expectedTags := []*domain.Tag{
		{
			ID:        1,
			Name:      "Go语言",
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

	// 设置Mock期望
	mockTagRepo.On("List", mock.Anything).Return(expectedTags, nil)

	// 执行测试
	ctx := context.Background()
	result, err := tagService.List(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedTags[0].Name, result[0].Name)
	assert.Equal(t, expectedTags[1].Name, result[1].Name)

	// 验证Mock调用
	mockTagRepo.AssertExpectations(t)
}

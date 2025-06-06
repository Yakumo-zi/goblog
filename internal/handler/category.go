package handler

import (
	"errors"
	"strconv"

	"goblog/internal/domain"
	"goblog/internal/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService domain.CategoryService
	validator       *validator.Validate
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validator:       validator.New(),
	}
}

// Create 创建分类
func (h *CategoryHandler) Create(c echo.Context) error {
	var req domain.CategoryCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	category, err := h.categoryService.Create(c.Request().Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Created(c, category)
}

// GetByID 获取单个分类
func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的分类ID")
	}

	category, err := h.categoryService.GetByID(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, category)
}

// Update 更新分类
func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的分类ID")
	}

	var req domain.CategoryUpdateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	category, err := h.categoryService.Update(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, category)
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的分类ID")
	}

	err = h.categoryService.Delete(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "分类删除成功"})
}

// List 获取分类列表
func (h *CategoryHandler) List(c echo.Context) error {
	categories, err := h.categoryService.List(c.Request().Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, categories)
}

// handleError 处理错误
func (h *CategoryHandler) handleError(c echo.Context, err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return response.NotFound(c, "分类不存在")
	}
	if errors.Is(err, domain.ErrInvalidInput) {
		return response.BadRequest(c, "无效的输入参数")
	}
	if errors.Is(err, domain.ErrDuplicateResource) {
		return response.BadRequest(c, "分类名称已存在")
	}
	return response.InternalServerError(c, "内部服务器错误")
}

package handler

import (
	"errors"
	"strconv"

	"goblog/internal/domain"
	"goblog/internal/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService domain.TagService
	validator  *validator.Validate
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagService domain.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		validator:  validator.New(),
	}
}

// Create 创建标签
func (h *TagHandler) Create(c echo.Context) error {
	var req domain.TagCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	tag, err := h.tagService.Create(c.Request().Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Created(c, tag)
}

// GetByID 获取单个标签
func (h *TagHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的标签ID")
	}

	tag, err := h.tagService.GetByID(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, tag)
}

// Update 更新标签
func (h *TagHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的标签ID")
	}

	var req domain.TagUpdateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	tag, err := h.tagService.Update(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, tag)
}

// Delete 删除标签
func (h *TagHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的标签ID")
	}

	err = h.tagService.Delete(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "标签删除成功"})
}

// List 获取标签列表
func (h *TagHandler) List(c echo.Context) error {
	tags, err := h.tagService.List(c.Request().Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, tags)
}

// handleError 处理错误
func (h *TagHandler) handleError(c echo.Context, err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return response.NotFound(c, "标签不存在")
	}
	if errors.Is(err, domain.ErrInvalidInput) {
		return response.BadRequest(c, "无效的输入参数")
	}
	if errors.Is(err, domain.ErrDuplicateResource) {
		return response.BadRequest(c, "标签名称已存在")
	}
	return response.InternalServerError(c, "内部服务器错误")
}

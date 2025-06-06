package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"goblog/internal/domain"
	"goblog/internal/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService domain.ArticleService
	validator      *validator.Validate
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(articleService domain.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		validator:      validator.New(),
	}
}

// Create 创建文章
func (h *ArticleHandler) Create(c echo.Context) error {
	var req domain.ArticleCreateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	article, err := h.articleService.Create(c.Request().Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Created(c, article)
}

// GetByID 获取单个文章
func (h *ArticleHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的文章ID")
	}

	article, err := h.articleService.GetByID(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, article)
}

// Update 更新文章
func (h *ArticleHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的文章ID")
	}

	var req domain.ArticleUpdateRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "无效的请求参数")
	}

	if err := h.validator.Struct(&req); err != nil {
		return response.BadRequest(c, "请求参数验证失败")
	}

	article, err := h.articleService.Update(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, article)
}

// Delete 删除文章
func (h *ArticleHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "无效的文章ID")
	}

	err = h.articleService.Delete(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "文章删除成功"})
}

// List 获取文章列表
func (h *ArticleHandler) List(c echo.Context) error {
	params := h.parseQueryParams(c)

	articles, total, err := h.articleService.List(c.Request().Context(), params)
	if err != nil {
		return h.handleError(c, err)
	}

	if params.Page > 0 && params.Limit > 0 {
		meta := response.PageMeta{
			Page:      params.Page,
			Limit:     params.Limit,
			Total:     total,
			TotalPage: int((total + int64(params.Limit) - 1) / int64(params.Limit)),
		}
		return response.SuccessPaged(c, articles, meta)
	}

	return response.Success(c, articles)
}

// ListByCategory 按分类获取文章
func (h *ArticleHandler) ListByCategory(c echo.Context) error {
	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		return response.BadRequest(c, "无效的分类ID")
	}

	params := h.parseQueryParams(c)

	articles, total, err := h.articleService.ListByCategory(c.Request().Context(), categoryID, params)
	if err != nil {
		return h.handleError(c, err)
	}

	if params.Page > 0 && params.Limit > 0 {
		meta := response.PageMeta{
			Page:      params.Page,
			Limit:     params.Limit,
			Total:     total,
			TotalPage: int((total + int64(params.Limit) - 1) / int64(params.Limit)),
		}
		return response.SuccessPaged(c, articles, meta)
	}

	return response.Success(c, articles)
}

// ListByTag 按标签获取文章
func (h *ArticleHandler) ListByTag(c echo.Context) error {
	tagID, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		return response.BadRequest(c, "无效的标签ID")
	}

	params := h.parseQueryParams(c)

	articles, total, err := h.articleService.ListByTag(c.Request().Context(), tagID, params)
	if err != nil {
		return h.handleError(c, err)
	}

	if params.Page > 0 && params.Limit > 0 {
		meta := response.PageMeta{
			Page:      params.Page,
			Limit:     params.Limit,
			Total:     total,
			TotalPage: int((total + int64(params.Limit) - 1) / int64(params.Limit)),
		}
		return response.SuccessPaged(c, articles, meta)
	}

	return response.Success(c, articles)
}

// parseQueryParams 解析查询参数
func (h *ArticleHandler) parseQueryParams(c echo.Context) domain.QueryParams {
	params := domain.QueryParams{}

	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil && page > 0 {
		params.Page = page
	}

	if limit, err := strconv.Atoi(c.QueryParam("limit")); err == nil && limit > 0 && limit <= 100 {
		params.Limit = limit
	} else if params.Page > 0 {
		params.Limit = 10 // 默认限制
	}

	if published := c.QueryParam("published"); published != "" {
		if published == "true" {
			b := true
			params.Published = &b
		} else if published == "false" {
			b := false
			params.Published = &b
		}
	}

	params.Search = c.QueryParam("search")

	return params
}

// handleError 处理错误
func (h *ArticleHandler) handleError(c echo.Context, err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return response.NotFound(c, "文章不存在")
	}
	if errors.Is(err, domain.ErrInvalidInput) {
		return response.BadRequest(c, "无效的输入参数")
	}
	if errors.Is(err, domain.ErrDuplicateResource) {
		return response.BadRequest(c, "资源已存在")
	}
	return response.InternalServerError(c, "内部服务器错误")
}

// Backup 备份所有文章
func (h *ArticleHandler) Backup(c echo.Context) error {
	// 生成备份数据
	backupData, err := h.articleService.BackupAll(c.Request().Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("articles_backup_%s.zip", timestamp)

	// 设置响应头
	c.Response().Header().Set("Content-Type", "application/zip")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(backupData)))

	// 直接返回ZIP数据
	return c.Blob(200, "application/zip", backupData)
}

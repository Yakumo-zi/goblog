package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PagedResponse 分页响应格式
type PagedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    PageMeta    `json:"meta"`
}

// PageMeta 分页元信息
type PageMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// Success 成功响应
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// Created 创建成功响应
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

// Error 错误响应
func Error(c echo.Context, code int, message string) error {
	return c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message)
}

// NotFound 404错误
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message)
}

// InternalServerError 500错误
func InternalServerError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message)
}

// SuccessPaged 分页成功响应
func SuccessPaged(c echo.Context, data interface{}, meta PageMeta) error {
	return c.JSON(http.StatusOK, PagedResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Meta:    meta,
	})
}

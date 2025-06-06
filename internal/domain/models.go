package domain

import "time"

// Article 文章领域模型
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Category  *Category `json:"category,omitempty"`
	Tags      []Tag     `json:"tags,omitempty"`
}

// Category 分类领域模型
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag 标签领域模型
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleCreateRequest 创建文章请求
type ArticleCreateRequest struct {
	Title      string `json:"title" validate:"required,min=1,max=200"`
	Content    string `json:"content" validate:"required,min=1"`
	Summary    string `json:"summary" validate:"max=500"`
	Published  bool   `json:"published"`
	CategoryID *int   `json:"category_id"`
	TagIDs     []int  `json:"tag_ids"`
}

// ArticleUpdateRequest 更新文章请求
type ArticleUpdateRequest struct {
	Title      string `json:"title" validate:"required,min=1,max=200"`
	Content    string `json:"content" validate:"required,min=1"`
	Summary    string `json:"summary" validate:"max=500"`
	Published  bool   `json:"published"`
	CategoryID *int   `json:"category_id"`
	TagIDs     []int  `json:"tag_ids"`
}

// CategoryCreateRequest 创建分类请求
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// CategoryUpdateRequest 更新分类请求
type CategoryUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// TagCreateRequest 创建标签请求
type TagCreateRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Color string `json:"color" validate:"hexcolor"`
}

// TagUpdateRequest 更新标签请求
type TagUpdateRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Color string `json:"color" validate:"hexcolor"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// QueryParams 查询参数
type QueryParams struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Published *bool  `query:"published"`
	Search    string `query:"search"`
}

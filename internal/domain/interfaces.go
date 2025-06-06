package domain

import "context"

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Create(ctx context.Context, article *Article) (*Article, error)
	GetByID(ctx context.Context, id int) (*Article, error)
	Update(ctx context.Context, id int, article *Article) (*Article, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params QueryParams) ([]*Article, int64, error)
	ListByCategory(ctx context.Context, categoryID int, params QueryParams) ([]*Article, int64, error)
	ListByTag(ctx context.Context, tagID int, params QueryParams) ([]*Article, int64, error)
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	Update(ctx context.Context, id int, category *Category) (*Category, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
}

// TagRepository 标签仓储接口
type TagRepository interface {
	Create(ctx context.Context, tag *Tag) (*Tag, error)
	GetByID(ctx context.Context, id int) (*Tag, error)
	Update(ctx context.Context, id int, tag *Tag) (*Tag, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*Tag, error)
	GetByName(ctx context.Context, name string) (*Tag, error)
	GetByIDs(ctx context.Context, ids []int) ([]*Tag, error)
}

// ArticleService 文章服务接口
type ArticleService interface {
	Create(ctx context.Context, req *ArticleCreateRequest) (*Article, error)
	GetByID(ctx context.Context, id int) (*Article, error)
	Update(ctx context.Context, id int, req *ArticleUpdateRequest) (*Article, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params QueryParams) ([]*Article, int64, error)
	ListByCategory(ctx context.Context, categoryID int, params QueryParams) ([]*Article, int64, error)
	ListByTag(ctx context.Context, tagID int, params QueryParams) ([]*Article, int64, error)
	BackupAll(ctx context.Context) ([]byte, error)
}

// CategoryService 分类服务接口
type CategoryService interface {
	Create(ctx context.Context, req *CategoryCreateRequest) (*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	Update(ctx context.Context, id int, req *CategoryUpdateRequest) (*Category, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*Category, error)
}

// TagService 标签服务接口
type TagService interface {
	Create(ctx context.Context, req *TagCreateRequest) (*Tag, error)
	GetByID(ctx context.Context, id int) (*Tag, error)
	Update(ctx context.Context, id int, req *TagUpdateRequest) (*Tag, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*Tag, error)
}

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (string, error)
	GenerateToken(ctx context.Context, username string) (string, error)
}

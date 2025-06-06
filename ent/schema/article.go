package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			NotEmpty().
			Comment("文章标题"),
		field.Text("content").
			NotEmpty().
			Comment("文章内容"),
		field.String("summary").
			Optional().
			Comment("文章摘要"),
		field.Time("created_at").
			Default(time.Now).
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
		field.Bool("published").
			Default(false).
			Comment("是否发布"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).
			Ref("articles").
			Unique(),
		edge.From("tags", Tag.Type).
			Ref("articles"),
	}
}

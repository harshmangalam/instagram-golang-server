package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("username").Unique().NotEmpty().Optional().Nillable(),
		field.String("password").NotEmpty().Sensitive(),
		field.String("bio").MaxLen(100).Optional(),
		field.String("profile_pic").Optional().Nillable(),
		field.Enum("gender").Values("male", "female", "custom", "not_prefer").Default("male"),
		field.Enum("role").Values("user", "admin").Default("user"),
		field.Bool("is_active").Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable().Optional().Nillable(),

		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Optional().Nillable(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("followings", User.Type).From("followers"),
		edge.To("posts", Post.Type),
		edge.From("posts_like", Post.Type).Ref("likes"),
	}
}

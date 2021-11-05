package basic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Schema struct {
	Key         string  `bson:"key" json:"key"`
	Label       string  `bson:"label" json:"label"`
	Kind        string  `bson:"kind" json:"kind"`
	Description string  `bson:"description,omitempty" json:"description"`
	System      bool    `bson:"system,omitempty" json:"system"`
	Fields      []Field `bson:"fields,omitempty" json:"fields"`
}

type Field struct {
	Key         string      `bson:"key" json:"key"`
	Label       string      `bson:"label" json:"label"`
	Type        string      `bson:"type" json:"type"`
	Description string      `bson:"description,omitempty" json:"description"`
	Default     string      `bson:"default,omitempty" json:"default,omitempty"`
	Unique      bool        `bson:"unique,omitempty" json:"unique,omitempty"`
	Required    bool        `bson:"required,omitempty" json:"required,omitempty"`
	Private     bool        `bson:"private,omitempty" json:"private,omitempty"`
	System      bool        `bson:"system,omitempty" json:"system,omitempty"`
	Option      FieldOption `bson:"option,omitempty" json:"option,omitempty"`
}

type FieldOption struct {
	// 数字类型
	Max interface{} `bson:"max,omitempty" json:"max,omitempty"`
	Min interface{} `bson:"min,omitempty" json:"min,omitempty"`
	// 枚举类型
	Values   bson.D `bson:"values,omitempty" json:"values,omitempty"`
	Multiple *bool  `bson:"multiple,omitempty" json:"multiple,omitempty"`
	// 引用类型
	Mode   string `bson:"mode,omitempty" json:"mode,omitempty"`
	Target string `bson:"target,omitempty" json:"target,omitempty"`
	To     string `bson:"to,omitempty" json:"to,omitempty"`
}

func GenerateSchema(ctx context.Context, db *mongo.Database) (err error) {
	collection := db.Collection("schema")
	if _, err = collection.InsertMany(ctx, []interface{}{
		Schema{
			Key:         "page",
			Label:       "动态页面",
			Kind:        "manual",
			Description: "",
			System:      true,
			Fields:      nil,
		},
		Schema{
			Key:    "role",
			Label:  "权限组",
			Kind:   "collection",
			System: true,
			Fields: []Field{
				{
					Key:      "key",
					Label:    "权限代码",
					Type:     "text",
					Required: true,
					Unique:   true,
					System:   true,
				},
				{
					Key:      "name",
					Label:    "权限名称",
					Type:     "text",
					Required: true,
					System:   true,
				},
				{
					Key:      "status",
					Label:    "状态",
					Type:     "bool",
					Required: true,
					System:   true,
				},
				{
					Key:    "description",
					Label:  "描述",
					Type:   "text",
					System: true,
				},
				{
					Key:     "pages",
					Label:   "页面",
					Type:    "reference",
					Default: "'[]'",
					System:  true,
					Option: FieldOption{
						Mode:   "manual",
						Target: "page",
					},
				},
			},
		},
		Schema{
			Label: "成员",
			Key:   "admin",
			Kind:  "collection",
			Fields: []Field{
				{
					Key:      "username",
					Label:    "用户名",
					Type:     "text",
					Required: true,
					Unique:   true,
					System:   true,
				},
				{
					Key:      "password",
					Label:    "密码",
					Type:     "password",
					Required: true,
					Private:  true,
					System:   true,
				},
				{
					Key:      "status",
					Label:    "状态",
					Type:     "bool",
					Required: true,
					System:   true,
				},
				{
					Key:      "roles",
					Label:    "权限",
					Type:     "reference",
					Required: true,
					Default:  "'[]'",
					System:   true,
					Option: FieldOption{
						Mode:   "many",
						Target: "role",
						To:     "key",
					},
				},
				{
					Key:    "name",
					Label:  "姓名",
					Type:   "text",
					System: true,
				},
				{
					Key:    "email",
					Label:  "邮件",
					Type:   "email",
					System: true,
				},
				{
					Key:    "phone",
					Label:  "联系方式",
					Type:   "text",
					System: true,
				},
				{
					Key:     "avatar",
					Label:   "头像",
					Type:    "media",
					Default: "'[]'",
					System:  true,
				},
			},
			System: true,
		},
	}); err != nil {
		return
	}
	if _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"key": 1,
		},
		Options: options.
			Index().
			SetUnique(true).
			SetName("key_idx"),
	}); err != nil {
		return
	}
	return
}

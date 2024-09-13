package article

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"not null" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`

	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`

	CreatedAt int64 `bson:"created_at,omitempty"`
	UpdatedAt int64 `bson:"updated_at,omitempty"`
	DeletedAt int64 `bson:"deletedAt,omitempty"`
}
type PublishedArticle Article

type PublishedArticleV1 struct {
	Article
}

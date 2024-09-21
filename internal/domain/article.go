package domain

import "time"

type Article struct {
	Id        int64
	Title     string
	Content   string
	Author    Author
	Status    ArticleStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnPublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (a Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return string(cs)
	}
	return string(cs[:100])
}
func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) NonPublished() bool {
	return s != ArticleStatusPublished
}

func (s ArticleStatus) String() string {
	return [...]string{"unknown", "unpublished", "published", "private"}[s.ToUint8()]
}

type Author struct {
	Id   int64
	Name string
}

package events

type ReadEvent struct {
	UserId    int64 `json:"user_id"`
	ArticleId int64 `json:"article_id"`
}

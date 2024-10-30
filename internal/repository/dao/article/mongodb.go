package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IDGenerator int64

type MongoDBDAO struct {

	//// 代表webook
	//database *mongo.Database

	// 代表制作库
	col *mongo.Collection

	//线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
	idGen   IDGenerator
}

func (m *MongoDBDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) FindByAuthorId(Dao context.Context, id int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}
func NewMongoDBDAOV1(db *mongo.Database, idGen IDGenerator) *MongoDBDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		idGen:   idGen,
	}
}
func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "created_at", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}
func (m *MongoDBDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreatedAt = now
	article.UpdatedAt = now
	id := m.node.Generate().Int64()
	article.Id = id
	_, err := m.col.InsertOne(ctx, article)
	return id, err
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, article Article) error {
	// 操作制作库
	filter := bson.M{"id": article.Id, "author_id": article.AuthorId}
	update := bson.D{bson.E{"$set", bson.M{
		"title":      article.Title,
		"content":    article.Content,
		"updated_at": time.Now().UnixMilli(),
		"status":     article.Status,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// 这边就是校验了 author_id 是不是正确的 ID
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

func (m *MongoDBDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, article)
	} else {
		id, err = m.Insert(ctx, article)
	}
	if err != nil {
		return 0, err
	}
	article.Id = id
	// 操作线上库了, upsert 语义
	now := time.Now().UnixMilli()
	//update := bson.E{"$set", art}
	//upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	article.UpdatedAt = now
	updateV1 := bson.M{
		// 更新，如果不存在，就是插入，
		"$set": PublishedArticle(article),
		// 在插入的时候，要插入 ctime
		"$setOnInsert": bson.M{"created_at": now},
	}
	filter := bson.M{"id": article.Id}
	_, err = m.liveCol.UpdateOne(ctx, filter,
		//bson.D{update, upsert},
		updateV1,
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBDAO) SyncV1(ctx context.Context, article Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, article Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

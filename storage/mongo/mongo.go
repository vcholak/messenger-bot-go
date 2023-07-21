package mongo

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vcholak/messenger-bot/lib/errp"
	"github.com/vcholak/messenger-bot/storage"
)

type MongoDBStorage struct {
	pages Pages
}

type Pages struct {
	*mongo.Collection
}

type Page struct {
	URL       string `bson:"url"`
	FirstName string `bson:"first_name"`
}

func New(connectString string, connectTimeout time.Duration) MongoDBStorage {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	pages := Pages{
		Collection: client.Database("read-adviser").Collection("pages"),
	}

	return MongoDBStorage{
		pages: pages,
	}
}

func (s MongoDBStorage) Save(ctx context.Context, page *storage.Page) error {
	_, err := s.pages.InsertOne(ctx, Page{
		URL:       page.URL,
		FirstName: page.FirstName,
	})
	if err != nil {
		return errp.Wrap("can't save page", err)
	}
	return nil
}

func (s MongoDBStorage) PickRandom(ctx context.Context, firstName string) (page *storage.Page, err error) {
	defer func() { err = errp.WrapIfErr("can't pick random page", err) }()

	pipe := bson.A{
		bson.M{"$sample": bson.M{"size": 1}},
	}

	cursor, err := s.pages.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var p Page

	cursor.Next(ctx)

	err = cursor.Decode(&p)
	switch {
	case errors.Is(err, io.EOF):
		return nil, storage.ErrNoSavedPages
	case err != nil:
		return nil, err
	}
	res := &storage.Page{
		URL:       p.URL,
		FirstName: p.FirstName,
	}
	return res, nil
}

func (s MongoDBStorage) Remove(ctx context.Context, storagePage *storage.Page) error {
	_, err := s.pages.DeleteOne(ctx, toPage(storagePage).Filter())
	if err != nil {
		return errp.Wrap("can't remove page", err)
	}
	return nil
}

func (s MongoDBStorage) IsExists(ctx context.Context, storagePage *storage.Page) (bool, error) {
	count, err := s.pages.CountDocuments(ctx, toPage(storagePage).Filter())
	if err != nil {
		return false, errp.Wrap("can't check if page exists", err)
	}
	return count > 0, nil
}

func toPage(p *storage.Page) Page {
	return Page{
		URL:       p.URL,
		FirstName: p.FirstName,
	}
}

func (p Page) Filter() bson.M {
	return bson.M{
		"url":        p.URL,
		"first_name": p.FirstName,
	}
}

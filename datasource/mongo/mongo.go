package mongo

import (
	"context"
	"instagram/logger"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection    = "users"
	photosCollection   = "photos"
	commentsCollection = "comments"
	likesCollection    = "likes"
)

type UsersManager struct {
	db *mongo.Database

	log logger.Logger
}

type PhotosManager struct {
	db *mongo.Database

	log logger.Logger
}

type СommentsManager struct {
	db *mongo.Database

	log logger.Logger
}

type LikesManager struct {
	db *mongo.Database

	log logger.Logger
}

func NewUsersManager(log logger.Logger, username, password string, dbHosts []string, database string) *UsersManager {
	client, err := connect(dbHosts, username, password)
	if err != nil {
		panic(err)
	}

	db := client.Database(database)

	return &UsersManager{
		log: log,
		db:  db,
	}
}

func (usersManager *UsersManager) Shutdown(ctx context.Context) error {
	if err := usersManager.db.Client().Disconnect(ctx); err != nil {
		return errors.Wrap(err, "usersManager shutdown")
	}

	return nil
}

func NewPhotosManager(log logger.Logger, username, password string, dbHosts []string, database string) *PhotosManager {
	client, err := connect(dbHosts, username, password)
	if err != nil {
		panic(err)
	}

	db := client.Database(database)

	return &PhotosManager{
		log: log,
		db:  db,
	}
}

func (photosManager *PhotosManager) Shutdown(ctx context.Context) error {
	if err := photosManager.db.Client().Disconnect(ctx); err != nil {
		return errors.Wrap(err, "photosManager shutdown")
	}

	return nil
}

func NewCommentsManager(log logger.Logger, username, password string, dbHosts []string, database string) *СommentsManager {
	client, err := connect(dbHosts, username, password)
	if err != nil {
		panic(err)
	}

	db := client.Database(database)

	return &СommentsManager{
		log: log,
		db:  db,
	}
}

func (commentsManager *СommentsManager) Shutdown(ctx context.Context) error {
	if err := commentsManager.db.Client().Disconnect(ctx); err != nil {
		return errors.Wrap(err, "commentsManager shutdown")
	}

	return nil
}

func NewLikeManager(log logger.Logger, username, password string, dbHosts []string, database string) *LikesManager {
	client, err := connect(dbHosts, username, password)
	if err != nil {
		panic(err)
	}

	db := client.Database(database)

	return &LikesManager{
		log: log,
		db:  db,
	}
}

func (likesManager *LikesManager) Shutdown(ctx context.Context) error {
	if err := likesManager.db.Client().Disconnect(ctx); err != nil {
		return errors.Wrap(err, "likesManager shutdown")
	}

	return nil
}

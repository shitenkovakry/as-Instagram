package mongo

import (
	"context"
	models "instagram/models/photos"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (mongo *PhotosManager) Read() (models.Photos, error) {
	collectionPhotos := mongo.db.Collection(photosCollection)
	filter := &bson.M{}

	cursor, err := collectionPhotos.Find(context.Background(), filter)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read photos")
	}

	var photos models.Photos

	if err := cursor.All(context.Background(), &photos); err != nil {
		return nil, errors.Wrap(err, "can not read cursor")
	}

	return photos, nil
}

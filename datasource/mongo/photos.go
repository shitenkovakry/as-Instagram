package mongo

import (
	"context"
	"instagram/errordefs"
	models "instagram/models/photos"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (photos *PhotosManager) obtainNextIDForPhoto() (int, error) {
	nextID := 0
	collectionPhotos := photos.db.Collection(photosCollection)

	// Определение этапов агрегации.
	pipeline := []bson.M{
		{"$sort": bson.M{"id_photo": -1}},
		{"$limit": 1},
		{"$project": bson.M{"_id": 0, "id_photo": 1}},
	}

	// Создание объекта Aggregation.
	agg, err := collectionPhotos.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}

	// Получение результатов агрегации.
	var resultForPhoto models.Photos

	if err := agg.All(context.Background(), &resultForPhoto); err != nil {
		panic(err)
	}

	// Вывод результата.
	if len(resultForPhoto) > 0 {
		nextID = resultForPhoto[0].IDPhoto
	}

	return nextID + 1, nil
}

func (photos *PhotosManager) findPhotoByFilter(filter *bson.M) (*models.Photo, error) {
	collectionPhoto := photos.db.Collection(photosCollection)
	result := collectionPhoto.FindOne(context.Background(), filter)

	err := result.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errordefs.ErrNoDocuments
	}

	if err != nil {
		return nil, errors.Wrap(err, "can not find by ID")
	}

	var foundPhoto *models.Photo
	if err := result.Decode(&foundPhoto); err != nil {
		return nil, errors.Wrap(err, "can not decode found photo")
	}

	return foundPhoto, nil
}

func (photos *PhotosManager) InsertForPhoto(photoOfUser *models.Photo) (*models.Photo, error) {
	collectionPhotos := photos.db.Collection(photosCollection)

	nextID, err := photos.obtainNextIDForPhoto()
	if err != nil {
		return nil, errors.Wrap(err, "can not find next id for photo")
	}

	photoOfUser.IDPhoto = nextID

	opts := options.InsertOne()

	result, err := collectionPhotos.InsertOne(context.Background(), photoOfUser, opts)
	if err != nil {
		return nil, errors.Wrap(err, "can not insert photo")
	}

	filter := &bson.M{
		"_id": result.InsertedID,
	}

	insertedPhoto, err := photos.findPhotoByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find inserted photo")
	}

	return insertedPhoto, nil
}

func (photos *PhotosManager) DeletePhoto(photoID int) (*models.Photo, error) {
	collectionPhotos := photos.db.Collection(photosCollection)

	filter := &bson.M{
		"photo_id": photoID,
	}

	deletedPhoto, err := photos.findPhotoByFilter(filter)
	if err != nil {
		return nil, errors.Wrap(err, "can not find deleted photo")
	}

	if _, err := collectionPhotos.DeleteOne(context.Background(), filter); err != nil {
		return nil, errors.Wrap(err, "can not delete photo")
	}

	return deletedPhoto, nil
}

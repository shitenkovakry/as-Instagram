package likes

type Like struct {
	IDUser  int `bson:"user_id"`
	IDPhoto int `bson:"photo_id"`
}

type Likes []*Like

package likes

type Like struct {
	IDUser  int `bson:"user_id"`
	IdPhoto int `bson:"photo_id"`
}

type Likes []*Like

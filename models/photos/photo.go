package photo

type Photo struct {
	IDPhoto int    `bson:"id_photo"`
	IDUser  int    `bson:"id_user"`
	Path    string `bson:"path"`
}

type Photos []*Photo

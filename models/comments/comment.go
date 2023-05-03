package comments

type Comment struct {
	IdComment int    `bson:"comment_id"`
	IdUser    int    `bson:"user_id"`
	IdPhoto   int    `bson:"photo_id"`
	Comment   string `bson:"comment"`
}

type Comments []*Comment

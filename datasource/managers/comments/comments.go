package comments

import (
	"instagram/logger"
	models "instagram/models/comments"

	"github.com/pkg/errors"
)

type DB interface {
	InsertForComment(commentOfUser *models.Comment) (*models.Comment, error)
}

type CommentsManager struct {
	log logger.Logger
	db  DB
}

func New(log logger.Logger, db DB) *CommentsManager {
	return &CommentsManager{
		log: log,
		db:  db,
	}
}

func (comment *CommentsManager) Add(commentOfUser *models.Comment) (*models.Comment, error) {
	insertedComment, err := comment.db.InsertForComment(commentOfUser)
	if err != nil {
		return nil, errors.Wrap(err, "cannot add comment")
	}

	return insertedComment, nil
}

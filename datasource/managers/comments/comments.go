package comments

import (
	"instagram/logger"
	models "instagram/models/comments"

	"github.com/pkg/errors"
)

type DB interface {
	InsertForComment(commentOfUser *models.Comment) (*models.Comment, error)
	DeleteComment(commentID int) (*models.Comment, error)
	UpdateComment(commentID int, newComment string) (*models.Comment, error)
	ReadComments(userID int, photoID int) (models.Comments, error)
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

func (comment *CommentsManager) Delete(commentID int) (*models.Comment, error) {
	deletedComment, err := comment.db.DeleteComment(commentID)
	if err != nil {
		return nil, errors.Wrap(err, "can not delete comment")
	}

	return deletedComment, nil
}

func (comment *CommentsManager) Update(commentID int, newComment string) (*models.Comment, error) {
	updatedComment, err := comment.db.UpdateComment(commentID, newComment)
	if err != nil {
		return nil, errors.Wrap(err, "can not update users comment")
	}

	return updatedComment, nil
}

func (comments *CommentsManager) ReadComments(userID int, photoID int) (models.Comments, error) {
	read, err := comments.db.ReadComments(userID, photoID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not read")
	}

	return read, nil
}

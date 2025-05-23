package quizzes

import (
	"wakuwaku_nihongo/internals/model"
	"wakuwaku_nihongo/internals/query"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type repo struct {
	*query.Query
}

func NewQuizRepo(db *gorm.DB) *repo {
	return &repo{
		query.Use(db),
	}
}

func (r *repo) GetQuizzes() (out *model.Quiz, err error) {
	q := r.Quiz
	quiz, err := q.Where(q.DeletedAt.IsNotNull(), q.Title.Eq(DEFAULT_QUIZ_TITLE)).
		Preload(q.Questions).First()
	if err != nil {
		log.Error().Err(err).Msg("error query")
		return
	}

	out = quiz
	return out, nil
}

func (r *repo) GetQuestionIDs(quizTitle string) (out []*string, err error) {
	q := r.Question
	qz := r.Quiz
	questions, err := q.Preload(q.Quiz.On(qz.Title.Eq(quizTitle))).
		Where(q.DeletedAt.IsNull()).Select(q.QuestionID).Find()
	if err != nil {
		log.Error().Err(err).Msg("error query")
		return
	}
	out = []*string{}
	for _, val := range questions {
		out = append(out, &val.QuestionID)
	}
	return out, nil
}

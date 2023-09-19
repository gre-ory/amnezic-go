package store

import (
	"context"
	"database/sql"
	"sync"

	"github.com/gre-ory/amnezic-go/internal/model"
)

// //////////////////////////////////////////////////
// themeQuestion memory store

func NewThemeQuestionMemoryStore() ThemeQuestionStore {
	return &themeQuestionMemoryStore{
		themeQuestions: make(map[model.ThemeQuestionId]*model.ThemeQuestion),
	}
}

type themeQuestionMemoryStore struct {
	themeQuestions     map[model.ThemeQuestionId]*model.ThemeQuestion
	themeQuestionsLock sync.RWMutex
}

var (
	NextThemeQuestionId = 0
)

func (s *themeQuestionMemoryStore) Create(ctx context.Context, _ *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	NextThemeQuestionId++
	themeQuestion.Id = model.ThemeQuestionId(NextThemeQuestionId)
	s.themeQuestions[themeQuestion.Id] = themeQuestion.Copy()
	return s.themeQuestions[themeQuestion.Id].Copy()
}

func (s *themeQuestionMemoryStore) Retrieve(ctx context.Context, _ *sql.Tx, id model.ThemeQuestionId) *model.ThemeQuestion {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	themeQuestion, found := s.themeQuestions[id]
	if !found {
		return nil
	}
	return themeQuestion.Copy()
}

func (s *themeQuestionMemoryStore) Update(ctx context.Context, _ *sql.Tx, themeQuestion *model.ThemeQuestion) *model.ThemeQuestion {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	_, found := s.themeQuestions[themeQuestion.Id]
	if !found {
		panic(model.ErrThemeQuestionNotFound)
	}
	s.themeQuestions[themeQuestion.Id] = themeQuestion.Copy()
	return s.themeQuestions[themeQuestion.Id].Copy()
}

func (s *themeQuestionMemoryStore) Delete(ctx context.Context, _ *sql.Tx, filter *model.ThemeQuestionFilter) {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	for id, question := range s.themeQuestions {
		if filter.IsMatching(question) {
			delete(s.themeQuestions, id)
		}
	}
}

func (s *themeQuestionMemoryStore) List(ctx context.Context, _ *sql.Tx, filter *model.ThemeQuestionFilter) []*model.ThemeQuestion {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	questions := make([]*model.ThemeQuestion, 0, len(s.themeQuestions))
	for _, question := range s.themeQuestions {
		if filter.IsMatching(question) {
			questions = append(questions, question.Copy())
		}
	}

	return questions
}

func (s *themeQuestionMemoryStore) CountByTheme(ctx context.Context, _ *sql.Tx) map[model.ThemeId]int {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	count := make(map[model.ThemeId]int, 0)
	for _, question := range s.themeQuestions {
		count[question.ThemeId]++
	}

	return count
}

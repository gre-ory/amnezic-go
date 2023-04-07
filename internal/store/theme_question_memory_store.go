package store

import (
	"context"
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

func (s *themeQuestionMemoryStore) Create(ctx context.Context, themeQuestion *model.ThemeQuestion) (*model.ThemeQuestion, error) {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	themeQuestionNumber := len(s.themeQuestions) + 1
	themeQuestion.Id = model.ThemeQuestionId(themeQuestionNumber)
	s.themeQuestions[themeQuestion.Id] = themeQuestion
	return s.themeQuestions[themeQuestion.Id], nil
}

func (s *themeQuestionMemoryStore) Retrieve(ctx context.Context, id model.ThemeQuestionId) (*model.ThemeQuestion, error) {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	themeQuestion, found := s.themeQuestions[id]
	if !found {
		return nil, model.ErrThemeQuestionNotFound
	}
	return themeQuestion, nil
}

func (s *themeQuestionMemoryStore) Update(ctx context.Context, themeQuestion *model.ThemeQuestion) (*model.ThemeQuestion, error) {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	_, found := s.themeQuestions[themeQuestion.Id]
	if !found {
		return nil, model.ErrThemeQuestionNotFound
	}
	s.themeQuestions[themeQuestion.Id] = themeQuestion
	return s.themeQuestions[themeQuestion.Id], nil
}

func (s *themeQuestionMemoryStore) Delete(ctx context.Context, filter *model.ThemeQuestionFilter) error {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	for id, question := range s.themeQuestions {
		if filter.IsMatching(question) {
			delete(s.themeQuestions, id)
		}
	}
	return nil
}

func (s *themeQuestionMemoryStore) List(ctx context.Context, filter *model.ThemeQuestionFilter) ([]*model.ThemeQuestion, error) {
	s.themeQuestionsLock.Lock()
	defer s.themeQuestionsLock.Unlock()

	questions := make([]*model.ThemeQuestion, 0, len(s.themeQuestions))
	for _, question := range s.themeQuestions {
		if filter.IsMatching(question) {
			questions = append(questions, question.Copy())
		}
	}

	return questions, nil
}

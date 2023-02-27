package model

// //////////////////////////////////////////////////
// game settings

type GameSettings struct {
	Seed       int64
	NbQuestion int
	NbAnswer   int
	NbPlayer   int
	Sources    []Source
}

const (
	MinNbPlayer = 2
	MaxNbPlayer = 99

	MinNbQuestion = 1
	MaxNbQuestion = 999

	MinNbAnswer = 2
	MaxNbAnswer = 99
)

func (o *GameSettings) Validate() error {
	if o.NbPlayer < MinNbPlayer || o.NbPlayer > MaxNbPlayer {
		return ErrInvalidNbPlayer
	}
	if o.NbQuestion < MinNbQuestion || o.NbQuestion > MaxNbQuestion {
		return ErrInvalidNbQuestion
	}
	if o.NbAnswer < MinNbAnswer || o.NbAnswer > MaxNbAnswer {
		return ErrInvalidNbAnswer
	}
	if len(o.Sources) == 0 {
		return ErrMissingSource
	}
	return nil
}

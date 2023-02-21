package model

// //////////////////////////////////////////////////
// game settings

type GameSettings struct {
	Seed       int64
	UseLegacy  bool
	NbQuestion int
	NbAnswer   int
	NbPlayer   int
}

const (
	MinNbPlayer = 1
	MaxNbPlayer = 9

	MinNbQuestion = 1
	MaxNbQuestion = 99

	MinNbAnswer = 2
	MaxNbAnswer = 9
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
	return nil
}

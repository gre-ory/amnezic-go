package model

import (
	"fmt"

	"github.com/gre-ory/amnezic-go/internal/util"
	"go.uber.org/zap/zapcore"
)

// //////////////////////////////////////////////////
// theme id

type ThemeId int64

func (i ThemeId) String() string {
	return fmt.Sprintf("%d", i)
}

func ToThemeId(value string) ThemeId {
	return ThemeId(util.StrToInt64(value))
}

// //////////////////////////////////////////////////
// theme

type Theme struct {
	Id     ThemeId
	Title  string
	ImgUrl string
	Labels map[string]string

	// consolidated data
	Questions []*ThemeQuestion
}

func (o *Theme) Copy() *Theme {
	return &Theme{
		Id:     o.Id,
		Title:  o.Title,
		ImgUrl: o.ImgUrl,
		Labels: util.CopyMap(o.Labels),
	}
}

func (o *Theme) GetInfo() *ThemeInfo {
	return &ThemeInfo{
		Id:     o.Id,
		Title:  o.Title,
		ImgUrl: o.ImgUrl,
		Labels: o.Labels,
	}
}

func (o *Theme) Equal(other *Theme) bool {
	if o == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	return (o.Id == other.Id) &&
		(o.Title == other.Title) &&
		(o.ImgUrl == other.ImgUrl) &&
		(util.EqualMap(o.Labels, other.Labels))
}

func (o *Theme) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("id", int64(o.Id))
	enc.AddString("title", o.Title)
	enc.AddString("img-url", o.ImgUrl)
	enc.AddObject("labels", zapcore.ObjectMarshalerFunc(o.MarshalLogLabels))
	if o.Questions != nil {
		enc.AddInt("nb-questions", len(o.Questions))
	}
	return nil
}

func (o *Theme) MarshalLogLabels(enc zapcore.ObjectEncoder) error {
	for key, value := range o.Labels {
		enc.AddString(key, value)
	}
	return nil
}

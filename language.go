package gateway

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log"
)

type Language interface {
	ShouldLocalize(lc *i18n.LocalizeConfig) string
	Localize(msgId, message string, params ...any) string
	Translate(msgId string, params ...any) string
}

type lang struct {
	bundle *i18n.Bundle
	*i18n.Localizer
}

func NewLanguage(bundle *i18n.Bundle, langs ...string) Language {
	l := &lang{
		bundle:    bundle,
		Localizer: i18n.NewLocalizer(bundle, langs...),
	}
	return l
}

func (l *lang) ShouldLocalize(lc *i18n.LocalizeConfig) string {
	result, err := l.Localizer.Localize(lc)
	if err != nil {
		log.Printf("error on localize, err : %v", err)
	}
	return result
}

func (l *lang) Localize(msgId, message string, params ...any) string {
	var p any
	if len(params) > 0 {
		p = params[0]
	}
	msg, err := l.Localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgId,
			Other: message,
		},
		TemplateData: p,
	})
	if err != nil {
		log.Printf("error on localize, err : %v", err)
	}
	return msg
}

func (l *lang) Translate(msgId string, params ...any) string {
	var p any
	if len(params) > 0 {
		p = params[0]
	}
	msg, err := l.Localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgId,
			Other: msgId,
		},
		TemplateData: p,
	})
	if err != nil {
		log.Printf("error on localize, err : %v", err)
	}
	return msg
}

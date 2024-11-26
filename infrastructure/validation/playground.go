package validation

import (
	"backend-agent-demo/adapter/validator"
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	go_playground "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	errTranslatorNotFound = errors.New("translator not found")
)

type playground struct {
	validator *go_playground.Validate
	translate ut.Translator
	err       error
	msg       []string
}

func NewGoPlayground() (validator.Validator, error) {
	var (
		language         = en.New()
		uni              = ut.New(language, language)
		translate, found = uni.GetTranslator("en")
	)

	if !found {
		return nil, errTranslatorNotFound
	}

	v := go_playground.New()
	if err := en_translations.RegisterDefaultTranslations(v, translate); err != nil {
		return nil, errTranslatorNotFound
	}

	return &playground{
		validator: v,
		translate: translate,
	}, nil
}

func (p *playground) Validate(i interface{}) error {
	if len(p.msg) > 0 {
		p.msg = nil
	}

	p.err = p.validator.Struct(i)
	if p.err != nil {
		return p.err
	}

	return nil
}

func (p *playground) Message() []string {
	if p.err != nil {
		for _, err := range p.err.(go_playground.ValidationErrors) {
			p.msg = append(p.msg, err.Translate(p.translate))
		}
	}

	return p.msg
}

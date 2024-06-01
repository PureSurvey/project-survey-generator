package minifier

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"regexp"
)

type Service struct {
	minifier *minify.M
}

func NewService() *Service {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	return &Service{minifier: m}
}

func (s *Service) Minify(string string) string {
	result, _ := s.minifier.String("text/html", string)
	return result
}

func (s *Service) MinifyJs(string string) string {
	result, _ := s.minifier.String("text/javascript", string)
	return result
}

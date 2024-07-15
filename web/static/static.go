package static

import (
	"github.com/a-h/templ"
)

var HttpErrorPages = map[int]templ.Component{
	400: ErrorPage(ErrorInfo400),
	401: ErrorPage(ErrorInfo401),
	403: ErrorPage(ErrorInfo403),
	404: ErrorPage(ErrorInfo404),
}

func GetHttpErrorPage(code int) templ.Component {
	component := HttpErrorPages[code]
	if component == nil {
		return nil
	}
	return component
}

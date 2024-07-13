package static

import (
	"github.com/a-h/templ"
)

var ErrorPages = map[int]templ.Component{
	400: HTTPError400(),
	401: HTTPError401(),
	403: HTTPError403(),
	404: HTTPError404(),
}

func GetHttpErrorPage(code int) templ.Component {
	component := ErrorPages[code]
	if component == nil {
		return nil
	}
	return component
}

// func createFile(path string) (*os.File, error) {
// 	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
// 		return nil, err
// 	}
// 	return os.Create(path)
// }
//
// func GetHttpErrorFilename(code int) string {
// 	return fmt.Sprintf("dist/%d.html", code)
// }
//
// func Build() {
// 	for code, component := range ErrorPages {
// 		name := GetHttpErrorFilename(code)
// 		f, err := createFile(name)
// 		if err != nil {
// 			log.Fatalf("failed to create output file: %v", err)
// 		}
//
// 		err = component.Render(context.Background(), f)
// 		if err != nil {
// 			log.Fatalf("failed to write output file: %v", err)
// 		}
// 	}
// }

package cyoa

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Story map[string]Chapter

// The html template fileName.
const TemplFileName string = "default.html"

type Chapter struct {
	Title   string   `json:"title"`
	Para    []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"chapter"`
}

// GetDefaultTempl reads the html template file, Returns it as a string
func GetDefaultTempl(fileName string) string {

	body, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}

	return string(body)
}

// Will return a new http.Handler
func GetHandler(s Story) http.Handler {
	return handler{s}
}

// handler interface for Story, Which implements the ServeHTTP
type handler struct {
	s Story
}

// ServeHTTP method of handler Story
func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tpl := template.Must(template.New("").Parse(GetDefaultTempl(TemplFileName)))

	// Parse the path and display the specific page
	path := strings.TrimSpace(req.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	fmt.Println(path)

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Chapter Not found
	http.Error(w, "Chapter Not Found", http.StatusNotFound)

}

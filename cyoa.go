package cyoa

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// The html template fileName.
const TemplFileName string = "default.html"

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Para    []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"chapter"`
}

// handler interface for Story, Which implements the ServeHTTP
type handler struct {
	s    Story
	tmpl string
}

// GetTempl reads the html template file, Returns it as a string
func GetTempl(fileName string) string {

	body, err := os.ReadFile(fileName)
	if err != nil {
		return ""
	}

	return string(body)
}

type HandlerOptions func(h *handler)

func WithHandlerOptions(tmplName string) HandlerOptions {
	return func(h *handler) {
		h.tmpl = tmplName
	}
}

// Will return a new http.Handler
func GetHandler(s Story, opts ...HandlerOptions) http.Handler {
	// The default handler will use the TemplFileName template, which is "default.html"
	h := handler{s, TemplFileName}
	for _, opt := range opts {
		// This calls the HandlreOptions function, which is returned by WithHandlerOptions func
		// Set the 'tmpl' as the user provided option.
		// user call should be something like below
		// 	handler := cyoa.GetHandler(story, cyoa.WithHandlerOptions("fileName"))
		opt(&h)
	}
	return h
}

// ServeHTTP method of handler Story
func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tpl := template.Must(template.New("").Parse(GetTempl(h.tmpl)))

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

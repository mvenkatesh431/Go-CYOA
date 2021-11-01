package main

import (
	"cyoa"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// The html template fileName.
//const templFName string = "default.html"
const custTemplFName string = "custTmpl-1.html"
const storyTemplFName string = "custStoryTmpl.html"

func main() {
	portNumber := flag.Int("port", 8000, "Port Number to run web server")
	fileName := flag.String("json", "gophers.json", "JSON file with the CYOA stories")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("Failed open the file ", *fileName)
	}

	story, err := JsonStory(file)
	if err != nil {
		panic(fmt.Sprintln("Failed to decode the JSON ", err))
	}

	// for index, value := range story {
	// 	fmt.Printf("%+v - %+v\n", index, value)
	// }

	// Get the http.Handler
	// To use the default template, don't pass the second argument
	handler := cyoa.GetHandler(story, cyoa.WithHandlerTmpl(storyTemplFName),
		cyoa.WithHandlerPathFunc(parseStoryPath))

	defaultHandler := cyoa.GetHandler(story, cyoa.WithHandlerTmpl(custTemplFName))

	// For the request like '/story/intro', If the user not providing the '/story', Then chapater won't found
	// So we will use the mux
	mux := http.NewServeMux()

	// To handle the default requests like '/' and '/intro
	mux.Handle("/", defaultHandler)

	// Now we are going to give the all requests with '/story/' or '/story/intro' to 'handler', Othre will get 404 not found
	mux.Handle("/story/", handler)

	fmt.Println("Starting the Server on port", *portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumber), mux))
}

// Parse json file and decode into the Story struct
func JsonStory(file io.Reader) (cyoa.Story, error) {
	d := json.NewDecoder(file)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// This function takes the request path as '/story/intro' type
func parseStoryPath(req *http.Request) string {
	// Parse the path and display the specific page
	path := strings.TrimSpace(req.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	log.Printf("Got Request for page '%s', Serving from StoryMux \n", path)

	// If the path won't contain the '/story/', We end up going 'slice out of bounds', So server will get run time error
	if strings.Contains(path, "/story/") {
		path = path[len("/story/"):]
	}
	return path
}

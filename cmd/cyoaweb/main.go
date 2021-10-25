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
)

// The html template fileName.
//const templFileName string = "default.html"
const templFileName string = "custTmpl-1.html"

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
	handler := cyoa.GetHandler(story, cyoa.WithHandlerOptions(templFileName))
	fmt.Println("Starting the Server on port", *portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portNumber), handler))
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

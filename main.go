package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})

	n := negroni.Classic()

	// Handler goes last
	n.UseHandler(mux)

	// Serve
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	n.Run(":" + port)
}

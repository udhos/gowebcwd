package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var knownPaths []string

func main() {

	currDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Getwd: %s", err)
	}

	http.HandleFunc("/", rootHandler) // default handler

	registerStatic("/www/", currDir)

	addr := ":8080"

	log.Printf("serving on port TCP %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

type staticHandler struct {
	innerHandler http.Handler
}

func registerStatic(path, dir string) {
	http.Handle(path, staticHandler{http.StripPrefix(path, http.FileServer(http.Dir(dir)))})
	knownPaths = append(knownPaths, path)
	log.Printf("registering static directory %s as www path %s", dir, path)
}

func (handler staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("staticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("root: URL=%s", r.URL.Path)
	log.Printf(msg)

	var paths string
	for _, p := range knownPaths {
		paths += fmt.Sprintf("<a href=\"%s\">%s</a> <br>", p, p)
	}

	rootStr :=
		`<!DOCTYPE html>

<html>
  <head>
    <title>gowebcwd root page</title>
  </head>
  <body>
    <h1>All known paths:</h1>
    %s
  </body>
</html>
`

	rootPage := fmt.Sprintf(rootStr, paths)

	io.WriteString(w, rootPage)
}

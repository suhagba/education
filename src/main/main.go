package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	//http.HandleFunc("/", HandleR)
	//http.ListenAndServe(":8000", nil)
	templates := populateTemplate()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		requestedFile := req.URL.Path[1:]
		templateWebPage := templates.Lookup(requestedFile + ".html")

		if templateWebPage != nil {
			templateWebPage.Execute(w, nil)
		} else {
			w.WriteHeader(404)
		}

	})

	http.HandleFunc("/images/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/dummy/", serveResource)
	http.HandleFunc("/fonts/", serveResource)
	http.HandleFunc("/js/", serveResource)
	http.HandleFunc("/sass/", serveResource)

	http.ListenAndServe(":8000", nil)
}
func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "text/javascript"
	} else {
		contentType = "text/plain"
	}
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}
func populateTemplate() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
		}
	}
	result.ParseFiles(*templatePaths...)

	return result

}

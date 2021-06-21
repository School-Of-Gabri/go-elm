package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func main() {
	//Run function to build elm file
	build_elm()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

func build_elm() {

	// get ELM execution path
	elm_executeable, _ := exec.LookPath("elm")

	// ELM make command
	cmd_make_elm := &exec.Cmd{
		Path:   elm_executeable,
		Args:   []string{elm_executeable, "make", "src/App.elm", "--output", "static/js/app.js"},
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}

	//Run ELM make command
	if err := cmd_make_elm.Run(); err != nil {
		fmt.Println("Error", err)
	}

}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		index := filepath.Join(fp, "index.html")
		info, err = os.Stat(index)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
		}
		fp = index
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

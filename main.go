package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var Dir string

type AudioFile struct {
	Name string
	Path string
}

func main() {
	flag.StringVar(&Dir, "dir", ".", "Directory to serve audio files from")
	flag.Parse()

	http.HandleFunc("/", serveIndex)
	http.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(Dir))))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	files := []AudioFile{}

	supported := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".ogg":  true,
		".aac":  true,
		".m4a":  true,
		".webm": true,
		".opus": true,
		".flac": true,
	}

	err := filepath.Walk(Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)
		if !info.IsDir() && supported[ext] {
			files = append(files, AudioFile{
				Name: info.Name(),
				Path: "/audio/" + info.Name(),
			})
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Error reading audio files", http.StatusInternalServerError)
		return
	}

	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Audio Player</title>
    </head>
    <body>
        <h1>Audio Files</h1>
        {{range .}}
            <div>
                <p>{{.Name}}</p>
                <audio controls style="width: 100%">
                    <source src="{{.Path}}" type="audio/mpeg">
                    Your browser does not support the audio element.
                </audio>
            </div>
        {{end}}
    </body>
    </html>
    `

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, files)
}

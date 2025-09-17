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

var Supported = map[string]bool{
	".aac":  true,
	".flac": true,
	".m4a":  true,
	".mp3":  true,
	".ogg":  true,
	".opus": true,
	".wav":  true,
}

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

	err := filepath.Walk(Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)
		if !info.IsDir() && Supported[ext] {
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
        <style>
            body {
                font-family: sans-serif;
                margin: 20px;
            }
            details {
                margin-bottom: 15px;
                border: 1px solid #ccc;
                border-radius: 5px;
                padding: 10px;
            }
            audio {
                width: 100%;
                padding-top: 10px;
            }
        </style>
    </head>
    <body>
        <h1>Audio Files</h1>
        {{range .}} <details>
            <summary>{{.Name}}</summary>
            <audio controls>
                <source src="{{.Path}}">
                Your browser does not support the audio element.
            </audio>
        </details>
        {{end}}
    </body>
    </html>
    `

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, files)
}

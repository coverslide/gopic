package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopic/internal/config"
	"gopic/web/template"
	htmlTemplate "html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var QuestionMarkSvg = []byte(`<?xml version="1.0" ?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 445 445"><g>
<path d="M 219.744 443.88C 341.103 443.88 439.488 344.51 439.488 221.94C 439.488 99.368 341.103 0 219.744 0C 98.387 0 0 99.368 0 221.94C 0 344.51 98.387 443.88 219.744 443.88z" style="stroke:none; fill:#000000"/>
<path d="M 219.744 221.94" style="stroke:none; fill:#000000"/>
</g>
<g>
<path d="M 219.744 392.714C 313.128 392.714 388.83 316.255 388.83 221.94C 388.83 127.623 313.128 51.166 219.744 51.166C 126.362 51.166 50.659 127.623 50.659 221.94C 50.659 316.255 126.362 392.714 219.744 392.714z" style="stroke:none; fill:#ffffff"/>
<path d="M 219.744 221.94" style="stroke:none; fill:#ffffff"/>
</g>
<g>
<path d="M 196.963 300.274L 246.494 300.172L 246.494 261.69C 246.494 251.252 251.36 241.39 264.38 232.849C 277.399 224.312 313.744 206.988 313.744 161.44C 313.744 115.89 275.577 84.582 243.494 77.94C 211.416 71.298 176.659 75.668 151.994 102.69C 129.907 126.887 125.253 146.027 125.253 188.255L 174.744 188.255L 174.744 178.44C 174.744 155.939 177.347 132.186 209.494 125.69C 227.04 122.144 243.488 127.648 253.244 137.19C 264.404 148.102 264.494 172.69 246.711 184.933L 218.815 203.912C 202.543 214.35 196.963 225.971 196.963 243.051L 196.963 300.274z" style="stroke:none; fill:#000000"/>
<g>
	<path d="M 196.638 370.692L 196.638 319.687L 246.85 319.687L 246.85 370.692L 196.638 370.692z" style="stroke:none; fill:#000000"/>
	<path d="M 221.744 345.19" style="stroke:none; fill:#000000"/>
</g></g></svg>`)

var templateDir *htmlTemplate.Template
var mainJs string

type TemplateData struct {
	Path       string
	FolderData htmlTemplate.HTML
	MainJs     string
}

type FolderData struct {
	FileName string `json:"filename"`
	IsDir    bool   `json:"isDir"`
}

func prepTemplates() {
	var err error
	templateDir, err = htmlTemplate.ParseFS(template.Content, "index.template.html", "main.template.js")
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = templateDir.ExecuteTemplate(&tpl, "main.template.js", nil)
	if err != nil {
		panic(err)
	}

	mainJs = tpl.String()
}

type Server struct {
	Server  *http.Server
	RootDir string
}

func NewServer(conf config.Config) *Server {
	prepTemplates()
	server := &Server{
		Server: &http.Server{
			Addr: conf.ListenAddr,
		},
		RootDir: conf.RootDir,
	}

	server.Server.Handler = server

	return server
}

func (s *Server) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := filepath.Join(s.RootDir, r.URL.Path)
	f, err := os.Open(fullPath)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stat, err := f.Stat()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if stat.IsDir() {
		if r.URL.Query().Get("json") == "true" {
			handleDirectoryJson(w, fullPath, f)
		} else {
			handleDirectoryPage(w, r.URL.Path)
		}
	} else if r.URL.Query().Get("thumbnail") == "true" {
		dirName := filepath.Dir(fullPath)
		baseName := filepath.Base(fullPath)
		thumbFilename := fmt.Sprintf("%s/.%s.png", dirName, baseName)
		thumbFile, err := os.Open(thumbFilename)
		if errors.Is(err, os.ErrNotExist) {
			err = createThumbNail(fullPath, thumbFilename)
			if err != nil {
				if errors.Is(err, os.ErrInvalid) {
					w.Header().Add("Content-Type", "image/svg+xml")
					w.Write(QuestionMarkSvg)
					return
				}
				log.Panic(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			thumbFile, err = os.Open(thumbFilename)
			if err != nil {
				log.Panic(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(w, thumbFile)
	} else {
		io.Copy(w, f)
	}
}

func handleDirectoryPage(w http.ResponseWriter, path string) {
	templateDir.ExecuteTemplate(w, "index.template.html", TemplateData{
		Path: path,
	})
}

func handleDirectoryJson(w http.ResponseWriter, fullPath string, f *os.File) {

	direntList, err := os.ReadDir(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fd := make([]FolderData, len(direntList))
	for i, dirent := range direntList {
		fd[i] = FolderData{
			FileName: dirent.Name(),
			IsDir:    dirent.Type().IsDir(),
		}
	}

	fdBytes, err := json.Marshal(fd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(fdBytes)
}

type FileType int

const (
	UNKNOWN FileType = iota
	IMAGE
	VIDEO
)

var fileTypeMap = map[string]FileType{
	".jpg":  IMAGE,
	".png":  IMAGE,
	".tif":  IMAGE,
	".3gpp": VIDEO,
	".3gp":  VIDEO,
	".mp4":  VIDEO,
	".m4v":  VIDEO,
}

func getFileTypeFromExtension(extension string) FileType {
	if fileType, ok := fileTypeMap[extension]; ok {
		return fileType
	}
	return UNKNOWN
}

func createThumbNail(fullPath, thumbFilename string) error {
	extension := filepath.Ext(fullPath)
	fileType := getFileTypeFromExtension(extension)
	log.Printf("extension: %s ft: %+v", extension, fileType)
	if fileType == IMAGE {
		cmd := exec.Command("convert", fullPath, "-auto-orient", "-resize", "50x50^", thumbFilename)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	} else if fileType == VIDEO {
		cmd := exec.Command("ffmpeg", "-loglevel", "warning", "-ss", "00:00:00", "-i", fullPath, "-vframes", "1", thumbFilename)

		stdout, _ := cmd.StdoutPipe()
		if err := cmd.Run(); err != nil {
			log.Println(err.Error())
			out, _ := io.ReadAll(stdout)
			log.Println(out)
			return err
		}
		return nil
	} else {
		return os.ErrInvalid
	}
}

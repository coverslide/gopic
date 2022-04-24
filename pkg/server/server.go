package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopic/pkg/config"
	"gopic/web/static"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

type FolderData struct {
	FileName string `json:"filename"`
	IsDir    bool   `json:"isDir"`
	ModTime  int64  `json:"mtime"`
}

type Server struct {
	Server  *http.Server
	RootDir string
}

func NewServer(conf config.Config) *Server {
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
	if strings.HasPrefix(r.URL.Path, "/_static/") {
		handleStatic(w, strings.TrimPrefix(r.URL.Path, "/_static/"))
		return
	}
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
		handleThumbnail(w, fullPath)
	} else if r.URL.Query().Get("info") == "true" {
		handleInfo(w, fullPath)
	} else if r.URL.Query().Get("ffprobe") == "true" {
		handleFFProbe(w, fullPath)
	} else if r.URL.Query().Get("identify") == "true" {
		handleIdentify(w, fullPath)
	} else {
		basename := filepath.Base(fullPath)
		http.ServeContent(w, r, basename, stat.ModTime(), f)
	}
}

func handleInfoCommand(w http.ResponseWriter, fullPath string, cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error: %s, output: %s", err.Error(), output)))
	} else {
		w.Write(output)
	}
}

func handleFFProbe(w http.ResponseWriter, fullPath string) {
	cmd := exec.Command("ffprobe", "-i", fullPath)
	handleInfoCommand(w, fullPath, cmd)
}

func handleIdentify(w http.ResponseWriter, fullPath string) {
	cmd := exec.Command("identify", "-verbose", fullPath)
	handleInfoCommand(w, fullPath, cmd)
}

func handleInfo(w http.ResponseWriter, path string) {
	fileExtension := filepath.Ext(path)
	mimeType := mime.TypeByExtension(fileExtension)
	w.Header().Add("X-Mime-Type", mimeType)

	mimeParts := strings.Split(mimeType, "/")
	if mimeParts[0] == "image" {
		handleIdentify(w, path)
	} else if mimeParts[0] == "video" {
		handleFFProbe(w, path)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unsupported type: %q filename: %q", mimeType, path)))
	}
}

func handleStatic(w http.ResponseWriter, path string) {
	staticFile, err := static.Content.Open(path)
	if err != nil {
		log.Default().Printf("static file error: %q", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer staticFile.Close()
	fileExtension := filepath.Ext(path)
	w.Header().Add("Content-Type", mime.TypeByExtension(fileExtension))
	io.Copy(w, staticFile)
}

func handleThumbnail(w http.ResponseWriter, fullPath string) {
	dirName := filepath.Dir(fullPath)
	baseName := filepath.Base(fullPath)
	thumbFilename := fmt.Sprintf("%s/.%s.png", dirName, baseName)
	thumbFile, err := os.Open(thumbFilename)
	if errors.Is(err, os.ErrNotExist) {
		err = createThumbNail(fullPath, thumbFilename)
		if err != nil {
			if errors.Is(err, os.ErrInvalid) {
				w.Header().Add("Content-Type", "image/svg+xml")
				questionMarkImage, err := static.Content.Open("images/question-mark.svg")
				if err != nil {
					log.Printf("error: %q", err.Error())
					w.WriteHeader(http.StatusNotFound)
					return
				}
				defer questionMarkImage.Close()
				io.Copy(w, questionMarkImage)
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
	defer thumbFile.Close()
	io.Copy(w, thumbFile)
}

func handleDirectoryPage(w http.ResponseWriter, path string) {
	fh, err := static.Content.Open("html/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fh.Close()
	io.Copy(w, fh)
}

func handleDirectoryJson(w http.ResponseWriter, fullPath string, f *os.File) {
	direntList, err := os.ReadDir(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fd := make([]FolderData, len(direntList))
	for i, dirent := range direntList {
		fileInfo, err := dirent.Info()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fd[i] = FolderData{
			FileName: fileInfo.Name(),
			IsDir:    fileInfo.IsDir(),
			ModTime:  fileInfo.ModTime().UnixMilli(),
		}
	}

	fdBytes, err := json.Marshal(fd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(fdBytes)
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
		defer stdout.Close()
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

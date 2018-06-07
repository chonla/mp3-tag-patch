package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chonla/format"
	id3 "github.com/mikkyang/id3-go"
)

const version = "0.3"

// FileID3 is file with id3 tag
type FileID3 struct {
	FileName string `json:"file-name"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Year     string `json:"year"`
	Genre    string `json:"genre"`
}

func main() {
	if len(os.Args) <= 1 {
		throwHelp()
	}

	cmd := os.Args[1]

	switch cmd {
	case "list":
		path := "./"
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		files := listFiles(path)
		b, _ := json.MarshalIndent(files, "", "    ")
		ioutil.WriteFile("mp3-tags.json", b, 0644)
	case "patch":
		path := "./"
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		files := loadList(path)
		patchID3Tags(path, files)
	case "info":
		if len(os.Args) <= 2 {
			throwError("no mp3 file to be read")
		}
		fname := os.Args[2]
		info, e := getInfo(fname)
		if e != nil {
			throwError(e.Error())
		}
		format.Printfln("File: %<filename>s\nTitle: %<title>s\nAlbum: %<album>s\nArtist: %<artist>s\nYear: %<year>s\nGenre: %<genre>s", map[string]interface{}{
			"filename": info.FileName,
			"title":    info.Title,
			"album":    info.Album,
			"artist":   info.Artist,
			"year":     info.Year,
			"genre":    info.Genre,
		})
	case "version":
		format.Printfln("v%<version>s", map[string]interface{}{
			"version": version,
		})
	default:
		fmt.Printf("unexpected command: %s\n", cmd)
		throwHelp()
	}
}

func getInfo(fname string) (f FileID3, e error) {
	if _, e = os.Stat(fname); e == nil {
		mp3, e := id3.Open(fname)
		defer mp3.Close()
		if e == nil {
			f = FileID3{
				FileName: fname,
				Artist:   strings.TrimSpace(strings.TrimRight(mp3.Artist(), "\u0000")),
				Album:    strings.TrimSpace(strings.TrimRight(mp3.Album(), "\u0000")),
				Title:    strings.TrimSpace(strings.TrimRight(mp3.Title(), "\u0000")),
				Year:     strings.TrimSpace(strings.TrimRight(mp3.Year(), "\u0000")),
				Genre:    strings.TrimSpace(strings.TrimRight(mp3.Genre(), "\u0000")),
			}
		}
	}
	return
}

func patchID3Tags(path string, files []FileID3) {
	if len(files) == 0 {
		throwError("no file to be patched")
	}
	for _, v := range files {
		fname := filepath.Clean(path + "/" + v.FileName)
		format.Printfln("Patching %<filename>s", map[string]interface{}{
			"filename": fname,
		})
		mp3, err := id3.Open(fname)

		if err != nil {
			mp3.Close()
			throwError(err.Error())
		}

		mp3.SetAlbum(v.Album)
		mp3.SetArtist(v.Artist)
		mp3.SetTitle(v.Title)
		mp3.SetGenre(v.Genre)
		mp3.SetYear(v.Year)

		mp3.Close()
	}
}

func loadList(path string) []FileID3 {
	out := []FileID3{}
	fname := filepath.Clean(path + "/mp3-tags.json")
	if _, err := os.Stat(fname); err == nil {
		content, err := ioutil.ReadFile(fname)
		if err != nil {
			throwError(err.Error())
		}

		err = json.Unmarshal(content, &out)
		if err != nil {
			throwError(err.Error())
		}
	} else {
		throwError("no mp3-tags.json found")
	}
	return out
}

func listFiles(path string) []FileID3 {
	path += fmt.Sprintf("%c", os.PathSeparator)

	out := []FileID3{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		throwError(err.Error())
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if strings.ToLower(ext) == ".mp3" {
			if !f.IsDir() {
				fname := filepath.Clean(path + f.Name())

				mp3, err := getInfo(fname)
				if err == nil {
					out = append(out, mp3)
				}
			}
		}
	}
	return out
}

func throwHelp() {
	throwError("usage:\n    mp3-tag-patch list [<path>]\n    mp3-tag-patch version\n    mp3-tag-patch info <filename>")
}

func throwError(msg string) {
	fmt.Println(msg)
	os.Exit(2)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"feodor.dk/linkyd/backend"
	"feodor.dk/linkyd/linky"
	"feodor.dk/linkyd/linky/link"
	"feodor.dk/linkyd/static"
)

var ErrInvalidPath = errors.New("invalid path")

var help bool
var port int = 8080
var loadFile string

type ParserState int

const (
	ParserStateFirst ParserState = iota
	ParserStateInvalid
	ParserStateFlag
	ParserStatePort
	ParserStateLoadFile
)

func parseCliArgs() error {
	var parserState ParserState = ParserStateFirst
	for _, arg := range os.Args {
		switch parserState {
		case ParserStateFirst:
			parserState = ParserStateFlag
			continue
		case ParserStateFlag:
			if s, err := parseFlag(arg); err != nil {
				return err
			} else {
				parserState = s
			}
		case ParserStatePort:
			if err := parsePort(arg); err != nil {
				return err
			}
		case ParserStateLoadFile:
			if err := validateLoadFile(arg); err != nil {
				return err
			}
		default:
			return errors.New("invalid parser state")
		}
	}

	return nil
}

func parseFlag(arg string) (ParserState, error) {
	switch arg {
	case "-p", "--port":
		return ParserStatePort, nil
	case "-l", "--load":
		return ParserStateLoadFile, nil
	case "-h", "--help":
		help = true
		return ParserStateFlag, nil
	default:
		return ParserStateInvalid, errors.New("invalid flag")
	}
}

func parsePort(arg string) error {
	if i, err := strconv.ParseInt(arg, 10, 32); err != nil {
		return err
	} else {
		port = int(i)
		return nil
	}
}

func validateLoadFile(arg string) error {
	if finfo, err := os.Stat(arg); err != nil {
		return err
	} else if finfo.IsDir() {
		return errors.New("load file is a directory")
	}

	loadFile = arg

	return nil
}

func main() {
	repo := link.NewInMemoryLinkRepository()
	linky := linky.New(repo)
	if err := parseCliArgs(); err != nil {
		printHelp()
		os.Exit(1)
	}

	if help {
		printHelp()
		os.Exit(0)
	}

	if loadFile != "" {
		if err := loadLinks(repo); err != nil {
			println("Error occured loading links from file:", err.Error())
			os.Exit(1)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(&linky, w, r)
		b.List()
	})

	http.HandleFunc("/as/", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(&linky, w, r)
		if asUser, err := getPathSegment(r, 2); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			b.As(asUser)
		}
	})

	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(&linky, w, r)

		switch r.Method {
		case "POST":
			b.Create()
		case "GET":
			b.List()
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		b := backend.Get(&linky, w, r)

		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, err := getPathSegment(r, 2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b.Delete(id)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("Content-Type", "image/x-icon")
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(static.Favicon)))
		w.Write(static.Favicon)
	})

	listenAddr := fmt.Sprintf(":%d", port)
	slog.Info("Starting HTTP listener", slog.String("address", listenAddr))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func getPathSegment(r *http.Request, argumentIndex int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != argumentIndex+1 {
		return "", ErrInvalidPath
	}

	return parts[argumentIndex], nil
}

type LinkDump = map[string]link.Link

func loadLinks(repo link.Repository) error {
	var dump LinkDump

	if data, err := os.ReadFile(loadFile); err != nil {
		return err
	} else if err := json.Unmarshal(data, &dump); err != nil {
		return err
	}

	for _, value := range dump {
		repo.Create(value)
	}

	return nil
}

func printHelp() {
	printlns(
		"Usage: linkyd [OPTION]...",
		"Run the linky daemon web server.",
		"",
		"  -h, --help              display this message and exit",
		"  -l, --load <DUMP FILE>  load a dump of links upon start-up",
		"  -p, --port <PORT>       specify the port, defaults to 8080",
	)
}

func printlns(lines ...string) {
	for _, line := range lines {
		println(line)
	}
}

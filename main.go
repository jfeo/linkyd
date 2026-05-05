package main

import (
	"log/slog"

	"feodor.dk/linkyd/cli"
	"feodor.dk/linkyd/http"
	"feodor.dk/linkyd/linky"
	"feodor.dk/linkyd/linky/link"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args, err := cli.ParseArguments(os.Args)
	if err != nil {
		println("Error: " + err.Error())
		cli.PrintHelp()
		os.Exit(1)
	}

	configureSlog(args)

	repo, err := link.NewSQLiteLinkRepository()
	if err != nil {
		slog.Error("could not create sqlite3 repository")
		os.Exit(1)
	}

	if args.Help {
		cli.PrintHelp()
		os.Exit(0)
	}

	linky := linky.New(repo)

	if args.LoadFile != "" {
		if err := loadDump(repo, args.LoadFile); err != nil {
			println("Error occured loading links from file:", err.Error())
			os.Exit(1)
		} else {
			println("Loaded links from", args.LoadFile)
		}
	}

	http.ListenAndServe(&linky, args.Port)
}

func configureSlog(args cli.Arguments) {
	logLevel := slog.LevelInfo
	if args.Verbose {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
}

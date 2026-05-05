package cli

import (
	"os"
	"strconv"
)

type ParserState int

const (
	ParserStateInitial ParserState = iota
	ParserStateInvalid
	ParserStateFlag
	ParserStatePort
	ParserStateLoadFile
)

func ParseArguments(args []string) (Arguments, error) {
	var res Arguments = ArgumentDefaults()

	var parserState ParserState = ParserStateInitial
	for _, arg := range args {
		switch parserState {
		case ParserStateInitial:
			parserState = ParserStateFlag
			continue
		case ParserStateFlag:
			if s, err := parseFlag(&res, arg); err != nil {
				return res, err
			} else {
				parserState = s
			}
		case ParserStatePort:
			if err := parsePort(&res, arg); err != nil {
				return res, err
			}
			parserState = ParserStateFlag
		case ParserStateLoadFile:
			if err := validateLoadFile(&res, arg); err != nil {
				return res, err
			}
			parserState = ParserStateFlag
		default:
			return res, ErrInvalidParserState
		}
	}

	return res, nil
}

func parseFlag(args *Arguments, arg string) (ParserState, error) {
	switch arg {
	case "-p", "--port":
		return ParserStatePort, nil
	case "-l", "--load":
		return ParserStateLoadFile, nil
	case "-h", "--help":
		args.Help = true
		return ParserStateFlag, nil
	case "-v", "--verbose":
		args.Verbose = true
		return ParserStateFlag, nil
	default:
		return ParserStateInvalid, ErrInvalidFlag
	}
}

func parsePort(args *Arguments, arg string) error {
	if i, err := strconv.ParseInt(arg, 10, 32); err != nil {
		return err
	} else {
		args.Port = int(i)
		return nil
	}
}

func validateLoadFile(args *Arguments, arg string) error {
	if finfo, err := os.Stat(arg); err != nil {
		if os.IsNotExist(err) {
			return ErrLoadFileDoesNotExist
		} else {
			return err
		}
	} else if finfo.IsDir() {
		return ErrLoadFileIsDirectory
	}

	args.LoadFile = arg

	return nil
}

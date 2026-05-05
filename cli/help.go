package cli

func PrintHelp() {
	printlns(
		"Usage: linkyd [OPTION]...",
		"Run the linky daemon web server.",
		"",
		"  -h, --help              display this message and exit",
		"  -l, --load <DUMP FILE>  load a dump of links upon start-up",
		"  -p, --port <PORT>       specify the port, defaults to 8080",
		"  -v, --verbose           print more verbose logs",
	)
}

func printlns(lines ...string) {
	for _, line := range lines {
		println(line)
	}
}

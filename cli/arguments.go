package cli

type Arguments struct {
	Help     bool
	Verbose  bool
	Port     int
	LoadFile string
}

func ArgumentDefaults() Arguments {
	return Arguments{
		Port: 8080,
	}
}

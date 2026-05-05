package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArguments(t *testing.T) {
	type testCase struct {
		name     string
		args     []string
		expected Arguments
	}

	for _, tc := range []testCase{
		{
			name:     "no arguments returns defaults",
			args:     []string{"name"},
			expected: Arguments{Help: false, Verbose: false, Port: 8080, LoadFile: ""},
		},
		{
			name:     "long help flag returns help true",
			args:     []string{"name", "--help"},
			expected: Arguments{Help: true, Verbose: false, Port: 8080, LoadFile: ""},
		},
		{
			name:     "short help flag returns help true",
			args:     []string{"name", "-h"},
			expected: Arguments{Help: true, Verbose: false, Port: 8080, LoadFile: ""},
		},
		{
			name:     "long port argument is used",
			args:     []string{"name", "--port", "1234"},
			expected: Arguments{Help: false, Verbose: false, Port: 1234, LoadFile: ""},
		},
		{
			name:     "short port argument is used",
			args:     []string{"name", "--port", "1234"},
			expected: Arguments{Help: false, Verbose: false, Port: 1234, LoadFile: ""},
		},
		{
			name:     "long verbose flag returns verbose true",
			args:     []string{"name", "--verbose"},
			expected: Arguments{Help: false, Verbose: true, Port: 8080, LoadFile: ""},
		},
		{
			name:     "short verbose flag returns verbose true",
			args:     []string{"name", "-v"},
			expected: Arguments{Help: false, Verbose: true, Port: 8080, LoadFile: ""},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parsed_args, err := ParseArguments(tc.args)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, parsed_args)
		})
	}

	t.Run("long load file with valid path", func(t *testing.T) {
		dir := t.TempDir()
		loadfile := filepath.Join(dir, "loadfile")
		err := os.WriteFile(loadfile, []byte("loadfile"), fs.ModePerm)
		require.NoError(t, err)

		expected := ArgumentDefaults()
		expected.LoadFile = loadfile

		args := []string{"name", "--load", loadfile}
		parsed_args, err := ParseArguments(args)

		assert.NoError(t, err)
		assert.Equal(t, expected, parsed_args)
	})

	t.Run("short load file with valid path", func(t *testing.T) {
		dir := t.TempDir()
		loadfile := filepath.Join(dir, "loadfile")
		err := os.WriteFile(loadfile, []byte("loadfile"), fs.ModePerm)
		require.NoError(t, err)

		expected := ArgumentDefaults()
		expected.LoadFile = loadfile

		args := []string{"name", "-l", loadfile}
		parsed_args, err := ParseArguments(args)

		assert.NoError(t, err)
		assert.Equal(t, expected, parsed_args)
	})

	t.Run("returns error if load file is directory", func(t *testing.T) {
		dir := t.TempDir()
		_, err := ParseArguments([]string{"name", "--load", dir})

		assert.ErrorIs(t, err, ErrLoadFileIsDirectory)
	})

	t.Run("returns error if load file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		_, err := ParseArguments([]string{"name", "--load", filepath.Join(dir, "non-existing-file")})

		assert.ErrorIs(t, err, ErrLoadFileDoesNotExist)
	})

	t.Run("returns error if invalid flag parsed", func(t *testing.T) {
		_, err := ParseArguments([]string{"name", "--invalid-flag"})

		assert.ErrorIs(t, err, ErrInvalidFlag)
	})
}

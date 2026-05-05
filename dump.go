package main

import (
	"encoding/json"
	"os"

	"feodor.dk/linkyd/linky/link"
)

type LinkDump = map[string]link.Link

func loadDump(repo link.Repository, loadFile string) error {
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

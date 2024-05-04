package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runConfig(binaryPath string, configPath string) {
	cmd := exec.Command(binaryPath+"/bin/stochadex", "--config", configPath)
	cmd.Dir = binaryPath
	err := cmd.Run()
	if err != nil {
		log.Fatal("config run err: ", err)
	}
}

func TestAllConfigs(t *testing.T) {
	t.Run(
		"test all of the configs in this repo run successfully",
		func(t *testing.T) {
			binaryPath := os.Getenv("STOCHADEX_PATH")
			matches, err := filepath.Glob("./*.yaml")
			if err != nil {
				log.Fatal(err)
			}
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			for _, configPath := range matches {
				runConfig(binaryPath, cwd+"/"+configPath)
			}
		},
	)
}

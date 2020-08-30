package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

func cjpeg(path string, quality int64) error {
	// replace ext with jpg for output path
	ext := filepath.Ext(path)
	name := filepath.Base(path)
	outfile := path[0:len(path)-len(ext)] + ".jpg"

	cmd := exec.Command("cjpeg", "-quality", strconv.FormatInt(quality, 10), "-optimize", "-progressive", "-outfile", outfile, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = cmd.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("converted: %s\n", name)
	return nil
}

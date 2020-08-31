// Copyright 2020 Tony Mackay.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var name = "png-to-jpeg"
var version = "dev"

func main() {
	path := flag.String("dir", ".", "Path to a directory containing PNG images to convert")
	workers := flag.Int("workers", runtime.NumCPU(), "Maximum amount of goroutines to use")
	quality := flag.Int64("quality", 75, "Image Quality, N between 5-95")
	showVersion := flag.Bool("version", false, "Print the version")

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s (runtime: %s)\n", name, version, runtime.Version())
		os.Exit(0)
	}

	var wg = sync.WaitGroup{}
	var guard = make(chan struct{}, *workers)
	walkDir(*path, quality, &wg, &guard)
	wg.Wait()
}

func usage() {
	fmt.Printf("usage: %s [options]\n\n", name)
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  %s -dir images\n", name)
	fmt.Printf("  %s -dir images -quality 60\n", name)
	fmt.Printf("  %s -dir images -quality 60 -workers 1\n\n", name)
}

func walkDir(dir string, quality *int64, wg *sync.WaitGroup, guard *chan struct{}) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, quality, wg, guard)
		} else {
			// only send PNG images
			path := filepath.Join(dir, entry.Name())
			if strings.ToLower(filepath.Ext(path)) == ".png" {
				*guard <- struct{}{}
				wg.Add(1)
				go func() {
					cjpeg(path, *quality)
					<-*guard
					wg.Done()
				}()
			}
		}
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil
	}
	return entries
}

func cjpeg(path string, quality int64) error {
	// replace ext with jpg for output path
	ext := filepath.Ext(path)
	//name := filepath.Base(path)
	outfile := path[0:len(path)-len(ext)] + ".jpg"
	q := strconv.FormatInt(quality, 10)
	cmd := exec.Command("cjpeg", "-quality", q, "-optimize", "-progressive", "-outfile", outfile, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = cmd.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("converted: %s to: %s\n", path, outfile)
	return nil
}

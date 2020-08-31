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

type sizes struct {
	original     int64
	originalPath string
	new          int64
	newPath      string
}

func main() {
	path := flag.String("dir", ".", "Path to a directory containing PNG images to convert")
	workers := flag.Int("workers", runtime.NumCPU(), "Maximum amount of goroutines to use")
	quality := flag.Int64("quality", 75, "Image Quality, N between 5-95")
	showVersion := flag.Bool("version", false, "Print the version")
	displayStats := flag.Bool("stats", false, "Display amount of converted images and size differences")

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s (runtime: %s)\n", name, version, runtime.Version())
		os.Exit(0)
	}

	var wg = sync.WaitGroup{}
	var guard = make(chan struct{}, *workers)
	fileSizes := make(chan sizes)

	go func() {
		walkDir(*path, quality, &wg, &guard, fileSizes)
		wg.Wait()
		close(fileSizes)
	}()

	var files, originalTotal, newTotal int64
	for s := range fileSizes {
		files++
		originalTotal += s.original
		newTotal += s.new
		fmt.Printf("converted: %s to: %s\n", s.originalPath, s.newPath)
	}

	if *displayStats {
		fmt.Printf("\nconverted: %v\n", files)
		fmt.Printf("old size: %s\nnew size: %s\n", byteCountIEC(originalTotal), byteCountIEC(newTotal))
		fmt.Printf("saved:    %0.2f%%\n", float64(originalTotal-newTotal)/float64(originalTotal)*100)
	}
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

func walkDir(dir string, quality *int64, wg *sync.WaitGroup, guard *chan struct{}, fileSizes chan<- sizes) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, quality, wg, guard, fileSizes)
		} else {
			// only send PNG images
			path := filepath.Join(dir, entry.Name())
			if strings.ToLower(filepath.Ext(path)) == ".png" {
				*guard <- struct{}{}
				go func() {
					wg.Add(1)
					defer wg.Done()
					defer func() { <-*guard }()
					cjpeg(path, *quality, fileSizes)
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

func cjpeg(path string, quality int64, fileSizes chan<- sizes) error {
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

	originalSize, _ := fileSize(path)
	newSize, _ := fileSize(outfile)

	fileSizes <- sizes{
		original:     originalSize,
		originalPath: path,
		new:          newSize,
		newPath:      outfile,
	}
	return nil
}

func fileSize(path string) (int64, error) {
	// get file size
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// convert a size in bytes to a human-readable string IEC (binary) format
// credit: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

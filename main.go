package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var sema = make(chan struct{}, runtime.NumCPU())

func main() {
	//start := time.Now()
	path := flag.String("dir", ".", "Path to a directory containing PNG images to convert")
	quality := flag.Int64("quality", 75, "Image Quality, N between 5-95")

	flag.Parse()

	filePaths := make(chan string)
	var n sync.WaitGroup
	n.Add(1)
	go walkDir(*path, &n, filePaths)

	go func() {
		n.Wait()
		close(filePaths)
	}()

	var nfiles int
	for path := range filePaths {
		nfiles++

		err := cjpeg(path, *quality)
		if err != nil {
			fmt.Printf("failed: %s\n", err)
		}

		fmt.Printf("%d files %s\n", nfiles, path)
	}
	//fmt.Println(time.Since(start))
}

func walkDir(dir string, n *sync.WaitGroup, filePaths chan<- string) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, n, filePaths)
		} else {
			// only send PNG images
			path := filepath.Join(dir, entry.Name())
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".png" {
				filePaths <- path
			}
		}
	}
}

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil
	}
	return entries
}

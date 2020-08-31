# Batch PNG to JPEG Converter
`png-to-jpeg` is a command line tool written in Go that recursively converts all PNG images in a directory to JPEGs using the Mozilla JPEG Encoder.

## Usage

```
usage: png-to-jpeg [options]

Options:
  -dir string
        Path to a directory containing PNG images to convert (default ".")
  -quality int
        Image Quality, N between 5-95 (default 75)
  -version
        Print the version
  -workers int
        Maximum amount of goroutines to use (default 4)

Examples:
  png-to-jpeg -dir images
  png-to-jpeg -dir images -quality 60
  png-to-jpeg -dir images -quality 60 -workers 1
```

## Building
This tool requires the MozJPEG package to be installed. For Mac users it can be installed with `brew install mozjpeg`.

Clone the repo then run the following commands:

```   
go build
```

To install it run:

```
go install
```

To assign a version when building use the `-ldflags` option of the `build` or `install` command. For example: 

```
go build -ldflags=-X=main.version=v1.0.0
```

## License
[MIT License](LICENSE)
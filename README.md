# Batch PNG to JPEG Converter
`png-to-jpeg` is a command line tool written in Go that recursively converts all PNG images in a directory to JPEGs using the Mozilla JPEG Encoder.

Update: it turns out this program is not required as the same can be done with a single line of bash:

```
find $ARG -name "*.png" -print0 | xargs -0 -I{} -P4 bash -c 'cjpeg -quality 75 -progressive -optimize -outfile "${1%.png}.jpg" $1' -- {}
```

However, I will leave this repo up because it does have the option to show `-stats`.

## Usage

```
usage: png-to-jpeg [options]

Options:
  -dir string
        Path to a directory containing PNG images to convert (default ".")
  -quality int
        Image Quality, N between 5-95 (default 75)
  -stats
        Display amount of converted images and size differences
  -version
        Print the version
  -workers int
        Maximum amount of goroutines to use (default 4)

Examples:
  png-to-jpeg -dir images
  png-to-jpeg -dir images -quality 60
  png-to-jpeg -dir images -quality 60 -workers 1
```

Example:

```
png-to-jpeg % ./png-to-jpeg -dir images -stats
converted: images/example.png to: images/example.jpg
converted: images/subdirectory/subdirectory-2/example-2.png to: images/subdirectory/subdirectory-2/example-2.jpg
converted: images/subdirectory/example.png to: images/subdirectory/example.jpg
converted: images/subdirectory/subdirectory-2/example.png to: images/subdirectory/subdirectory-2/example.jpg

converted: 4
old size: 1.8 MiB
new size: 514.7 KiB
saved:    71.58%
```

## Building
This tool requires the MozJPEG package to be installed. For Mac users it can be installed with `brew install mozjpeg`.

To make a release, clone the repo and run the following command:

```   
./make.sh release
```

The binaries will be placed in the build folder. 


## License
[MIT License](LICENSE)
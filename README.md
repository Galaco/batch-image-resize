# Batch Image Resizer

Resizes all JPEG and PNG images in a directory. Preserves filenames.

This does not check child directories, only the specified directory.

This will preserve aspect ratio.



### Usage
First build the binary:
`go build ./main.go`

Run the binary with the following options:
```-dir="./pictures-to-resize" -output="./pictures-resized" -maxWidth=300 -maxHeight=300```

**NOTE: The output directory must already exist**

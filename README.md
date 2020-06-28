# Batch Image Resizer

Resizes all JPEG and PNG images in a directory. Preserves filenames.

This does not check child directories, only the specified directory.

This will preserve aspect ratio.


### Usage
Super simple; run the binary with the following options:
```bash
./batch-image-resizer -dir="./pictures-to-resize" -output="./pictures-resized" -maxWidth=300 -maxHeight=300
```

**NOTE: The output directory must already exist**
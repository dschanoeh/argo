[![Build Status](https://travis-ci.org/dschanoeh/argo.svg?branch=master)](https://travis-ci.org/dschanoeh/argo)
# argo

Asset resizing in go.

argo reads an input folder with picture assets, processes them, and writes
the resulting images to an output folder.
Possible options for the processing are:

* Resizing
* Quality adjustments
* Progressive JPEG

Different variants of the same image can be generated to allow optimization for
multiple screen resolutions.

1. [Installation](#installation)
1. [Usage](#usage)

## Installation
Ubuntu:
Assuming that you already installed go and set your $GOPATH, you'll just need to:

```
sudo apt install imagemagick
go install github.com/dschanoeh/argo
```

## Usage
Edit the config file:

```
# folder with input files
inputFolder = "./input/"
# folder with output files
outputFolder = "./output/"
# suffixes to be assigned to the different sizes
suffixes = ["", "@1x"]
# widths for the suffixes
widths = [1520, 760]
# JPEG qualities to be used when saving
qualities = [85, 75]
# Don't overwrite existing files
noOverwrite = false
# Progressive JPEG
progressive = true
```

And then run argo:

```
$ argo config.toml
Reading from input folder /home/dschanoeh/go/src/github.com/dschanoeh/argo/input
Processing file DSCF0849.jpg [1 / 1]
        Original dimensions: 1520 x 1013
        Original size: 388 kb
        → /home/dschanoeh/go/src/github.com/dschanoeh/argo/output/DSCF0849.jpg (172 kb) 44.37%
        → /home/dschanoeh/go/src/github.com/dschanoeh/argo/output/DSCF0849@1x.jpg (47 kb) 12.16%
```



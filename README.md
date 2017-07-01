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
Edit the config file and then point argo to it.
```
argo config.toml
```



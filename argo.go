package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bclicn/color"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// Config represents the argo configuration read from the config file
type Config struct {
	InputFolder  string
	OutputFolder string
	MinWidth     uint
	Suffixes     []string
	Widths       []uint
	Qualities    []uint
	NoOverwrite  bool
	Progressive  bool
}

var config Config
var l log.Logger

func main() {
	l := log.New(os.Stderr, "", 0)

	if len(os.Args) != 2 {
		usage()
		l.Fatal("No config file provided\n")
	}

	_, err := toml.DecodeFile(os.Args[1], &config)
	if err != nil {
		l.Fatal("Could not parse config file\n", err)
	}

	config.InputFolder, err = filepath.Abs(config.InputFolder)
	if err != nil {
		l.Fatal("Could not parse input folder name\n", err)
	}

	config.OutputFolder, err = filepath.Abs(config.OutputFolder)
	if err != nil {
		l.Fatal("Could not parse output folder name\n", err)
	}

	fmt.Println("Reading from input folder " + config.InputFolder)

	files, err := ioutil.ReadDir(config.InputFolder)
	if err != nil {
		l.Fatal(err)
	}

	imagick.Initialize()
	defer imagick.Terminate()

	for i, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		fmt.Printf("Processing file %s [%d / %d]\n", file.Name(), i+1, len(files))
		handleFile(file, config.InputFolder)
	}
}

func usage() {
	fmt.Printf("Usage: argo config.toml\n")
}

func handleFile(file os.FileInfo, path string) {
	// Skip files without known extension
	if !strings.HasSuffix(file.Name(), ".jpg") {
		fmt.Printf("Skipping: %s", file.Name())
		return
	}

	filename := path + string(os.PathSeparator) + file.Name()

	fileHandle, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open file %s. Skipping...", file.Name())
		return
	}

	mw := imagick.NewMagickWand()
	err = mw.PingImageFile(fileHandle)
	if err != nil {
		fmt.Printf("Wasn't able to read file %s: %s", file.Name(), err.Error())
		return
	}
	fileHandle.Close()

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	fmt.Printf("\tOriginal dimensions: %d x %d\n", width, height)

	size := getSize(filename)
	fmt.Printf("\tOriginal size: %d kb\n", size/1024)

	for i := range config.Suffixes {
		suffix := config.Suffixes[i]
		targetWidth := config.Widths[i]
		quality := config.Qualities[i]

		if width < targetWidth {
			fmt.Printf(color.Red("\tWarning: File is smaller than minimum width (%d px).\n"), targetWidth)
		}

		fileNameIn, fileNameOut := fullFileNames(file.Name(), suffix)

		// Skip it the output file exists
		if config.NoOverwrite && fileExists(fileNameOut) {
			fmt.Printf("\t✓ skipping (exists) %s", fileNameOut)
			continue
		}

		fmt.Printf("\t→ %s", fileNameOut)

		err := writeToOutput(fileNameIn, fileNameOut, targetWidth, quality)

		if err != nil {
			fmt.Printf(color.Red(" ✘ failed\n"))
			continue
		}

		newSize := getSize(fileNameOut)
		percentage := float64(newSize) / float64(size) * 100
		fmt.Printf(" (%d kb) ", newSize/1024)
		if percentage > 100 {
			fmt.Printf(color.Red("%.2f%%\n"), percentage)
		} else {
			fmt.Printf(color.Green("%.2f%%\n"), percentage)
		}
	}
}

func fullFileNames(fileName string, suffix string) (string, string) {
	splits := strings.Split(fileName, ".")
	name := splits[0]
	extension := splits[1]
	fileNameIn := config.InputFolder + string(os.PathSeparator) + name + "." + extension
	fileNameOut := config.OutputFolder + string(os.PathSeparator) + name + suffix + "." + extension
	return fileNameIn, fileNameOut
}

func writeToOutput(fileNameIn string, fileNameOut string, width uint, quality uint) error {
	mw := imagick.NewMagickWand()

	fileHandle, err := os.Open(fileNameIn)
	defer fileHandle.Close()
	if err != nil {
		fmt.Printf("Could not open file to read %s. Skipping...", fileNameIn)
		return err
	}

	fileHandleOut, err := os.Create(fileNameOut)
	defer fileHandleOut.Close()
	if err != nil {
		fmt.Printf("Could not open file to write %s. Skipping...", fileNameIn)
		return err
	}

	err = mw.ReadImageFile(fileHandle)
	if err != nil {
		fmt.Printf("Wasn't able to read file %s: %s", fileNameIn, err.Error())
		return err
	}

	originalWidth := mw.GetImageWidth()
	originalHeight := mw.GetImageHeight()
	scaleFactor := float64(width) / float64(originalWidth)
	height := uint(scaleFactor * float64(originalHeight))

	mw.ResizeImage(width, height, imagick.FILTER_LANCZOS)
	//mw.SetSamplingFactors([]float64{1.0, 1.0, 1.0})
	mw.SetColorspace(imagick.COLORSPACE_RGB)

	if config.Progressive {
		mw.SetImageInterlaceScheme(imagick.INTERLACE_JPEG)
	}
	mw.SetCompression(imagick.COMPRESSION_JPEG)
	mw.SetCompressionQuality(quality)

	err = mw.WriteImageFile(fileHandleOut)
	if err != nil {
		fmt.Printf("Error during write: %s", err.Error())
		return err
	}

	return err
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getSize(fileName string) int64 {
	fi, e := os.Stat(fileName)
	if e != nil {
		return -1
	}
	return fi.Size()
}

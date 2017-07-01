package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bclicn/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Config struct {
	InputFolder  string
	OutputFolder string
	MinWidth     int
	Suffixes     []string
	Widths       []int
	Qualities    []int
	NoOverwrite  bool
	Progressive  bool
}

var config Config

func main() {
	_, err := toml.DecodeFile(os.Args[1], &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading from input folder " + config.InputFolder)

	files, err := ioutil.ReadDir(config.InputFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		handleFile(file, config.InputFolder)
	}
}

func handleFile(file os.FileInfo, path string) {
	// Skip files without known extension
	if !strings.HasSuffix(file.Name(), ".jpg") {
		fmt.Printf("Skipping: %s", file.Name())
		return
	}

	fmt.Println("File: " + file.Name())

	filename := path + file.Name()

	width, height, _ := getDimensions(filename)
	fmt.Printf("\tDimensions: %d x %d\n", width, height)
	size := getSize(filename)
	fmt.Printf("\tSize: %d kb\n", size/1024)

	for i, _ := range config.Suffixes {
		suffix := config.Suffixes[i]
		targetWidth := config.Widths[i]
		quality := config.Qualities[i]

		if width < config.MinWidth {
			fmt.Printf(color.Red("Warning: File is smaller than minimum width (%d kb)\n"), config.MinWidth)
		}

		fileNameIn, fileNameOut := fullFileNames(file.Name(), suffix)

		// Skip it the output file exists
		if config.NoOverwrite && fileExists(fileNameOut) {
			fmt.Printf("\t-x- %s", fileNameOut)
			continue
		}

		fmt.Printf("\t--> %s", fileNameOut)
		err := writeToOutput(fileNameIn, fileNameOut, targetWidth, quality)

		if err != nil {
			log.Fatal(err)
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
	fileNameIn := config.InputFolder + name + "." + extension
	fileNameOut := config.OutputFolder + name + suffix + "." + extension
	return fileNameIn, fileNameOut
}

func writeToOutput(fileNameIn string, fileNameOut string, width int, quality int) error {
	params := []string{fileNameIn,
		"-resize", strconv.Itoa(width),
		"-sampling-factor", "4:2:0",
		"-colorspace", "RGB",
		"-quality", strconv.Itoa(quality)}
	if config.Progressive {
		params = append(params, "-interlace")
		params = append(params, "JPEG")
	}
	params = append(params, fileNameOut)
	_, err := exec.Command("convert", params...).Output()
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

func getDimensions(filename string) (int, int, error) {
	out, err := exec.Command("identify", []string{"-ping", "-format", "\"%[w]:%[h]\"", filename}...).Output()
	if err != nil {
		log.Fatal(err)
	}
	outs := strings.Replace(string(out), "\"", "", -1)
	s := strings.Split(outs, ":")
	width, err := strconv.Atoi(s[0])
	height, err := strconv.Atoi(s[1])
	return width, height, err
}

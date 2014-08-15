package main

import (
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

var source = flag.String("source", "", "Directory containing image assets")
var destination = flag.String("destination", "", "Destination for image assets")

func main() {
	flag.Parse()
	out, err := exec.Command("bash", "-c", "/usr/bin/find "+*source+" -name \\*.png -print").CombinedOutput()
	if err != nil {
		log.Fatal(err, string(out))
	}
	imagePaths := strings.Split(string(out), "\n")
	for _, imagePath := range imagePaths {
		paths := strings.Split(imagePath, "/")
		fileName := paths[len(paths)-1]
		if len(fileName) < 4 {
			continue
		}
		imageName := fileName[0 : len(fileName)-4]
		mkdirOut, err := exec.Command("bash", "-c", "mkdir -p "+*destination+"/"+imageName+".imageset").CombinedOutput()
		if err != nil {
			log.Fatal(err, string(mkdirOut))
		}

		newPath := *destination + "/" + imageName + ".imageset"
		cpOut, err := exec.Command("bash", "-c", "cp "+imagePath+" "+newPath).CombinedOutput()
		if err != nil {
			log.Fatal(err, string(cpOut))
		}
		writeErr := ioutil.WriteFile(newPath+"/Contents.json", []byte(contentsFileTemplate(fileName)), 0644)
		if writeErr != nil {
			log.Fatal(writeErr)
		}

	}
}

func contentsFileTemplate(imageName string) string {
	return `{
		"images" : [
		{
			"idiom" : "universal",
			"scale" : "2x",
			"filename" : "` + imageName + `"
			}
			],
		"info" : {
				"version" : 1,
				"author" : "xcode"
			}
	}`
}

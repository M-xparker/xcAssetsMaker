package main

import (
	"encoding/json"
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
		scaleSplit := strings.Split(imageName, "@")
		imageName = scaleSplit[0]
		var scale string
		if len(scaleSplit) > 1 {
			scale = string(scaleSplit[1][0])
		}

		if scale == "" {
			scale = "1x"
		} else {
			scale += "x"
		}
		mkdirOut, err := exec.Command("bash", "-c", "mkdir -p "+*destination+"/"+imageName+".imageset").CombinedOutput()
		if err != nil {
			log.Fatal(err, string(mkdirOut))
		}

		newPath := *destination + "/" + imageName + ".imageset"
		cpOut, err := exec.Command("bash", "-c", "cp "+imagePath+" "+newPath).CombinedOutput()
		if err != nil {
			log.Fatal(err, string(cpOut))
		}

		var metadata *ImageSet
		data, err := ioutil.ReadFile(newPath + "/Contents.json")

		if len(data) > 0 {
			json.Unmarshal(data, &metadata)
		} else {
			metadata = new(ImageSet)
			metadata.Info.Author = "xcode"
			metadata.Info.Version = 1
			metadata.Images = []Image{}
		}
		image := Image{Idiom: "universal", Scale: scale, Filename: fileName}
		metadata.Images = append(metadata.Images, image)

		contents, err := json.Marshal(metadata)
		writeErr := ioutil.WriteFile(newPath+"/Contents.json", contents, 0644)
		if writeErr != nil {
			log.Fatal("Here", writeErr)
		}

	}
}

type SetInfo struct {
	Version int    `json:"version"`
	Author  string `json:"xcode"`
}
type Image struct {
	Idiom    string `json:"idiom"`
	Scale    string `json:"scale"`
	Filename string `json:"filename"`
}
type ImageSet struct {
	Info   SetInfo `json:"info"`
	Images []Image `json:"images"`
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Generates 50 character random string
func randomName() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	output := make([]rune, 50)
	for i := 0; i < 50; i++ {
		output[i] = letters[rand.Intn(len(letters))]
	}
	return string(output)
}

// Parses command line arguments with the flag module
func parseFlags() (path string, files bool, dirs bool, depth int) {
	filesP := flag.Bool("files", true, "Rename files")
	dirsP := flag.Bool("dirs", false, "Rename directories")
	depthP := flag.Int("depth", 1, "Depth of the rename (not yet implemented)")

	flag.Parse()

	path = flag.Arg(0)
	if len(path) == 0 {
		path = "./"
	}

	return path, *filesP, *dirsP, *depthP
}

func main() {
	path, files, dirs, depth := parseFlags()
	fmt.Println(path, files, dirs, depth)

	fileList, err := ioutil.ReadDir(".")
	handle(err)

	tempFile, err := ioutil.TempFile(os.TempDir(), "mva-")
	handle(err)
	for _, f := range fileList {
		tempFile.Write([]byte(f.Name() + "\n"))
	}
	handle(tempFile.Close())

	cmd := exec.Command("vim", tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	handle(cmd.Run())

	fileOutput, err := ioutil.ReadFile(tempFile.Name())
	handle(err)

	renamedList := strings.Split(string(fileOutput), "\n")
	renamedList = renamedList[:len(renamedList)-1]

	operations := make(map[string]string) // [original]new names

	if len(fileList) != len(renamedList) {
		log.Fatal("Original list and renamed list have different lengths")
	}
	for i := 0; i < len(renamedList); i++ {
		for j := i + 1; j < len(renamedList); j++ {
			if renamedList[i] == renamedList[j] {
				log.Fatal("Renamed list has duplicates")
			}
		}
		operations[fileList[i].Name()] = renamedList[i]
	}

	for len(operations) != 0 {
		for original, new := range operations {
			fmt.Println(original, new)
			if original == new {
				delete(operations, original)
				continue
			}
			for otherOriginal, otherNew := range operations {
				if new == otherOriginal {
					random := otherOriginal
					for _, nameIsInUse := operations[random]; nameIsInUse; { // while name is in use
						random = randomName()
					}
					fmt.Println(otherOriginal, "->", random)
					handle(os.Rename(otherOriginal, random))
					operations[random] = otherNew
					delete(operations, otherOriginal)
					break
				}
			}
			fmt.Println(original, "->", new)
			os.Rename(original, new)
		}
	}

	os.Remove(tempFile.Name())
}

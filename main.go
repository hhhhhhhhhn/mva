package main

import (
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

func main() {
	fileList, err := ioutil.ReadDir(".")
	handle(err)

	tempFile, err := ioutil.TempFile(os.TempDir(), "mva-")
	handle(err)
	for _, f := range fileList {
		tempFile.Write([]byte(f.Name() + "\n"))
	}
	handle(tempFile.Close())

	cmd := exec.Command("sh", "-c", "$EDITOR " + tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		cmd := exec.Command("vim", tempFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		handle(cmd.Run())
	}

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
			if original == new {
				delete(operations, original)
				continue
			}
			for otherOriginal, otherNew := range operations {
				if new == otherOriginal {
					fmt.Println("Name in use:", new)
					random := otherOriginal
					for true { // while name is in use
						random = randomName()
						_, nameUsed := operations[random]
						if !nameUsed {
							break
						}
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
			delete(operations, original)
		}
	}

	os.Remove(tempFile.Name())
}

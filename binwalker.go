package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {
	filepath := flag.String("file", "", "File to analyze. It calls binwalk3")
	binwalkerFolder := flag.String("output", fmt.Sprintf("binwalker-%d", time.Now().Unix()), "Directory for the binwalker output")
	extractedDir := flag.String("extracted", ".outtmp", "Folder of binwalk output")
	help := flag.Bool("help", false, "Print the help message")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	//filepath is required
	if len(*filepath) == 0 {
		log.Println("[ERROR] -file argument must be used")
		flag.Usage()
		return
	}

	log.Println("[OK] Binwalker started")

	log.Println("[INFO] Executing binwalk3")
	if len(*filepath) != 0 {
		//extract
		command := exec.Command("binwalk3", "-a", "-M", "-C", *extractedDir, "-e", *filepath)
		_, outErr := command.Output()
		if outErr != nil {
			fmt.Println("[ERROR] Binwalk3:", outErr.Error())
			return
		}
	}

	log.Println("[INFO] Creating output directory:", *binwalkerFolder)
	// Create binwalker output dir
	if err := os.MkdirAll(*binwalkerFolder, 0770); err != nil {
		log.Println("[ERROR] Couldn't create directory:", err.Error())
		return
	}

	log.Println("[INFO] Reading content of", *extractedDir)
	// read the output temp dir
	entries, readDirErr := os.ReadDir(*extractedDir)
	if readDirErr != nil {
		fmt.Println("[ERROR] Couldn't read directory:", readDirErr.Error())
		return
	}

	for _, element := range entries {
		// Check the current folder that contains the extracted files
		if strings.HasSuffix(element.Name(), ".extracted") {

			log.Println("[INFO] Analyzing", element.Name())
			// entering .outtmp/binary.bin.extracted/
			newPath := path.Join(*extractedDir, element.Name())
			newEntries, readNewDirErr := os.ReadDir(newPath)

			if readNewDirErr != nil {
				fmt.Printf("[ERROR] Coudln't read %s: %s\n", newPath, readNewDirErr.Error())
				return
			}

			// Inside of it, enumerate all the folders available
			// enumerating inside .outtmp/binary.bin.extracted/
			for _, filesDir := range newEntries {
				dirName := filesDir.Name()
				// converting filenames from hex to int in order to sort them
				n := new(big.Int)
				n.SetString(dirName, 16)

				// Add padding to the number to avoid messing up the sorting
				valueFormat := fmt.Sprintf("%09d", n)

				originalPath := path.Join(newPath, filesDir.Name())

				// enter inside the folder containing the file
				log.Printf("[INFO] Analyzing file (%s): %s\n", valueFormat, originalPath)
				lastDir, lastDirErr := os.ReadDir(originalPath)
				if lastDirErr != nil {
					fmt.Printf("[ERROR] accessing %s: %s\n", originalPath, lastDirErr.Error())
					return
				}
				// Let's move the files to our new output folder and use the following format: 123456789-filename.extension
				for _, fileName := range lastDir {
					filePath := path.Join(originalPath, fileName.Name())
					new := path.Join(*binwalkerFolder, valueFormat+fileName.Name())
					log.Println("[INFO] Moving", filePath, "to", new)
					if err := os.Rename(filePath, new); err != nil {
						log.Printf("[ERROR] Moving file: %s\n", err.Error())
					}
				}

			}
		}
	}

	// Remove leftovers
	if err := os.RemoveAll(*extractedDir); err != nil {
		log.Println("[ERROR] remove .outtmp:", err.Error())
		return
	}

	log.Println("[OK] Completed. Extracted files are in:", *binwalkerFolder)
}

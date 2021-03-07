package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"keyholders/config"
	"keyholders/helpers"
	"keyholders/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/peterhellberg/lossypng"
)

var dbError error
var files []string
var fileBlacklist []string
var deletedFiles []string
var fileCycle int32
var fileContent string

func main() {
	fileCycle = 0
	fileBlacklist = []string{"120x90", "150x150", "202x300", "404x600", "496x737", "584x438", "592x444", "690x1024", "758x564", "768x1140", "800x785", "783x600", "1024x784", "1170x785", "300x230", "496x380", "768x588", "584x309", "592x386", "496x239", "300x145", "592x309", "300x150", "496x249", "584x426", "592x426", "758x426", "768x385"}
	// For db connection
	config.DB, dbError = gorm.Open("mysql", config.DbURL(config.BuildDBConfig()))
	if dbError != nil {
		fmt.Println("Status:", dbError)
	} else {
		fmt.Println("Connection Successfully")
	}

	defer config.DB.Close()

	posts, err := models.GetAllList()
	if err != nil {
		log.Fatal(err)
	}

	postMetas, err := models.GetAllListMeta()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, item := range posts {
		fileWrite(f, item.GUID)
	}

	for _, item := range postMetas {
		fileWrite(f, item.MetaValue)
	}

	b, err := ioutil.ReadFile("data.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileContent = string(b)

	//Images has been searched
	root := helpers.DotEnvVariable("IMAGE_PATH")
	fileListOnFolder(root)
}

// For image list on folder
func fileListOnFolder(root string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			imageNames := strings.Split(info.Name(), ".")
			imageName := imageNames[len(imageNames)-1]

			if strings.Contains(imageName, "jpg") ||
				strings.Contains(imageName, "png") ||
				strings.Contains(imageName, "jpeg") ||
				strings.Contains(imageName, "gif") ||
				strings.Contains(imageName, "svg") ||
				strings.Contains(imageName, "tif") ||
				strings.Contains(imageName, "tiff") ||
				strings.Contains(imageName, "bmp") ||
				strings.Contains(imageName, "eps") {
				fileCheckOnDb(path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\n")
	fmt.Println("Process completed")
	fmt.Printf("%d file deleted...\n", len(deletedFiles))
}

func fileCheckOnDb(file string) {
	count := 0
	fileCycle++
	imagePaths := strings.Split(file, "uploads/")
	imagePath := imagePaths[len(imagePaths)-1]
	fmt.Printf("%d --> %s is checking...\n", fileCycle, imagePath)

	for _, item := range fileBlacklist {
		if strings.Contains(imagePath, "-"+item) {
			imagePath = strings.Replace(imagePath, "-"+item, "", -1)
			break
		}
	}

	// //check whether s contains substring text
	if strings.Contains(fileContent, imagePath) {
		count++
		files = append(files, file)
		fmt.Printf("%d --> ***** File is using in the post table or post_meta table...\n", fileCycle)
	}

	if count == 0 {
		fmt.Printf("%d --> ***** File is deleting...\n", fileCycle)
		e := os.Remove(file)
		if e != nil {
			log.Fatal(e)
		} else {
			deletedFiles = append(deletedFiles, file)
		}
	} else {
		fileOptimize(file)
	}

	// To search in the post table
	// postResult, _ := models.GetSearchOnPost(imagePath)
	// if len(postResult) > 0 {
	// 	count++
	// 	files = append(files, file)
	// 	fmt.Printf("%d --> ***** File is using in the post table...\n", fileCycle)
	// }
	// // To search in the post_meta table
	// metaResult, _ := models.GetSearchOnMeta(imagePath)
	// if len(metaResult) > 0 {
	// 	count++
	// 	files = append(files, file)
	// 	fmt.Printf("%d --> ***** File is using in the post_meta table...\n", fileCycle)
	// }

	// if count == 0 {
	// 	fmt.Printf("%d --> ***** File is deleting...\n", fileCycle)
	// 	e := os.Remove(file)
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	} else {
	// 		deletedFiles = append(deletedFiles, file)
	// 	}
	// } else {
	// 	fileOptimize(file)
	// }
}

func fileOptimize(file string) {
	fmt.Printf("%d --> ***** File is optimizing...\n", fileCycle)

	tempFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	imageNames := strings.Split(tempFile.Name(), ".")
	imageName := imageNames[len(imageNames)-1]

	var image image.Image

	if imageName == "jpg" || imageName == "jpeg" {
		image, err = jpeg.Decode(tempFile)
		if err != nil {
			fmt.Println(err)
		}
	}
	if imageName == "png" {
		image, err = png.Decode(tempFile)
		if err != nil {
			fmt.Println(err)
		}
	}
	if image != nil {
		resultImg := lossypng.Optimize(image, lossypng.RGBAConversion, 10)
		defer tempFile.Close()

		if imageName == "jpg" {
			jpeg.Encode(tempFile, resultImg, nil)
		}

		if imageName == "png" {
			png.Encode(tempFile, resultImg)
		}
	}
}

func fileWrite(f *os.File, text string) {
	if strings.Contains(text, ".jpg") ||
		strings.Contains(text, ".png") ||
		strings.Contains(text, ".jpeg") ||
		strings.Contains(text, ".gif") ||
		strings.Contains(text, ".svg") ||
		strings.Contains(text, ".tif") ||
		strings.Contains(text, ".tiff") ||
		strings.Contains(text, ".bmp") ||
		strings.Contains(text, ".eps") {
		f.WriteString(text + "\n")
	}
}

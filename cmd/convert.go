package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func convertPDF2CBZ(source string) (result string) {
	output := "tempcbz"
	err := prepFolder(output)
	if err != nil {
		return fmt.Sprintf("Error %f", err)
	}

	fmt.Println("Begining Extractions...")
	err = api.ExtractImagesFile(source, output, nil, nil)
	if err != nil {
		return fmt.Sprintf("Error %f", err)
	}

	fmt.Println("Creating CBZ...")
	filename := strings.Split(source, ".")
	zipname := filename[0] + ".cbz"
	archive, err := os.Create(zipname)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	files, err := os.ReadDir(output)
	if err != nil {
		return fmt.Sprintf("Error %f", err)
	}
	if len(files) < 1 {
		return "Error no files found"
	}

	for _, file := range files {
		err = addToZip(output, file, zipWriter)
		if err != nil {
			return fmt.Sprintf("Error %f", err)
		}
	}

	zipWriter.Close()
	_ = os.RemoveAll(output)

	return zipname
}

func prepFolder(path string) error {
	_ = os.RemoveAll(path)
	err := os.Mkdir(path, 0777)
	if err != nil {
		return err
	}
	return nil
}

func addToZip(output string, file os.DirEntry, zipWriter *zip.Writer) error {
	if strings.Contains(file.Name(), "thumb") {
		return nil
	}
	if strings.Contains(file.Name(), "_Fm0.") {
		return nil
	}

	filepath := output + "/" + file.Name()
	// fmt.Println("file ", output, "/", file.Name())
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zipWriter.Create(file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}

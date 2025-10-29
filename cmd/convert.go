package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func convertPDF2CBZ(source string) (result string) {
	err := validateFile(source)
	if err != nil {
		return fmt.Sprintf("Error %s", err)
	}

	output := "tempcbz"
	err = prepFolder(output)
	if err != nil {
		return fmt.Sprintf("Error %s", err)
	}

	fmt.Println("Begining Extractions...")
	err = api.ExtractImagesFile(source, output, nil, nil)
	if err != nil {
		return fmt.Sprintf("Error %s", err)
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
		return fmt.Sprintf("Error %s", err)
	}
	if len(files) < 1 {
		return "Error no files found"
	}

	for _, file := range files {
		err = addToZip(output, file, zipWriter)
		if err != nil {
			return fmt.Sprintf("Error %s", err)
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
	// check if we should filter out the file
	err := imgFileFilter(file)
	if err != nil {
		return nil
	}

	filepath := output + "/" + file.Name()
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

func imgFileFilter(file os.DirEntry) error {
	if strings.Contains(file.Name(), "thumb") {
		return errors.New("filtred file")
	}
	if strings.Contains(file.Name(), "_Fm0.") {
		return errors.New("filtred file")
	}
	filename := strings.Split(file.Name(), ".")
	fileext := strings.ToLower(filename[len(filename)-1])
	ext := []string{"png", "jpg", "jpeg", "tif", "tiff", "webp"}
	if !contains(ext, fileext) {
		return errors.New("incorrect file format")
	}

	return nil
}

func validateFile(file string) error {
	filename := strings.Split(file, ".")

	// Check the file to be processed is a PDF
	if strings.ToLower(filename[len(filename)-1]) != "pdf" {
		return errors.New("incorrect file format")
	}

	// Test is the file can be found
	_, err := os.Stat(file)
	if err != nil {
		return errors.New("file not found")
	} else {
	}

	return nil
}

func contains(list []string, term string) bool {
	for _, item := range list {
		if item == term {
			return true
		}
	}
	return false
}

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

type SourceFile struct {
	FullPath string
	Name     string
	Path     string
	Ext      string
}

func convertPDF2CBZ(source string) (result string) {
	sourceFile, err := validateFile(source)
	if err != nil {
		return fmt.Sprintf("Error %s", err)
	}

	// fmt.Println("sourceFile ", sourceFile.Name, " Path ", sourceFile.Path, " ext ", sourceFile.Ext)
	output := "tempcbz"
	err = prepFolder(output)
	if err != nil {
		return fmt.Sprintf("Error %s", err)
	}

	switch sourceFile.Ext {
	case "pdf":
		fmt.Println("Begining PDF Extractions...")
		err = api.ExtractImagesFile(source, output, nil, nil)
		if err != nil {
			return fmt.Sprintf("Error %s", err)
		}
	case "epub":
		fmt.Println("Begining EPUB Extractions...")
		err = ExtractImagesFileFromEPUB(&sourceFile, output)
		if err != nil {
			return fmt.Sprintf("Error %s", err)
		}
	}

	fmt.Println("Creating CBZ...")
	_ = os.Mkdir("output", 0777)
	zipname := "output\\" + sourceFile.Name + ".cbz"
	fmt.Printf("Creating zipname... %s \n", zipname)
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
	if !inList(ext, fileext) {
		return errors.New("incorrect file format")
	}

	return nil
}

func validateFile(file string) (SourceFile, error) {
	tName := ""
	tPath := ""
	tExt := ""
	if strings.Contains(file, "\\") {
		fileSplit := strings.Split(file, "\\")
		fileSplit2 := strings.Split(fileSplit[len(fileSplit)-1], ".")
		tName = fileSplit2[0]
		tExt = strings.ToLower(fileSplit2[len(fileSplit2)-1])
		pathLen := len(file) - len(fileSplit[len(fileSplit)-1])
		tPath = file[:pathLen]
	} else {
		fileExtSplit := strings.Split(file, ".")
		tExt = strings.ToLower(fileExtSplit[len(fileExtSplit)-1])
		nameLen := len(file) - (len(tExt) + 1)
		tName = file[:nameLen]
		tPath = ""
	}

	// fmt.Printf("Creating FullPath... %s \n", file)
	// fmt.Printf("Creating Path... %s \n", tPath)
	// fmt.Printf("Creating Name... %s \n", tName)
	// fmt.Printf("Creating Ext... %s \n", tExt)

	sourceFile := SourceFile{
		FullPath: file,
		Name:     tName,
		Path:     tPath,
		Ext:      tExt,
	}

	// Check the file to be processed is a PDF
	if sourceFile.Ext != "pdf" && sourceFile.Ext != "epub" {
		return sourceFile, errors.New("incorrect file format")
	}

	// Test is the file can be found
	_, err := os.Stat(file)
	if err != nil {
		return sourceFile, errors.New("file not found")
	}

	return sourceFile, nil
}

func inList(list []string, term string) bool {
	for _, item := range list {
		if item == term {
			return true
		}
	}
	return false
}

func ExtractImagesFileFromEPUB(sourceFile *SourceFile, outputDir string) error {
	arch, err := zip.OpenReader(sourceFile.FullPath)
	if err != nil {
		return fmt.Errorf("OpenReader error %s", err)
	}
	defer arch.Close()

	for _, zfile := range arch.File {
		if strings.Contains(strings.ToLower(zfile.Name), "oebps/image") && zfile.UncompressedSize64 > 0 {
			// fmt.Printf("Name: %s | Size: %d bytes\n", zfile.Name, zfile.UncompressedSize64)
			curZFile, err := zfile.Open()
			if err != nil {
				return fmt.Errorf("open zfile error %s", err)
			}
			defer curZFile.Close()

			filename := strings.Split(zfile.Name, "/")
			outFileName := outputDir + "/" + strings.ToLower(filename[len(filename)-1])
			outFile, err := os.Create(outFileName)
			if err != nil {
				return fmt.Errorf("error %s", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, curZFile)
			if err != nil {
				return fmt.Errorf("error %s", err)
			}
		}
	}

	fmt.Println("Extraction complete.")
	return nil
}

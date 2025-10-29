# PDF2CBZ
Converts PDF files of images into a CBZ by extracxting images. 
This application will not generate new images from pages in the PDF.

## Instructions
pdf2cbz -help
\
pdf2cbz convert "path to pdf with filename"
\
pdf2cbz batch "path to folder containing 1 or more pdf"

### filtered file names
- thumb
- _Fm0.

## To Do
 - add config file to add exclusions too
 - add option to dump inages only
 - add support for epub
 - add better help text

## Releases
 - 0.1.0 - Converts PDF, or a folder of PDFs, to CBZ

## Dev notes
1. Requires
   - github.com/spf13/cobra
   - github.com/pdfcpu/pdfcpu
2. Run main using `go run main.go ARGS`
3. Build exe using `go build`

## Author
Max Glick \
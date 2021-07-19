package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/disintegration/imaging"
)

//go:embed commit.txt
var gitCommit string // Git hash, set by build pipeline
//go:embed version.txt
var buildVersion string // human readable version, set by build pipeline
var port = "3000"

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/upload.php", uploadFile)
	http.HandleFunc("/test", test)
	// http.Handle("/public/", http.StripPrefix("/public/", fs))
	log.Println(buildVersion, "listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func test(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("tests/input.jpg")
	var ior io.Reader = file
	if err != nil {
		log.Println("file not found")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	convertImage(w, &ior)
}

func convertImage(w io.Writer, r *io.Reader) error {
	// config
	autorotate := true

	// a5, dpi 300
	targetHeight := 1748
	targetWidth := 2480
	targetRatio := float64(targetWidth) / float64(targetHeight)

	src, err := imaging.Decode(*r)
	srcWidth := src.Bounds().Max.X
	srcHeight := src.Bounds().Max.Y
	srcRatio := float64(srcWidth) / float64(srcHeight)
	if err != nil {
		log.Printf("failed to decode image: %v", err)
		return err
	}
	dst := imaging.New(targetWidth, targetHeight, color.White)

	if autorotate && srcHeight > srcWidth {
		src = imaging.Rotate90(src)
		srcWidth = src.Bounds().Max.X
		srcHeight = src.Bounds().Max.Y
		srcRatio = float64(srcWidth) / float64(srcHeight)
	}

	if targetRatio < srcRatio {
		img1 := imaging.Resize(src, targetWidth, 0, imaging.Lanczos)
		img1Height := img1.Bounds().Max.Y
		offsetTop := (targetHeight - img1Height) / 2
		dst = imaging.Paste(dst, img1, image.Pt(0, offsetTop))
	} else {
		img1 := imaging.Resize(src, 0, targetHeight, imaging.Lanczos)
		img1Width := img1.Bounds().Max.X
		offsetLeft := (targetWidth - img1Width) / 2
		dst = imaging.Paste(dst, img1, image.Pt(offsetLeft, 0))
	}

	err = imaging.Encode(w, dst, imaging.JPEG)
	if err != nil {
		log.Printf("failed to save image: %v", err)
		return err
	}
	return nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile` it also returns the FileHeader so we can get the Filename, the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer file.Close()
	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	//fmt.Printf("File Size: %+v\n", handler.Size)

	filename := fileNameWithoutExtTrimSuffix(handler.Filename)
	filename = filename + "_a5.jpg"
	w.Header().Set("Content-Type", "image/jpeg")
	// attachment offers the file to download, inline shows it in browser
	w.Header().Add("Content-Disposition", "inline; filename="+filename)
	var ior io.Reader = file
	convertImage(w, &ior)
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

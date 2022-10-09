package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Descargamos el documento PDF firmado mediante GET
func DownloadPdfSigned(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, _ := url.QueryUnescape(vars["file"])
	dirPdf, _ := url.QueryUnescape(vars["dir"])
	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[R].pdf")

	fmt.Println(filePdfSigned)

	// Open file
	f, err := os.Open(filePdfSigned)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-type", "application/pdf")
	w.Header().Add("Content-disposition", "filename="+namePdf+"[R].pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

}

// Descargamos el docuemento PDF firmado en Base64 mediante POST
func DownloadPdfSignedBase64(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, _ := url.QueryUnescape(vars["file"])
	dirPdf, _ := url.QueryUnescape(vars["dir"])
	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[R].pdf")

	fmt.Println(filePdfSigned)

	data, err := os.ReadFile(filePdfSigned)
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(500)
		return
	}

	//mimeType := http.DetectContentType(data)
	//fmt.Println(mimeType)

	previewRespuesta := map[string]interface{}{
		"codigo": 2000, //codigo interno
		"data":   base64.StdEncoding.EncodeToString(data),
	}

	respuesta, _ := json.Marshal(previewRespuesta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write(respuesta)

}

package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Exclusivamente utilizado por ReFirmaPCX para descargar los documentos (sin firmar) que esta comprimidos con 7z
func Download7z(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Descargando 7z sin firmar...")
	documentName7z := r.URL.Query().Get("documentName") + ".7z"

	filename := filepath.Join(os.TempDir(), "upload", documentName7z)

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("Content-disposition", "attachment; filename="+documentName7z)

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}

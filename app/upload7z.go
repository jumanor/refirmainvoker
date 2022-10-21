package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jumanor/refirmainvoker/logging"
	"github.com/jumanor/refirmainvoker/util"
)

// Exclusivamente utilizado por ReFirmaPCX para subir los documentos (firmados) que esta comprimidos con 7z
func Upload7z(w http.ResponseWriter, r *http.Request) {
	logging.Log().Trace().Msg("Inicio Subiendo 7z firmado...")
	// 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20) //32 Megas
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}

	var clave string
	for key := range r.MultipartForm.File {

		clave = key
	}

	if err := util.VerificarJWT(clave); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Recuperamos el archivo 7z
	f, h, err := r.FormFile(clave)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}

	fileName, _ := url.QueryUnescape(h.Filename)
	fileNameWithOutExtension := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	signedPath := filepath.Join(os.TempDir(), "upload", "signed")
	os.MkdirAll(signedPath, os.ModePerm)
	filePathSigned := filepath.Join(signedPath, fileName)

	file, err := os.OpenFile(filePathSigned, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}

	// Persistimos al disco duro
	bytes, err := io.Copy(file, f)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}
	f.Close()
	file.Close()

	// Descomprimimos 7z
	c := exec.Command("7z", "x", filePathSigned, "-o"+filepath.Join(signedPath, fileNameWithOutExtension))

	if err := c.Run(); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}

	fmt.Printf("The number of bytes are: %d\n", bytes)
	logging.Log().Debug().Str("7z", filePathSigned).Msg("Archivo [R]7z signed descomprimido satisfactoriamente")

	w.WriteHeader(200) //SUCCESS
	w.Write([]byte(""))
}

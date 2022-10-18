package app

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jumanor/refirmainvoker/logging"
)

// Descargamos el documento PDF firmado mediante GET
func DownloadPdfSigned(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, err := url.QueryUnescape(vars["file"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de " + vars["file"]))
		return
	}
	dirPdf, err := url.QueryUnescape(vars["dir"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de  " + vars["dir"]))
		return
	}
	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[R].pdf")

	logging.Log().Trace().Msgf("descargando pdf %s", filePdfSigned)

	// Open file
	f, err := os.Open(filePdfSigned)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo abrir el archivo" + namePdf))
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-type", "application/pdf")
	w.Header().Add("Content-disposition", "filename="+namePdf+"[R].pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo copiar el archivo " + namePdf + " en flujo de envio "))
	}

}

// Descargamos el docuemento PDF firmado en Base64 mediante POST
func DownloadPdfSignedBase64(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, err := url.QueryUnescape(vars["file"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de " + vars["file"]))
		return
	}
	dirPdf, _ := url.QueryUnescape(vars["dir"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de  " + vars["dir"]))
		return
	}
	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[R].pdf")

	logging.Log().Trace().Msgf("descargando pdf b64 %s", filePdfSigned)

	data, err := os.ReadFile(filePdfSigned)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo leer el archivo" + namePdf))
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

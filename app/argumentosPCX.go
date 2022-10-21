package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jumanor/refirmainvoker/logging"
	"github.com/jumanor/refirmainvoker/util"
)

var CLIENT_ID string = ""
var CLIENT_SECRET string = ""
var MAX_FILE_SIZE_7Z string

// Cadena de argumentos que se envia a refirma PCX
func paramWeb(documentName string, fileDownloadUrl string, fileDownloadLogoUrl string, fileDownloadStampUrl string,
	fileUploadUrl string, posx string, posy string, reason string, token string) ([]byte, error) {

	param := make(map[string]string)

	param["app"] = "pcx"
	param["mode"] = "lot-p"
	param["clientId"] = CLIENT_ID
	param["clientSecret"] = CLIENT_SECRET
	param["idFile"] = token
	param["type"] = "W"
	param["protocol"] = "T"
	param["fileDownloadUrl"] = fileDownloadUrl
	param["fileDownloadLogoUrl"] = fileDownloadLogoUrl
	param["fileDownloadStampUrl"] = fileDownloadStampUrl
	param["fileUploadUrl"] = fileUploadUrl
	param["contentFile"] = documentName
	param["reason"] = reason
	param["isSignatureVisible"] = "true"
	param["stampAppearanceId"] = "0"
	param["pageNumber"] = "0"
	param["posx"] = posx
	param["posy"] = posy
	param["fontSize"] = "7"
	param["dcfilter"] = ".*FIR.*|.*FAU.*"
	param["signatureLevel"] = "0"
	param["outputFile"] = "out-" + documentName
	param["maxFileSize"] = MAX_FILE_SIZE_7Z //Por defecto será 5242880 5MB - Maximo 100MB

	respuesta, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("No se pudo parsear a json la cadena de argumentos")
	}

	logging.Log().Debug().Msg(string(respuesta))

	return respuesta, nil
}

type ResultCanalDescarga struct {
	Message string
	Error   error
}

// Descargamos y creamos el documento PDF en el HD mediante una gorutina
func downloadPdfAndPersist(rutaMain string, pdf struct {
	URL  string "json:\"url\""
	Name string "json:\"name\""
}, ch chan ResultCanalDescarga) {

	fmt.Println(pdf)

	out, err := os.Create(filepath.Join(rutaMain, pdf.Name+".pdf"))
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo crear archivo: " + pdf.Name)}
		return
	}

	client := http.Client{
		Timeout: 60 * time.Second, //timeout 60 segundos
	}
	resp, err := client.Get(pdf.URL)
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo descargar: " + pdf.URL)}
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("Error : " + resp.Status + "  " + pdf.URL)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo descargar: " + pdf.URL)}
		return

	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo copiar a disco: " + pdf.URL)}
		return
	}

	out.Close()
	resp.Body.Close()

	ch <- ResultCanalDescarga{Message: "OK " + pdf.Name, Error: nil}
}

// Descargamos todos los documentos PDF concurrentemente.
// Todos los documentos deben ser persistidos caso contrario se lanza en error
func downloadAllPdfAndPersistConcurrency(rutaMain string, urls Pdf) error {

	ch := make(chan ResultCanalDescarga)

	for _, pdf := range urls {
		go downloadPdfAndPersist(rutaMain, pdf, ch) //usamos go rutinas
	}

	for range urls {

		result := <-ch //bloqueamos a la espera de la respuesta

		if result.Error != nil {
			fmt.Println(result.Error)

			return result.Error
		}

		//fmt.Println(result.Message)

	}

	return nil
}

// Creamos un archivo 7z con los PDFs descargados
func createFile7z(urls Pdf) (string, error) {

	nameUUID := uuid.New().String()

	rutaMain := filepath.Join(os.TempDir(), "upload", nameUUID)

	if err := os.MkdirAll(rutaMain, os.ModePerm); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", errors.New("No se puede crear el directorio " + rutaMain)
	}

	if err := downloadAllPdfAndPersistConcurrency(rutaMain, urls); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", err
	}

	file7z := filepath.Join(rutaMain, "..", nameUUID+".7z")
	c := exec.Command("7z", "a", file7z, rutaMain+string(filepath.Separator)+".")

	if err := c.Run(); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", errors.New("No se pudo comprimir a 7z")
	}

	logging.Log().Debug().Str("7z", file7z).Msg("Archivo 7z creado satisfactoriamente")
	return nameUUID, nil
}

type Pdf []struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
type DatoArgumentos struct {
	Pdfs  Pdf `json:"pdfs"`
	Firma struct {
		Posx        int    `json:"posx"`
		Posy        int    `json:"posy"`
		Reason      string `json:"reason"`
		StampSigned string `json:"stampSigned"`
	} `json:"firma"`
}

// URI llamado por el Cliente para contruir Cadena de Argumentos en BASE64
func ArgumentsServletPCX(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("x-access-token")
	if err := util.VerificarJWT(token); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Acceso no autorizado"))
		return
	}

	var inputParameter DatoArgumentos
	err := json.NewDecoder(r.Body).Decode(&inputParameter)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte("No se pudo parsear a json los parametros de entrada"))
		return
	}

	serverURL := "http://" + r.Host

	documentNameUUID, err := createFile7z(inputParameter.Pdfs)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte(err.Error()))
		return
	}

	documentName7z := documentNameUUID + ".7z"
	fileDownloadUrl := serverURL + "/download7z?documentName=" + url.QueryEscape(documentNameUUID) + "&token=" + url.QueryEscape(token)
	fileDownloadLogoUrl := serverURL + "/public/iLogo.png"
	fileDownloadStampUrl := serverURL + "/public/iFirma.png"
	fileUploadUrl := serverURL + "/upload7z"
	posx := strconv.Itoa(inputParameter.Firma.Posx)
	posy := strconv.Itoa(inputParameter.Firma.Posy)
	reason := inputParameter.Firma.Reason

	if inputParameter.Firma.StampSigned != "" {
		fileDownloadStampUrl = inputParameter.Firma.StampSigned
	}

	param, err := paramWeb(documentName7z, fileDownloadUrl, fileDownloadLogoUrl,
		fileDownloadStampUrl, fileUploadUrl, posx, posy, reason, token)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte(err.Error()))
		return
	}

	argumentosEnc := base64.StdEncoding.EncodeToString(param)

	urlBasePDFDownloadSigned := serverURL + "/downloadPdfSigned/" + url.QueryEscape(documentNameUUID)

	previewRespuesta := map[string]interface{}{
		"codigo": 2000, //codigo interno
		"data": map[string]interface{}{
			"argumentosBase64": argumentosEnc,
			"urlBase":          urlBasePDFDownloadSigned,
		},
	}

	respuesta, _ := json.Marshal(previewRespuesta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write(respuesta)

} //////////////////////////////////////////////////////////////////////////////////////

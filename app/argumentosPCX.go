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

	"github.com/google/uuid"
)

var CLIENT_ID string = ""
var CLIENT_SECRET string = ""

// Cadena de argumentos que se envia a refirma PCX
func paramWeb(documentName string, fileDownloadUrl string, fileDownloadLogoUrl string, fileDownloadStampUrl string,
	fileUploadUrl string, posx string, posy string) []byte {

	param := make(map[string]string)

	param["app"] = "pcx"
	param["mode"] = "lot-p"
	param["clientId"] = CLIENT_ID
	param["clientSecret"] = CLIENT_SECRET
	param["idFile"] = "001"
	param["type"] = "W"
	param["protocol"] = "T"
	param["fileDownloadUrl"] = fileDownloadUrl
	param["fileDownloadLogoUrl"] = fileDownloadLogoUrl
	param["fileDownloadStampUrl"] = fileDownloadStampUrl
	param["fileUploadUrl"] = fileUploadUrl
	param["contentFile"] = documentName
	param["reason"] = "soy autor del documento"
	param["isSignatureVisible"] = "true"
	param["stampAppearanceId"] = "0"
	param["pageNumber"] = "0"
	param["posx"] = posx
	param["posy"] = posy
	param["fontSize"] = "7"
	param["dcfilter"] = ".*FIR.*|.*FAU.*"
	param["signatureLevel"] = "0"
	param["outputFile"] = "out-" + documentName
	param["maxFileSize"] = "5242880"

	respuesta, _ := json.Marshal(param)

	fmt.Println(string(respuesta))

	return respuesta
}

// Creamos un archivo 7z con los PDFs descargados
func createFile7z(urls Pdf) (string, error) {

	nameUUID := uuid.New().String()

	rutaMain := filepath.Join(os.TempDir(), "upload", nameUUID)
	fmt.Println(rutaMain)
	os.MkdirAll(rutaMain, os.ModePerm)

	for _, v := range urls { //descargar los PDFs

		fmt.Println(v)

		out, err := os.Create(filepath.Join(rutaMain, v.Name+".pdf"))
		if err != nil {
			fmt.Println(err)
			return "", errors.New("No se pudo crear archivo: " + v.Name)
		}

		resp, err := http.Get(v.URL)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("No se pudo descargar: " + v.URL)
		}

		io.Copy(out, resp.Body)

		out.Close()
		resp.Body.Close()

	}

	file7z := filepath.Join(rutaMain, "..", nameUUID+".7z")
	c := exec.Command("7z", "a", file7z, rutaMain+string(filepath.Separator)+".")

	if err := c.Run(); err != nil {
		fmt.Println(err)
		return "", errors.New("No se pudo comprimir a 7z")
	}

	return nameUUID, nil
}

type Pdf []struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
type DatoArgumentos struct {
	Pdfs  Pdf `json:"pdfs"`
	Firma struct {
		Posx int `json:"posx"`
		Posy int `json:"posy"`
	} `json:"firma"`
}

// URI llamado por el Cliente para contruir Cadena de Argumentos en BASE64
func ArgumentsServletPCX(w http.ResponseWriter, r *http.Request) {

	var argumentos DatoArgumentos
	err := json.NewDecoder(r.Body).Decode(&argumentos)
	if err != nil {
		panic(err)
	}

	serverURL := "http://" + r.Host

	documentNameUUID, err := createFile7z(argumentos.Pdfs)
	if err != nil {
		panic(err)
	}

	documentName7z := documentNameUUID + ".7z"
	fileDownloadUrl := serverURL + "/download7z?documentName=" + url.QueryEscape(documentNameUUID)
	fileDownloadLogoUrl := serverURL + "/public/iLogo.png"
	fileDownloadStampUrl := serverURL + "/public/iFirma.png"
	fileUploadUrl := serverURL + "/upload7z"
	posx := strconv.Itoa(argumentos.Firma.Posx)
	posy := strconv.Itoa(argumentos.Firma.Posy)

	argumentosEnc := base64.StdEncoding.EncodeToString(paramWeb(documentName7z, fileDownloadUrl, fileDownloadLogoUrl,
		fileDownloadStampUrl, fileUploadUrl, posx, posy))

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

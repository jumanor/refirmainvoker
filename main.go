package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/jumanor/refirmainvoker/app"
	"github.com/jumanor/refirmainvoker/util"
)

var SERVER_ADDRESS string = "0.0.0.0:9091"

func identificadores() {

	if len(os.Args) > 1 { //Leemos de argumentos
		app.CLIENT_ID = os.Args[1]
		app.CLIENT_SECRET = os.Args[2]
		SERVER_ADDRESS = os.Args[3]

	} else { //Leemos de archivo properties

		abs_fname, _ := filepath.Abs("./")
		ruta := abs_fname + string(filepath.Separator) + "config.properties"

		properties, err := util.ReadPropertiesFile(ruta)
		if err != nil {
			log.Fatal(err)
		}
		app.CLIENT_ID, _ = properties["clientId"]
		app.CLIENT_SECRET, _ = properties["clientSecret"]
		SERVER_ADDRESS, _ = properties["serverAddress"]
	}

}

func main() {
	identificadores()

	enrutador := mux.NewRouter()

	fs := http.FileServer(http.Dir("./public/"))
	enrutador.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	enrutador.HandleFunc("/argumentsServletPCX", app.ArgumentsServletPCX).Methods("POST")
	enrutador.HandleFunc("/download7z", app.Download7z).Methods("GET")
	enrutador.HandleFunc("/upload7z", app.Upload7z).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}", app.DownloadPdfSignedBase64).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}", app.DownloadPdfSigned).Methods("GET")

	enrutador.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)

	enrutador.Use(util.EnableCors)

	servidor := &http.Server{
		Handler: enrutador,
		Addr:    SERVER_ADDRESS,
		// Timeouts para evitar que el servidor se quede "colgado" por siempre
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Escuchando en %s. Presiona CTRL + C para salir", SERVER_ADDRESS)

	log.Fatal(servidor.ListenAndServe())

} //////////////////////////////////////////////////////////////////////////////////////

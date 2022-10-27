package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jumanor/refirmainvoker/app"
	"github.com/jumanor/refirmainvoker/logging"
	"github.com/jumanor/refirmainvoker/util"
)

var SERVER_ADDRESS string
var CERTIFICATE_FILE_TLS string
var PRIVATE_KEY_FILE_TLS string

func init() {

	util.TIME_EXPIRE_TOKEN = 5        //5 Minutos
	app.MAX_FILE_SIZE_7Z = "10485760" //10 Megas

	if len(os.Args) > 1 { //Leemos de argumentos
		app.CLIENT_ID = os.Args[1]
		app.CLIENT_SECRET = os.Args[2]
		SERVER_ADDRESS = os.Args[3]
		util.SECRET_KEY_JWT = os.Args[4]
		app.USER_ACCESS_API = os.Args[5]

		if len(os.Args) >= 7 {
			if n, err := strconv.ParseInt(os.Args[6], 10, 64); err != nil {
				panic(err)
			} else {
				util.TIME_EXPIRE_TOKEN = n
			}
		}
		if len(os.Args) == 8 {
			app.MAX_FILE_SIZE_7Z = os.Args[7]
		}

	} else { //Leemos de archivo properties

		abs_fname, _ := filepath.Abs("./")
		ruta := abs_fname + string(filepath.Separator) + "config.properties"

		properties, err := util.ReadPropertiesFile(ruta)
		if err != nil {
			panic(err)
		}
		app.CLIENT_ID = properties["clientId"]
		app.CLIENT_SECRET = properties["clientSecret"]
		SERVER_ADDRESS = properties["serverAddress"]
		util.SECRET_KEY_JWT = properties["secretKeyJwt"]
		app.USER_ACCESS_API = properties["userAccessApi"]

		if properties["timeExpireToken"] != "" {
			if exp, err := strconv.ParseInt(properties["timeExpireToken"], 10, 64); err != nil {
				panic(err)
			} else {
				util.TIME_EXPIRE_TOKEN = exp
			}
		}
		if properties["maxFileSize7z"] != "" {
			app.MAX_FILE_SIZE_7Z = properties["maxFileSize7z"]
		}
		if properties["certificateFileTls"] != "" {
			CERTIFICATE_FILE_TLS = properties["certificateFileTls"]
		}
		if properties["privateKeyFileTls"] != "" {
			PRIVATE_KEY_FILE_TLS = properties["privateKeyFileTls"]
		}

	}

	logging.Log().Trace().Str("CLIENT_ID", app.CLIENT_ID).Str("CLIENT_SECRET", app.CLIENT_SECRET).
		Str("SERVER_ADDRESS", SERVER_ADDRESS).Str("SECRET_KEY_JWT", util.SECRET_KEY_JWT).
		Str("USER_ACCESS_API", app.USER_ACCESS_API).Int64("TIME_EXPIRE_TOKEN", util.TIME_EXPIRE_TOKEN).
		Str("MAX_FILE_SIZE_7Z", app.MAX_FILE_SIZE_7Z).Send()
}
func main() {

	enrutador := mux.NewRouter()

	fs := http.FileServer(http.Dir("./public/"))
	enrutador.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	enrutador.HandleFunc("/argumentsServletPCX", app.ArgumentsServletPCX).Methods("POST")
	enrutador.HandleFunc("/download7z", app.Download7z).Methods("GET")
	enrutador.HandleFunc("/upload7z", app.Upload7z).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}", app.DownloadPdfSignedBase64).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}/{token}", app.DownloadPdfSigned).Methods("GET")
	enrutador.HandleFunc("/autenticacion", app.Autenticacion).Methods("POST")

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

	if CERTIFICATE_FILE_TLS != "" && PRIVATE_KEY_FILE_TLS != "" {

		logging.Log().Info().Str("Scheme", "https").Msgf("Escuchando en %s. Presiona CTRL + C para salir", SERVER_ADDRESS)
		err := servidor.ListenAndServeTLS(CERTIFICATE_FILE_TLS, PRIVATE_KEY_FILE_TLS)
		logging.Log().Fatal().Err(err).Send()

	} else {

		logging.Log().Info().Str("Scheme", "http").Msgf("Escuchando en %s. Presiona CTRL + C para salir", SERVER_ADDRESS)
		err := servidor.ListenAndServe()
		logging.Log().Fatal().Err(err).Send()
	}

} //////////////////////////////////////////////////////////////////////////////////////

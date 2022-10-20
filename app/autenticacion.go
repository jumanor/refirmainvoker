package app

import (
	"encoding/json"
	"net/http"

	"github.com/jumanor/refirmainvoker/logging"
	"github.com/jumanor/refirmainvoker/util"
)

var USER_ACCESS_API string

type InputParametrosAutenticacion struct {
	UsuarioAccesoApi string `json:"usuarioAccesoApi"`
}

// URI para autenticacion
func Autenticacion(w http.ResponseWriter, r *http.Request) {

	var inputParameter InputParametrosAutenticacion
	err := json.NewDecoder(r.Body).Decode(&inputParameter)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte("No se pudo parsear a json los parametros de entrada"))
		return
	}

	//Una muy simple forma de validar acceso a las API
	if inputParameter.UsuarioAccesoApi != USER_ACCESS_API {
		logging.Log().Info().Msgf("Usuario %s incorrecto", inputParameter.UsuarioAccesoApi)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario incorrecto"))
		return
	}

	tokenString, err := util.GenerarJWT()
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	previewRespuesta := map[string]interface{}{
		"codigo": 2000, //codigo interno
		"data":   tokenString,
	}

	respuesta, _ := json.Marshal(previewRespuesta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write(respuesta)
}

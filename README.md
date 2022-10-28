![Go](https://img.shields.io/badge/Golang-1.19-blue.svg?logo=go&longCache=true&style=flat)
# Refirma Invoker - Reniec
Implementación del Motor de Firma Digital - Refirma Invoker - del [RENIEC](https://dsp.reniec.gob.pe/refirma_suite/main/web/main.jsf)

***Refirma Invoker*** es parte de la ***ReFirma Suite*** y su uso es gratuito para las Entidades Públicas del Perú, su funcionamiento lo puede ver [acá](https://drive.google.com/file/d/1S-FrH2HX6vawsO4oXESUHwDsSQjJOGMj/view?usp=sharing)

La versión [v1.1.0](https://github.com/jumanor/refirmainvoker/tree/v1.1.0) es el último lanzamiento

Para mayor información de esta implementación puede ver en el siguiente [video](https://www.youtube.com/watch?v=aOto5CStZNA)

# Requisitos

Para ejecutar Refirma Invoker es necesario que se contacte con la [RENIEC](https://dsp.reniec.gob.pe/refirma_suite/main/web/main.jsf) para que le brinde los identificadores **[clientId]** y **[clientSecret]**

Esta implementación usa [7-zip](https://www.7-zip.org/) que normalmente ya viene instalada en **LINUX**; sin embargo, en **WINDOWS** tendra que instalar manualmente y verificar que se puede acceder desde el terminal.

Windows
    
    C:\Users\Jumanor>7z i

Linux

    jumanor@ubuntu:~$7z i

# Instalación

Esta disponible un video de la instalación en el siguiente [enlace](https://www.youtube.com/watch?v=7q4dS8y3Sws)

Se compilo Refirma Invoker para Windows y Linux, y estan disponibles en los [releases](https://github.com/jumanor/refirmainvoker/releases/tag/v1.1.0).

1. Descargue el ejecutable
   
   Windows: [main.exe](https://github.com/jumanor/refirmainvoker/releases/download/v1.1.0/main.exe)
   
   Linux:   [main](https://github.com/jumanor/refirmainvoker/releases/download/v1.1.0/main)

2. Copia la carpeta **public** del repositorio esta contiene 2 imagenes: iFirma.png e iLogo.png
3. Crea un archivo **config.properties** con los siguientes parametros :
    ``` bash
    # Identificador proporcionado por RENIEC
    clientId=K57845459hkj
    # Identificador proporcionado por RENIEC
    clientSecret=TYUOPDLDFDG
    # Direccion Ip y Puerto de escucha ReFirma Invoker
    serverAddress=0.0.0.0:9091
    # Clave secreta para generar Tokens
    secretKeyJwt=muysecretokenjwt
    # Usuario que accedera a la API
    userAccessApi=usuarioAccesoApi
    # Tiempo de expiración del Token en minutos. Ejemplo 5 minutos (Opcional)
    timeExpireToken=5
    # Maximo tamaño del archivo 7z en bytes. Ejemplo 10 megas (Opcional)
    maxFileSize7z=10485760
    # Certificado SSL/TLS (Opcional)
    #certificateFileTls=C:\Users\jumanor\cert.pem
    #certificateFileTls=/home/jumanor/cert.pem
    # Clave Privada SSL/TLS  (Opcional)
    #privateKeyFileTls=C:\Users\jumanor\key.pem
    #privateKeyFileTls=/home/jumanor/key.pem
    ``` 
4. En caso desee habilitar protocolo **https** es necesario que ingrese los siguientes parametros :
    ``` bash
    # Certificado SSL/TLS (Opcional)
    certificateFileTls=/home/jumanor/cert.pem
    # Clave Privada SSL/TLS (Opcional)
    privateKeyFileTls=/home/jumanor/key.pem
    ```
5. Ejecuta ReFirma Invoker

    Windows

        main.exe

    Linux

        ./main

# Funcionamiento

Esta disponible un video del funcionamiento en el siguiente [enlace](https://youtu.be/7q4dS8y3Sws?t=218)

Copia la carpeta [example](https://github.com/jumanor/refirmainvoker/tree/master/example) de este repositorio en un Servidor Web y ejecuta **test.html**.

A pesar que todos los archivos del ejemplo son **estaticos** es necesario usar un Servidor Web debido a que se esta usando **ES6**.

En caso use **Visual Studio Code** instale el plugin [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) que habilita un Servidor Web.

En el ejemplo se usa **Microsoft Click Once**; por lo tanto, si esta usando navegador **Chome** o **Firefox** instales los siguientes plugins:

- Chrome instale este [plugin](https://chrome.google.com/webstore/detail/clickonce-for-google-chro/kekahkplibinaibelipdcikofmedafmb) 

- Firefox instale este [plugin](https://addons.mozilla.org/es/firefox/addon/meta4clickoncelauncher/?utm_source=addons.mozilla.org&utm_medium=referral&utm_content=search)

En caso use el navegador **Edge** no es necesario instalar nada adicional.

Esta implementacion de Refirma Invoker usa JSON Web Tokens ([JWT](https://jwt.io/)).

``` javascript
//Listamos los documentos que se desean firmar digitalmente
let pdfs=[];
pdfs[0]={url:"http://miservidor.com/docs1.pdf",name:"doc1"};
pdfs[1]={url:"http://miservidor.com/docs2.pdf",name:"doc2"};

//Enviamos la posicion en donde se ubicara la representación gráfica de la firma digital
let firmaParam={};
firmaParam.posx=10;
firmaParam.posy=12;
firmaParam.reason="Soy el autor del documento pdf";
firmaParam.stampSigned="http://miservidor.com/estampillafirma.png";//parametro opcional

//Llamamos a Refirma Invoker con la dirección ip en donde se ejecuta main.exe o main
let firma=new RefirmaInvoker("http://192.168.1.10:9091");
//Importante:
//El Sistema de Gestion Documental se encarga de la autenticación y envía un token al Cliente
//Este método se usa solo como demostración no se debe de usar en el Cliente
let token=await firma.autenticacion("usuarioAccesoApi");
//Realiza el proceso de Firma Digital
let url_base=await firma.ejecutar(pdfs,firmaParam,token);

//En este caso obtenemos los documentos firmados digitalmente y los enviamos a un frame
document.getElementById("frame1").src=url_base+"/"+encodeURI("doc1")+"/"+encodeURI(token);
document.getElementById("frame2").src=url_base+"/"+encodeURI("doc2")+"/"+encodeURI(token);
```          

El *Sistema de Gestión Documental* autentica a los Usuarios normalmente contra una Base de Datos,
despues de la autencación satisfactoria se debe de consumir  el API REST /autenticacion de ReFirma Invoker 
y enviar el **token** al Cliente.

![a link](https://drive.google.com/uc?export=view&id=1cvW39tUvRkCctHetWvhqJskCnXQwRWKn)

Ejemplo en Python
``` python
import requests
import json
api_url = "http://127.0.0.1:9091/autenticacion"
param={"usuarioAccesoApi":"usuarioAccesoApi"}
response = requests.post(api_url,json=param)
if response.status_code == 200:
	token=response.json().get("data")
	print(token)

```
Ejemplo en Php
``` php
$params=array("usuarioAccesoApi"=>"usuarioAccesoApi");
$postdata=json_encode($params);
$opts = array('http' =>
    array(
    'method' => 'POST',
    'header' => 'Content-type: application/json',
    'content' => $postdata
    )
);
$context = stream_context_create($opts);
@$response = file_get_contents("http://127.0.0.1:9091/autenticacion", false, $context);
if(isset($http_response_header) == true){
    
    $status_line = $http_response_header[0];
    preg_match('{HTTP\/\S*\s(\d{3})}', $status_line, $match);
    $status = $match[1];

    if ($status == 200){
        $obj=json_decode($response,true);
        $token=$obj["data"];
        echo $token;
    }    
}
```

# Contribución

Por favor contribuya usando [Github Flow](https://guides.github.com/introduction/flow/). Crea un fork, agrega los commits, y luego abre un [pull request](https://github.com/fraction/readme-boilerplate/compare/).

# License
Copyright © 2022 [Jorge Cotrado](https://github.com/jumanor). <br />
This project is [MIT](https://github.com/jumanor/refirmainvoker/blob/master/License) licensed.

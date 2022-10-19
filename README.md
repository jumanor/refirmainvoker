# Refirma Invoker - Reniec
Implementación del Motor de Firma Digital - Refirma Invoker - del [RENIEC](https://dsp.reniec.gob.pe/refirma_suite/main/web/main.jsf)

***Refirma Invoker*** es parte de la ***ReFirma Suite*** y su uso es gratuito para las Entidades Públicas del Perú, su funcionamiento lo puede ver [acá](https://drive.google.com/file/d/1S-FrH2HX6vawsO4oXESUHwDsSQjJOGMj/view?usp=sharing)

La versión [v2.0.0-alpha](https://github.com/jumanor/refirmainvoker/tree/v2.0.0-alpha) es el último lanzamiento

Para mayor información de esta implementación la puede ver en el siguiente [video](https://www.youtube.com/watch?v=aOto5CStZNA)

***Importante***: Este software esta aún en desarrollo, no use en entornos de produción

# Instalación
Para ejecutar Refirma Invoker es necesario que se contacte con la [RENIEC](https://dsp.reniec.gob.pe/refirma_suite/main/web/main.jsf) para que le brinde los identificadores correspondientes

Esta implementación usa [7-zip](https://www.7-zip.org/) que normalmente ya viene instalada en **LINUX**; sin embargo, en **WINDOWS** tendra que instalar manualmente y verificar que se puede acceder desde el terminal.

*ejemplo:*
    
    C:\Users\Jumanor>7z i

Se compilo Refirma Invoker para Windows y Linux, y estan disponibles en los [releases](https://github.com/jumanor/refirmainvoker/releases/tag/v2.0.0-alpha), tambien puede descargar los ejecutables [main.exe](https://github.com/jumanor/refirmainvoker/releases/download/v2.0.0-alpha/main.exe) y [main](https://github.com/jumanor/refirmainvoker/releases/download/v2.0.0-alpha/main) siguiendo los enlaces correspondientes.

**Windows**

    main.exe [clientId] [clientSecret] [serverAddress]

*ejemplo:*

    main.exe K57845459hkj TYUOPDLDFDG 192.168.1.10:9091

**Linux**

    ./main [clientId] [clientSecret] [serverAddress]

*ejemplo:*

    ./main K57845459hkj TYUOPDLDFDG 192.168.1.10:9091

Tambien puede crear el archivo **config.properties** en la misma ubicación donde se encuentra el ejecutable y ejecutar directamente sin parametros.

    clientId=K57845459hkj
    clientSecret=TYUOPDLDFDG
    serverAddress=192.168.1.10:9091

# Funcionamiento
Copia la carpeta [example](https://github.com/jumanor/refirmainvoker/tree/master/example) de este repositorio en un Servidor Web y ejecuta **test.html**.

A pesar que todos los archivos del ejemplo son **estaticos** es necesario usar un Servidor Web debido a que se esta usando **ES6**.

En caso use **Visual Studio Code** instale el plugin [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) que habilita un Servidor Web.

En el ejemplo se usa **Microsoft Click Once**; por lo tanto, si esta usando navegador **Chome** o **Firefox** instales los siguientes plugins:

- Chrome instale este [plugin](https://chrome.google.com/webstore/detail/clickonce-for-google-chro/kekahkplibinaibelipdcikofmedafmb) 

- Firefox instale este [plugin](https://addons.mozilla.org/es/firefox/addon/meta4clickoncelauncher/?utm_source=addons.mozilla.org&utm_medium=referral&utm_content=search)

En caso use el navegador **Edge** no es necesario instalar nada adicional.

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
let url_base=await firma.ejecutar(pdfs,firmaParam);

//En este caso obtenemos los documentos firmados digitalmente y los enviamos a un frame
document.getElementById("frame1").src=url_base+"/"+encodeURI("doc1")
document.getElementById("frame2").src=url_base+"/"+encodeURI("doc2")
```          
# Contribución

Por favor contribuya usando [Github Flow](https://guides.github.com/introduction/flow/). Crea un fork, agrega los commits, y luego abre un [pull request](https://github.com/fraction/readme-boilerplate/compare/).

# License
Copyright © 2022 [Jorge Cotrado](https://github.com/jumanor). <br />
This project is [MIT](https://github.com/jumanor/refirmainvoker/blob/master/License) licensed.

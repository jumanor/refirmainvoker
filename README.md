# Refirma Invoker - Reniec
Implementación del Motor de Firma Digital - Refirma Invoker - del [RENIEC](https://dsp.reniec.gob.pe/refirma_suite/main/web/main.jsf)

***Importante***: Este software esta aún en desarrollo, no use en entornos de produción

# Instalación
Para ejecutar Refirma Invoker es necesario que se contacte con la **RENIEC** para que le brinde los identificadores correspondientes

Esta implementación usa [7-zip](https://www.7-zip.org/) que normalmente ya viene instalada en **LINUX**; sin embargo, en **WINDOWS** tendra que instalar manualmente y verificar que se puede acceder desde el terminal.

*ejemplo:*
    
    C:\Users\Jumanor>7z i

Se compilo Refirma Invoker para Windows y Linux, y estan disponibles en los [releases](), tambien puede descargar los ejecutables [main.exe]() y [main]() siguiendo los enlaces correspondientes.

**Windows**

    main.exe [clientId] [clientSecret] [ip]

*ejemplo:*

    main.exe K57845459hkj TYUOPDLDFDG 192.168.1.10:9091

**Linux**

    ./main [clientId] [clientSecret] [ip]

*ejemplo:*

    ./main K57845459hkj TYUOPDLDFDG 192.168.1.10:9091

# Funcionamiento
Copia la carpeta [example]() de este repositorio en un Servidor Web y ejecuta **test.html**.

A pesar que todos los archivos del ejemplo son **estaticos** es necesario usar un Servidor Web debido a que se esta usando **ES6**.

En caso use **Visual Studio Code** instale el plugin [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) que habilita un Servidor Web.
``` javascript
//Listamos los documentos que se desean firmar digitalmente
let pdfs=[];
pdfs[0]={url:"http://miservidor.com/docs1.pdf",name:"doc1"};
pdfs[1]={url:"http://miservidor.com/docs2.pdf",name:"doc2"};

//Enviamos la posicion en donde se ubicara la representación gráfica de la firma digital
let firmaParam={};
firmaParam.posx=10;
firmaParam.posy=12;

//Llamamos a Refirma Invoker con la dirección ip en donde se ejecuta main.exe o main
let firma=new RefirmaInvoker("192.168.1.10:9091");
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

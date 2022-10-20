export class RefirmaInvoker{

    ARGUMENTOS=null;
    URL_BASE=null;

    EVENT_SUCCESS_INVOKER=null;
    EVENT_ERROR_INVOKER=null;
    //Direccion de Servidor Refirma Invoker
    URL_SERVER_REFIRMA_INVOKER=null;
    
    constructor(url_server_refirma_invoker){
        this.URL_SERVER_REFIRMA_INVOKER=url_server_refirma_invoker;
    }
    event_getArguments(){
        dispatchEventClient('sendArguments', this.ARGUMENTOS);
    }
    event_invokerOk(){
        this.EVENT_SUCCESS_INVOKER(this.URL_BASE);
    }
    event_invokerCancel(e){
        console.log(e)
        this.EVENT_ERROR_INVOKER("FIRMA CANCELADA")
    }
    precarga(){
            
        this.event_getArguments = this.event_getArguments.bind(this);
        this.event_invokerOk = this.event_invokerOk.bind(this);
        this.event_invokerCancel = this.event_invokerCancel.bind(this);

        window.removeEventListener("getArguments", this.event_getArguments);
        window.removeEventListener("invokerOk",this.event_invokerOk);
        window.removeEventListener("invokerCancel", this.event_invokerCancel);

        window.addEventListener('getArguments', this.event_getArguments);
        window.addEventListener('invokerOk', this.event_invokerOk);
        window.addEventListener('invokerCancel',this.event_invokerCancel);
    }////////////////////////////////////////////////////////////////////////////////////
    // Solamente el SGD debe usar /autenticacion.
    // No utilize este metodo en producción
    async autenticacion(usuarioAccesoApi){
        let response=await fetch(this.URL_SERVER_REFIRMA_INVOKER+"/autenticacion",{
            method:'POST',
            body:JSON.stringify({usuarioAccesoApi:usuarioAccesoApi}),
            headers: {
                        'Content-Type': 'application/json; charset=UTF-8',      
                    },
        });
        if(!response.ok){//200-299
            console.log(response.statusText)
            let tt=await response.text()
            throw Error(tt);
        }
        
        let result=await response.json();
        return result.data
    }////////////////////////////////////////////////////////////////////////////////////
    async ejecutar(urlPdfs,parametros,token){
            this.precarga();
        
            let params={};
            //params.pdfs=[];
            //params.pdfs[0]={url:"http://127.0.0.1:5500/01.pdf",name:"doc0"};
            //params.pdfs[1]={url:"http://127.0.0.1:5500/02.pdf",name:"doc1"};
            params.pdfs=urlPdfs
            //params.firma={};
            //params.firma.posx=10;
            //params.firma.posy=12;
            params.firma=parametros;
             
                let response=await fetch(this.URL_SERVER_REFIRMA_INVOKER+"/argumentsServletPCX",{
                    method:'POST',
                    body:JSON.stringify(params),
                    headers: {
                                'Content-Type': 'application/json; charset=UTF-8', 
                                'x-access-token' : token    
                            },
                });

                if(!response.ok){//200-299
                    console.log(response.statusText)
                    let tt=await response.text()
                    throw Error(tt);
                }
                
                let result=await response.json();
                console.log(result);

                this.ARGUMENTOS=result.data.argumentosBase64;
                this.URL_BASE=result.data.urlBase;

                initInvoker('W');//REFIRMA
           
            return new Promise((resolve,reject)=>{

                this.EVENT_SUCCESS_INVOKER=resolve;
                this.EVENT_ERROR_INVOKER=reject;

            });

    }//////////////////////////////////////////////////////////////////////////////
}
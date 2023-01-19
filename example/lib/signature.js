export class Signature{
    //////////////////////////////////////////////////////////////////////PUBLIC///////////////////////////////////////////////////////////
    signaturesByPage(pdfTexto,pagNumber){
        
        if(this._formatoCorrecto(pdfTexto)===false){
            return [];
        }
      
        if(pdfTexto.includes('/T (Signature')){ 

            let mapSignatures = new Map()

            let arrayDeCadenas = pdfTexto.split('/T (Signature');
            arrayDeCadenas.shift();                    
            arrayDeCadenas.forEach(element => {
                
                let indexIni=element.indexOf("/P")
                let indexFin=element.indexOf("0 R",indexIni)

                let pageObj=element.slice(indexIni,indexFin).split(" ")[1]
                let rect = element.slice(element.indexOf("[")+1,element.indexOf("]")).split(" ");

                if(rect[2]<=1 || rect[3]<=1){//firma sin representacion grafica [0.0 0.0 0.0 0.0]
                    return;
                }
                
                if(mapSignatures.size==0){
                    mapSignatures.set(pageObj,[rect])
                }
                else{
                    let arr=mapSignatures.get(pageObj)
                    if(arr==undefined){
                        mapSignatures.set(pageObj,[rect])
                    }
                    else{
                        arr.push(rect)
                    }
                }

            });

            //console.log(mapSignatures)

            let arregloDeSignatures=[];
            let pag=this._getPagesObj(pdfTexto);//los page object ya estan ordenados
            for (let i=0;i<pag.length;i++){
                if(pagNumber==i+1){

                    let pageObj=pag[i];

                    let map=mapSignatures.get(pageObj);
                    if(map==undefined){

                        arregloDeSignatures=[];
                    }
                    else{

                        arregloDeSignatures=map;
                    }

                }
            }
            //console.log(arregloDeSignatures)
            return arregloDeSignatures
        }
        else{   //El documento no tienes firmas digitales

            return []
        }
    }
    //////////////////////////////////////////////////////////////////////PRIVATE///////////////////////////////////////////////////////////
    _formatoCorrecto(pdfTexto){
        let kid = pdfTexto.split('/Kids');//por ahora solo funciona con un Kids
        if(kid.length==2)
            return true
        else 
            return false
    }
    _getPagesObj(pdfTexto){

        let arrayTemp=pdfTexto.split('/Kids')
        arrayTemp=arrayTemp[1]
        arrayTemp.tr
   
        let rect = arrayTemp.slice(arrayTemp.indexOf("[")+1,arrayTemp.indexOf("]")).trim().split(" ");
        let pages=[];
        
        for(let i=0;i<rect.length;i++){
            if(i%3==0){
                pages.push(rect[i])
            }
        }
        return pages;//el primer elemento del array es la primera pagina
    }
}
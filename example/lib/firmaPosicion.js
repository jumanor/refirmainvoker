export class FirmaPosicion{

    CANVAS=null
    IS_DOWN=false;

    LAST_POINTER_X;
    LAST_POINTER_Y;

    LAST_RECTANGLE_X;
    LAST_RECTANGLE_Y;

    WIDTH_RECTANGLE=168;
    HEIGHT_RECTANGLE=62;

    COLOR;
    SCALE=1;

    SIGNATURES=[];
    IMAGEN_INICIAL_SOBRE_SIGNATURE=false;

     ///////////////////////////////////////PUBLIC///////////////////////////////////////////////////////////
    constructor(canvas){
        
        this.CANVAS=canvas

        this.mouseup = this.mouseup.bind(this);
        this.CANVAS.addEventListener('mouseup', this.mouseup);

        this.mousedown = this.mousedown.bind(this);
        this.CANVAS.addEventListener('mousedown',this.mousedown);

        this.mousemove = this.mousemove.bind(this);
        this.CANVAS.addEventListener('mousemove', this.mousemove);

        this.mouseover = this.mouseover.bind(this);
        this.CANVAS.addEventListener('mouseover', this.mouseover);
        
    }
    //opcional
    setScale(scale){
        this.SCALE=scale;
    }
    sizeCanvas(width,height){

        this.CANVAS.width=width;
        this.CANVAS.height=height
    }
    dibujarInicio(x,y,color){

        x=x*this.SCALE
        y=y*this.SCALE

        if(this.restrictSignature(x,y,this.getWithRectangleScale(),this.getHeightRectangleScale())){
            //Por ahora una excepcion en caso la imagen inicial este sobre un signature
            this.IMAGEN_INICIAL_SOBRE_SIGNATURE=true;
        }

        this.LAST_RECTANGLE_X=x;
        this.LAST_RECTANGLE_Y=y;
        this.COLOR=color;
        this.dibujar(x,y,color);
    }
    getX(){

        return this.LAST_RECTANGLE_X/this.SCALE;
    }
    getY(){

        return this.LAST_RECTANGLE_Y/this.SCALE;
    }
    dibujarSignatures(canvasContentPDF,signatures,color){

        signatures.forEach(signRect =>{

               //  Rect [10.0 718.5 178.0 780.0]
               
               let alto=(signRect[3]-signRect[1])*this.SCALE
               let ancho=(signRect[2]-signRect[0])*this.SCALE
               let x=signRect[0]*this.SCALE
               let y=(this.CANVAS.height-signRect[1]*this.SCALE)-alto
               
               this.dibujarRect(canvasContentPDF,x,y,ancho,alto,color) 
               
               this.SIGNATURES.push([x,y,ancho,alto])//valors escalados
        })
    }
    ///////////////////////////////////////PRIVATE///////////////////////////////////////////////////////////
    getWithRectangleScale(){

        var width=this.WIDTH_RECTANGLE*this.SCALE
        //if (width % 1 !== 0)
        //    throw new Error("El ancho del rectangulo tiene que ser entero")

        return width
    }
    getHeightRectangleScale(){

        var height=this.HEIGHT_RECTANGLE*this.SCALE
        //if (height % 1 !== 0)
        //    throw new Error("El ancho del rectangulo tiene que ser entero")
                
        return height
    }
    dibujarRect(canvas,x,y,ancho,alto,color){
        
        var ctx = canvas.getContext("2d");
        ctx.fillStyle = color
        ctx.fillRect(x, y, ancho, alto);
    }
    dibujar(x,y,color){
        var ctx = this.CANVAS.getContext("2d");
        ctx.fillStyle = color
        ctx.fillRect(x, y, this.getWithRectangleScale(), this.getHeightRectangleScale());
    }
    clear(x,y){
        var ctx = this.CANVAS.getContext("2d");
        //un rectangulo mas grande por el redondeo a entero en caso se escale
        ctx.clearRect(x-2,y-2, this.getWithRectangleScale()+4, this.getHeightRectangleScale()+4)
    }
    restrictSignature(xR,yR,anchoR,altoR){

        for (let i=0;i<this.SIGNATURES.length;i++){

            let signRect=this.SIGNATURES[i];

            let alto=signRect[3]
            let ancho=signRect[2]
            let x=signRect[0]
            let y=signRect[1]

            if( xR<=(x+ancho)  &&  xR>=x  &&  yR<=(y+alto) && yR>=y){
                return true;
            }
            if( (xR+anchoR)<=(x+ancho)  &&  (xR+anchoR)>=x  &&  yR<=(y+alto) && yR>=y){
                return true;
            }
            if( (xR+anchoR)<=(x+ancho)  &&  (xR+anchoR)>=x  &&  (yR+altoR)<=(y+alto) && (yR+altoR)>=y){
                return true;
            }
            if( xR<=(x+ancho)  &&  xR>=x  &&  (yR+altoR)<=(y+alto) && (yR+altoR)>=y){
                return true;
            }


            if( x<=(xR+anchoR)  &&  x>=xR  &&  y<=(yR+altoR) && y>=yR){
                return true;
            }
            if( (x+ancho)<=(xR+anchoR)  &&  (x+ancho)>=xR  &&  y<=(yR+altoR) && y>=yR){
                return true;
            }
            if( (x+ancho)<=(xR+anchoR)  &&  (x+ancho)>=xR  &&  (y+alto)<=(yR+altoR) && (y+alto)>=yR){
                return true;
            }
            if( x<=(xR+anchoR)  &&  x>=xR  &&  (y+alto)<=(yR+altoR) && (y+alto)>=yR){
                return true;
            }

        };

        return false;
    }
    mouseover(e){
        this.CANVAS.style.cursor="grab"  
    }
    mouseup(e){

        e.preventDefault();
        e.stopPropagation();
    
        this.IS_DOWN = false;
        this.CANVAS.style.cursor="grab"
    }
    mousedown(e){
        e.preventDefault();
        e.stopPropagation();
        
        this.LAST_POINTER_X=e.offsetX;
        this.LAST_POINTER_Y=e.offsetY;
        
        this.IS_DOWN = true;
        this.CANVAS.style.cursor="move"
        
    }
    mousemove(e){

        if (!this.IS_DOWN) {
            return;
        }
        
       
        // tell the browser we'll handle this event
        e.preventDefault();
        e.stopPropagation();
    
        var dx = e.offsetX - this.LAST_POINTER_X;
        var dy = e.offsetY - this.LAST_POINTER_Y;

        if(this.restrictSignature(this.LAST_RECTANGLE_X+dx, this.LAST_RECTANGLE_Y+dy,this.getWithRectangleScale(),this.getHeightRectangleScale())){
            if(this.IMAGEN_INICIAL_SOBRE_SIGNATURE===false){
                return
            }
        }
        else{//Sin restricciones
            if(this.IMAGEN_INICIAL_SOBRE_SIGNATURE===true){//Solo si la imagen de inicio esta restriginda
                this.IMAGEN_INICIAL_SOBRE_SIGNATURE=false
            }
        }
            
        
        var _canvas=this.CANVAS; 
       
        if(this.LAST_RECTANGLE_X+dx>0 && this.LAST_RECTANGLE_Y+dy>0 && this.LAST_RECTANGLE_X+dx < _canvas.width-this.getWithRectangleScale() && this.LAST_RECTANGLE_Y+dy < _canvas.height-this.getHeightRectangleScale()){
           
            this.clear(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y);
    
            this.LAST_POINTER_X=e.offsetX;
            this.LAST_POINTER_Y=e.offsetY;
    
            this.dibujar(this.LAST_RECTANGLE_X+dx, this.LAST_RECTANGLE_Y+dy,this.COLOR)
    
            this.LAST_RECTANGLE_X += dx;
            this.LAST_RECTANGLE_Y += dy;
    
           
        }
        else{
    
            if(this.LAST_RECTANGLE_X+dx<=0 && (this.LAST_RECTANGLE_Y+dy>0 && this.LAST_RECTANGLE_Y+dy < _canvas.height-this.getHeightRectangleScale())){
            
                this.clear(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y);
    
    
                //this.LAST_POINTER_X=e.offsetX;
                this.LAST_POINTER_Y=e.offsetY;
    
                this.dibujar(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y+dy,this.COLOR)
    
        
                //this.LAST_RECTANGLE_X += dx;
                this.LAST_RECTANGLE_Y += dy;
               
    
            }
            else if(this.LAST_RECTANGLE_Y+dy<=0 && (this.LAST_RECTANGLE_X+dx>0 && this.LAST_RECTANGLE_X+dx<_canvas.width-this.getWithRectangleScale())){
                
                
                this.clear(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y);
    
    
                this.LAST_POINTER_X=e.offsetX;
                // this.LAST_POINTER_Y=e.offsetY;
    
                this.dibujar(this.LAST_RECTANGLE_X+dx, this.LAST_RECTANGLE_Y,this.COLOR)
    
                
                this.LAST_RECTANGLE_X += dx;
                // this.LAST_RECTANGLE_Y += dy;
                
            }
            
            else if(this.LAST_RECTANGLE_X+dx>=_canvas.width-this.getWithRectangleScale() && (this.LAST_RECTANGLE_Y+dy>0 && this.LAST_RECTANGLE_Y+dy < _canvas.height-this.getHeightRectangleScale())){
            
                this.clear(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y);
    
    
                //his.LAST_POINTER_X=e.offsetX;
                this.LAST_POINTER_Y=e.offsetY;
    
                this.dibujar(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y+dy,this.COLOR)
    

                //this.LAST_RECTANGLE_X += dx;
                this.LAST_RECTANGLE_Y += dy;
            
    
            }
            else if(this.LAST_RECTANGLE_Y+dy>=_canvas.height-this.getHeightRectangleScale() && (this.LAST_RECTANGLE_X+dx>0 && this.LAST_RECTANGLE_X+dx<_canvas.width-this.getWithRectangleScale())){
                
                this.clear(this.LAST_RECTANGLE_X, this.LAST_RECTANGLE_Y);
    
    
                this.LAST_POINTER_X=e.offsetX;
                //this.LAST_POINTER_Y=e.offsetY;
                
                this.dibujar(this.LAST_RECTANGLE_X+dx, this.LAST_RECTANGLE_Y,this.COLOR)
    
                this.LAST_RECTANGLE_X += dx;
                //this.LAST_RECTANGLE_Y += dy;
            
            }
            else{//Zona restringida
    
            
            }
        }
    }
}
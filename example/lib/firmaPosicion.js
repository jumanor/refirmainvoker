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
    ///////////////////////////////////////PRIVATE///////////////////////////////////////////////////////////
    getWithRectangleScale(){

        var width=this.WIDTH_RECTANGLE*this.SCALE
        if (width % 1 !== 0)
            throw new Error("El ancho del rectangulo tiene que ser entero")

        return width
    }
    getHeightRectangleScale(){

        var height=this.HEIGHT_RECTANGLE*this.SCALE
        if (height % 1 !== 0)
            throw new Error("El ancho del rectangulo tiene que ser entero")
                
        return height
    }
    dibujar(x,y,color){
        var ctx = this.CANVAS.getContext("2d");
        ctx.fillStyle = color
        ctx.fillRect(x, y, this.getWithRectangleScale(), this.getHeightRectangleScale());
    }
    clear(x,y){
        var ctx = this.CANVAS.getContext("2d");
        ctx.clearRect(x,y, this.getWithRectangleScale(), this.getHeightRectangleScale())
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
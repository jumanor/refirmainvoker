export class Dialog{

    static prompt(mensaje){

        return new Promise((resolve, reject) => {
    
            window.bootbox.prompt(mensaje, function(result){
                resolve(result);
            });
              
        });
    }////////////////////////////////////////////////////////
    static alert(mensaje){

        return new Promise((resolve, reject) => {
    
            window.bootbox.alert(mensaje, function(){
                resolve("");
            });
              
        });
    }////////////////////////////////////////////////////////    
}

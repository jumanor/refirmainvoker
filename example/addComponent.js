//Refirma Invoker necesita este div oara funcionar
let componente=document.getElementById("addComponent")
if(componente==null){

    let elemDiv = document.createElement('div');
    elemDiv.setAttribute("id", "addComponent");
    document.body.appendChild(elemDiv);
}
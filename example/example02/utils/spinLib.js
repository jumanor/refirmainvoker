let modal=document.getElementById("myModalSpin")
if(modal==null){

    let elemDiv = document.createElement('div');
    elemDiv.setAttribute("id", "myModalSpin");
    elemDiv.setAttribute("class", "modalSpin");

    let elemDiv1 = document.createElement('div');

    let elemDiv2 = document.createElement('div');
    elemDiv2.setAttribute("class", "spinner");
    
    document.body.appendChild(elemDiv);
    elemDiv.appendChild(elemDiv1);
    elemDiv1.appendChild(elemDiv2);
    
    /*
    <div id="myModalSpin" class="modalSpin">
        <div>
            <div class="spinner"></div>
        </div>
    </div>
    */
}   
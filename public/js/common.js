// Ugis Germanis
// My javascript

var music = {
    selected: "byartist",
    sort: new Event('sort'),
};

document.addEventListener("sort", function(e){
    console.log("sort event: ", music.selected);
}, false)

window.addEventListener('WebComponentsReady', function(e) {
    console.log('webcomponents are ready!!!');
    //Init the default
    document.getElementById(music.selected).active=true;
    document.dispatchEvent(music.sort);
    //Nifty button stuff
    var sorter = document.getElementsByClassName("sorter");
    for (i=0; i< sorter.length; i++) {
        sorter[i].addEventListener("change", function(e){
            music.selected = e.path[0].id
            for(j=0; j<sorter.length; j++){
                sorter[j].active=false;
            }
            document.getElementById(music.selected).active=true;
            document.dispatchEvent(music.sort);
        }, false);
    }
});
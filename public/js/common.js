// Ugis Germanis
// My javascript

var music = {
    selected: "byartist",
    sort: new Event('sort'),
};

document.addEventListener("sort", function(e){
    console.log("Raw: ", music.selected);
}, false)

window.addEventListener('WebComponentsReady', function(e) {
    console.log('components ready');
    byartist.addEventListener("change", function(e){
        if (byartist.active == true){
            bytitle.active=false;
            byalbum.active=false;
            music.selected="byartist";
            document.dispatchEvent(music.sort);
        }
    }, false);
    bytitle.addEventListener("change", function(e){
        if (bytitle.active == true) {
            byartist.active=false;
            byalbum.active=false;
            music.selected="bytitle";
            document.dispatchEvent(music.sort);
        }
    }, false);
    byalbum.addEventListener("change", function(e){
        if (byalbum.active == true) {
            byartist.active=false;
            bytitle.active=false;
            music.selected="byalbum";
            document.dispatchEvent(music.sort);    
        }
    }, false);
});
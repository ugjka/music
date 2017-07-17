// Ugis Germanis
// My javascript

var music = {
    selected: "byartist",
};

var sort = new Event('sort');

document.addEventListener("sort", function(e){
    console.log("Raw: ", music.selected);
}, false)

window.addEventListener('WebComponentsReady', function(e) {
    console.log('components ready');
    artist.addEventListener("change", function(e){
        if (artist.active == true){
            song.active=false;
            album.active=false;
            music.selected="byartist";
            document.dispatchEvent(sort);
        }
    }, false);
    song.addEventListener("change", function(e){
        if (song.active == true) {
            artist.active=false;
            album.active=false;
            music.selected="bysong";
            document.dispatchEvent(sort);
        }
    }, false);
    album.addEventListener("change", function(e){
        if (album.active == true) {
            artist.active=false;
            song.active=false;
            music.selected="byalbum";
            document.dispatchEvent(sort);    
        }
    }, false);
});
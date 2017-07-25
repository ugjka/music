// Ugis Germanis
// My javascript
window.addEventListener('WebComponentsReady', function(e) {
    console.log('webcomponents are ready!!!');
    //Slider events
    music.slider = document.getElementById("slider");
    slider.addEventListener("mousedown", function(){
        music.sliderDrag = true;
    });
    slider.addEventListener("mouseup", function(){
        music.sliderDrag = false;
    });
    slider.addEventListener("change", function(e){
        var estimate = soundManager.getSoundById(["sound", music.current].join("")).durationEstimate;
        soundManager.getSoundById(["sound", music.current].join("")).setPosition((estimate/1000) * e.target.value);
        music.sliderDrag = false;
    });
    document.getElementById("menubutton").addEventListener("click", function(e){
        startDrawer.toggle();
    });
    //Sorter
    document.getElementById("startDrawer").addEventListener("change", function(e){
        if(e.target && e.target.nodeName == "PAPER-BUTTON") {
            var sorters = document.getElementsByClassName("sorter");
            for (i=0; sorters.length > i; i++){
                sorters[i].active = false;
            }
            music.selected = e.target.id.replace("post-", "");
            e.target.active = true;
            document.dispatchEvent(music.sort);
        }
    }, false);
    //Playback controls
    document.getElementById("next").addEventListener("click", function(e){
        playnext();
    });
    document.getElementById("previous").addEventListener("click", function(e){
        playprevious();
    });
    document.getElementById("playit").addEventListener("click", function(e){
        play();
    });
    document.getElementById("pauseit").addEventListener("click", function(e){
        pause();
    });
    //Sort event handler
    document.addEventListener("sort", function(e){
        music.request.open("GET", [music.url, "/api?sort=", music.selected].join(""));
        music.request.send();
    }, false);
    //Default sort
    document.getElementById(music.selected).active=true;
    document.dispatchEvent(music.sort);
});
//Main object
var music = {
    selected: "byartist",
    sort: new Event('sort'),
    playlist: null,
    request: new XMLHttpRequest(),
    url: [document.location.protocol, "//", window.location.host].join(""),
    current: 0,
    slider: null,
    sliderDrag: false,
    render: function(){
        $("#playlist").empty();
        var playlist = [""];
        for (i=0; i<this.playlist.length; i++) {
            playlist.push("<li id='sound", i, "' class='song' index='",i ,"'>", this.playlist[i].Title, " - ", this.playlist[i].Artist, "</li>");
        }
        $("#playlist").append(playlist.join(""));
        document.getElementById("playlist").addEventListener('click', function(e){
            if(e.target && e.target.nodeName == "LI") {
                music.current= e.target.getAttribute("index");
                playSong(music.current);
            }
        });
    },
};

music.request.onloadend = function(){
    if (this.readyState = 4 && this.status == 200) {
        music.playlist = JSON.parse(this.responseText);
        $("#playlist").attr('loading', true);
        music.render();
        $("#playlist").attr('loading', false);
    }
}
//Play song by id
function playSong(id){
    soundManager.stopAll();
    soundManager.createSound({
        id: ["sound", id].join(""),
        url: [music.url, "/stream?id=", music.playlist[id].ID].join(""),
        multiShot:false,
        onfinish: function() {
            if (music.current == (music.playlist.length - 1)) {
                music.current = 0;
            } else {
                music.current++;
            }
            playSong(music.current);
            $(["#sound", id].join("")).attr("playing", false);
            this.destruct();
        },
        onstop: function(){
            $(["#sound", id].join("")).attr("playing", false);
            this.destruct();
        },
        whileplaying: function(){
            if (music.sliderDrag == false){
                music.slider.value = (this.position / this.durationEstimate) * 1000;
            }
        },
    });
    $(["#sound", id].join("")).attr("playing", true);
    soundManager.play(["sound", id].join(""));
    document.getElementsByTagName("body")[0].style.backgroundImage = ["url(", music.url, "/art?id=", music.playlist[id].ID, ")"].join("");
}
//Playback functions
function playnext(){
    if (music.current == (music.playlist.length - 1)) {
        music.current = 0;
    } else {
        music.current++;
    }
    playSong(music.current);
}

function playprevious(){
    if (music.current == 0) {
        music.current = music.playlist.length - 1
    } else {
        music.current--;
    }
    playSong(music.current);
}

function play(){
    if (soundManager.getSoundById(["sound", music.current].join("")) === undefined){
        playSong(music.current);
    } else {
        soundManager.getSoundById(["sound", music.current].join("")).resume();
    }
}

function pause(){
    soundManager.getSoundById(["sound", music.current].join("")).pause();
}
// Ugis Germanis
// My javascript

var music = {
    selected: "byartist",
    sort: new Event('sort'),
    playlist: null,
    request: new XMLHttpRequest(),
    url: document.location.protocol+ "//"+ window.location.host,
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

document.addEventListener("sort", function(e){
    console.log("sort event: ", music.selected);
    music.request.open("GET", music.url + "/api?sort=" + music.selected);
    music.request.send();
}, false)

window.addEventListener('WebComponentsReady', function(e) {
    console.log('webcomponents are ready!!!');
    //Init the default
    music.slider = document.getElementById("slider");
    slider.addEventListener("mousedown", function(){
        music.sliderDrag = true;
    });
    slider.addEventListener("mouseup", function(){
        music.sliderDrag = false;
    });
    slider.addEventListener("change", function(e){
        var estimate = soundManager.getSoundById("sound"+music.current).durationEstimate;
        soundManager.getSoundById("sound"+music.current).setPosition((estimate/1000) * e.target.value);
        music.sliderDrag = false;
    });
    document.getElementById(music.selected).active=true;
    document.dispatchEvent(music.sort);
    //Nifty button stuff
    var sorter = document.getElementsByClassName("sorter");
    for (i=0; i< sorter.length; i++) {
        sorter[i].addEventListener("change", function(e){
            for(j=0; j<sorter.length; j++){
                if (sorter[j].id == e.target.id) {
                    continue
                }
                sorter[j].active=false;
            }
            e.target.active=true;
            if (music.selected == e.target.id){
                return
            }
            music.selected = e.target.id
            document.dispatchEvent(music.sort);
        }, false);
    }
});

function playSong(id){
    soundManager.stopAll();
    soundManager.createSound({
        id: "sound"+id,
        url: music.url+ "/stream?id="+ music.playlist[id].ID,
        multiShot:false,
        onfinish: function() {
            if (music.current == (music.playlist.length - 1)) {
                music.current = 0;
            } else {
                music.current++;
            }
            playSong(music.current);
            $("#sound"+id).attr("playing", false);
            this.destruct();
        },
        onstop: function(){
            $("#sound"+id).attr("playing", false);
            this.destruct();
        },
        whileplaying: function(){
            if (music.sliderDrag == false){
                music.slider.value = (this.position / this.durationEstimate) * 1000;
            }
        },
    });
    $("#sound"+id).attr("playing", true);
    soundManager.play("sound"+id);
    document.getElementsByTagName("body")[0].style.backgroundImage = "url("+music.url+"/art?id="+music.playlist[id].ID+")";
}

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
    if (soundManager.getSoundById("sound"+music.current) === undefined){
        playSong(music.current);
    } else {
        soundManager.getSoundById("sound"+music.current).resume();
    }
}

function pause(){
    soundManager.getSoundById("sound"+music.current).pause();
}
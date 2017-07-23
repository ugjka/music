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
        for (i=0; i<this.playlist.length; i++) {
            $("#playlist").append("<li class='song'></li>");
            $(".song").last().attr('id', i);
            $(".song").last().append(this.playlist[i].Title+ " - "+ this.playlist[i].Artist);
            document.getElementById(i).addEventListener('click', function(e){
                music.current= e.target.id;
                playSong(music.current);
            });
        }
    },
};

music.request.onloadend = function(){
    if (this.readyState = 4 && this.status == 200) {
        music.playlist = JSON.parse(this.responseText);
        music.render();
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
            $("#"+id).attr("playing", false);
            this.destruct();
        },
        onstop: function(){
            $("#"+id).attr("playing", false);
            this.destruct();
        },
        whileplaying: function(){
            if (music.sliderDrag == false){
                music.slider.value = (this.position / this.durationEstimate) * 1000;
            }
        },
    });
    $("#"+id).attr("playing", true);
    soundManager.play("sound"+id);
}

function playnext(){
    music.current++;
    playSong(music.current);
}

function playprevious(){
    music.current--;
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
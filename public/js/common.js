//ugjka <esesmu@protonmail.com>
// My javascript

/*
                    .-'-'`-.               .-'`-`-.
                 .-'-       \             /       -`-.
                /    -.    - \           / -    .-    \
               /_/___ __.'-\  \        /  /-`.__ ___\_\
              /./ \_ \___ / \  |      |  / \ ___/ _/ \.\
              \/      /_ \  /  |      |  \  / _\      \/
               |  __\/   /./   |      |   \.\   \/__  |
               \  \_\  .'      |      |      `.  /_/  /
             _/ \____.'       /        \       `.____/ \_
           /'        \      /'|        |`\      /        `\
          |           \    /  |        |  \    /           |
          \/           \  \   |        |   /  /           \/
         /              `\ \   \      /   / /'              \
        /                 | \  |      |  / |                 \
       /   /              \    |      |    /              \   \
    .-'  ,'       \        \   \      /   /        /       `,  `-.
   |'  ./          `\       \   |    |   /       /'          \.  `|
   \   |'            \      `   |    |   '      /            `|   /
    `-/\     ./      |`\     \  |    |  /     /'|      \.     /\-'
      ||`---'        |  \     . |    | .     /  |        `---'||
      \\             /   \    | \    / |    /   \             //
       ||           /\_   \    \|    |/    /   _/\           ||
       \ \         /|  `\_|\    '    `    /|_/'  |\         / /
        | `\        |      |    \    /    |      |        /' |
        |   \       |      |     \  /     |      |       /   |
        |'  |        \      \    |  |    /      /        |  `|
        \   |         `-.    \   |  |   /    .-'         |   /
  -----._\__._           \   |   |  |   |   /           _.__/_.-----
              `-.         |  |   |  |   |  |         .-'
                          |  \   |  |   /  |
                          |  _|  |  |  |_  |
                          / /||..|  |..||\ \
       :F_P:            .' |/||||/  \||||\| `.            :F_P:
      ---._________.---'   ` \|''    ``|/ '   `---._________.---
*/

// Set up stuff when webcomponents are ready
window.addEventListener('WebComponentsReady', function() {
  console.log('webcomponents are ready!!!');
  //Slider events
  music.slider = document.getElementById("slider");
  slider.addEventListener("mousedown", function() {
    music.sliderDrag = true;
  });
  slider.addEventListener("mouseup", function() {
    music.sliderDrag = false;
  });
  slider.addEventListener("change", function(e) {
    var estimate = music.mainSound.durationEstimate;
    music.mainSound.setPosition((estimate / 1000) * e.target.value);
    music.sliderDrag = false;
  });
  document.getElementById("menubutton").addEventListener("click", function() {
    startDrawer.toggle();
  });
  //Sorter
  document.getElementById("startDrawer").addEventListener("change", function(e) {
    if (e.target && e.target.nodeName == "PAPER-BUTTON") {
      var sorters = document.getElementsByClassName("sorter");
      for (i = 0; sorters.length > i; i++) {
        sorters[i].active = false;
      }
      music.selected = e.target.id.replace("post-", "");
      e.target.active = true;
      document.dispatchEvent(music.sort);
    }
  }, false);
  //Playback controls
  document.getElementById("next").addEventListener("click", function() {
    playnext();
  });
  document.getElementById("previous").addEventListener("click", function() {
    playprevious();
  });
  document.getElementById("playit").addEventListener("click", function() {
    play();
  });
  document.getElementById("pauseit").addEventListener("click", function() {
    pause();
  });
  document.getElementById("favoriteit").addEventListener("click", function() {
    //Set or Unset favorite
    $.get(
      [music.url, "/like"].join(""),
      { "like": music.playlist[music.current].ID },
      function(data) {
        document.getElementById("favoriteit").setAttribute("favorited", data);
      },
      "json"
    );
  });
  // Focus to playing handler
  document.getElementById("focus").addEventListener("click", function() {
    document.getElementById("sound" + music.current).scrollIntoView({ block: "center" });
  });
  // Logout handler
  document.getElementById("logout").addEventListener("click", function() {
    $.removeCookie('secret');
    location.reload();
  });
  //Sort event handler
  //Get the playlist
  document.addEventListener("sort", function() {
    $.get(
      [music.url, "/api"].join(""),
      { "sort": music.selected },
      function(data) {
        music.playlist = data;
        $("#playlist").attr('loading', true);
        music.render();
        $("#playlist").attr('loading', false);
      },
      "json"
    );
  }, false);
  //Init sound object
  soundManager.setup({
    url: 'bower_components/SoundManager2/swf/soundmanager2_flash9.swf',
    flashVersion: 9,
    preferFlash: false,
    onready: function() {
      soundManager.createSound({
        id: "main",
        url: "",
        multiShot: false,
        onfinish: function() {
          music.previous = music.current;
          $(["#sound", music.current].join("")).attr("playing", false);
          if (music.current == (music.playlist.length - 1)) {
            music.current = 0;
          } else {
            music.current++;
          }
          playSong(music.current);
          //Count the play
          $.get(
            [music.url, "/count"].join(""),
            { "id": music.playlist[music.previous].ID }
          );
        },
        onstop: function() {
          $(["#sound", music.previous].join("")).attr("playing", false);
        },
        //Update the slider on playback
        whileplaying: function() {
          if (music.sliderDrag === false) {
            music.slider.value = (this.position / this.durationEstimate) * 1000;
          }
        },
      });
      music.mainSound = soundManager.getSoundById("main");
    }
  });
  //Default sort
  document.getElementById(music.selected).active = true;
  document.dispatchEvent(music.sort);
  //Set up Media Session API controls
  if ('mediaSession' in navigator) {
    navigator.mediaSession.setActionHandler('nexttrack', function() {
      playnext();
    });
    navigator.mediaSession.setActionHandler('previoustrack', function() {
      playprevious();
    });
    navigator.mediaSession.setActionHandler('play', function() {
      play();
    });
    navigator.mediaSession.setActionHandler('pause', function() {
      pause();
    });
  }
  //Handle login
  $.post([music.url, "/login"].join(""), { password: "", secret: $.cookie("secret") },
    function(data) {
      if (data === false) {
        login.open();
      } else if ($.cookie("secret") == "pp9zzKI6msXItWfcGFp1bpfJghZP4lhZ4NHcwUdcgKYVshI68fX5TBHj6UAsOsVY9QAZnZW20-MBdYWGKB3NJg==") {
        document.getElementById("logout").style.display = "none";
      }
    },
    "json"
  );
  document.getElementById("passwordInput").addEventListener("keydown", function(e) {
    if (e.key === "Enter") {
      loginNow();
    }
  });
  document.getElementById("passwordEnter").addEventListener("click", function() {
    loginNow();
  });

  function loginNow() {
    var pass = passwordInput.value;
    if (pass === undefined) {
      $("#passwordWrong").empty();
      $("#passwordWrong").append("Enter the password");
      return;
    }
    $.post([music.url, "/login"].join(""), { password: pass },
      function(data) {
        if (data === false) {
          login.open();
          $("#passwordWrong").empty();
          $("#passwordWrong").append("Wrong password");
        } else {
          login.close();
        }
      },
      "json"
    );
  }

  document.getElementById("playlist").addEventListener('click', function(e) {
    var closest = e.target.closest("TR");
    if (closest) {
      music.previous = music.current;
      music.current = closest.getAttribute("index");
      playSong(music.current);
    }
  });
});
//Main object
var music = {
  selected: "byartist",
  sort: new Event('sort'),
  playlist: null,
  mainSound: null,
  url: [document.location.protocol, "//", window.location.host].join(""),
  current: 0,
  previous: 0,
  slider: null,
  sliderDrag: false,
  //Render the playlist
  render: function() {
    $("#playlist").empty();
    var playlist = [""];
    for (i = 0; i < this.playlist.length; i++) {
      var albumPrint = "";
      if (this.playlist[i].Track != 0) {
        albumPrint = ["<td class='album'>[", this.playlist[i].Track, "]", " ", this.playlist[i].Album, "</td>"].join("");
      } else {
        albumPrint = "<td class='album'></td>";
      }
      playlist.push(
        "<tr id='sound", i, "' class='song' index='", i, "'>",
        "<td class='title'>", this.playlist[i].Title, "</td>",
        "<td class='artist'>", this.playlist[i].Artist, "</td>",
        albumPrint,
        "</tr>"
      );
    }
    $("#playlist").append(playlist.join(""));
  },
};

//Play song by id
function playSong(id) {
  if (music.mainSound.playState == 1) {
    music.mainSound.stop();
  }
  music.mainSound.unload();
  $(["#sound", id].join("")).attr("playing", true);
  music.mainSound.url = [music.url, "/stream?id=", music.playlist[id].ID].join("");
  music.mainSound.play();
  document.getElementsByTagName("body")[0].style.backgroundImage = ["url(", music.url, "/art?id=", music.playlist[id].ID, ")"].join("");
  //Set up media session data
  if ('mediaSession' in navigator) {
    navigator.mediaSession.metadata = new MediaMetadata({
      title: music.playlist[music.current].Title,
      artist: music.playlist[music.current].Artist,
      album: music.playlist[music.current].Album,
      artwork: [{ src: [music.url, "/art?id=", music.playlist[id].ID].join(""), }],
    });
  }
  document.title = [music.playlist[music.current].Title, "by", music.playlist[music.current].Artist].join(" ");
  //Get if liked or not
  $.get(
    [music.url, "/like"].join(""),
    { "id": music.playlist[id].ID },
    function(data) {
      document.getElementById("favoriteit").setAttribute("favorited", data);
    },
    "json"
  );
}

//Playback functions
function playnext() {
  music.previous = music.current;
  if (music.current >= (music.playlist.length - 1)) {
    music.current = 0;
  } else {
    music.current++;
  }
  playSong(music.current);
}

function playprevious() {
  music.previous = music.current;
  if ((music.current == 0) || (music.current > music.playlist.length - 1)) {
    music.current = music.playlist.length - 1;
  } else {
    music.current--;
  }
  playSong(music.current);
}

function play() {
  if (music.mainSound.paused) {
    music.mainSound.resume();
  } else {
    playSong(music.current);
  }
}

function pause() {
  music.mainSound.pause();
}

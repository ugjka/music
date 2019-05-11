# music

[![Donate](https://dl.ugjka.net/Donate-PayPal-green.svg)](https://www.paypal.me/ugjka)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fugjka%2Fmusic.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fugjka%2Fmusic?ref=badge_shield)

Stream your god damned music hassle free

## 2019.05 Update

* Use html tables for the playlist
* Logout button
* Focus on currently playing track button
* Track number and album title on right side if the screen is big enough
* Some ui improvements
* Coincidentally KDE browser integration plugin works OOB for those who use it

## Note

Use `-enableFlac` flag for flac support (not every browser supports flac file streaming).

Use `-password yourpass` flag to protect the website with password.

## Caveats

Make sure your music is properly tagged, otherwise it will all be wonky.

Build needs Go with Go module support

## Install instructions

```bash
export GO111MODULE=on
git clone https://github.com/ugjka/music.git
cd music
go generate #requires Bower https://bower.io/
go build
./music -path /path/to/mp3/collection -port 8080
```

Navigate in your browser to: `http://127.0.0.1:8080/`

**Desktop view:**

![desktop](https://i.imgur.com/ceqI4LJ.png)

![desktop](https://i.imgur.com/OAIx2Vz.png)

**Mobile lockscreen:**

![mobile](https://img.ugjka.net/XPdyMKUk.png)

**Optional Password Protection:**

![password](https://i.imgur.com/SLsa9EQ.png)

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fugjka%2Fmusic.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fugjka%2Fmusic?ref=badge_large)

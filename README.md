# music

[![Donate](https://dl.ugjka.net/Donate-PayPal-green.svg)](https://www.paypal.me/ugjka)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fugjka%2Fmusic.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fugjka%2Fmusic?ref=badge_shield)

Stream your god damned mp3 (and flac?) music hassle free

## Demo

https://demo.ugjka.net/

## Note

Use `-enableFlac` flag for flac support (not every browser supports flac file streaming).

Use `-password yourpass` flag to protect the website with password.

## Caveats

Make sure your music is properly tagged, otherwise it will all be wonky.

Build needs Go with Go module support and node.js NPM

## Install instructions

```bash
export GO111MODULE=on
git clone https://github.com/ugjka/music.git
cd music
go generate #requires NPM, see https://nodejs.org
go build
./music -path /path/to/mp3/collection -port 8080
```

Navigate in your browser to: `http://127.0.0.1:8080/`

**Desktop view:**

![desktop](https://i.imgur.com/FeleWvb.png)

![desktop](https://i.imgur.com/AAXh2vn.png)

**Mobile lockscreen:**

![mobile](https://img.ugjka.net/XPdyMKUk.png)

**Optional Password Protection:**

![password](https://i.imgur.com/MOFEkvc.png)

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fugjka%2Fmusic.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fugjka%2Fmusic?ref=badge_large)

## The future of this project?

When I get time and inspiration I'll try to work on the features listed in this link: https://github.com/ugjka/music/projects/1

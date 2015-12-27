![Logo](http://rpiai.com/musicsaur/musicsaur1.png)

# MusicSAUR (formerly the un-googleable "μsic")
## Music Synchronization And Uniform Relaying
[![Version 1.3](https://img.shields.io/badge/version-1.3-brightgreen.svg)]()
[![Join the chat at https://gitter.im/schollz/musicsaur](https://badges.gitter.im/schollz/music.svg)](https://gitter.im/schollz/musicsaur?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

---

**Update 12/26/2015: Version 1.3 released! Now includes Go and Python versions!**

**Major features include (for all versions):**

- Go version with [precompiled libraries](http://www.musicsaur.com) for folks that want to unzip and run!
- Mute button
- All files served locally (no internet connection needed!)
- Replaying works now (it was broken in 1.2)
- Playlist sorted by Artist/Album/Track (in order)

---

Want to sync up your music on multiple computers? This accomplishes exactly that - using only a simple Python server and a browser. Simply run the Python server, and open up a browser on each computer you want to sync up - that's it!

This program is powered by [the excellent howler.js library from goldfire](https://github.com/goldfire/howler.js/). Essentially all the client computers [sync their clocks](http://www.mine-control.com/zack/timesync/timesync.html) and then try to start a song at the same time. Any dissimilarities between playback are also fixed, because the clients will automatically seek to the position of the server.

# Setup

If you don't want to install *anything*, just download the [compiled version](http://www.musicsaur.com/download-binary/).

## Python

If you're interested in installing the Python version, follow these instructions. First install the required packages using

```bash
sudo python setup.py install
```

then copy the configuration file for the Python version

```bash
cp config-python.cfg config.cfg
```

and edit line #42 with the locations of your music folderse. There are other parameters to edit, if you feel so inclined, but you needn't just to get started. Then simply run with

```bash
python syncmusic.py
```

## Golang

If you're interested in installing the Golang version, follow these instructions. First install the required packages

```bash
go get github.com/mholt/caddy/caddy
go get github.com/tcolgate/mp3
go get github.com/bobertlo/go-id3/id3
go get github.com/BurntSushi/toml
go get gopkg.in/tylerb/graceful.v1
```

Then copy the configuration file 

```bash
cp config-go.cfg config.cfg
```

and edit line #5 with your music folders. Then simpily use 

```bash
go run *.go
```

to start up the server!

### Some notes

- If you don't hear anything, the client is probably trying to synchronize. The browser automatically mutes when it goes out of sync to avoid the headache caused by mis-aligned audio. You can see synchronization progress in [your browser console](https://webmasters.stackexchange.com/questions/8525/how-to-open-the-javascript-console-in-different-browsers). 
- If you still dont' hear anything, and you're using Chrome browser on Android [you need change one of the flags in chrome to allow audio without gestures](http://android.stackexchange.com/questions/59134/enable-autoplay-html5-video-in-chrome). To do this, copy and paste this into your Chrome browser:

```bash
chrome://flags/#disable-gesture-requirement-for-media-playback
```

- If you want to play music from a Raspberry Pi, just type this command (works on headless):

```bash
xinit /usr/bin/midori -a http://W.X.Y.Z:5000
```

OR

```bash
xinit /usr/bin/luakit -u http://W.X.Y.Z:5000
```


## Screenshot

![Screenshot](http://rpiai.com/musicsaur/screenshot2.png)

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

# History

- 12/26/2015: Version 1.3 Release
- 12/19/2015: Version 1.2 Release
- 12/18/2015 (evening): Version 1.1 Release
- 12/18/2015 (morning): Version 1.0 Release

## Credits

* James Simpson and [Goldfire studios](http://goldfirestudios.com/blog/104/howler.js-Modern-Web-Audio-Javascript-Library) for their amazing [howler.js library](https://github.com/goldfire/howler.js/)
* Zach Simpson for [his paper on simple clock synchronization](http://www.mine-control.com/zack/timesync/timesync.html)
* Everyone on the [/r/raspberry_pi](https://www.reddit.com/r/raspberry_pi/comments/3xc8kq/simple_python_script_to_allow_multiple_raspberry/) and [/r/python](https://www.reddit.com/r/Python/comments/3xc8mj/simple_python_script_to_allow_multiple_computers/) threads for great feature requests and bug reports!
* [ClkerFreeVectorImages](https://pixabay.com/en/users/ClkerFreeVectorImages-3736/) and [OpenClipartVectors](https://pixabay.com/en/users/OpenClipartVectors-30363/) for the Public Domain vectors
- [mholt](github.com/mholt) for the invaluable ```Caddy``` 
- [tcolgate](http://github.com/tcolgate) for the ```mp3``` package
- [bobertlo](http://github.com/bobertlo) for the ```go-id3``` package
- [BurntSushi](http://github.com/BurntSushi) for their ```toml``` library

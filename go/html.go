package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func getTime() (curTime int64) {
	//curTime = time.Now().UnixNano() / 1000000
	curTime = time.Since(time.Date(2015, 6, 1, 12, 0, 0, 0, time.UTC)).Nanoseconds() / 1000000
	return
}

func songControl(millisecondWait int64, is_playing bool, text string, song string, start_next bool) {
	time.Sleep(time.Duration(millisecondWait) * time.Millisecond)
	if song == statevar.CurrentSong {
		log.Printf(song + " " + text)
		statevar.IsPlaying = is_playing
		if start_next == true {
			skipTrack(-1)
		}
	}
}

func getPlaylistHTML() (playlist_html string) {
	playlist_html = ""
	for i, k := range statevar.SongList {
		if statevar.CurrentSong != k {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'>" + k + "</a><br>\n"
		} else {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'><b>" + k + "</b></a><br>\n"

		}
	}
	return
}

func getPlaybackPositionInSeconds() float64 {
	position := float64(getTime()-statevar.SongStartTime) / 1000.0
	if statevar.IsPlaying == true && position > 0 {
		return position
	} else {
		return 0.0
	}
}

func SyncRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//current_song := r.FormValue("current_song")
		client_timestamp, _ := strconv.Atoi(r.FormValue("client_timestamp"))
		data := SyncJSON{
			Current_song:     statevar.CurrentSong,
			Client_timestamp: int64(client_timestamp),
			Server_timestamp: getTime(),
			Is_playing:       statevar.IsPlaying,
			Song_time:        getPlaybackPositionInSeconds(),
			Song_start_time:  statevar.SongStartTime,
		}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(b))
	}
}

func NextSongRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		skip, _ := strconv.Atoi(r.FormValue("skip"))
		skipTrack(skip)
		data := SyncJSON{
			Current_song:     "None",
			Client_timestamp: 0,
			Server_timestamp: 0,
			Is_playing:       false,
			Song_time:        0,
			Song_start_time:  0,
		}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(b))
	}
}

func skipTrack(song_index int) {
	if song_index < 0 {
		statevar.CurrentSongIndex += song_index + 2
	} else {
		statevar.CurrentSongIndex = song_index
	}
	song := statevar.SongList[statevar.CurrentSongIndex]
	rawSongData, _ = ioutil.ReadFile(statevar.SongMap[song].Path)
	statevar.CurrentSong = song
	statevar.SongStartTime = getTime() + 11000
	statevar.IsPlaying = false
	b, _ := json.Marshal(statevar)
	ioutil.WriteFile("state.json", b, 0644)
	go songControl(statevar.SongStartTime-getTime()-3000, false, "3", song, false)
	go songControl(statevar.SongStartTime-getTime()-2000, false, "2", song, false)
	go songControl(statevar.SongStartTime-getTime()-1000, false, "1", song, false)
	go songControl(statevar.SongStartTime-getTime(), true, "Playing "+song, song, false)
	go songControl(statevar.SongStartTime-getTime()+statevar.SongMap[song].Length, false, "Stopping "+song, song, true)
}

func cleanup() {
	fmt.Println("cleanup")
}

func main() {

	// Load configuration parameters
	if _, err := toml.DecodeFile("./config.cfg", &conf); err != nil {
		// handle error
	}
	fmt.Printf("%v", conf)

	// Load state
	if _, err := os.Stat("state.json"); err == nil {
		dat, err := ioutil.ReadFile("state.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(dat, &statevar)
		fmt.Println("\n*******")
		fmt.Println(statevar.CurrentSong)
		fmt.Println("*******\n")
		statevar.IsPlaying = false
		statevar.SongList = []string{}
	} else {
		statevar = State{
			SongMap:          make(map[string]Song),
			SongList:         []string{},
			PathList:         make(map[string]bool),
			SongStartTime:    0,
			IsPlaying:        false,
			CurrentSong:      "None",
			CurrentSongIndex: 0,
		}
	}

	// Load Mp3s
	loadMp3s(conf.ServerParameters.MusicFolder)

	// Load song list
	for k, _ := range statevar.SongMap {
		statevar.SongList = append(statevar.SongList, k)
	}
	statevar.SongList.Sort()

	skipTrack(statevar.CurrentSongIndex)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /")
		html_response := index_html
		html_response = strings.Replace(html_response, "{{ data['random_integer'] }}", strconv.Itoa(rand.Intn(10000)), -1)
		html_response = strings.Replace(html_response, "{{ data['check_up_wait_time'] }}", strconv.Itoa(1700), -1)
		html_response = strings.Replace(html_response, "{{ data['max_sync_lag'] }}", strconv.Itoa(50), -1)
		html_response = strings.Replace(html_response, "{{ data['message'] }}", "Syncing...", -1)
		html_response = strings.Replace(html_response, "https://cdnjs.cloudflare.com/ajax/libs/mathjs/2.5.0/math.min.js", "/math.js", -1)
		html_response = strings.Replace(html_response, "https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js", "/jquery.js", -1)
		html_response = strings.Replace(html_response, "/static/howler.js", "/howler.js", -1)
		html_response = strings.Replace(html_response, "/static/normalize.css", "/normalize.css", -1)
		html_response = strings.Replace(html_response, "/static/skeleton.css", "/skeleton.css", -1)
		html_response = strings.Replace(html_response, "{{ data['playlist_html'] | safe }}", getPlaylistHTML(), -1)
		fmt.Fprintf(w, html_response)
	})

	mux.HandleFunc("/sound.mp3", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /sound.mp3")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte(rawSongData))
	})
	mux.HandleFunc("/howler.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /howler.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(howler_js))
	})
	mux.HandleFunc("/math.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /math.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(jquery_js))
	})
	mux.HandleFunc("/jquery.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /jquery.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(math_js))
	})
	mux.HandleFunc("/skeleton.css", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /skeleton.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(skeleton_css))
	})
	mux.HandleFunc("/normalize.css", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /normalize.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(normalize_css))
	})
	mux.HandleFunc("/sync", SyncRequest)
	mux.HandleFunc("/nextsong", NextSongRequest)
	//http.ListenAndServe(":5000", nil)
	graceful.Run(":5000", 10*time.Second, mux)
}
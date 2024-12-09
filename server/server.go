package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
)

func Run(conn net.Conn) error {
	sm := http.NewServeMux()

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "get_property", "playlist")
		if err != nil {
			log.Panic(err)
		}

		var reply mpv.Reply
		trimmedBuf := bytes.TrimRight(buf, "\x00")
		err = json.Unmarshal(trimmedBuf, &reply)
		if err != nil {
			log.Panic(err)
		}

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "text/html")

		tmpl.Execute(w, reply.Data)
	})

	sm.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("parseform err: ", err)
		}

		query := r.FormValue("query")
		fmt.Printf("query: %v\n", query)
		searchResult, err := client.Search(query)

		link := fmt.Sprintf("https://youtube.com/watch?v=%v", searchResult.Id.VideoId)

		mpv.Send(conn, mpv.PlaylistManipLoadFile, link, mpv.PlaylistManipLoadfileFlagAppend+"-play")

		w.Header().Add("HX-Refresh", "true")
	})

	sm.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		mpv.Send(conn, mpv.PlaybackControlStop)
	})

	sm.HandleFunc("/toggle-pause", func(w http.ResponseWriter, r *http.Request) {
		mpv.Send(conn, mpv.PropertyManipCycle, mpv.PropertyNamePause)
	})

	sm.HandleFunc("/playlist-next", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "playlist-next", "force")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
	})

	sm.HandleFunc("/playlist-prev", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "playlist-prev", "force")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
	})

	sm.HandleFunc("/playlist-clear", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "playlist-clear")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
	})

	sm.HandleFunc("/playlist", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "get_property", "playlist")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))

		trimmedBuf := bytes.TrimRight(buf, "\x00")

		var reply mpv.Reply

		if err = json.Unmarshal(trimmedBuf, &reply); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("reply: %v\n", reply)

		for _, song := range reply.Data {
			fmt.Printf("song: %v\n", song)
		}
	})

	sm.HandleFunc("/video-yes", func(w http.ResponseWriter, r *http.Request) {
		if err := mpv.Send(conn, "set_property", "video", "auto"); err != nil {
			log.Fatal(err)
		}
	})

	sm.HandleFunc("/video-no", func(w http.ResponseWriter, r *http.Request) {
		if err := mpv.Send(conn, "set_property", "video", "no"); err != nil {
			log.Fatal(err)
		}

	})

	sm.HandleFunc("/fullscreen-yes", func(w http.ResponseWriter, r *http.Request) {
		if err := mpv.Send(conn, "set_property", "fullscreen", "yes"); err != nil {
			log.Fatal(err)
		}

	})

	sm.HandleFunc("/fullscreen-no", func(w http.ResponseWriter, r *http.Request) {
		if err := mpv.Send(conn, "set_property", "fullscreen", "no"); err != nil {
			log.Fatal(err)
		}

	})

	server := http.Server{
		Addr:    "192.168.3.112:8080",
		Handler: sm,
		// ReadTimeout:       10 * time.Second,
		// ReadHeaderTimeout: 10 * time.Second,
		// WriteTimeout:      10 * time.Second,
		// IdleTimeout:       10 * time.Second,
	}

	fmt.Printf("server.Addr: http://%v\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

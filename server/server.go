package server

import (
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
		err := mpv.Send(conn, "get_property", "playlist")
		if err != nil {
			log.Panic(err)
		}

		// var reply mpv.Reply
		//
		// trimmedBuf := bytes.TrimRight(buf, "\x00")
		//
		// err = json.Unmarshal(trimmedBuf, &reply)
		// if err != nil {
		// 	log.Panic(err)
		// }

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "text/html")

		tmpl.Execute(w, nil)
	})

	sm.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\"this ran\": %v\n", "this ran")
		err := r.ParseForm()
		if err != nil {
			log.Fatal("parseform err: ", err)
		}

		query := r.FormValue("query")
		fmt.Printf("query: %v\n", query)
		searchResult, err := client.Search(query)

		link := fmt.Sprintf("https://youtube.com/watch?v=%v", searchResult.Id.VideoId)

		mpv.Send(conn, mpv.PlaylistManipLoadFile, link, mpv.PlaylistManipLoadfileFlagAppend+"-play")
	})

	sm.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		mpv.Send(conn, mpv.PlaybackControlStop)
	})

	sm.HandleFunc("/toggle-pause", func(w http.ResponseWriter, r *http.Request) {
		mpv.Send(conn, mpv.PropertyManipCycle, mpv.PropertyNamePause)
	})

	sm.HandleFunc("/playlist-clear", func(w http.ResponseWriter, r *http.Request) {
		buf, err := mpv.SendAndReceive(conn, "playlist-clear")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("string(buf): %v\n", string(buf))
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

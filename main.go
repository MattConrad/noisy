package main

import (
	"flag"
	"fmt"
	"github.com/jsgoecke/attspeech"
	"github.com/julienschmidt/httprouter"
	"github.com/twinj/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var port int
var wavfilesdir string
var attclient *attspeech.Client

func init() {
	wavfilesdir = "./wavfiles/"

	flag.IntVar(&port, "port", 8080, "The port number running the web server interface for Noisy.")
	attId := flag.String("attid", "", "The application id used in the AT&T speech API.")
	attSecret := flag.String("attsecret", "", "The application secret used in the AT&T speech API.")

	flag.Parse()
	attclient = attspeech.New(*attId, *attSecret, "")
}

func main() {
	router := httprouter.New()

	router.GET("/", index)
	router.GET("/failed", failed)
	router.POST("/run/:wavname", run)
	router.POST("/speak", speak)

	addr := fmt.Sprintf(":%d", port)

	log.Fatal(http.ListenAndServe(addr, router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./index.html")
}

func failed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "./failed.html")
}

func run(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	play(wavfilesdir + params.ByName("wavname") + ".wav")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func speak(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	r.ParseForm()
	sentence := r.PostFormValue("sentence")
	if strings.TrimSpace(sentence) == "" {
		log.Println("no sentence provided, request aborted")
		http.Redirect(w, r, "/failed", http.StatusSeeOther)
		return
	}

	attclient.SetAuthTokens()
	apiRequest := attclient.NewAPIRequest(attspeech.TTSResource)
	apiRequest.ContentType = "text/plain"
	apiRequest.Accept = "audio/x-wav"
	apiRequest.VoiceName = "mike"
	apiRequest.Text = sentence

	data, err := attclient.TextToSpeech(apiRequest)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/failed", http.StatusSeeOther)
		return
	}

	filename := wavfilesdir + uuid.NewV4().String() + ".wav"
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/failed", http.StatusSeeOther)
		return
	}
	defer os.Remove(filename)

	play(filename)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//written for Windows machines only. requires PowerShell (most Win machines will have this already).
func play(wavfile string) {
	cmd := exec.Command("powershell", "-c", "(New-Object Media.SoundPlayer \""+wavfile+"\").PlaySync()")
	cmd.Run()
}

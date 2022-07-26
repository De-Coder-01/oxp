package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/man0xff/oxp"
	. "github.com/man0xff/oxp"
)

func ServeRequests() {
	fmt.Println("Serving now...")
	http.HandleFunc("/getMeaning/", getMeaningHandler)
	http.HandleFunc("/getSentence/", getSentenceHandler)
	http.ListenAndServe("127.0.0.1:9090", nil)
}

func getMeaningHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	word := strings.TrimPrefix(r.URL.Path, "/getMeaning/")

	if word == "" {
		jsonWordMissing, _ := json.Marshal("word missing in url path")
		w.Write(jsonWordMissing)
	} else {
		res := getRes(word, "meaning")
		jsonRes, _ := json.Marshal(res)
		w.Write(jsonRes)
	}

}

func getSentenceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	word := strings.TrimPrefix(r.URL.Path, "/getSentence/")

	if word == "" {
		jsonWordMissing, _ := json.Marshal("word missing in url path")
		w.Write(jsonWordMissing)
	} else {
		res := getRes(word, "sentence")
		jsonRes, _ := json.Marshal(res)
		w.Write(jsonRes)
	}

}

func getRes(word, what string) string {

	p := oxp.NewClient()
	resp, _ := p.Search(context.Background(), word)

	var entries = []*Entry(resp.([]*Entry))

	if what == "meaning" {
		res := entries[0].Senses[0].Def
		return res
	} else {
		res := entries[0].Senses[0].Examples[0]
		return res
	}

}
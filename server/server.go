package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/plain")
	word := strings.TrimPrefix(r.URL.Path, "/getMeaning/")

	if word == "" {
		jsonWordMissing, _ := json.Marshal("word missing in url path")
		w.Write(jsonWordMissing)
	} else {
		res := getRes(word, "meaning")
		w.Write([]byte(res))
	}

}

func getSentenceHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/plain")
	word := strings.TrimPrefix(r.URL.Path, "/getSentence/")

	if word == "" {
		jsonWordMissing, _ := json.Marshal("word missing in url path")
		w.Write(jsonWordMissing)
	} else {
		res := getRes(word, "sentence")
		w.Write([]byte(res))
	}

}

func getRes(word, what string) string {

	f, err := os.OpenFile("usage.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	p := oxp.NewClient()
	resp, _ := p.Search(context.Background(), word)

	if resp == nil {
		if what == "meaning" {
			return "no meaning"
		}
		return "no sentence"
	}

	var entries = []*Entry(resp.([]*Entry))

	if what == "meaning" {
		log.Println("Got request for meaning of : ", word)
		if len(entries[0].Senses) == 0 {
			return "no meaning"
		}
		res := entries[0].Senses[0].Def
		return res
	} else {
		log.Println("Got request for sentence of : ", word)

		/*
			Return if no sense
		*/
		if len(entries[0].Senses) == 0 {
			return "no sentence"
		}

		/*
			Iterative search for example, append with |
		*/
		resultantExample := ""
		for _, sens := range entries[0].Senses {
			if len(sens.Examples) != 0 {
				for _, exp := range sens.Examples {
					resultantExample += (exp + " | ")
				}
			}
		}

		/*
			Remove last stray " |"
		*/
		resultantExample = resultantExample[:len(resultantExample)-3]

		return resultantExample
	}

}

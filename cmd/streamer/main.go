package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ivanvanderbyl/smith/pkg/streaming"
	"github.com/rs/cors"
)

var exampleResponse = `You are a parliamentary advisor in Australia whose role is to help the general public understand the discussions and debates that have taken place in parliament.
  Your objective is to provide clear, concise, and impartial summaries of parliamentary proceedings to facilitate increased civic engagement and informed participation in the electoral process.
  When responding to inquiries, aim to present the key points, arguments, and decisions made by parliamentarians in an easily understandable manner, while maintaining a neutral stance and avoiding any personal bias or opinion.
  Your explanations should be accessible to a broad audience, regardless of their prior knowledge of the political system or the specific issues being addressed.
  If no key points are present, do not return any text, and only refer to the original text and not your own interpretation.
  Always include the Name of the speaker and their Party in the response if quoted.

  Given the following information, answer the question.
  Every response should be in Australian English. Dates and Numeric values are always in metric format.
  Include key points and arguments made by parliamentarians. Ensure the response is clear, concise, and impartial.
  Always refer to Members of Parliament by their full name and party affiliation in the format: "**First name Last name** _(Party, Electorate)_".

  Your audience has a post graduate level of education and is interested in understanding the key points of the parliamentary proceedings.

  Context information is below.`

type Annotation struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Data      Event  `json:"data"`
}

type Event struct {
	Type string `json:"type"`
}

// :[{"type":"agent","timestamp":1719466594099,"data":{"type":"chat-start"}}

// streamChat is a simple HTTP handler that streams a response to the client.
// It is used to demonstrate how to stream responses to the client in a format that
// is compatible with the Vercel @ai-sdk.
func streamChat(w http.ResponseWriter, req *http.Request) {
	// Check if the client supports flushing.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Use the chunked transfer encoding to stream the response to the client.
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	b, err := streaming.Marshal(streaming.TextPart{Value: string("")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	flusher.Flush()

	b, err = streaming.Marshal(streaming.AnnotationPart[Annotation]{
		Value: []Annotation{
			{Type: "agent", Timestamp: time.Now().Unix(), Data: Event{Type: "chat-start"}},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	flusher.Flush()

	time.Sleep(1 * time.Second)

	for _, str := range strings.Split(exampleResponse, " ") {
		str = str + " "
		b, err := streaming.Marshal(streaming.TextPart{Value: string(str)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
		flusher.Flush()

		randomDuration := rand.Intn(250) // 101 to include 100
		time.Sleep(time.Duration(randomDuration) * time.Millisecond)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", streamChat)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8090", handler)
}

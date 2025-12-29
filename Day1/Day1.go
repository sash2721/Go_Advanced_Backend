package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Time struct {
	CurrentTime string `json:"currentTime"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()                                   // find current time
	formattedTime := currentTime.Format("2006/01/02, 03:04 PM") // formatting time in string

	jsonString := Time{CurrentTime: formattedTime} // creating a json string
	jsonData, err := json.Marshal(jsonString)      // converting in JSON

	if err != nil {
		fmt.Println("Error while marshaling:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error while processing the data!"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonData))
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// check if the request is a POST request or not
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Only POST methods allowed"}`))
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error while reading the request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Failed to read the request body!"}`))
		return
	}
	defer r.Body.Close() // close the connection at the end

	// echo the json back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Return response as Hello!
		w.Write([]byte("Hello!"))
	})

	// Adding the /health route
	http.HandleFunc("/health", HealthHandler)

	// Adding the /time route
	http.HandleFunc("/time", TimeHandler)

	// Adding the echo route
	http.HandleFunc("/echo", EchoHandler)

	var port string = ":8080"
	fmt.Printf("Starting the Server on PORT %s\n", port)

	// Start the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("ListenAndServer Error:", err)
	}
}

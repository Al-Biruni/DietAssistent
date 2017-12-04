package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	cors "github.com/heppu/simple-cors"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
	"bytes"
	"os"
)

var key string = "qHwatDFKxKxnzGTNLIFT2paDyuQvkMNsHo193M3o"
//var key string = "y0eBExLZ0UmM8l9ycH9O8Ap7BRd28m6ywIoTNcre"
var fullReport bool = false
type listItem struct {
	List struct {
		Q     string `json:"q"`
		Sr    string `json:"sr"`
		Ds    string `json:"ds"`
		Start int    `json:"start"`
		End   int    `json:"end"`
		Total int    `json:"total"`
		Group string `json:"group"`
		Sort  string `json:"sort"`
		Item  []struct {
			Offset int    `json:"offset"`
			Group  string `json:"group"`
			Name   string `json:"name"`
			Ndbno  string `json:"ndbno"`
			Ds     string `json:"ds"`
		} `json:"item"`
	} `json:"list"`
}

type resReport struct {
	Report struct {
		Sr   string `json:"sr"`
		Type string `json:"type"`
		Food struct {
			Ndbno     string `json:"ndbno"`
			Name      string `json:"name"`
			Ds        string `json:"ds"`
			Manu      string `json:"manu"`
			Ru        string `json:"ru"`
			Nutrients []struct {
				NutrientID string `json:"nutrient_id"`
				Name       string `json:"name"`
				Derivation string `json:"derivation"`
				Group      string `json:"group"`
				Unit       string `json:"unit"`
				Value      string `json:"value"`
				Measures   []struct {
					Label string  `json:"label"`
					Eqv   float64 `json:"eqv"`
					Eunit string  `json:"eunit"`
					Qty   float64 `json:"qty"`
					Value string  `json:"value"`
				} `json:"measures"`
			} `json:"nutrients"`
		} `json:"food"`
		Footnotes []interface{} `json:"footnotes"`
	} `json:"report"`
}

func searchFor(items string) string {
	tList := listItem{}

	searchItemsURL := []string{}
	searchItemsURL = append(searchItemsURL, "https://api.nal.usda.gov/ndb/search/?", "api_key=", key, "&format=json", "&ds=Standard%20Reference")
	searchItemsURL = append(searchItemsURL, "&q=", items)
	searchItemsURLstring := strings.Join(searchItemsURL, "")
	fmt.Printf(searchItemsURLstring)
	resp, err := http.Get(searchItemsURLstring)
	if err != nil {
		fmt.Printf("response err")
		panic(err.Error())
	} else {
		defer resp.Body.Close()
		rList := new(bytes.Buffer)
		rList.ReadFrom(resp.Body)
		data := rList.String()
		//fmt.Printf(data)
		err := json.Unmarshal([]byte(data), &tList)
		if err != nil {
			fmt.Printf("decode err")
			panic(err.Error())
		}
		//  fmt.Println(tList)
		items := tList.List.Item
		return getReport(items[0].Ndbno)
	}
	//return tList
}

func getReport(num string) string {
	reportURL := []string{}
	tReport := resReport{}
	reportURL = append(reportURL, "https://api.nal.usda.gov/ndb/reports/?", "api_key=", key, "&format=json", "&type=b")
	//for _,num :=range dataNum{
	reportURL = append(reportURL, "&ndbno=", num)
	//  }

	reportURLstring := strings.Join(reportURL, "")
	//  fmt.Println(reportURLstring)
	resp, err := http.Get(reportURLstring)
	if err != nil {
		fmt.Println("response err")
		panic(err.Error())
	} else {
		defer resp.Body.Close()

		rReport := new(bytes.Buffer)
		rReport.ReadFrom(resp.Body)
		data := rReport.String()
		//fmt.Printf(data)
		err := json.Unmarshal([]byte(data), &tReport)
		if err != nil {
			fmt.Println("decode err")
			panic(err.Error())
		}
		//fmt.Println(tReport.Report.Food.Nutrients[2].Name)
		if(fullReport){
			fR := []string{}
			for _, n := range tReport.Report.Food.Nutrients {
				fR=append(fR,n.Name ," " , n.Value ," ",n.Unit ," ")
			}
			fRString := strings.Join(fR,"")
			return fRString
		}else{

		for _, f := range tReport.Report.Food.Nutrients {
			if f.Name == "Energy" {
				return f.Value + "   " + f.Unit + "  per 100g"
			}
		}

	}
}
	return ""

}

var (
	// WelcomeMessage A constant to hold the welcome message
	WelcomeMessage = "Welcome,if you want to inqury about food enter food name;for chat history type history  what food do u wanna know about ?"

	// sessions = {String welcomeUrl
	//   "uuid1" = Session{...},
	//   ...
	// }
	sessions = map[string]Session{}

	processor = sampleProcessor
)

type (
	// Session Holds info about a session
	Session map[string]interface{}

	// JSON Holds a JSON object
	JSON map[string]interface{}

	// Processor Alias for Process func
	Processor func(session Session, message string) (string, error)
)

func sampleProcessor(session Session, message string) (string, error) {
	// Make sure a history key is defined in the session which points to a slice of strings
	_, historyFound := session["history"]
	if !historyFound {
		session["history"] = []string{}
	}

	// Fetch the history from session and cast it to an array of strings
	history, _ := session["history"].([]string)

	// Make sure the message is unique in history
	/*
	for _, m := range history {
		if strings.EqualFold(m, message) {
			return "", fmt.Errorf("You've already ordered %s before!", message)
		}
	}*/

	// Add the message in the parsed body to the messages in the session
	history = append(history, message)

	// Form a sentence out of the history in the form Message 1, Message 2, and Message 3
	l := len(history)
	wordsForSentence := make([]string, l)
	copy(wordsForSentence, history)
	if l > 1 {
		wordsForSentence[l-1] = "and " + wordsForSentence[l-1]

	}
	//sentence := strings.Join(wordsForSentence, ", ")

	// Save the updated history to the session
	session["history"] = history
	messageA:=[]string{}
	messageA = strings.Split(message," ")

	for _,m := range messageA{
		fmt.Printf(m)
if strings.EqualFold("fullReport",m){
	fullReport =true;
	sI :=strings.Join(messageA[1:],"")
	cal := searchFor(sI)
	fullReport =false;
	if(cal==""){
		cal:="no such thing"
		return fmt.Sprintf("invalid ",cal),nil
	}
	return fmt.Sprintf(" %s , have %s! What else?", message,cal), nil
}
}


	if strings.EqualFold("history",message){
		 hist:=[]string{}
		for _, m := range history{
		hist = append(hist ," ", m)
	}
historyString := strings.Join(hist,"")

		return fmt.Sprintf("you have asked about  %s", historyString), nil
	}
	cal := searchFor(message)
	if(cal==""){
		cal:="no such thing"
		return fmt.Sprintf("invalid ",cal),nil
	}
	return fmt.Sprintf(" %s , have %s! What else?", message,cal), nil
}

// withLog Wraps HandlerFuncs to log requests to Stdout
func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := httptest.NewRecorder()
		fn(c, r)
		log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)
	}
}

// writeJSON Writes the JSON equivilant for data into ResponseWriter w
func writeJSON(w http.ResponseWriter, data JSON) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ProcessFunc Sets the processor of the chatbot
func ProcessFunc(p Processor) {
	processor = p
}

// handleWelcome Handles /welcome and responds with a welcome message and a generated UUID
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	// Generate a UUID.
	hasher := md5.New()
	hasher.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	uuid := hex.EncodeToString(hasher.Sum(nil))

	// Create a session for this UUID
	sessions[uuid] = Session{}
	fmt.Sprintf(uuid)
	// Write a JSON containg the welcome message and the generated UUID

	writeJSON(w, JSON{
		"uuid":    uuid,
		"message": WelcomeMessage,
	})
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// Make sure only POST requests are handled
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Make sure a UUID exists in the Authorization header
	uuid := r.Header.Get("Authorization")
	if uuid == "" {
		http.Error(w, "Missing or empty Authorization header.", http.StatusUnauthorized)
		return
	}

	// Make sure a session exists for the extracted UUID
	session, sessionFound := sessions[uuid]
	if !sessionFound {
		http.Error(w, fmt.Sprintf("No session found for: %v.", uuid), http.StatusUnauthorized)
		return
	}

	// Parse the JSON string in the body of the request
	data := JSON{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Couldn't decode JSON: %v.", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Make sure a message key is defined in the body of the request
	_, messageFound := data["message"]
	if !messageFound {
		http.Error(w, "Missing message key in body.", http.StatusBadRequest)
		return
	}

	// Process the received message
	message, err := processor(session, data["message"].(string))
	if err != nil {
		http.Error(w, err.Error(), 422 /* http.StatusUnprocessableEntity */)
		return
	}

	// Write a JSON containg the processed response
	writeJSON(w, JSON{
		"message": message,
	})
}

// handle Handles /
func handle(w http.ResponseWriter, r *http.Request) {
	body :=
		"<!DOCTYPE html><html><head><title>Chatbot</title></head><body><pre style=\"font-family: monospace;\">\n" +
			"Available Routes:\n\n" +
			"  GET  /welcome -> handleWelcome\n" +
			"  POST /chat    -> handleChat\n" +
			"  GET  /        -> handle        (current)\n" +
			"</pre></body></html>"
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(w, body)
}

// Engage Gives control to the chatbot
func Engage(addr string) error {
	// HandleFuncs
	mux := http.NewServeMux()
	mux.HandleFunc("/welcome", withLog(handleWelcome))
	mux.HandleFunc("/chat", withLog(handleChat))
	mux.HandleFunc("/", withLog(handle))

	// Start the server
	return http.ListenAndServe(addr, cors.CORS(mux))
}

func main() {
	// Uncomment the following lines to customize the chatbot

	//WelcomeMessage = "What's your gender?" + getSymptoms(tokennn)

	//fmt.Println("What's your gender?")
	//ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(Engage(":" + port))
}

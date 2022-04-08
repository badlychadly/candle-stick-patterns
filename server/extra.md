package main


import (
	"fmt"
	"net/http"  
	"io/ioutil"
	"io"
	"strconv"
	"strings"
	// "net/url"
	// "math"
	"golang.org/x/net/html" 
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"time" 
)



var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
  WriteBufferSize: 1024,

  // We'll need to check the origin of our connection
  // this will allow us to make requests from our React
  // development server to here.
  // For now, we'll do no checking and just allow any connection
  CheckOrigin: func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
    // print out that message for clarity
        fmt.Printf("p= : %v\n", string(p))

        if err := conn.WriteMessage(messageType, []byte("In Server!")); err != nil {
            log.Println(err)
            return
        }

    }
}

func writer(conn *websocket.Conn) {
	// defer func() {
	// 	conn.Close()
	// }()
	start := time.Now()
	allStocks := []Stock{}
	scrapeUrl := "https://finviz.com/screener.ashx?v=111&f=sh_curvol_o5000,sh_price_u1"
	resp, err := http.Get(scrapeUrl)
		
		if err != nil {
			if newErr := conn.WriteMessage(1, []byte(err.Error())); newErr != nil {
				log.Println(err)
				return
			}
		}
		defer resp.Body.Close()
		colSlice, total := GetHtmlTable(resp.Body)
		fmt.Printf("writerTotal: %v, length: %v\n", total, len(colSlice))

		allStocks = append(allStocks, colSlice...)

		// if err := conn.WriteJSON(colSlice); err != nil {
		// 	log.Println(err)
		// 	return
		// }
		
		for countedStocks := len(colSlice); countedStocks < total; { 
			nextUrl := fmt.Sprintf("%v&r=%v", scrapeUrl, countedStocks+1)
			fmt.Printf("nextUrl: %v\n", nextUrl)
			newResp, err := http.Get(nextUrl)
			if err != nil {
				if err := conn.WriteMessage(1, []byte("request failed")); err != nil {
					log.Println(err)
					return
				}
			}
			defer resp.Body.Close()
			newSlice, _ := GetHtmlTable(newResp.Body)
			countedStocks += len(newSlice)

			allStocks = append(allStocks, newSlice...)

			// if err := conn.WriteJSON(newSlice); err != nil {
			// 	log.Println(err)
			// 	return
			// }

		}
		duration := time.Since(start)
		fmt.Printf("first operation: %v\n", duration.Nanoseconds())

		start = time.Now()

		for i, v := range allStocks {
			// encodedValue := "https://query1.finance.yahoo.com/v8/finance/chart/GNUS?symbol=GNUS&period1=1542092400&period2=1595621341&interval=1d&includePrePost=true&events=div%7Csplit%7Cearn&lang=en-US&region=US&crumb=RXvYwZiv1bi&corsDomain=finance.yahoo.com"
			
			// decodedValue, err := url.QueryUnescape(encodedValue)
			// if err != nil {
			// 	log.Fatal(err)
			// 	return
			// }
			// fmt.Println(decodedValue)
			now := 	time.Now().UnixNano() / 1000000
			candleUrl := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%v?symbol=%v&period1=1591018200&period2=%v&interval=1d&includePrePost=true&events=div|split|earn&lang=en-US&region=US&crumb=RXvYwZiv1bi&corsDomain=finance.yahoo.com", v.Symbol, v.Symbol, now)
			// fmt.Printf("path: %v\n", now)

			newResp, err := http.Get(candleUrl)
			if err != nil {
				if err := conn.WriteMessage(1, []byte("request failed")); err != nil {
					log.Println(err)
					return
				}
			}
			defer newResp.Body.Close()
			body, _ := ioutil.ReadAll(newResp.Body)
			allStocks[i].CandleData = string(body)
			// if err := conn.WriteJSON(v); err != nil {
			// 	log.Println(err)
			// 	return
			// }
			// fmt.Printf("stock: %v\n", v)
			if err := conn.WriteJSON(allStocks[i]); err != nil {
				log.Println(err)
				return
			}
		}
		duration = time.Since(start)
		fmt.Printf("SECOND operation: %v\n",duration.Nanoseconds())
		if err := conn.WriteMessage(1, []byte("finished")); err != nil {
			log.Println(err)
			return
		}
		// fmt.Printf("allstocks: %v\n", allStocks)
		
		// fmt.Printf("allStocks: %v\n", allStocks)
}


func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	w.Header().Set("Access-Control-Allow-Origin", "*")

  // upgrade this connection to a WebSocket
  // connection
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
  }
  // listen indefinitely for new messages coming
  // through on our WebSocket connection
  go writer(ws)
    reader(ws)
}




func HomePage(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		resp, err := http.Get("https://finviz.com/screener.ashx?v=111&f=ta_candlestick_d&o=-volume")
		defer resp.Body.Close()
		colSlice, _ := GetHtmlTable(resp.Body)
		// body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "could not get request")
		}
		json.NewEncoder(w).Encode(colSlice)
		// fmt.Fprintf(w, "HomePage")
	
		// fmt.Fprintf(w, "Homepage")
}



func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/ws", serveWs)
	return router
}

type Stock struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	Volume string `json:"volume"`
	CandleData string `json:"candleData"`
}

func CreateStock(Symbol, Price, Volume string) Stock {
	fmt.Printf("Symbol: %v, Price: %v, Volume: %v\n", Symbol, Price, Volume)
	return Stock{Symbol: Symbol, Price: Price, Volume: Volume}
}


func GetHtmlTable(httpBody io.Reader) ([]Stock, int) {
	z := html.NewTokenizer(httpBody)
	var content = []Stock{}
	getNextTextToken := false
	var total int

	for {
		tt := z.Next()
		switch tt {
			case html.ErrorToken:
				return content, total
			case html.TextToken:
				text := (string)(z.Text())
				// fmt.Printf("textToken: %v\n", text)
				if getNextTextToken {
					// total = 0
					tString := strings.TrimSpace(strings.Split(text, "#")[0])
					
					total, _ = strconv.Atoi(tString)
					getNextTextToken = false
					// fmt.Printf("textToken: %v, tString: %v\n", total, tString)
				}
				
				if strings.TrimSpace(text) == "Total:"{
					getNextTextToken = true
					// fmt.Printf("textToken: %v\n", text)

				}
			case html.CommentToken:
				text := (string)(z.Text())
				if text[:1] != "<" {
	
					stocks := strings.Split(text, "\n")
					// formatted := []string{}
					for _, v := range stocks {
						if strings.Contains(v, "|") {
							sd := strings.Split(v, "|")
							stock := Stock{Symbol: sd[0], Price: sd[1], Volume: sd[2]}
							// stock := CreateStock(sd...)
							// fmt.Printf("value: %v\n", stock)
							content = append(content, stock)
						}
	
					}
				}
		}

	}
	
	return content, total
}




func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3001", router))
  } 
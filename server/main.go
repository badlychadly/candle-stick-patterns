package main


import (
	"fmt"
	"net/http"  
	"io/ioutil"
	// "io"
	// "strconv"
	// "strings"
	// "net/url"
	"sync"
	// "math"
	"bytes"
	// "golang.org/x/net/html" 
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"time" 
	"stock-helper/server/patterns"
	"stock-helper/server/structs"
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

type Settings struct {
	From int64 `json:"from"`
	Interval string `json:"interval"`
}

func reader(conn *websocket.Conn) {
    for {
	// read in a message
	// fmt.Printf("Conn: %+T\n", conn)
        messageType, _, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
		}
		// var settings Settings
		// if err := json.Unmarshal(p, &settings); err != nil{
		// 	fmt.Printf("errr: %v\n", err)
		// }
		// fmt.Printf("p= : %+v\n", settings)
		// if settings.From != "" {
		// 	writer(conn, settings)
		// }

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
	allStocks := []structs.Stock{}
	scrapeUrl := "https://scanner.tradingview.com/america/scan"
		postBody := []byte(`{"filter":[{"left":"volume","operation":"nempty"},{"left":"volume","operation":"egreater","right":200000},{"left":"close","operation":"in_range","right":[0.40,3]}],"options":{"lang":"en"},"symbols":{"query":{"types":[]},"tickers":[]},"columns":["logoid","name","close","change","volume"],"sort":{"sortBy":"market_cap_basic","sortOrder":"desc"},"range":[0,300]}`)

	
	resp, err := http.Post(scrapeUrl, "application/x-www-form-urlencoded", bytes.NewReader(postBody))
		
		if err != nil {
			if err := conn.WriteMessage(1, []byte("request failed")); err != nil {
				log.Println(err)
				return
			}
		}
		defer resp.Body.Close()
		textBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ioutil Error: %+v\n", err)
		}

		var tradingViewData structs.TradingViewJson
		err = json.Unmarshal(textBody, &tradingViewData)
		if err != nil {
			fmt.Printf("Cant Unmashall err: %v\n", err)
			// return
			log.Println(err)
		}

		for _, stockData := range tradingViewData.Data {
			symbol, sok := stockData.D[1].(string)
			price, pok := stockData.D[2].(float64)
			volume, vok := stockData.D[4].(float64)
			// fmt.Printf("volume %T ok: %v\n", volume, vok)
			if sok && pok && vok {
				allStocks = append(allStocks, structs.Stock{Symbol: symbol, Price: price, Volume: volume})
			} else {
				fmt.Println("wrong dataTypes for Stock")
			}
			
		} 


		fmt.Printf("allStocks: %+v\n", len(allStocks))

		
		// colSlice, total := GetHtmlTable(resp.Body)
		// fmt.Printf("writerTotal: %v, length: %v\n", total, len(colSlice))

		// allStocks = append(allStocks, colSlice...)

		
		// for countedStocks := len(colSlice); countedStocks < total; { 
		// 	nextUrl := fmt.Sprintf("%v&r=%v", scrapeUrl, countedStocks+1)
		// 	fmt.Printf("nextUrl: %v\n", nextUrl)
		// 	newResp, err := http.Get(nextUrl)
		// 	if err != nil {
		// 		if err := conn.WriteMessage(1, []byte("request failed")); err != nil {
		// 			log.Println(err)
		// 			return
		// 		}
		// 	}
		// 	defer resp.Body.Close()
		// 	newSlice, _ := GetHtmlTable(newResp.Body)
		// 	countedStocks += len(newSlice)

		// 	allStocks = append(allStocks, newSlice...)

		// }
		duration := time.Since(start)
		fmt.Printf("first operation: %v\n", duration.Nanoseconds())

		start = time.Now()
	



		getPattern := func(stock structs.Stock, outerChan chan structs.Stock, failed chan structs.Stock, wg *sync.WaitGroup) {
			now := 	time.Now().UnixNano() / 1000000
						candleUrl := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%v?symbol=%v&period1=1641080470&period2=%v&interval=1d&includePrePost=true&events=div|split|earn&lang=en-US&region=US&crumb=RXvYwZiv1bi&corsDomain=finance.yahoo.com", stock.Symbol, stock.Symbol, now)
						// fmt.Printf("path: %v\n", now)
			
						newResp, err := http.Get(candleUrl)
						if err != nil {
	
							// fmt.Printf("ln 116::: http.Get: %v\n", err.Error())
							// return
							log.Println(err)
							failed <-stock
							// fmt.Printf("newResp: %+v\n", newResp)
							// fmt.Printf("Header: %+v\n", newResp.Header)
							// newResp.Body.Close()
							// wg.Done()
							return
						}
						// fmt.Printf("Header: %+v\n", newResp.Header)
						
						body, err := ioutil.ReadAll(newResp.Body)
						if err != nil {
							fmt.Printf("ioutil: %v\n", err)
						}
						var chart structs.ChartData
						err = json.Unmarshal(body, &chart)
						if err != nil {
							fmt.Printf("Cant Unmashall err: %v\n", err)
							// return
							log.Println(err)
						}
						newResp.Body.Close()
						
						if len(chart.Chart.Result) < 1 {
							wg.Done()
							return
							// fmt.Printf("Chart: %+v\n", len(chart.Chart.Result))
						}
						timestamps := chart.Chart.Result[0].Timestamp
						quote := chart.Chart.Result[0].Indicators.Quote[0]
						opens := quote.Open 
						closes := quote.Close
						lows := quote.Low
						highs := quote.High
						for i := 0; i < len(timestamps); i++ {
							// fmt.Printf("timestamp %v\n", timestamps[i])
							dataPoint := structs.DataPoint{X: timestamps[i] * 1000, Y: []float64{opens[i], highs[i], lows[i], closes[i]}}
							stock.DataPoints = append(stock.DataPoints, dataPoint)
						}

						if len(stock.DataPoints) < 6 {
							wg.Done()
							return
						}
						if pattern := patterns.LadderBottom(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.BullThreeLineStrike(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.BearThreeLineStrike(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.RisingThreeMethods(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.FallingThreeMethods(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.BullishMatHold(stock.DataPoints[len(stock.DataPoints) - 5:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.ThreeWhiteSoldiers(stock.DataPoints[len(stock.DataPoints) - 4:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.ThreeBlackCrows(stock.DataPoints[len(stock.DataPoints) - 4:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						} 
						if pattern := patterns.ThreeStarsInTheSouth(stock.DataPoints[len(stock.DataPoints) - 3:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.ThreeInsideUp(stock.DataPoints[len(stock.DataPoints) - 3:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						if pattern := patterns.BearishAbandonedBaby(stock.DataPoints[len(stock.DataPoints) - 3:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
						// if pattern := patterns.EveningStar(stock.DataPoints[len(stock.DataPoints) - 3:]...); pattern.Name != "" {
						// 	stock.Pattern = pattern
						// 	fmt.Printf("Sending stock!: %v\n", stock.Symbol)
						// 	// stockChan <- stock
						// 	outerChan <-stock
						// 	// return
						// 	wg.Done()
						// 	return
						// }
						// if pattern := patterns.MorningStar(stock.DataPoints[len(stock.DataPoints) - 3:]...); pattern.Name != "" {
						// 	stock.Pattern = pattern
						// 	fmt.Printf("Sending stock!: %v\n", stock.Symbol)
						// 	// stockChan <- stock
						// 	outerChan <-stock
						// 	// return
						// 	wg.Done()
						// 	return
						// } 
						if pattern := patterns.Hammer(stock.DataPoints[len(stock.DataPoints) - 2:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						} 
						if pattern := patterns.PiercingLine(stock.DataPoints[len(stock.DataPoints) - 2:]...); pattern.Name != "" {
							stock.Pattern = pattern
							fmt.Printf("Sending stock!: %v\n", stock.Symbol)
							// stockChan <- stock
							outerChan <-stock
							// return
							wg.Done()
							return
						}
	
						
						// if pattern, err := patterns.IsDoji(stock.DataPoints[len(stock.DataPoints) - 1]); err == nil{
							
						// 	stock.Pattern = pattern
						// 	// fmt.Printf("Symbol %+v\n", stock.Symbol)
						// 	fmt.Printf("Sending stock!: %v\n", stock.Symbol)
						// 	// stockChan <- stock
						// 	outerChan <-stock
						// 	// return
						// 	wg.Done()
						// 	return
						// } 
						// else {
							fmt.Printf("No pattern for stock %v\n", stock.Symbol)
							wg.Done()
							return
						// }
		}




		innerFunc := func() <-chan structs.Stock {
			innerChan := make(chan structs.Stock, 20)
			go func() {
				for _, stock := range allStocks {
					innerChan <-stock
				}
				close(innerChan)
			}()
			
			return innerChan
		}




		outerFunc := func(in <-chan structs.Stock) <-chan structs.Stock {
			outerChan := make(chan structs.Stock, 20)
			failed := make(chan structs.Stock, 20)
			var wg sync.WaitGroup
			
			wg.Add(len(allStocks))
			go func() {
				for stock := range in {
					// fmt.Printf("wg: %v\n", &wg)
					go getPattern(stock, outerChan, failed, &wg)
				}
				fmt.Println("Finished with range")
				
			}()

			go func() {
				for stock := range failed {
					fmt.Printf("recieved failed stock: %+v\n", stock)
					go getPattern(stock, outerChan, failed, &wg)
				}
				fmt.Println("Finished with failed range")
				
			}()

			go func() {
				wg.Wait()
				close(outerChan)
			}()
			return outerChan
		}




		for stock := range outerFunc(innerFunc()) {
			// fmt.Printf("Recieved Stock: %v\n", stock.Symbol)
			if err := conn.WriteJSON(stock); err != nil {
				log.Println(err)
				// return
			}
		}



		
		fmt.Println("Now closing the connection")
					// conn.Close()
					if err := conn.WriteMessage(1, []byte("finished")); err != nil {
						log.Println(err)
						// return
					}
					duration = time.Since(start)
					fmt.Printf("SECOND operation: %v\n",duration.Nanoseconds())
					fmt.Printf("allStocks len: %v\n", len(allStocks))
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




// func HomePage(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		resp, err := http.Get("https://finviz.com/screener.ashx?v=111&f=ta_candlestick_d&o=-volume")
// 		defer resp.Body.Close()
// 		colSlice, _ := GetHtmlTable(resp.Body)
// 		// body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Fprintf(w, "could not get request")
// 		}
// 		json.NewEncoder(w).Encode(colSlice)

// }



func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/ws", serveWs)
	return router
}



// func CreateStock(Symbol, Price, Volume string) structs.Stock {
// 	fmt.Printf("Symbol: %v, Price: %v, Volume: %v\n", Symbol, Price, Volume)
// 	return structs.Stock{Symbol: Symbol, Price: Price, Volume: Volume}
// }


// func GetHtmlTable(httpBody io.Reader) ([]structs.Stock, int) {
// 	z := html.NewTokenizer(httpBody)
// 	var content = []structs.Stock{}
// 	getNextTextToken := false
// 	var total int

// 	for {
// 		tt := z.Next()
// 		switch tt {
// 			case html.ErrorToken:
// 				return content, total
// 			case html.TextToken:
// 				text := (string)(z.Text())
// 				// fmt.Printf("textToken: %v\n", text)
// 				if getNextTextToken {
// 					// total = 0
// 					tString := strings.TrimSpace(strings.Split(text, "#")[0])
					
// 					total, _ = strconv.Atoi(tString)
// 					getNextTextToken = false
// 					// fmt.Printf("textToken: %v, tString: %v\n", total, tString)
// 				}
				
// 				if strings.TrimSpace(text) == "Total:"{
// 					getNextTextToken = true
// 					// fmt.Printf("textToken: %v\n", text)

// 				}
// 			case html.CommentToken:
// 				text := (string)(z.Text())
// 				if text[:1] != "<" {
	
// 					stocks := strings.Split(text, "\n")
// 					// formatted := []string{}
// 					for _, v := range stocks {
// 						if strings.Contains(v, "|") {
// 							sd := strings.Split(v, "|")
// 							stock := structs.Stock{Symbol: sd[0], Price: sd[1], Volume: sd[2]}
// 							// stock := CreateStock(sd...)
// 							// fmt.Printf("value: %v\n", stock)
// 							content = append(content, stock)
// 						}
	
// 					}
// 				}
// 		}

// 	}
	
// 	return content, total
// }




func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3001", router))
  } 
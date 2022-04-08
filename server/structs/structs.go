package structs


type ChartData struct {
	Chart struct {
		Result []struct {
			
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Volume []int     `json:"volume"`
					Close  []float64 `json:"close"`
					Low    []float64 `json:"low"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}


type TradingViewJson struct {
	Data       []Data `json:"data"`
	TotalCount int    `json:"totalCount"`
}
type Data struct {
	S string        `json:"s"`
	D []interface{} `json:"d"`
}




type DataPoint struct {
	X int64 `json:"x"`
	Y []float64 `json:"y"`
}


type Stock struct {
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
	Volume float64 `json:"volume"`
	DataPoints []DataPoint `json:"dataPoints"`
	Pattern StockPattern `json:"pattern"`
}

type StockPattern struct {
	Name       string `json:"name"`
	StripLines []struct {
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	} `json:"stripLines"`
}

type StripLineData map[string]string
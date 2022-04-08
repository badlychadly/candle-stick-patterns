package patterns


import(
	"math"
	"time"
	"fmt"
	"stock-helper/server/structs"
)

type Day struct {
	Open float64
	High float64
	Close float64
	Low float64
	Ct float64
	Cb float64
	Uw float64
	Lw float64
	StartTime int64
	EndTime int64
}

func CreateDay(dayData structs.DataPoint) Day {
	open := dayData.Y[0]
	high := dayData.Y[1]
	low := dayData.Y[2]
	close := dayData.Y[3]
	ct := math.Abs(high - low)
	cb := math.Abs(open - close)
	date := time.Unix(int64(dayData.X / 1000), 0)
	// fmt.Printf("Date: %v\n", date)
	year, month, day := date.Date()
	denver, err := time.LoadLocation("America/Denver")
    if err != nil {
        fmt.Println(err)
        // return
	}
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, denver)
	endOfDay := time.Date(year, month, day, 23, 0, 0, 0, denver)

	createdDay := Day{Open: open, High: high, Close: close, Low: low, Ct: ct, Cb: cb, StartTime: startOfDay.Unix() * 1000, EndTime: endOfDay.Unix() * 1000}
	if Bullish(createdDay){
		createdDay.Uw = createdDay.High - createdDay.Close
		createdDay.Lw = createdDay.Open - createdDay.Open
	} else {
		createdDay.Uw = createdDay.High - createdDay.Open
		createdDay.Lw = createdDay.Close - createdDay.Low
	}
	return createdDay
}

func Bullish(day Day) bool {
	if day.Close > day.Open {
		return true
	}
	return false
}

func LadderBottom(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])
	dayFour := CreateDay(days[3])
	lastDay := CreateDay(days[4])

	if (Bullish(dayOne)) || (Bullish(dayTwo)) || (Bullish(dayThree)) || (Bullish(dayFour)) {
		return
	}
	if (dayTwo.Open >= dayOne.Open) || (dayThree.Open >= dayTwo.Open) || (dayFour.Open >= dayThree.Open) {
		return
	}
	if (dayTwo.Close >= dayOne.Close) || (dayThree.Close >= dayTwo.Close) || (dayFour.Close >= dayThree.Close) {
		return
	}
	if (dayFour.Lw > 0) || (dayFour.Uw < (dayFour.Cb / 2)) {
		return
	}
	if (!Bullish(lastDay)) || (lastDay.Open <= dayFour.Close) || (lastDay.Cb <= dayFour.Cb) {
		return
	}
	

	pattern = structs.StockPattern{Name: "Ladder Bottom"}
		pattern.StripLines = append(pattern.StripLines, struct{
			StartValue int64     `json:"startValue"`
			EndValue   int64     `json:"endValue"`
			Color      string  `json:"color"`
			Label      string  `json:"label"`
			Opacity    float64 `json:"opacity"`
		}{StartValue: dayOne.StartTime, EndValue: lastDay.EndTime, Color: "#fcff4d", Label: "Ladder Bottom", Opacity: 0.4})
		return

}


func BullThreeLineStrike(days ...structs.DataPoint) (pattern structs.StockPattern) {
	initDay := CreateDay(days[0])
	dayOne := CreateDay(days[1])
	dayTwo := CreateDay(days[2])
	dayThree := CreateDay(days[3])
	dayFour := CreateDay(days[4])

	if (dayOne.Close <= dayOne.Open) || (dayTwo.Close <= dayTwo.Open) || (dayThree.Close <= dayThree.Open) {
		return
	}
	if (dayOne.Close <= initDay.Close) || (dayTwo.Close <= dayOne.Close) || (dayThree.Close <= dayTwo.Close) {
		return
	}
	if (dayOne.Low < initDay.Low) || (dayTwo.Low <= dayOne.Low) || (dayThree.Low <= dayTwo.Low){
		return
	}
	if ((dayOne.High <= initDay.High) || (dayTwo.High <= dayOne.High)) || ((dayThree.High <= dayTwo.High) || (dayFour.High < dayThree.High)) {
		return
	}
	if (dayFour.Close >= dayFour.Open) || (dayOne.Low < dayFour.Low) || (dayOne.Close < dayFour.Close) {
		return
	}

	pattern = structs.StockPattern{Name: "Bull Three Line Strike"}
		pattern.StripLines = append(pattern.StripLines, struct{
			StartValue int64     `json:"startValue"`
			EndValue   int64     `json:"endValue"`
			Color      string  `json:"color"`
			Label      string  `json:"label"`
			Opacity    float64 `json:"opacity"`
		}{StartValue: dayOne.StartTime, EndValue: dayFour.EndTime, Color: "#fcff4d", Label: "Bull Three Line Strike", Opacity: 0.4})
		return

}

func BearThreeLineStrike(days ...structs.DataPoint) (pattern structs.StockPattern) {
	initDay := CreateDay(days[0])
	dayOne := CreateDay(days[1])
	dayTwo := CreateDay(days[2])
	dayThree := CreateDay(days[3])
	dayFour := CreateDay(days[4])

	if (dayOne.Close >= dayOne.Open) || (dayTwo.Close >= dayTwo.Open) || (dayThree.Close >= dayThree.Open) {
		return
	}
	if (dayOne.Close >= initDay.Close) || (dayTwo.Close >= dayOne.Close) || (dayThree.Close >= dayTwo.Close) {
		return
	}
	if (dayOne.Low > initDay.Low) || (dayTwo.Low >= dayOne.Low) || (dayThree.Low >= dayTwo.Low) || (dayFour.Low > dayThree.Low){
		return
	}	
	if (dayOne.High >= initDay.High) || (dayTwo.High >= dayOne.High) || (dayThree.High >= dayTwo.High) {
		return
	}
	if (dayFour.Close <= dayFour.Open) || (dayOne.High > dayFour.High) || (dayOne.Close > dayFour.Close) {
		return
	}

	pattern = structs.StockPattern{Name: "Bear Three Line Strike"}
		pattern.StripLines = append(pattern.StripLines, struct{
			StartValue int64     `json:"startValue"`
			EndValue   int64     `json:"endValue"`
			Color      string  `json:"color"`
			Label      string  `json:"label"`
			Opacity    float64 `json:"opacity"`
		}{StartValue: dayOne.StartTime, EndValue: dayFour.EndTime, Color: "#fcff4d", Label: "Bear Three Line Strike", Opacity: 0.4})
		return

}



func ThreeWhiteSoldiers(days ...structs.DataPoint) (pattern structs.StockPattern) {
	initDay := CreateDay(days[0])
	dayOne := CreateDay(days[1])
	dayTwo := CreateDay(days[2])
	dayThree := CreateDay(days[3])
	upperWick := dayThree.High - dayThree.Close 
	

	if (Bullish(initDay)) || (!Bullish(dayOne)) {
		return
	}

	if (!Bullish(dayTwo)) || (dayTwo.Cb < dayOne.Cb) {
		return
	}
	if (dayTwo.Open <= dayOne.Open) || (dayTwo.Close <= dayOne.Close) {
		return
	}
	if (dayTwo.High <= dayOne.High) || (!Bullish(dayThree)) {
		return
	}
	if (dayThree.Open <= dayTwo.Open) || (dayThree.Close <= dayTwo.Close) {
		return
	}
	if (dayThree.High <= dayTwo.High) || (dayThree.Cb <= upperWick) {
		return
	}

	

	pattern = structs.StockPattern{Name: "Three White Soldiers"}
		pattern.StripLines = append(pattern.StripLines, struct{
			StartValue int64     `json:"startValue"`
			EndValue   int64     `json:"endValue"`
			Color      string  `json:"color"`
			Label      string  `json:"label"`
			Opacity    float64 `json:"opacity"`
		}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Three White Soldiers", Opacity: 0.4})
		return

	// fmt.Printf("initDay: %v\n", InitDay)
}

func ThreeBlackCrows(days ...structs.DataPoint) (pattern structs.StockPattern) {
	initDay := CreateDay(days[0])
	dayOne := CreateDay(days[1])
	dayTwo := CreateDay(days[2])
	dayThree := CreateDay(days[3])
	// upperWick := dayThree.High - dayThree.Close 
	

	if (!Bullish(initDay)) || (Bullish(dayOne)) {
		return
	}

	if (Bullish(dayTwo)) || (dayTwo.Cb < dayOne.Cb) || (Bullish(dayThree)) {
		return
	}
	if (dayOne.Lw >= dayOne.Cb) || (dayTwo.Lw >= dayTwo.Cb) || (dayThree.Lw >= dayThree.Cb) {
		return
	}
	if (dayTwo.Open >= dayOne.Open) || (dayTwo.Close >= dayOne.Close) || (dayTwo.High >= dayOne.High) {
		return
	}
	
	if (dayThree.Open >= dayTwo.Open) || (dayThree.Close >= dayTwo.Close) {
		return
	}
	if (dayThree.High >= dayTwo.High) || (dayThree.Cb <= dayThree.Uw) {
		return
	}

	

	pattern = structs.StockPattern{Name: "Three Black Crows"}
		pattern.StripLines = append(pattern.StripLines, struct{
			StartValue int64     `json:"startValue"`
			EndValue   int64     `json:"endValue"`
			Color      string  `json:"color"`
			Label      string  `json:"label"`
			Opacity    float64 `json:"opacity"`
		}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Three Black Crows", Opacity: 0.4})
		return

	// fmt.Printf("initDay: %v\n", InitDay)
}


func MorningStar(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])
	dayOneCbCap := dayTwo.Cb * 0.20
	// upperWickLimit := dayOne.Cb * 0.11

	if Bullish(dayOne) {
		return
	}
	if (dayTwo.Open >= dayOne.Close) || (dayTwo.Close >= dayOne.Open) {
		return
	}
	if dayTwo.Cb > dayOneCbCap {
		return
	}
	
	// if dayTwo.Uw > upperWickLimit {
	// 	return
	// }
	if (dayThree.Open <= dayTwo.Open) || (dayThree.Open <= dayTwo.Close) {
		return
	}
	if dayThree.Cb <= (dayTwo.Cb * 2) {
		return
	}

	pattern = structs.StockPattern{Name: "Morning Star"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Morning Star", Opacity: 0.4})
	return

}

func EveningStar(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])
	dayOneCbCap := dayOne.Cb * 0.20
	// upperWickLimit := dayOne.Cb * 0.11

	if !Bullish(dayOne) {
		return
	}
	if Bullish(dayTwo) {
		return
	}
	if (dayTwo.Cb >= dayOne.Cb) || (dayTwo.Cb > dayOneCbCap) {
		return
	}
	
	if (dayThree.Open >= dayTwo.Open) || (dayThree.Open >= dayTwo.Close) {
		return
	}
	if dayThree.Cb <= (dayTwo.Cb * 2) {
		return
	}

	pattern = structs.StockPattern{Name: "Evening Star"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Evening Star", Opacity: 0.4})
	return

}


func RisingThreeMethods(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	firstBear := CreateDay(days[1])
	secondBear := CreateDay(days[2])
	thirdBear := CreateDay(days[3])
	lastDay := CreateDay(days[4])
	paddedCbLimit := dayOne.Cb * 0.80

	if (dayOne.Close <= dayOne.Open) || (dayOne.Cb <= firstBear.Cb){
		return
	}
	if (firstBear.Close >= firstBear.Open) || (secondBear.Close >= secondBear.Open) || (thirdBear.Close >= thirdBear.Open) {
		return
	}
	if (dayOne.Close < firstBear.Open) || (dayOne.Close < secondBear.Open) || (dayOne.Close < thirdBear.Open) {
		return
	}
	if (dayOne.Low < firstBear.Low) || (dayOne.Low < secondBear.Low) || (dayOne.Low < thirdBear.Low){
		return
	}
	if (firstBear.Close < dayOne.Open) || (secondBear.Close < dayOne.Open) || (thirdBear.Close < dayOne.Open) {
		return
	}
	if (paddedCbLimit < firstBear.Cb) || (paddedCbLimit < secondBear.Cb) || (paddedCbLimit < thirdBear.Cb) {
		return
	}
	if (lastDay.Cb < firstBear.Cb) || (lastDay.Cb < secondBear.Cb) || (lastDay.Cb < thirdBear.Cb) {
		return
	}
	if (lastDay.Close <= firstBear.Open) || (lastDay.Close <= secondBear.Open) || (lastDay.Close <= thirdBear.Open) {
		return
	}

	pattern = structs.StockPattern{Name: "Rising Three Methods"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: lastDay.EndTime, Color: "#fcff4d", Label: "Rising Three Methods", Opacity: 0.4})
	return
}

func FallingThreeMethods(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	firstBull := CreateDay(days[1])
	secondBull := CreateDay(days[2])
	thirdBull := CreateDay(days[3])
	lastDay := CreateDay(days[4])
	paddedCbLimit := dayOne.Cb * 0.80

	if (dayOne.Close >= dayOne.Open) || (dayOne.Cb <= firstBull.Cb){
		return
	}
	if (firstBull.Close <= firstBull.Open) || (secondBull.Close <= secondBull.Open) || (thirdBull.Close <= thirdBull.Open) {
		return
	}
	if (dayOne.Close > firstBull.Open) || (dayOne.Close > secondBull.Open) || (dayOne.Close > thirdBull.Open) {
		return
	}
	if (dayOne.High < firstBull.High) || (dayOne.High < secondBull.High) || (dayOne.High < thirdBull.High){
		return
	}
	if (firstBull.Close > dayOne.Open) || (secondBull.Close > dayOne.Open) || (thirdBull.Close > dayOne.Open) {
		return
	}
	if (paddedCbLimit < firstBull.Cb) || (paddedCbLimit < secondBull.Cb) || (paddedCbLimit < thirdBull.Cb) {
		return
	}
	if (lastDay.Cb < firstBull.Cb) || (lastDay.Cb < secondBull.Cb) || (lastDay.Cb < thirdBull.Cb) {
		return
	}
	if (lastDay.Close >= firstBull.Open) || (lastDay.Close >= secondBull.Open) || (lastDay.Close >= thirdBull.Open) {
		return
	}

	pattern = structs.StockPattern{Name: "Rising Three Methods"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: lastDay.EndTime, Color: "#fcff4d", Label: "Rising Three Methods", Opacity: 0.4})
	return
}



func BullishMatHold(days ...structs.DataPoint) (pattern structs.StockPattern) {
	firstDay := CreateDay(days[0])
	secondDay := CreateDay(days[1])
	thirdDay := CreateDay(days[2])
	fourthDay := CreateDay(days[3])
	lastDay := CreateDay(days[4])

	if (firstDay.Open >= firstDay.Close) || (lastDay.Open >= lastDay.Close) {
		return
	}
	if (secondDay.Close >= secondDay.Open)  || (fourthDay.Close >= fourthDay.Open) {
		return
	}
	openLimit := firstDay.Close + (firstDay.Cb * 0.10)
	bodyMin := firstDay.Cb * 0.20
	if (secondDay.Open <= openLimit) || (secondDay.Close <= firstDay.Close) {
		return
	}
	if (secondDay.Cb < bodyMin) || (thirdDay.Cb < bodyMin) || (fourthDay.Cb < bodyMin) {
		return
	}
	if (secondDay.Cb >= firstDay.Cb) || (thirdDay.Cb >= firstDay.Cb) || (fourthDay.Cb >= firstDay.Cb){
		return
	}
	if thirdDay.Open > thirdDay.Close {
		if (thirdDay.Open >= secondDay.Open) || (thirdDay.Close >= secondDay.Close) {
			return
		}
		if (lastDay.Open >= thirdDay.Open) || (lastDay.Close >= thirdDay.Close) {
			return
		}

	} else {
		if (thirdDay.Open >= secondDay.Close) || (thirdDay.Close >= secondDay.Open) {
			return
		}
		if (fourthDay.Open >= thirdDay.Close) || (fourthDay.Close >= thirdDay.Open) {
			return
		}
	}
	if (secondDay.Open <= firstDay.Low) || (thirdDay.Open <= firstDay.Low) || (fourthDay.Open <= firstDay.Low){
		return
	}
	if (secondDay.Close <= firstDay.Low) || (thirdDay.Close <= firstDay.Low) || (fourthDay.Close <= firstDay.Low) {
		return
	}
	lastDayBodyMin := firstDay.Cb * 0.70
	if (lastDay.Close <= secondDay.High) || (lastDay.Cb <= lastDayBodyMin) {
		return
	}
	pattern = structs.StockPattern{Name: "Bullish Mat Hold"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: firstDay.StartTime, EndValue: lastDay.EndTime, Color: "#fcff4d", Label: "Bullish Mat Hold", Opacity: 0.4})
	return


}

func BearishAbandonedBaby(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])

	if (dayOne.Close < dayOne.Open) || (!IsDoji(dayTwo)) || (dayThree.Open < dayThree.Close) {
		return
	}
	dayOneUpperWick := dayOne.High - dayOne.Close
	dayOneLowerWick := dayOne.Open - dayOne.Low
	wickMax := dayOne.Cb * 0.15
	if (dayOneUpperWick > wickMax) || (dayOneLowerWick > wickMax) {
		return
	}
	dayThreeUpperWick := dayThree.High - dayThree.Close
	dayThreeLowerWick := dayThree.Open - dayThree.Low
	if (dayThreeUpperWick > wickMax) || (dayThreeLowerWick > wickMax) {
		return
	}
	if dayTwo.Open <= dayOne.Close {
		return
	}
	pattern = structs.StockPattern{Name: "Bearish Abnd. Baby"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Bearish Abnd. Baby", Opacity: 0.4})
	return

}

func ThreeStarsInTheSouth(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])

	if (dayOne.Close >= dayOne.Open) || (dayTwo.Close >= dayTwo.Open) || (dayThree.Close >= dayTwo.Open) {
		return
	}
	dayOneLowerWick := dayOne.Close - dayOne.Low
	lowerWickMin := func (cb float64) float64 {
		return cb * 0.30
	}
	dayTwoLowerWick := dayTwo.Close - dayTwo.Low
	
	if (dayOneLowerWick < lowerWickMin(dayOne.Cb)) || (dayTwoLowerWick < lowerWickMin(dayTwo.Cb)) {
		return
	}
	if (dayTwo.Open >= dayOne.Open) || (dayThree.Open >= dayTwo.Open) || (dayThree.High >= dayTwo.High) {
		return
	}
	if (dayTwo.Low <= dayOne.Low) || (dayThree.Low <= dayTwo.Low) || (dayThree.Cb >= dayTwo.Cb) {
		return
	}
	if (dayOne.High > (dayOne.Open + (dayOne.Open * 0.01))) || (dayTwo.High > (dayTwo.Open + (dayTwo.Open * 0.02))) || (dayThree.High > dayThree.Open) {
		return
	}
	if (dayTwo.Cb >= dayOne.Cb) || (dayOne.Low >= dayTwo.Low) || (dayThree.Low < dayThree.Close) {
		return
	}


	pattern = structs.StockPattern{Name: "Three stars In South"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Three stars In South", Opacity: 0.4})
	return

}


func ThreeInsideUp(days ...structs.DataPoint) (pattern structs.StockPattern) {
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	dayThree := CreateDay(days[2])

	if (Bullish(dayOne)) || (!Bullish(dayTwo)) || (!Bullish(dayThree)) {
		return
	}
	if (dayTwo.Close < (dayOne.Open + (dayOne.Cb / 2))) || (dayThree.Close < dayOne.High) {
		return
	}
	if (dayTwo.Close > dayOne.Open) || (dayTwo.Open > dayOne.Open) {
		return
	}


	pattern = structs.StockPattern{Name: "Three Inside Up"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayThree.EndTime, Color: "#fcff4d", Label: "Three Inside Up", Opacity: 0.4})
	return

}




func PiercingLine(days ...structs.DataPoint) (pattern structs.StockPattern){
	dayOne := CreateDay(days[0])
	dayTwo := CreateDay(days[1])
	

	if dayOne.Open < dayOne.Close {
		return
	}
	if dayTwo.Close < dayTwo.Open {
		return
	}
	if dayOne.Close <=  dayTwo.Open {
		return
	}
	closeLimit := dayOne.Close + (dayOne.Cb / 2)
	if dayTwo.Close <= closeLimit {
		return
	}
	if dayTwo.Close >= dayOne.Open {
		return
	}

	pattern = structs.StockPattern{Name: "Piercing Line"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: dayOne.StartTime, EndValue: dayTwo.EndTime, Color: "#fcff4d", Label: "Piercing Line", Opacity: 0.4})
	return
}




func Hammer(days ...structs.DataPoint) (pattern structs.StockPattern) {
	initDay := CreateDay(days[0])
	hammer := CreateDay(days[1])
	// lowerWick := math.Abs(hammer.Open - hammer.Close)
	

	if initDay.Open > initDay.Close {
		return
	} 
	var upperWick float64
	var lowerWick float64
	if hammer.Close < hammer.Open {
		upperWick = hammer.High - hammer.Open
		lowerWick = hammer.Close - hammer.Low
	} else {
		upperWick = hammer.High - hammer.Close
		lowerWick = hammer.Open - hammer.Low
	}
	upperWickLimit := hammer.Ct * 0.05
	if (hammer.Cb >= (lowerWick / 2.1)) || (upperWick >= upperWickLimit) {
		return
	}

	pattern = structs.StockPattern{Name: "Hammer"}
	pattern.StripLines = append(pattern.StripLines, struct{
		StartValue int64     `json:"startValue"`
		EndValue   int64     `json:"endValue"`
		Color      string  `json:"color"`
		Label      string  `json:"label"`
		Opacity    float64 `json:"opacity"`
	}{StartValue: hammer.StartTime, EndValue: hammer.EndTime, Color: "#fcff4d", Label: "Hammer", Opacity: 0.4})
	return

}




func IsDoji(day Day) bool {

	// dayOne := CreateDay(dayData)

	if day.Open == day.Close || day.Cb <= (day.Ct * 0.06) {
		// pattern = structs.StockPattern{Name: "Doji"}
		// pattern.StripLines = append(pattern.StripLines, struct{
		// 	StartValue int64     `json:"startValue"`
		// 	EndValue   int64     `json:"endValue"`
		// 	Color      string  `json:"color"`
		// 	Label      string  `json:"label"`
		// 	Opacity    float64 `json:"opacity"`
		// }{StartValue: dayOne.StartTime, EndValue: dayOne.EndTime, Color: "#fcff4d", Label: "Doji", Opacity: 0.4})
		// fmt.Printf("in IsDoji Stipline:: %+v\n", pattern.StripLines)
		return true
	}
	// err = fmt.Errorf("No patterns found")
	 return false
}



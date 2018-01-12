package utils

import (
	"fmt"
    "github.com/bitfinexcom/bitfinex-api-go/v2"
	"time"
    "bufio"
	"os"
	"strings"
)

var monnaies []string = []string{"tBTCUSD","tDSHUSD","tBCHUSD","tSANUSD"}
const nbmonnaies int = 4
var index int = 3


func RecupMonnaies() {

	// monnaies := []string{"tBTCUSD","tDSHUSD","tBCHUSD","tSANUSD"}
	var filename string
	var file *os.File
	var err error
	var w *bufio.Writer
	var monnaie string
	var Tick bitfinex.Ticker

	tickers,err := Clientv2.Tickers.Get(strings.Join(monnaies,","))

	if err != nil {
		fmt.Println("Error : ", err)
		fmt.Printf("Tickers - %s - %s\n",
					monnaies,time.Now().Format("2006-01-02;15:04:05.999"))
		return
	}
	
	nbtickers:=len(tickers)
	// fmt.Printf("nbtickers : %v \n",nbtickers)
	
	for i:=nbtickers-1;i>=0; i-- {
		// fmt.Printf("Ticker : %v \n",tickers[i])
		Tick = tickers[i]
		monnaie = strings.Split(Tick.Symbol,"t")[1]
		// fmt.Printf("monnaie : %v \n",monnaie)

		filename = fmt.Sprintf("TickerFile_%s_%s.log",monnaie,time.Now().Format("2006-01-02"))
		file, err = os.OpenFile(filename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("File does not exists or cannot be created - %s - %s\n",
						filename,time.Now().Format("2006-01-02;15:04:05.999"))
			os.Exit(1)
		}

		w = bufio.NewWriter(file)

		fmt.Fprintf(w,"%s;%s;%v;%v;%v;%v\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),
		Tick.Bid,Tick.Ask,Tick.LastPrice,Tick.Volume)
		
		w.Flush()
		file.Close()


	}

}
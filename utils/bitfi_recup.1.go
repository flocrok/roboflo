package utils

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"time"
    "bufio"
	"os"
)

var monnaies [4]string = [...]string{"BTCUSD","DSHUSD","BCHUSD","SANUSD"}
const nbmonnaies int = 4
var index int = 3


func RecupMonnaies() {

//	monnaies = [...]string{"BTCUSD","DSHUSD","BCHUSD","SANUSD"}
	var filename string
	var file *os.File
	var err error
	var w *bufio.Writer
	var monnaie string
	var Tick bitfinex.Tick

	index++
	if index == nbmonnaies {
		index = 0
	}
	monnaie = monnaies[index]

	Tick, err = Client.Ticker.Get(monnaie)
	
	if err != nil {
		fmt.Println("Error : ", err)
		fmt.Printf("Ticker - %s - %s\n",
					monnaie,time.Now().Format("2006-01-02;15:04:05.999"))
	}

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
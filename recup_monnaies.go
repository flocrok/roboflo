package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"time"
    "bufio"
	"os"
)

const sleeptime time.Duration = 15

func main() {

	monnaies := [...]string{"BTCUSD", "XRPUSD","DSHUSD","BCHUSD","SANUSD"}
	nbmonnaies := len(monnaies)

	var index int
	var filename string
	var file *os.File
	var err error
	var w *bufio.Writer
	var monnaie string
	var retry int
	var Tick bitfinex.Tick
	
	client := bitfinex.NewClient().Auth("fsmDUPTPmEuiBuanTsAWQ61GTyisqsAmpLfq9VVYJQ2", "CJ4gJ3W9mpKBYXEsVvhEwO5hiY5gdYqHNNoCBbTdWS7")
	info, errinfo := client.Account.Info()

	if errinfo != nil {
		fmt.Println("Error : ", errinfo)
	} else {
		fmt.Println("Info : ", info)
	}


	for true {

		// Ici, je recupere l'heure de cette iteration
		IterationTime := time.Now()

		index = 0
		for index < nbmonnaies {
			monnaie = monnaies[index]
			retry = 0
			for true {
				// TODO : creer une fonction dediee ?
				Tick, err = client.Ticker.Get(monnaie)
				
				if err != nil {
					fmt.Println("Error : ", err)
					fmt.Printf("Ticker - %s - %s\n",
								monnaie,time.Now().Format("2006-01-02;15:04:05.999"))
								
					time.Sleep(3 * time.Second)
					retry++
					if retry == 5 {// si trop d'erreur, on attend 2 minutes avant d'essayer a nouveau
						time.Sleep(120 * time.Second)
						IterationTime = time.Now()
						retry = 0
					}

				} else {break} 
			} /*else {
				fmt.Printf("%s ",Tick.Bid)
			}*/

			filename = fmt.Sprintf("output_%s_%s.csv",monnaie,time.Now().Format("2006-01-02"))
			file, err = os.OpenFile(filename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Printf("File does not exists or cannot be created - %s - %s\n",
							filename,time.Now().Format("2006-01-02;15:04:05.999"))
				os.Exit(1)
			}
			defer file.Close()
		
			w = bufio.NewWriter(file)

			fmt.Fprintf(w,"%s;%s;%s;%s;%s;%s\n",time.Now().Format("2006-01-02;15:04:05.999"),
							Tick.Timestamp,Tick.Bid,Tick.Ask,Tick.LastPrice,Tick.Volume)
			w.Flush()
			file.Close()
			index++

		}

		time.Sleep(time.Until(IterationTime.Add(sleeptime * time.Second)))
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999"))
}
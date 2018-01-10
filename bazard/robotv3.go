package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"../utils"
	"bytes"
	"strings"
)

var Client *bitfinex.Client
var monnaie string = "XRPUSD"

func main() {

	
	monnaies := [...]string{"BTCUSD","BCHUSD","BTGUSD","DSHUSD","ZECUSD","XRPUSD","ETHUSD"}
	var str bytes.Buffer
	for _,pair := range monnaies {
		pair = strings.ToUpper(pair)
		str.WriteString("t")
		str.WriteString(pair)
		str.WriteString(",")
	}

	fmt.Printf("debug 1\n")
	ret:=utils.Init(monnaie,0,4)
	if ret == 1 {return}

	fmt.Printf("debug 2\n")
	Client=utils.Clientv2
	fmt.Printf("debug 3 : %v\n",str.String())

	Tickers,err := Client.Ticker.Gets(str.String())
	if err != nil {
		fmt.Printf("Err : %v\n",err)
	} else {
		fmt.Printf("Ticks : %v\n",Tickers)		
	}
	for _,Ticks := range Tickers{
		
		fmt.Printf("Symbol : %v\n",Ticks.Symbol)	
	}

	_=err
//	1 : récupérer les valeurs des 7 monnaies
	// Bitfi : ok
	// CEX : TODO

// 2 : pour chaque "couple", calculer le gain potentiel (en fonction de la somme à dépenser)
// 3 choisir le couple le plus rentable
// 4: acheter sur bitfinex
// 5 : transférer sur cex
// 6 : scuter l arriver sur cex
// 7: récupérer le $
// 8 : acheter la 2ème monnaie
// 9 : transférer vers bitfinex
// 10 : récupérer les $
// 11 : constater le gain :D


	

	


}
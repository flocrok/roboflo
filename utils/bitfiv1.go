package utils

import (
	"strconv"
	"../trc"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
)

var Client *bitfinex.Client

func Initv1() int {

	Client = bitfinex.NewClient().Auth(Param.Clef, Param.Secret)
	info, errinfo := Client.Account.Info()

	if errinfo != nil {
		trc.Error.Println("Error : ", errinfo)
		return 1
	} else {
		trc.Info.Println("Info : ", info)
	}
	available,err := GetBalancev1("exchange","usd")
	trc.Info.Println("Available : ",available," ",err)

	return 0
}


func GetTicker(monnaie string) (MyTicker,error) {

	var montick MyTicker

	tick,err:=Client.Ticker.Get(monnaie)

	if err != nil {
		return montick,err
	}

	montick.Ask,err = strconv.ParseFloat(tick.Ask,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.Ask) ;return montick,err}
	montick.Bid,err = strconv.ParseFloat(tick.Bid,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.Bid) ;return montick,err}
	montick.High,err = strconv.ParseFloat(tick.High,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.High) ;return montick,err}
	montick.LastPrice,err = strconv.ParseFloat(tick.LastPrice,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.LastPrice);return montick,err}
	montick.Low,err = strconv.ParseFloat(tick.Low,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.Low) ;return montick,err}
	montick.Symbol = monnaie
	montick.Volume,err = strconv.ParseFloat(tick.Volume,64)
	if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",tick.Volume) ;return montick,err}
	
	return montick,err
}

func GetBalancev1(wallettype string, currency string) (float64,error) {

	Balance,err := Client.Balances.All()
	if err != nil {
		trc.Error.Println ("Balance Error : ",err)
		return 0.0,err
	} else {
		trc.Info.Println ("Balance info : ",Balance)
	}

	nbwallet:=len(Balance)
	
	for i:=0;i<nbwallet; i++ {
		if Balance[i].Type == wallettype && Balance[i].Currency == currency {
			
			available,err := strconv.ParseFloat(Balance[i].Available,64)
			if err !=nil { trc.Error.Println("Erreur de conversion (", err,") de ",Balance[i].Available) ;return 0.0,err}
			return available,err
		}
	}
// positionner une erreur
	//err = Error("Balance non trouvee"
	return 0.0,err
}

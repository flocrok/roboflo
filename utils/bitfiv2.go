package utils

import (
	"../trc"
	"time"
	"github.com/bitfinexcom/bitfinex-api-go/v2"
	"../sto"
)

var Clientv2 *bitfinex.Client

func Initv2(monnaie string,LogStoch int,) int {
	var retour int
	var chaine string

	// init du client 
	Clientv2 = bitfinex.NewClient().Credentials(Param.Clef, Param.Secret)
	//pROD / "fsmDUPTPmEuiBuanTsAWQ61GTyisqsAmpLfq9VVYJQ2", "CJ4gJ3W9mpKBYXEsVvhEwO5hiY5gdYqHNNoCBbTdWS7"
	// TEST : "ClPhRBhaSS2n4uqmcOebBE7CeeuHDGeapf2vQXwRUn0", "sjmX1WrSXtQi3Mf6aVbQg92YJCUAAkVu45Lgf6SqE3L"

	etat, err :=Clientv2.Platform.Status()
	if etat == true {
		trc.Info.Println("Plaform OK")
	} else {
		trc.Error.Println("Plaform KO")
		return 1
	}

	bougies,err := Clientv2.Candles.Get(monnaie,"1m","hist?limit=500")
	if err != nil {
		trc.Error.Println("Get Candles KO : ",err)
		return 1
	}
	nbbougies:=len(bougies)

	sto.Periode=Param.Periode
	sto.Freqstoch=Param.Frequence
	sto.MoyD=Param.Moyd
	
	for i:=nbbougies-1;i>=0; i-- {
		// Base sur 1 mesure toutes les 15 secondes, on stocke les 4 donnees de la minute, 1 pour chaque 15 secondes
		retour,chaine = sto.Stoch(bougies[i].Open)
		retour,chaine = sto.Stoch(bougies[i].Low)
		retour,chaine = sto.Stoch(bougies[i].High)
		retour,chaine = sto.Stoch(bougies[i].Close)
		_=retour
		_=chaine
		if retour < 1 && LogStoch == 1 {
			trc.Stochlog.Printf("%s;%s;%s\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),chaine)
		}
    }
	
	return 0
}
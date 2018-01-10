package main

import (
	"fmt"
	"time"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"./sto"
	"./utils"
	"./trc"
	"strings"
)

var jour string

const sleeptime time.Duration = 15
var Client *bitfinex.Client
var LogStoch = 1

var monnaie string = "XRPUSD"
var Param utils.Config

var Kutile,Dutile float64
var signalpresent int = 0
var signalvente int = 0
var signalachat int = 0
var etat int = 0
var ilfautacheter int = 0
var ilfautvendre int = 0
var Prixactuel float64
var PrixAchat float64
var PrixVente float64

func main() {

	var err error
	var retry int
	var Tick utils.MyTicker

	var Margevente float64
	var Margetheorique float64
	var MargeCumulee float64 = 0.0
	var Prixactuel float64
	var PrixAchat float64
	var PrixVente float64
	var nbventes int = 0
	var sommeinitiale float64 = 0.0
	var sommeavantachat float64 = 0.0
	var sommeapresvente float64 = 0.0
	var amount float64
	var ret int

	ret=utils.Init(monnaie,LogStoch,4)
	if ret == 1 {return}

	jour = time.Now().Format("2006-01-02")
	//Param=utils.Param
	monnaie=utils.Param.Monnaie
	
	crypto := strings.ToLower(string([]rune(monnaie)[0:3]))

	for true {
		
		// Ici, je recupere l'heure de cette iteration
		IterationTime := time.Now()

		retry = 0
		for true {
			// TODO : creer une fonction dediee ?
			Tick, err =utils.GetTicker(monnaie)
			
			if err != nil {
				fmt.Println("Get Ticker Error : ", err)
				trc.Error.Println("Get Ticker : ", err)
							
				time.Sleep(3 * time.Second)
				retry++
				if retry == 5 {// si trop d'erreur, on attend 2 minutes avant d'essayer a nouveau
					time.Sleep(120 * time.Second)
					IterationTime = time.Now()
					retry = 0
				}

			} else {break} 
		} 

		Prixactuel = Tick.Bid
		// appel de Stoch avec la nouvelle valeur
		retour,chaine := sto.Stoch(Prixactuel)
		
		trc.Tickerlog.Printf("%s;%s;%v;%v;%v;%v\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),
		Tick.Bid,Tick.Ask,Tick.LastPrice,Tick.Volume)
		
		//Kutile=sto.Faststoch[sto.Indice]
		//Dutile=sto.Fastmean[sto.Indice]
		Kutile=sto.Slowstoch[sto.Indice]
		Dutile=sto.Slowmean[sto.Indice]
		
		if retour < 1 {

			decision()
			
			if ilfautacheter == 1 {
				// j'achete au prix de vente le plus bas tick.ask
				PrixAchat = Tick.Ask
				etat = 1
				ilfautacheter = 0
				sommeavantachat,err=utils.GetBalancev1("exchange","usd")
				_=err // on ne fait rien de l'erreur
				if sommeinitiale == 0.0 {sommeinitiale = sommeavantachat}
				amount=sommeavantachat/Tick.Ask
				//_=amount
				utils.Ordre("achat",monnaie,Tick,amount)
				fmt.Printf("ACHAT:%v;%v;%v - %s\n",PrixAchat,sommeavantachat,amount,time.Now().Format("2006-01-02;15:04:05.999"))
				trc.Trace.Printf("ACHAT:%v;%v;%v - %s\n",PrixAchat,sommeavantachat,amount,time.Now().Format("2006-01-02;15:04:05.999"))
			}
	
			if ilfautvendre == 1 {
				// je vends au prix d'achat le plus haut tick.Bid
				PrixVente = Tick.Bid
				etat = 0
				ilfautvendre = 0
				//Ordre(sens string, monnaie string, Tick MyTicker, amount float64)
				amount,err=utils.GetBalancev1("exchange",crypto)
				utils.Ordre("vente",monnaie,Tick,-1.0*amount)
				
				sommeapresvente,err=utils.GetBalancev1("exchange","usd")

				Margevente = (sommeapresvente - sommeavantachat)/sommeavantachat
				MargeCumulee = (sommeapresvente - sommeinitiale)/sommeinitiale
				Margetheorique= (PrixVente*0.998 - PrixAchat*1.002)/(PrixAchat*1.002)
				trc.Trace.Printf("VENTE:Prix:%v - Amount:%v - Somme Recuperee:%v - Gain:%v\n",
					PrixVente,amount,sommeapresvente,sommeapresvente-sommeavantachat)
				trc.Trace.Printf("MARGE:Réelle:%v - Théorique:%v - Cumulee:%v\n",Margevente,Margetheorique,MargeCumulee)

				fmt.Printf("VENTE:Prix:%v - Amount:%v - Somme Recuperee:%v - Gain:%v - Date:%s\n",
					PrixVente,amount,sommeapresvente,sommeapresvente-sommeavantachat,time.Now().Format("2006-01-02;15:04:05.999"))
				fmt.Printf("MARGE:Réelle:%v - Théorique:%v - Cumulee:%v\n",Margevente,Margetheorique,MargeCumulee)

				MargeCumulee +=  Margevente/(PrixAchat*1.002)
				nbventes++
			}
			trc.Trace.Printf("Stoch :%s;%s;%s;%d\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),chaine,etat)
			
			if LogStoch == 1 {
				trc.Stochlog.Printf("%s;%s;%s;%d\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),chaine,etat)
				
			}

		}

		// Recuperation des donnees d'autres monnaies 
		utils.RecupMonnaies()

		// Mise en sommeil x secondes
		time.Sleep(time.Until(IterationTime.Add(sleeptime * time.Second)))

		// Si changement de jour, on change le logger
		newjour := time.Now().Format("2006-01-02")
		if newjour != jour {
			jour = newjour
			trc.Rotate(monnaie)
		}
	}
	
}
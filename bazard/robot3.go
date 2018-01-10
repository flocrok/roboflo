package main

import (
	"bufio"
	"fmt"
	"time"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"../sto"
	"../utils"
)

const sleeptime time.Duration = 15
var Client *bitfinex.Client
var LogStoch = 1
var wstoch *bufio.Writer

var monnaie string = "XRPUSD"
var nbusd float64 = 45.0

func main() {

	var err error
	var retry int
	var Tick utils.MyTicker

	var Margevente float64
	var MargeCumulee float64
	var Prixactuel float64
	var PrixAchat float64
	var PrixVente float64
	var nbventes int = 0
	var premiereligne int = 1
	var premierprix float64
	var etat int = 0
	var ilfautacheter int = 0
	var ilfautvendre int = 0
	var signalpresent int = 0
	var signalvente int = 0
	var signalachat int = 0
	var amount float64
	var ret int

	ret,wstoch=utils.Initv2(monnaie,LogStoch)
	if ret == 1 {return}

	ret=utils.Initv1()
	if ret == 1 {return}

	for true {
		
		// Ici, je recupere l'heure de cette iteration
		IterationTime := time.Now()

		retry = 0
		for true {
			// TODO : creer une fonction dediee ?
			Tick, err =utils.GetTicker(monnaie)
			
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

		Prixactuel = Tick.Bid
		if premiereligne == 1 {
			premiereligne = 0
			premierprix = Prixactuel
			_ = premierprix
		}

		// appel de Stoch avec la nouvelle valeur
		retour,chaine := sto.Stoch(Prixactuel)
		
		//fmt.Printf(" retour de stoch : %d\n",retour)
		Kutile:=sto.Faststoch[sto.Indice]
		Dutile:=sto.Fastmean[sto.Indice]

	    if signalachat == 1 || signalvente == 1 {
			signalpresent = 1
		}
		if retour < 1 {
			if etat == 0 { // etat 0, j'attend pour acheter
				// si les 2 sont en dessous de 20
				if Kutile < 20 && Dutile < 20 {
					// si K > D : j'achete
					if Kutile > Dutile {
						if signalachat == 1 {
							ilfautacheter = 1
							signalachat = 0
						} else {signalachat = 1}}
						// une regle plus elaboree pourrait être d'attendre que K repasse au dessus de 20
					// si K < D : je suis pret a acheter
				}
				// si D < 20 et K > 20
					// D etant la moyenne de K, a priori K etait < 20 donc on devrait acheter aussi
				if  Kutile > 20 && Dutile < 20 {
					if signalachat == 1 {
						ilfautacheter = 1
						signalachat = 0
					} else {signalachat = 1}}

				if signalachat == 1 && Kutile > Dutile {
						ilfautacheter = 1
						signalachat = 0
					}
			//	if  Dutile > 80 && Kutile > Dutile { ilfautacheter = 1}
				// Si D> 20 et K< 20
					// alors on a pas la condition d'attente du signale 
					// car si K repasse > 20 alors D ne sera pas passe en dessous de 20
					
			//	if Kutile > 50 && Dutile < (Kutile-10)  { ilfautacheter = 1}
			}
			if etat == 1 {
				if Kutile > 80 && Dutile > 80 {
					// si K > D : j'achete
					if Kutile < Dutile {
						if signalachat == 1 {
							ilfautvendre = 1
							signalvente = 0
						} else {signalvente = 1}}
						// une regle plus elaboree pourrait être d'attendre que K repasse au dessus de 20
					// si K < D : je suis pret a acheter
				}
				if  Kutile < 80 && Dutile > 80 {
					if signalachat == 1 {
						ilfautvendre = 1
						signalvente = 0
					} else {signalvente = 1}}

				
				if signalvente == 1 && Kutile < Dutile {
						ilfautvendre = 1
						signalvente = 0
					}
			//	if  Dutile < 20 && Kutile < Dutile { ilfautvendre = 1}
			//	if Kutile < 50 && Dutile > (Kutile+10)  { ilfautvendre = 1}

			}

			if signalpresent == 1 && signalvente == 1 { signalvente = 0}
			if signalpresent == 1 && signalachat == 1 { signalachat = 0}
			
			if ilfautacheter == 1 {
				// j'achete au prix de vente le plus bas tick.ask
				PrixAchat = Tick.Ask
				etat = 1
				ilfautacheter = 0
				amount=nbusd/Tick.Ask
				_=amount
				utils.Ordre("achat",monnaie,Tick,amount)
				fmt.Printf(" achat  : %v;%v;%v - %s\n",PrixAchat,nbusd,amount,time.Now().Format("2006-01-02;15:04:05.999"))
			}
	
			if ilfautvendre == 1 {
				// je vends au prix d'achat le plus haut tick.Bid
				PrixVente = Tick.Bid
				etat = 0
				ilfautvendre = 0
				//Ordre(sens string, monnaie string, Tick MyTicker, amount float64)
				utils.Ordre("vente",monnaie,Tick,-1.0*amount)
				Margevente = (PrixVente*0.998 - PrixAchat*1.002)
				fmt.Printf(" vente  : %v;%v - %s\n",PrixVente,amount,time.Now().Format("2006-01-02;15:04:05.999"))
				fmt.Printf(" marge  : %v\n",Margevente/PrixAchat)
				MargeCumulee +=  Margevente/(PrixAchat*1.002)
				nbventes++
			}
	
			if LogStoch == 1 {
				fmt.Printf("%s;%s;%s;%d\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),chaine,etat)
				fmt.Fprintf(wstoch,"%s;%s;%s;%d\n",monnaie,time.Now().Format("2006-01-02;15:04:05.999"),chaine,etat)
				wstoch.Flush()
			}

		}

		fmt.Printf("Ticker - %s - %s\n",
			monnaie,time.Now().Format("2006-01-02;15:04:05.999"))

		fmt.Printf("%v;%v;%v;%v;%v\n",time.Now().Format("2006-01-02;15:04:05.999"),
						Tick.Bid,Tick.Ask,Tick.LastPrice,Tick.Volume)

		time.Sleep(time.Until(IterationTime.Add(sleeptime * time.Second)))
	}
	
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999"))
}
package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"time"
	"strconv"	
    "bufio"
	"os"
	"./sto"
	"strings"
	"./trc"
)

// On prend des mesures toutes les x secondes. Defaut à 10 secondes
var sleeptime time.Duration = 10
var LogStoch = 0

var Vitesse string = "fast"
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

func simu(monnaie string)  (float64, float64) {

	var Margevente float64 = 0.0
	var MargeCumulee float64 = 0.0
	var Tick bitfinex.Tick
	//var err error
	var monsplit []string
	var nbventes int = 0
	var premiereligne int = 1
	var premierprix float64
	var inputfile string
	var wstoch *bufio.Writer
	var heurefic string
	var datefic string

	etat = 0
	ilfautacheter = 0
	ilfautacheter = 0
	signalpresent = 0
	signalvente = 0
	signalachat = 0
	
	//for jour := 6;jour < 11; jour++ {
	//	inputfile = fmt.Sprintf("data/output_%s_2017-12-%.2d.csv",monnaie,jour)
		inputfile = fmt.Sprintf("TickerFile_XRPUSD_2018-01-02.csv")
		//fmt.Println("fichier : ", inputfile)

		// Ouverture du fichier contenant les donnees de test
		file, err := os.Open(inputfile)
		if err != nil {
			fmt.Println("Error : ", err)
			return 0.0,0.0
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		if LogStoch == 1 {
			filenamestoch := fmt.Sprintf("stoch_%s_%d_%d_%d.csv",time.Now().Format("2006-01-02"),sto.Periode,sto.Freqstoch,sto.MoyD)
			filestoch, err := os.OpenFile(filenamestoch,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println("File does not exists or cannot be created")
				os.Exit(1)
			}
			defer filestoch.Close()

			wstoch = bufio.NewWriter(filestoch)
		}
		
		numligne := 0
		for scanner.Scan() {
			numligne++
			monsplit = strings.Split(scanner.Text(), ";")
			//fmt.Println(monsplit[3])

			//Tick, err := client.Ticker.Get("BTCUSD")
			datefic = monsplit[1]
			heurefic = monsplit[2]
			Tick.Bid = monsplit[3]
			Tick.Ask = monsplit[4]
			Prixactuel,err = strconv.ParseFloat(Tick.Bid, 64)
			if premiereligne == 1 {
				premiereligne = 0
				premierprix = Prixactuel
			}

			// appel de Stoch avec la nouvelle valeur
			retour,chaine := sto.Stoch(Prixactuel)
			sto.Macd(Prixactuel)
			
			//fmt.Printf(" retour de stoch : %d\n",retour)
			if Vitesse == "fast" {
				Kutile=sto.Faststoch[sto.Indice]
				Dutile=sto.Fastmean[sto.Indice]
			} else
			{
				Kutile=sto.Slowstoch[sto.Indice]
				Dutile=sto.Slowmean[sto.Indice]
			}

			if retour < 1 {
				
				decision()

				if ilfautacheter == 1 {
					// j'achete au prix de vente le plus bas tick.ask
					PrixAchat, err = strconv.ParseFloat(Tick.Ask, 64)
					etat = 1
					ilfautacheter = 0
					//fmt.Printf(" achat  : %v - %s\n",PrixAchat,monsplit[1])
				}
		
				if ilfautvendre == 1 {
					// je vends au prix d'achat le plus haut tick.Bid
					PrixVente, err = strconv.ParseFloat(Tick.Bid, 64)
					etat = 0
					ilfautvendre = 0
					Margevente = (PrixVente*0.998 - PrixAchat*1.002)
				//	fmt.Printf(" vente  : %v - %s\n",PrixVente,monsplit[1])
				//	fmt.Printf(" marge  : %v %s %s\n",Margevente/PrixAchat,monsplit[1], monsplit[2] )
					MargeCumulee +=  Margevente/(PrixAchat*1.002)
					nbventes++
				}
		
				if LogStoch == 1 {
					fmt.Fprintf(wstoch,"%s;%s;%s;%s;%d;%v;%v;%v;%v\n","SIMU",
						datefic,heurefic,chaine,etat,sto.Mme1,sto.Mme2,sto.Macd1,sto.MmeMacd)
					wstoch.Flush()
				}

			}

		//	time.Sleep(sleeptime * time.Second)

		}
//	}

	//	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999"))

	filename := fmt.Sprintf("result_simu_%s.csv",monnaie)
	fileout, errout := os.OpenFile(filename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if errout != nil {
		fmt.Printf("File does not exists or cannot be created - %s - %s\n",
					filename,time.Now().Format("2006-01-02;15:04:05.999"))
		os.Exit(1)
	}
	defer fileout.Close()

	wout := bufio.NewWriter(fileout)

	fmt.Fprintf(wout,"%v;%d;%d;%d;%d;%v\n",MargeCumulee,sto.Periode,sto.Freqstoch,sto.MoyD,nbsignalcible,Vitesse)
	wout.Flush()
	fileout.Close()
//	fmt.Printf("%v;%d;%d;%d;%d;%d;%v\n",MargeCumulee,nbventes,sto.Periode,sto.Freqstoch,sto.MoyD,nbsignalcible,Vitesse)
	//fmt.Printf(" nb de ventes : %d\n",nbventes)
	// progression sur la periode
	_ = premierprix
//	fmt.Printf(" premierprix : %v\n",premierprix)
//	fmt.Printf(" Prixactuel : %v\n",Prixactuel)
	//fmt.Printf(" progression sans rien faire : %v\n",(Prixactuel - premierprix)/premierprix)
	return MargeCumulee,((Prixactuel - premierprix)/premierprix)

}

func main() {

	monnaies := [...]string{"XRPUSD"}
//	monnaies := [...]string{"BTCUSD", "XRPUSD","DSHUSD","BCHUSD","SANUSD"}
//	monnaies := [...]string{"XRPUSD","DSHUSD","BCHUSD","SANUSD"}
	nbmonnaies := len(monnaies)

	ret := trc.Init(0,"")
	if ret == 1 {return }

   //periode : 6, 10, 14, 20, 24, 28, 32
   //frequence : 1, 2, 3, 5, 10
   //moyD : 2, 3, 4, 5, 6
   frequencemax := 0
   periodemax := 0
   moydmax := 0
   margemax := 0.0
   marge := 0.0
   margeglobal := 0.0
   margestd := 0.0
   vitessemax := ""

   index := 0
   LogStoch = 1
	frequencemax = 0
	periodemax = 0
	moydmax = 0
	marge = 0.0
	margestd = 0.0
	margemax = -1
	vitessetab := []string {"fast","slow"}
	vitessetab = []string {"slow"}

	//periode : 8 à 36 par 2
	for periode := 28;periode < 29;periode=periode+3 {
			sto.Periode = periode
			for moyD := 4;moyD < 5;moyD++ {
				sto.MoyD = moyD
				// frequence : 4 à 40 par 4
				for frequence := 20;frequence < 22;frequence=frequence+8 {
					sto.Freqstoch = frequence
					for _,vitesse := range vitessetab {
						Vitesse = vitesse
						margeglobal=0
						index=0
						for index < nbmonnaies {
							sto.Init()
						//	fmt.Printf("simu : %s - %d;%d;%d\n", monnaies[index],periode,moyD,frequence)
							marge,margestd = simu(monnaies[index])
							_=margestd
							margeglobal+=marge
							index++
						}
						
						if margeglobal > margemax {
							margemax=margeglobal
							frequencemax=frequence
							periodemax=periode
							moydmax=moyD
							vitessemax=vitesse
							
	//	fmt.Printf("%s;%v;%d;%d;%d;%v;%v\n","toutesmonnaies",margemax,periodemax,frequencemax,moydmax,vitessemax,margestd)
						}
					}
				}
			}
		}
		fmt.Printf("%s;marge : %v - marge std %v\n periode :%d; frequence : %d; moyenne : %d; vitesse %v",
			"toutesmonnaies",margemax,margestd,periodemax,frequencemax,moydmax,vitessemax)
	
	
}
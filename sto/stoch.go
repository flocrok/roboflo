package sto

import (
	"fmt"
	"../trc"
)

const Maxvalues int = 10000
var Periode int = 33
var MoyD int = 4
var MoyK int = 3
// on calcule la stochastique sur des fréquences de n mesures
//   cela permet d'avoir des min et max valables et une stochastique qui n'est pas trop souvent à 0 ou 100%
var Freqstoch int = 29
var indicefreq int = 0


var Hightab [Maxvalues]float64
var Lowtab [Maxvalues]float64
var Tickertab [Maxvalues]float64
var Faststoch [Maxvalues]float64
var Fastmean [Maxvalues]float64
var Slowstoch [Maxvalues]float64
var Slowmean [Maxvalues]float64
var Indice int = 0
var attenteDebut = 1

func Init() {
	indicefreq = 0
	Indice = 0
	attenteDebut = 1
	indicefreq = 0
}

func lomvalue(sens string,speed string) float64 {
	i :=Indice
	tabref:=Tickertab
	if speed == "slow" { tabref=Faststoch }
	
//	fmt.Printf("\n je suis dans lomvalue - %s - %s\n",sens, speed)
	valret := tabref[Indice]
	indicemin := Indice - (Periode*Freqstoch)
	if indicemin < 0 {indicemin = indicemin + Maxvalues}
	for i != indicemin {
		
//		fmt.Printf(" je suis dans lomvalue : %v - %d \n",tabref[i],i)
		if sens == "max" {if tabref[i] > valret {valret=tabref[i]}
		}else {if (tabref[i] < valret && tabref[i] != 0){valret=tabref[i]}}
		i--
		if i < 0 {i=i+Maxvalues}
	}
	return valret
}

func pourcent(speed string, moy int) float64 {
	
	tabref:=Faststoch
	if speed == "slowD" { tabref=Slowstoch }
	if speed == "slowK" { tabref=Faststoch }
	i :=Indice
	indicemin := Indice - moy*Freqstoch
	if indicemin < 0 {indicemin = indicemin + Maxvalues}

	valret := 0.0
	for i != indicemin {
		//fmt.Printf(" je suis dans pourcentd : %v - %d \n",valret,i)
		valret+=tabref[i]
		//fmt.Printf(" je suis dans pourcentd : %v - %d \n",valret,i)
		i-=Freqstoch
		if i < 0 {i=i+Maxvalues}
	}
	valret=valret/float64(moy)
	//fmt.Printf(" fin pourcentd : %v - %d \n",valret,i)
	return valret
}


func Stoch(value float64) (int,string) {

	var ret int = 0
	Indice++
	if Indice == Maxvalues { Indice = 0}
	indicefreq++
	if indicefreq == Freqstoch { indicefreq = 0}
//	fmt.Printf(" je suis dans stoche : %v\n",value)
	Tickertab[Indice]=value
	Hightab[Indice] = lomvalue("max","fast")
	Lowtab[Indice] = lomvalue("min","fast")

	if indicefreq == 0 {
		if Hightab[Indice] == Lowtab[Indice] {Faststoch[Indice]=50.0
		}else {
			Faststoch[Indice] = 100 * (Tickertab[Indice]-Lowtab[Indice])/(Hightab[Indice]-Lowtab[Indice])
		}
		Fastmean[Indice]=pourcent("fastK",MoyD)
		/*minstoch:=lomvalue("min","slow")
		maxstoch:=lomvalue("max","slow")
		if minstoch==maxstoch {Slowstoch[Indice]=50.0
			}else {
				Slowstoch[Indice] = 100 * (Faststoch[Indice]-minstoch)/(maxstoch-minstoch)
			}*/
		Slowstoch[Indice]=pourcent("slowK",MoyK)
		Slowmean[Indice]=pourcent("slowD",MoyD)
	} else { ret = 1}


	chaine := fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%d",
				value,Hightab[Indice],Lowtab[Indice],
				Faststoch[Indice],Fastmean[Indice],Slowstoch[Indice],Slowmean[Indice],Indice)
	/*fmt.Printf("%v;%v;%v;%v;%v;%v;%v;%d",
					value,Hightab[Indice],Lowtab[Indice],
					Faststoch[Indice],Fastmean[Indice],Slowstoch[Indice],Slowmean[Indice],Indice)*/
//	fmt.Printf(" chaine : %s\n",chaine)

	if attenteDebut == 1 {
		if Indice == (Periode)*Freqstoch {
			attenteDebut = 0
			trc.Info.Println("Stoch : ok, assez de donnees")
		} else {ret = 2}
	}
	
	return ret,chaine

}




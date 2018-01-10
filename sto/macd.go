package sto

import (
	"fmt"
)

var Macd1 float64 = 0.0
var MmeMacd float64 = 0.0
var Mme1 float64 = 0.0
var Mme2 float64 = 0.0
var MmePeriode1 int = 12
var MmePeriode2 int = 36
var MmeMacdPeriode int = 12
var Macdtab [Maxvalues]float64

func Mms(periode int, tab string) float64 {

	tabref:=Tickertab
	if tab == "macd" {
		tabref=Macdtab
	}
	i :=Indice
	indicemin := Indice - periode*Freqstoch
	if indicemin < 0 {indicemin = indicemin + Maxvalues}

	valret := 0.0
	for i != indicemin {
		//fmt.Printf(" je suis dans pourcentd : %v - %d \n",valret,i)
		valret+=tabref[i]
		//fmt.Printf(" je suis dans pourcentd : %v - %d \n",valret,i)
		i-=Freqstoch
		if i < 0 {i=i+Maxvalues}
	}
	valret=valret/float64(periode)
	//fmt.Printf(" fin pourcentd : %v - %d \n",valret,i)
	return valret
}

func Mme(value float64) int {

	var a float64
	if Mme1 == 0.0 {
		// si on a assez de donnees, on fait la moyenne
		if Indice == (MmePeriode1)*Freqstoch {
			Mme1=Mms(MmePeriode1,"")
		}
	} else {
		oldMme1 := Mme1
		a = 2/(float64(MmePeriode1)+1)
		Mme1 = oldMme1 * (1-a) + a * value
		fmt.Printf("a : %v - mm1 : %v, old : %v, ticker : %v\n",a, Mme1, oldMme1,value)
	}

	if Mme2 == 0.0 {
		// si on a assez de donnees, on fait la moyenne
		if Indice == (MmePeriode2)*Freqstoch {
			Mme2=Mms(MmePeriode2,"")
		}
	} else {
		oldMme2 := Mme2
		a = 2/(float64(MmePeriode2)+1)
		Mme2 = oldMme2 * (1-a) + a * value
	}

	return 0
}

func Macd(value float64) float64 {

	var a float64

	if indicefreq != 0 { return 0.0}

	Mme(value)
	//oldMacd := Macd1
	if Indice < ( MmePeriode2*Freqstoch) {
		Macd1 = 0.0
	} else {
		Macd1 = Mme1 - Mme2
	}
	Macdtab[Indice]=Macd1

	if MmeMacd == 0 {
		if Indice == (MmeMacdPeriode - 1)*Freqstoch + MmePeriode2*Freqstoch {
			// On fait la moyenne des premieres Macd
			MmeMacd=Mms(MmeMacdPeriode,"macd")
		}
	} else {
		oldMmeMacd := MmeMacd
		a = 2/(float64(MmeMacdPeriode)+1)
		MmeMacd = oldMmeMacd * (1-a) + a * Macd1
	}

	return 0.0
}

func MacdInit()  {
	Macd1 = 0.0
	MmeMacd = 0.0
	Mme1 = 0.0
	Mme2 = 0.0
}

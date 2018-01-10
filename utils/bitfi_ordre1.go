package utils

import (
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"time"
	"strconv"	
	"../trc"
)

const sleeptime time.Duration = 15

func Ordre(sens string, monnaie string, Tick MyTicker, amount float64) {
	
	var err error
	var retry int
	var PrixCible float64
	var remainingamount float64
	var id int64
	var Order *bitfinex.Order

	if sens != "achat" && sens != "vente" {return}

	if Param.Simu {
		trc.Info.Println("Mode Simulation, pas d'achat ou vente")
		return
	}

	trc.Trace.Printf("sens : %s \n",sens)

	//trc.Trace.Printf("%v;%v;%v;%v;%v\n",time.Now().Format("2006-01-02;15:04:05.999"),
	//				Tick.Bid,Tick.Ask,Tick.LastPrice,Tick.Volume)

	PrixVente:=Tick.Ask
	PrixAchat:=Tick.Bid
	if sens == "achat" {
		// d'abord je passe un ordre a un prix convenable : Valeur de vente la plus basse - diff(entre vente et achat)/4
		PrixCible=PrixVente-(PrixVente-PrixAchat)/4
	}
	
	if sens == "vente" {
		// d'abord je passe un ordre a un prix convenable : Valeur de vente la plus basse - diff(entre vente et achat)/4
		PrixCible=PrixAchat+(PrixVente-PrixAchat)/6
	}

// on passe l'ordre
	retry = 0
	for true {
		Order, err = Client.Orders.Create(monnaie,amount,PrixCible,"exchange limit")
		if err != nil {
			trc.Error.Println("Erreur de creation d'ordre : ", err)
			time.Sleep(1 * time.Second)
			retry++
			if retry == 5 {
				// si trop d'erreur, on abandonne s'il s'agit d'un achat
				// S'il s'agit d'une vente, pour le moment on abandonne ...
				return
			}
		} else {
			trc.Info.Printf("Info de l'ordre : %d;%s;%s;%s\n", Order.ID,Order.Price, Order.RemainingAmount, Order.Side)
			break
		}

	}
	waitorder := 0
	id = Order.ID
// on attend que l'ordre soit resolu
	for waitorder < 60 {
		remainingamount,err =strconv.ParseFloat(Order.RemainingAmount, 64)
		if remainingamount == 0.0 {break}
		time.Sleep(1*time.Second)
		waitorder++
		// on recupere le status de l'ordre
		*Order, err = Client.Orders.Status(id)
		if err != nil {
			trc.Error.Println("Erreur de status d'ordre : ", err)
		} else {
			trc.Info.Printf("Status de l'ordre : %d;%s;%s;%s\n", Order.ID,Order.Price, Order.RemainingAmount, Order.Side)
		}
	}
	
	if remainingamount == 0.0 {
		trc.Info.Printf("Ordre passe et resolu avec succes\n")
	} else {
		trc.Warning.Printf("ATTENTION : ordre limit non resolu => ordre market!\n")
		// D'abord, on annule le precedent
		retry = 0
		for true {
			err = Client.Orders.Cancel(id)
			if err != nil {
				trc.Error.Println("Erreur d'annulation d'ordre : ", err)
				time.Sleep(1 * time.Second)
				retry++
				if retry == 5 {
					// si trop d'erreur, on abandonne s'il s'agit d'un achat
					// S'il s'agit d'une vente, pour le moment on abandonne ...
					return
				}
			} else {break}
		}
		// Ensuite, on recupere la somme qui restait sur l'ordre
		retry = 0
		for true {
			*Order, err = Client.Orders.Status(id)
			if err != nil {
				trc.Error.Println("Erreur de recuperation de status : ", err)
				time.Sleep(1 * time.Second)
				retry++
				if retry == 5 {
					// si trop d'erreur, on abandonne s'il s'agit d'un achat
					// S'il s'agit d'une vente, pour le moment on abandonne ...
					return
				}
			} else {break}
		}

		// On attend que la somme soit de nouveau disponible après le cancel
		time.Sleep(10 * time.Second)

		// ensuite on passe l'ordre market de la somme restante
		remainingamount,err =strconv.ParseFloat(Order.RemainingAmount, 64)
			// le prix sera plus eleve, on prevoit de pouvoir en acheter un peu moins que le prix cible initial
		Order, err = Client.Orders.Create(monnaie,remainingamount*0.9,PrixCible,"exchange market")
		if err != nil {
			trc.Error.Println("Erreur de creation d'ordre : ", err)
			return
		} else {
			trc.Info.Printf("Info de l'ordre : %d;%s;%s;%s\n", Order.ID,Order.Price, Order.RemainingAmount, Order.Side)
		}
		// on attend que l'ordre soit resolu
		waitorder = 0
		id = Order.ID
		for waitorder < 60 {
			remainingamount,err =strconv.ParseFloat(Order.RemainingAmount, 64)
			if remainingamount == 0.0 {break}
			time.Sleep(1*time.Second)
			waitorder++
			// on recupere le status de l'ordre
			*Order, err = Client.Orders.Status(id)
			if err != nil {
				trc.Error.Println("Erreur de status d'ordre : ", err)
			} else {
				trc.Info.Printf("Status de l'ordre : %d;%s;%s;%s\n", Order.ID,Order.Price, Order.RemainingAmount, Order.Side)
			}
		}
	}
	if remainingamount != 0.0 {
		trc.Error.Printf("On a tout essayé, mais ça n'est pas passé !\n")
		return
	}
	// On peut aussi passer un ordre pour vendre si ça descend trop

}
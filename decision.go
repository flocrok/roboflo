package main

var diffKD float64 = 0.0
var nbsignalcible int = 1
var regle1 bool = true
var regle2 bool = false
var regle3 bool = false
var regle4 bool = false
var RatioProtect float64 = 0.98
var seuilhautD bool = false
var seuilbasD bool = false

var Difftendance float64  = 0.0

func decision () {

// 1ere règle basée sur le croisement des courbes K et D
	if regle1 {
		if etat == 0 { // etat 0, j'attend pour acheter
			//Si j'ai deja un signal d'achat, alors je regarde si K est toujours superieu a D
			if signalachat > 0 {
			//	fmt.Println(" signal d'achat' : ",signalachat)
				if Kutile > Dutile {
					signalachat++
			//		fmt.Println("nouveau signal d'achat' : ",signalachat)
				} else { signalachat = 0}
			} else {
				if Kutile > Dutile + diffKD && Dutile < 20 {signalachat = 1}
			}
			if signalachat == nbsignalcible  {
			//	fmt.Println("donc il faut acheter' : ",signalachat)
					ilfautacheter = 1
					signalachat = 0
			}
		}
		if etat == 1 {
			if signalvente > 0 {
				if Kutile < Dutile {
					signalvente++
				//	fmt.Println("nouveau signal de vente : ",signalvente)
				} else {
					signalvente = 0
				}
			} else {
				if Kutile < Dutile - diffKD && Dutile > 80  {signalvente = 1}
			}
			if signalvente == nbsignalcible  {
					ilfautvendre = 1
					signalvente = 0
			}
		}
	}

// 2ème règle : protection contre les baisses alors que l'on a acheté
// => les simulations indiquent que cette règle n'est pas très utile
// cela est probablement lié au fait que si on vend à cause d'une baisse, on va quand même racheter à la a
// légère hausse suivante alors que la tendance est toujours à la baisse.
	if regle2 {
		if etat == 1 && ilfautvendre == 0 {
			if Prixactuel < (PrixAchat * RatioProtect) && Kutile < Dutile{
			ilfautvendre = 1
			signalvente = 0
			//fmt.Println("regle2 : ",PrixAchat,Prixactuel,Kutile,Dutile)
			}
		}
		
	}

// 3ème règle : on vend si on descend en dessous de 80 et on achete si on passe au dessus de 20
// TODO : ajouter la notion de nombre de fois ou ça se produit
	if regle3 {
		if etat == 0 && ilfautacheter == 0 {
			if Dutile < 20 {
				seuilbasD = true
			}
			if seuilbasD && Kutile > 20 {
				ilfautacheter = 1
				seuilbasD = false
			}
		}
		
		if etat == 1 && ilfautvendre == 0 {
			if Dutile  > 80 {
				seuilhautD = true
			}
			if seuilhautD && Kutile < 80 {
				ilfautvendre = 1
				seuilhautD = false
			}
		}
	}

// 4ème règle : si franchissement
	if regle4 {
		if etat == 0 && ilfautacheter == 0 {
			if Dutile > 50 - Difftendance  {
				ilfautacheter = 1
			}
		}

		if etat == 1 && ilfautvendre == 0 {
			if Dutile < 50 + Difftendance {
				ilfautvendre = 1
			}
		}
	}


}
package main

import (
	"fmt"
)


func main() {
	

		fmt.Println("debut ")
		s := []string      {"no error", "Eio",  "invalid argument"}
	for index,valeur := range s {
		
		fmt.Printf("index : %d \n ",index)
		fmt.Printf("indevaleurx : %s \n ",valeur)

	}

	for periode := 0;periode < 40;periode = periode+4 {
		
		fmt.Printf("periode : %d \n ",periode)

	}
}
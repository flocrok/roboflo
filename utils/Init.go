package utils

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"../trc"
)

type Config struct  {
	Clef string
	Secret string
	Run string
	Monnaie string
	Frequence int
	Periode int
	Moyd int
	Simu bool
}

var Param Config

func Init(monnaie string,LogStoch int,level int) int {

	var ret int

	ret = LectureParam() 
	if ret == 1 {return 1}

	ret = trc.Init(level,monnaie)
	if ret == 1 {return 1}

	ret=Initv2(monnaie,LogStoch)
	if ret == 1 {return 1}

	ret = Initv1()
	if ret == 1 {return 1}
	
	trc.Info.Println("Init terminée")
	fmt.Println("Init terminée")

	return 0
}

func LectureParam() int {
	
	var monsplit []string
	Param.Simu = false

	ParamFile, err := os.Open("Robot.ini")
	if err != nil {return 1}
	
	scanner := bufio.NewScanner(ParamFile)
	
	for scanner.Scan() {
		
		//fmt.Println("ligne : ",scanner.Text())
		monsplit = strings.Split(scanner.Text(), "=")
		switch monsplit[0] {
		case "run" : 
			fmt.Println("run : ",monsplit[1])
			Param.Run=monsplit[1]
		case "periode" : 
			fmt.Println("periode : ",monsplit[1])
			Param.Periode,err=strconv.Atoi(monsplit[1])
		case "frequence" : 
			fmt.Println("frequence : ",monsplit[1])
			Param.Frequence,err=strconv.Atoi(monsplit[1])
		case "moyd" : 
				fmt.Println("moyd : ",monsplit[1])
				Param.Moyd,err=strconv.Atoi(monsplit[1])
		case "monnaie" : 
				fmt.Println("monnaie : ",monsplit[1])
				Param.Monnaie=monsplit[1]
		case "clef" : 
			//	fmt.Println("clef : ",monsplit[1])
				Param.Clef=monsplit[1]
		case "secret" : 
			//	fmt.Println("secret : ",monsplit[1])
				Param.Secret=monsplit[1]
		case "simu" : 
			if monsplit[1] == "true" || monsplit[1] == "O" || monsplit[1] == "1" {
				Param.Simu=true
				fmt.Println("simu : ",Param.Simu)
			}
		}
	}
	ParamFile.Close()
	if Param.Clef=="" || Param.Secret == "" {
		fmt.Printf("Erreur : Clef [%s] - Secret[%s]\n",Param.Clef,Param.Secret)
		return 1
	}
	return 0
}
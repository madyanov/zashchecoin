package bullshit

import "log"

func FailIf(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func WarnIf(err error) {
	if err != nil {
		log.Println(err)
	}
}

func PacnicIf(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

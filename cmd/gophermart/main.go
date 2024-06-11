package main

import (
	"github.com/rs/zerolog/log"

	"github.com/andranikuz/gophermart/internal/application"
)

func main() {
	a, err := application.NewApplication()
	if err != nil {
		log.Panic().Msg(err.Error())
		panic(err)
	}
	if err := a.Run(); err != nil {
		log.Panic().Msg(err.Error())
		panic(err)
	}
}

package handler

import "github.com/rs/zerolog/log"

func logErrorIfExists(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

package handler

import (
	expedition "goravel/app/models/expeditions"
)

func HandleExpediton(expedition_type string, resi string) expedition.Response {
	switch expedition_type {
	case "spx":
		return HandleSpx(resi)
	default:
		return expedition.Response{}
	}
}

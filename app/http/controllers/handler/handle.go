package handler

import (
	expedition "goravel/app/models/expeditions"
)

func HandleExpediton(expedition_type string, resi string) expedition.Response {
	switch expedition_type {
	case "spx":
		return HandleSpx(resi)
	case "tokopedia":
		return HandleTokopedia(resi)
	case "jnt":
		return HandleJNT(resi)
	case "jnt-cargo":
		return HandleJNTCargo(resi)
	case "jne":
		return HandleJNE(resi)
	case "sicepat":
		return HandleSicepat(resi)
	default:
		return expedition.Response{}
	}
}

package handlers

import "encoding/json"

func generateErrorResp(errorMsg string) []byte {
	type errorResp struct {
		Error string `json:"error"`
	}
	resp := errorResp{
		Error: errorMsg,
	}
	rawResp, _ := json.Marshal(resp)
	return rawResp
}

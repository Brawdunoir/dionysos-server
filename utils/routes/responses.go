package utils

import "github.com/Brawdunoir/dionysos-server/constants"

type ErrorResponse struct {
	Error string `json:"error"`
}

type URIResponse struct {
	URI string `json:"uri" example:"/api/v0/model/id"`
}

func CreateURIResponse(uri string) *URIResponse {
	return &URIResponse{
		URI: constants.BasePath + uri,
	}
}

func CreateErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{
		Error: error,
	}
}

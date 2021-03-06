package webapp

import (
	"net/http"
	"encoding/json"
	"gosha/webapp/types"
)

func ErrResponse(w http.ResponseWriter, err string, status int) {

	response := types.APIError{}

	response.Error = true
	response.ErrorMessage = err

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

	return
}

func ValidResponse(w http.ResponseWriter, data interface{}) {

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)

	return
}

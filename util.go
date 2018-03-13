package httputil

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/form"
)

type response struct {
	StatusCode int         `json:"status_code"`
	Messages   []string    `json:"messages"`
	Data       interface{} `json:"data"`
}

func marshalJSONResponse(statusCode int, messages []string, data interface{}) ([]byte, error) {
	resp := response{
		statusCode,
		messages,
		data,
	}
	bs, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func AcceptAllRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func DecodeFormRequest(r *http.Request, req interface{}) error {
	decoder := form.NewDecoder()
	r.ParseForm()
	err := decoder.Decode(&req, r.Form)
	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	}
	return nil
}

func DecodeJSONRequest(r *http.Request, req interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	switch {
	case err == io.EOF:
		return nil
	case err != nil:
		return err
	}
	return nil
}

func EncodeJSONResponse(resp interface{}) ([]byte, error) {
	return json.Marshal(&resp)
}

func SetContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetContentHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte(ErrInternalServerError.Error()))
}

func WriteDecodeRequestError(w http.ResponseWriter, exampleReq interface{}) {
	WriteResponse(w, 400, []string{ErrDecodeRequest.Error()}, exampleReq)
}

func WriteRedirectResponse(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(301)
}

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int, messages []string, data interface{}) {
	if err == ErrInternalServerError {
		WriteInternalServerError(w)
		return
	}
	WriteResponse(w, statusCode, messages, data)
}

func WriteResponse(w http.ResponseWriter, statusCode int, messages []string, data interface{}) {
	resp, err := marshalJSONResponse(statusCode, messages, data)
	if err != nil {
		WriteInternalServerError(w)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(resp)
}

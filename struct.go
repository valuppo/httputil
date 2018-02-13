package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ContentType int

const (
	JSON ContentType = iota
	Form
	HTML
)

type jsonResponse struct {
	StatusCode int         `json:"status_code"`
	Messages   []string    `json:"messages"`
	Data       interface{} `json:"data"`
}

type httputil struct {
	requestContentType ContentType
	appError           error
	isAcceptAllRequest bool
}

func (hu *httputil) SetApplicationError(err error) {
	hu.appError = err
}

func (hu *httputil) SetRequestContentType(contentType ContentType) {
	hu.requestContentType = contentType
}

func (hu *httputil) AcceptAllRequest(isAcceptAllRequest bool) {
	hu.isAcceptAllRequest = isAcceptAllRequest
}

func (hu *httputil) DecodeRequest(r *http.Request, req interface{}) {
	switch hu.requestContentType {
	case JSON:
		DecodeJSONRequest(r, req)
	case Form:
		DecodeFormRequest(r, req)
	}
}

func (hu *httputil) JSON(w http.ResponseWriter, err error, statusCode int, messages []string, data interface{}) {
	if hu.isAcceptAllRequest {
		AcceptAllRequest(w)
	}

	SetContentJSON(w)

	switch {
	case err == hu.appError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(hu.appJsonError())
	case err != nil:
		w.WriteHeader(statusCode)
		w.Write(hu.encodeJSONResponse(statusCode, []string{err.Error()}, nil))
	case err == nil:
		w.WriteHeader(http.StatusOK)
		w.Write(hu.encodeJSONResponse(statusCode, messages, data))
	}
}

func (hu *httputil) encodeJSONResponse(statusCode int, messages []string, data interface{}) []byte {
	resp := jsonResponse{
		StatusCode: statusCode,
		Messages:   []string{"Internal Server Error"},
		Data:       nil,
	}

	bs, err := json.Marshal(&resp)
	if err != nil {
		logrus.Error(err)
		return hu.appJsonError()
	}
	return bs
}

func (hu *httputil) appJsonError() []byte {
	resp := jsonResponse{
		StatusCode: http.StatusInternalServerError,
		Messages:   []string{"Internal Server Error"},
		Data:       nil,
	}

	bs, _ := json.Marshal(&resp)
	return bs
}

func New() *httputil {
	return &httputil{
		requestContentType: JSON,
		appError:           ErrInternalServerError,
		isAcceptAllRequest: true,
	}
}

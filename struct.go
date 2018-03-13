package httputil

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
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

type Util struct {
	requestContentType  ContentType
	responseContentType ContentType
	appError            error
	decodeRequestError  error
}

func (u *Util) SetApplicationError(err error) {
	u.appError = err
}

func (u *Util) SetDecodeRequestError(err error) {
	u.decodeRequestError = err
}

func (u *Util) SetRequestContentType(contentType ContentType) {
	u.requestContentType = contentType
}

func (u *Util) DecodeRequest(r *http.Request, req interface{}) error {
	switch u.requestContentType {
	case JSON:
		if err := DecodeJSONRequest(r, req); err != nil {
			logrus.Error(err)
			return u.decodeRequestError
		}
		return nil
	case Form:
		if err := DecodeFormRequest(r, req); err != nil {
			logrus.Error(err)
			return u.decodeRequestError
		}
		return nil
	default:
		return nil
	}
}

func (u *Util) DecodeValidateRequest(r *http.Request, req interface{}) (bool, error) {
	if err := u.DecodeRequest(r, req); err != nil {
		logrus.Error(err)
		return false, err
	}
	isValid, err := govalidator.ValidateStruct(req)
	if err != nil {
		logrus.Error(err)
		return isValid, err
	}
	return isValid, nil
}

func (u *Util) EncodeResponse(resp interface{}) ([]byte, error) {
	switch u.responseContentType {
	case JSON:
		return EncodeJSONResponse(resp)
	default:
		return nil, nil
	}
}

func (u *Util) ErrorJSON(w http.ResponseWriter, err error, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case err == u.appError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(u.appJsonError())
	case err == u.decodeRequestError:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(u.encodeJSONResponse(http.StatusBadRequest, []string{err.Error()}, data))
	default:
		w.WriteHeader(statusCode)
		w.Write(u.encodeJSONResponse(statusCode, []string{err.Error()}, data))
	}
}

func (u *Util) JSON(w http.ResponseWriter, statusCode int, messages []string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	w.Write(u.encodeJSONResponse(statusCode, messages, data))
}

func (u *Util) encodeJSONResponse(statusCode int, messages []string, data interface{}) []byte {
	resp := jsonResponse{
		StatusCode: statusCode,
		Messages:   messages,
		Data:       data,
	}

	bs, err := json.Marshal(&resp)
	if err != nil {
		log.Println(err)
		return u.appJsonError()
	}
	return bs
}

func (u *Util) appJsonError() []byte {
	resp := jsonResponse{
		StatusCode: http.StatusInternalServerError,
		Messages:   []string{"Internal Server Error"},
		Data:       nil,
	}

	bs, _ := json.Marshal(&resp)
	return bs
}

func New() *Util {
	return &Util{
		requestContentType:  JSON,
		responseContentType: JSON,
		appError:            ErrInternalServerError,
		decodeRequestError:  ErrDecodeRequest,
	}
}

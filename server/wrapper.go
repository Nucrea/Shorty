package server

import (
	"encoding/json"
	"net/http"
	"shorty/src/common/logging"

	"github.com/gin-gonic/gin"
)

type DummyInput struct{}

var _ error = (*ErrorBadRequest)(nil)
var _ error = (*ErrorInternal)(nil)

type ErrorBadRequest struct {
	Message string
}

func (e *ErrorBadRequest) Error() string {
	return e.Message
}

type ErrorInternal struct {
	Message string
}

func (e *ErrorInternal) Error() string {
	return e.Message
}

type Handler[Input, Output interface{}] func(ctx *gin.Context, input Input) (Output, error)

type ResponseOk struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}

type ResponseError struct {
	Status string `json:"status"`
	Error  struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	} `json:"error"`
}

func wrap[In, Out any](logger logging.Logger, handler Handler[In, Out]) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.WithContext(c)

		var input In
		contentType := c.Request.Header.Get("Content-Type")
		if contentType == "application/json" {
			if err := c.ShouldBindJSON(&input); err != nil {
				response := ResponseError{
					Status: "error",
					Error: struct {
						Id      string `json:"id"`
						Message string `json:"message"`
					}{
						Id:      "WrongBody",
						Message: err.Error(),
					},
				}

				body, err := json.Marshal(response)
				if err != nil {
					log.Error().Err(err).Msg("bind request body error")
				}
				c.Data(400, "application/json", body)
				return
			}
		}

		var response interface{}

		status := http.StatusOK
		output, err := handler(c, input)
		if err != nil {
			switch err.(type) {
			case *ErrorBadRequest:
				status = http.StatusBadRequest
			case *ErrorInternal:
				status = http.StatusInternalServerError
			}

			log.Error().Err(err).Msg("error in request handler")
			response = ResponseError{
				Status: "error",
				Error: struct {
					Id      string `json:"id"`
					Message string `json:"message"`
				}{
					Id:      "-",
					Message: err.Error(),
				},
			}
		} else {
			var empty Out
			if interface{}(output) == interface{}(empty) {
				status = http.StatusNoContent
			}
			response = ResponseOk{
				Status: "success",
				Result: output,
			}
		}

		body, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("marshal response error")
			c.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		c.Data(status, "application/json", body)
	}
}

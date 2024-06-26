package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrNotFound            = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
	ErrBadRequest          = &ErrResponse{HTTPStatusCode: 400, StatusText: "Bad request"}
	ErrInternalServerError = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal Server Error"}
)

func Error(msg string, code int) ErrResponse {
	return ErrResponse{
		HTTPStatusCode: code,
		StatusText:     msg,
	}
}

//
//func (e *ErrResponse) JSON(w http.ResponseWriter, r *http.Request) {
//	render.JSON(w, r, e)
//}

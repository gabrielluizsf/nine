package server

import "net/http"

type Error struct {
	StatusCode  int
	ContentType string
	Err         error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.Err != nil {
		w.Header().Set("Content-Type", e.ContentType)

		if e.ContentType == "application/json" {
			if e.StatusCode >= 100 {
				w.WriteHeader(e.StatusCode)
			}
			b, err := JSON{
				"err": e.Err.Error(),
			}.Bytes()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}

		http.Error(w, e.Err.Error(), e.StatusCode)
		return
	}
}

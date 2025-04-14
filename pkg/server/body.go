package server

import "github.com/i9si-sistemas/nine/internal/json"

// Body decodes the JSON body of an HTTP request into a provided variable.
//
//	var body bodyType
//	if err := nine.Body(req, &body); err != nil {
//	    return res.Status(http.StatusBadRequest).JSON(nine.JSON{
//			"message": "invalid body"
//		})
//	}
func Body[T any](req *Request, v *T) error {
	return json.Decode(req.Body().Bytes(), v)
}

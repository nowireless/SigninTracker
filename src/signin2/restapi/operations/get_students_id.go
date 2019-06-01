// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetStudentsIDHandlerFunc turns a function with the right signature into a get students ID handler
type GetStudentsIDHandlerFunc func(GetStudentsIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetStudentsIDHandlerFunc) Handle(params GetStudentsIDParams) middleware.Responder {
	return fn(params)
}

// GetStudentsIDHandler interface for that can handle valid get students ID params
type GetStudentsIDHandler interface {
	Handle(GetStudentsIDParams) middleware.Responder
}

// NewGetStudentsID creates a new http.Handler for the get students ID operation
func NewGetStudentsID(ctx *middleware.Context, handler GetStudentsIDHandler) *GetStudentsID {
	return &GetStudentsID{Context: ctx, Handler: handler}
}

/*GetStudentsID swagger:route GET /students/{id} getStudentsId

Get a particular student

*/
type GetStudentsID struct {
	Context *middleware.Context
	Handler GetStudentsIDHandler
}

func (o *GetStudentsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetStudentsIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
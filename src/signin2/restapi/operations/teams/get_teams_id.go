// Code generated by go-swagger; DO NOT EDIT.

package teams

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetTeamsIDHandlerFunc turns a function with the right signature into a get teams ID handler
type GetTeamsIDHandlerFunc func(GetTeamsIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTeamsIDHandlerFunc) Handle(params GetTeamsIDParams) middleware.Responder {
	return fn(params)
}

// GetTeamsIDHandler interface for that can handle valid get teams ID params
type GetTeamsIDHandler interface {
	Handle(GetTeamsIDParams) middleware.Responder
}

// NewGetTeamsID creates a new http.Handler for the get teams ID operation
func NewGetTeamsID(ctx *middleware.Context, handler GetTeamsIDHandler) *GetTeamsID {
	return &GetTeamsID{Context: ctx, Handler: handler}
}

/*GetTeamsID swagger:route GET /teams/{id} Teams getTeamsId

Get a particular team

*/
type GetTeamsID struct {
	Context *middleware.Context
	Handler GetTeamsIDHandler
}

func (o *GetTeamsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTeamsIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
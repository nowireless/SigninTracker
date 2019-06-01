// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PatchStudentsIDHandlerFunc turns a function with the right signature into a patch students ID handler
type PatchStudentsIDHandlerFunc func(PatchStudentsIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PatchStudentsIDHandlerFunc) Handle(params PatchStudentsIDParams) middleware.Responder {
	return fn(params)
}

// PatchStudentsIDHandler interface for that can handle valid patch students ID params
type PatchStudentsIDHandler interface {
	Handle(PatchStudentsIDParams) middleware.Responder
}

// NewPatchStudentsID creates a new http.Handler for the patch students ID operation
func NewPatchStudentsID(ctx *middleware.Context, handler PatchStudentsIDHandler) *PatchStudentsID {
	return &PatchStudentsID{Context: ctx, Handler: handler}
}

/*PatchStudentsID swagger:route PATCH /students/{id} patchStudentsId

Update a particular student

*/
type PatchStudentsID struct {
	Context *middleware.Context
	Handler PatchStudentsIDHandler
}

func (o *PatchStudentsID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPatchStudentsIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
// Code generated by go-swagger; DO NOT EDIT.

package teams

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	validate "github.com/go-openapi/validate"

	models "signin2/models"
)

// GetTeamsHandlerFunc turns a function with the right signature into a get teams handler
type GetTeamsHandlerFunc func(GetTeamsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTeamsHandlerFunc) Handle(params GetTeamsParams) middleware.Responder {
	return fn(params)
}

// GetTeamsHandler interface for that can handle valid get teams params
type GetTeamsHandler interface {
	Handle(GetTeamsParams) middleware.Responder
}

// NewGetTeams creates a new http.Handler for the get teams operation
func NewGetTeams(ctx *middleware.Context, handler GetTeamsHandler) *GetTeams {
	return &GetTeams{Context: ctx, Handler: handler}
}

/*GetTeams swagger:route GET /teams Teams getTeams

Returns a collection of teams

Collection of teams registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.

*/
type GetTeams struct {
	Context *middleware.Context
	Handler GetTeamsHandler
}

func (o *GetTeams) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTeamsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetTeamsOKBody get teams o k body
// swagger:model GetTeamsOKBody
type GetTeamsOKBody struct {

	// at meta id
	// Required: true
	AtMetaID models.ID `json:"@meta.id"`

	// members
	// Required: true
	Members []*models.Team `json:"Members"`
}

// Validate validates this get teams o k body
func (o *GetTeamsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAtMetaID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateMembers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTeamsOKBody) validateAtMetaID(formats strfmt.Registry) error {

	if err := o.AtMetaID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("getTeamsOK" + "." + "@meta.id")
		}
		return err
	}

	return nil
}

func (o *GetTeamsOKBody) validateMembers(formats strfmt.Registry) error {

	if err := validate.Required("getTeamsOK"+"."+"Members", "body", o.Members); err != nil {
		return err
	}

	for i := 0; i < len(o.Members); i++ {
		if swag.IsZero(o.Members[i]) { // not required
			continue
		}

		if o.Members[i] != nil {
			if err := o.Members[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getTeamsOK" + "." + "Members" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetTeamsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetTeamsOKBody) UnmarshalBinary(b []byte) error {
	var res GetTeamsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

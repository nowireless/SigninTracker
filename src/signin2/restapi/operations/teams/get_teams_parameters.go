// Code generated by go-swagger; DO NOT EDIT.

package teams

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTeamsParams creates a new GetTeamsParams object
// with the default values initialized.
func NewGetTeamsParams() GetTeamsParams {

	var (
		// initialize parameters with default values

		expandDefault = bool(true)
	)

	return GetTeamsParams{
		Expand: &expandDefault,
	}
}

// GetTeamsParams contains all the bound params for the get teams operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetTeams
type GetTeamsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Expand the links
	  In: query
	  Default: true
	*/
	Expand *bool
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetTeamsParams() beforehand.
func (o *GetTeamsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qExpand, qhkExpand, _ := qs.GetOK("expand")
	if err := o.bindExpand(qExpand, qhkExpand, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindExpand binds and validates parameter Expand from query.
func (o *GetTeamsParams) bindExpand(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetTeamsParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("expand", "query", "bool", raw)
	}
	o.Expand = &value

	return nil
}

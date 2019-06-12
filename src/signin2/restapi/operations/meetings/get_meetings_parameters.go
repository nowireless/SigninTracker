// Code generated by go-swagger; DO NOT EDIT.

package meetings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetMeetingsParams creates a new GetMeetingsParams object
// with the default values initialized.
func NewGetMeetingsParams() GetMeetingsParams {

	var (
		// initialize parameters with default values

		expandDefault = bool(true)
	)

	return GetMeetingsParams{
		Expand: &expandDefault,
	}
}

// GetMeetingsParams contains all the bound params for the get meetings operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetMeetings
type GetMeetingsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Meetings after this date
	  In: query
	*/
	AfterDate *strfmt.Date
	/*Meetings before this date
	  In: query
	*/
	BeforeDate *strfmt.Date
	/*Expand the links
	  In: query
	  Default: true
	*/
	Expand *bool
	/*Limit meetings to a specific team. Use database ID
	  In: query
	*/
	Teamid *int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetMeetingsParams() beforehand.
func (o *GetMeetingsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qAfterDate, qhkAfterDate, _ := qs.GetOK("afterDate")
	if err := o.bindAfterDate(qAfterDate, qhkAfterDate, route.Formats); err != nil {
		res = append(res, err)
	}

	qBeforeDate, qhkBeforeDate, _ := qs.GetOK("beforeDate")
	if err := o.bindBeforeDate(qBeforeDate, qhkBeforeDate, route.Formats); err != nil {
		res = append(res, err)
	}

	qExpand, qhkExpand, _ := qs.GetOK("expand")
	if err := o.bindExpand(qExpand, qhkExpand, route.Formats); err != nil {
		res = append(res, err)
	}

	qTeamid, qhkTeamid, _ := qs.GetOK("teamid")
	if err := o.bindTeamid(qTeamid, qhkTeamid, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAfterDate binds and validates parameter AfterDate from query.
func (o *GetMeetingsParams) bindAfterDate(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date
	value, err := formats.Parse("date", raw)
	if err != nil {
		return errors.InvalidType("afterDate", "query", "strfmt.Date", raw)
	}
	o.AfterDate = (value.(*strfmt.Date))

	if err := o.validateAfterDate(formats); err != nil {
		return err
	}

	return nil
}

// validateAfterDate carries on validations for parameter AfterDate
func (o *GetMeetingsParams) validateAfterDate(formats strfmt.Registry) error {

	if err := validate.FormatOf("afterDate", "query", "date", o.AfterDate.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindBeforeDate binds and validates parameter BeforeDate from query.
func (o *GetMeetingsParams) bindBeforeDate(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date
	value, err := formats.Parse("date", raw)
	if err != nil {
		return errors.InvalidType("beforeDate", "query", "strfmt.Date", raw)
	}
	o.BeforeDate = (value.(*strfmt.Date))

	if err := o.validateBeforeDate(formats); err != nil {
		return err
	}

	return nil
}

// validateBeforeDate carries on validations for parameter BeforeDate
func (o *GetMeetingsParams) validateBeforeDate(formats strfmt.Registry) error {

	if err := validate.FormatOf("beforeDate", "query", "date", o.BeforeDate.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindExpand binds and validates parameter Expand from query.
func (o *GetMeetingsParams) bindExpand(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetMeetingsParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("expand", "query", "bool", raw)
	}
	o.Expand = &value

	return nil
}

// bindTeamid binds and validates parameter Teamid from query.
func (o *GetMeetingsParams) bindTeamid(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("teamid", "query", "int64", raw)
	}
	o.Teamid = &value

	return nil
}
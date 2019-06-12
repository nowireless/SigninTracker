// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Day day
// swagger:model Day
type Day struct {

	// day
	// Maximum: 31
	// Minimum: 1
	Day int64 `json:"Day,omitempty"`

	// month
	// Maximum: 12
	// Minimum: 1
	Month int64 `json:"Month,omitempty"`

	// year
	Year int64 `json:"Year,omitempty"`
}

// Validate validates this day
func (m *Day) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDay(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMonth(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Day) validateDay(formats strfmt.Registry) error {

	if swag.IsZero(m.Day) { // not required
		return nil
	}

	if err := validate.MinimumInt("Day", "body", int64(m.Day), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("Day", "body", int64(m.Day), 31, false); err != nil {
		return err
	}

	return nil
}

func (m *Day) validateMonth(formats strfmt.Registry) error {

	if swag.IsZero(m.Month) { // not required
		return nil
	}

	if err := validate.MinimumInt("Month", "body", int64(m.Month), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("Month", "body", int64(m.Month), 12, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Day) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Day) UnmarshalBinary(b []byte) error {
	var res Day
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
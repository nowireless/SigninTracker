// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// IDRef A reference to a resource.
// swagger:model idRef
type IDRef struct {

	// at meta id
	AtMetaID ID `json:"@meta.id,omitempty"`
}

// Validate validates this id ref
func (m *IDRef) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAtMetaID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *IDRef) validateAtMetaID(formats strfmt.Registry) error {

	if swag.IsZero(m.AtMetaID) { // not required
		return nil
	}

	if err := m.AtMetaID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("@meta.id")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IDRef) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IDRef) UnmarshalBinary(b []byte) error {
	var res IDRef
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

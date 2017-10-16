// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// APILabelPair api label pair
// swagger:model apiLabelPair

type APILabelPair struct {

	// Label name
	Name string `json:"name,omitempty"`

	// Label value
	Value string `json:"value,omitempty"`
}

/* polymorph apiLabelPair name false */

/* polymorph apiLabelPair value false */

// Validate validates this api label pair
func (m *APILabelPair) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *APILabelPair) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APILabelPair) UnmarshalBinary(b []byte) error {
	var res APILabelPair
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

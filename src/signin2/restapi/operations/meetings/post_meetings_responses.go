// Code generated by go-swagger; DO NOT EDIT.

package meetings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "signin2/models"
)

// PostMeetingsCreatedCode is the HTTP code returned for type PostMeetingsCreated
const PostMeetingsCreatedCode int = 201

/*PostMeetingsCreated Meeting successfully created

swagger:response postMeetingsCreated
*/
type PostMeetingsCreated struct {
}

// NewPostMeetingsCreated creates PostMeetingsCreated with default headers values
func NewPostMeetingsCreated() *PostMeetingsCreated {

	return &PostMeetingsCreated{}
}

// WriteResponse to the client
func (o *PostMeetingsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(201)
}

/*PostMeetingsDefault Unexpected error

swagger:response postMeetingsDefault
*/
type PostMeetingsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostMeetingsDefault creates PostMeetingsDefault with default headers values
func NewPostMeetingsDefault(code int) *PostMeetingsDefault {
	if code <= 0 {
		code = 500
	}

	return &PostMeetingsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post meetings default response
func (o *PostMeetingsDefault) WithStatusCode(code int) *PostMeetingsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post meetings default response
func (o *PostMeetingsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post meetings default response
func (o *PostMeetingsDefault) WithPayload(payload *models.Error) *PostMeetingsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post meetings default response
func (o *PostMeetingsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostMeetingsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
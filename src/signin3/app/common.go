package app

import (
	"io/ioutil"
	"net/http"
	"signin3/database"
	"signin3/models"
	"signin3/tags"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
)

type handlerContext struct {
	w http.ResponseWriter
	r *http.Request

	db *database.Database

	marshalJSON   func(obj models.Model) []byte
	unmarshalJSON func(raw []byte) (models.Model, error)
}

func handlePost(ctx handlerContext) {
	// API Start (1)
	// TODO: Somwhere check headers of the POST request for JSON
	requestBody, err := ioutil.ReadAll(ctx.r.Body)
	if err != nil {
		panic(err)
	}
	// requestValue := reflect.Zero(reflect.TypeOf(objType))
	// requestRaw := requestValue.Interface()
	request, err := ctx.unmarshalJSON(requestBody)
	// log.Info("RequestRaw Type: ", fmt.Sprintf("%T ", requestRaw), reflect.ValueOf(requestRaw).Type())
	// err = json.Unmarshal(requestBody, &requestRaw)
	if err != nil {
		log.Error(err)
		MalformedJSON(ctx.w, ctx.r)
		return
	}
	// log.Info("RequestRaw Type: ", fmt.Sprintf("%T ", requestRaw), reflect.ValueOf(requestRaw).Type())

	spew.Dump(request)
	// API End (1)

	// Application LOGIC start
	missing := tags.CheckRequiredOnCreate(request)
	if len(missing) > 0 {
		MissingRequiredOnCreate(ctx.w, ctx.r, missing)
		return
	}

	// Note: request needs to be a pointer to model.*
	err = ctx.db.Create(request)
	if err != nil {
		// TODO the following should maybe move to database package?
		// Create a custom error struct for friendlier error handle
		if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.Code == "23505" {
			e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
			writeError(ctx.w, ctx.r, e)
			return
		}
		InternalError(ctx.w, ctx.r, err, "Database Error")
		return
	}

	requestModel := request.(models.Model)
	result, err := ctx.db.Get(request, requestModel.GetDatabaseID())
	if err != nil {
		InternalError(ctx.w, ctx.r, err, "Database Error")
		return
	}

	// Application LOGIC end

	// API Start (2)
	// Does a header need to created to store the new resources location?
	writeStruct(ctx.w, http.StatusCreated, result)
	// API End (2)
}

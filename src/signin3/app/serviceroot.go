package app

import (
	"net/http"
	"signin3/constants"
	"signin3/models"
)

func (app *App) ServiceRoot(w http.ResponseWriter, r *http.Request) {
	root := map[string]interface{}{}
	root["@uri"] = r.RequestURI
	root["People"] = models.Link{URI: constants.PeopleCollection}
	root["Meetings"] = models.Link{URI: constants.MeetingsCollection}
	root["Teams"] = models.Link{URI: constants.TeamsCollection}

	writeStruct(w, http.StatusOK, root)
}

package internal

import (
	"database/sql"
	"fmt"
	"signin3/models"

	"github.com/jackc/pgx/pgtype"
)

func makeIDLinks(collectionURI string, ids pgtype.Int4Array) []models.Link {
	result := []models.Link{}
	for _, id := range ids.Elements {
		if id.Status != pgtype.Present {
			panic("Element not present")
		}

		link := models.Link{}
		link.URI = fmt.Sprintf("%s/%d", collectionURI, id.Int)

		result = append(result, link)
	}

	return result
}

func makeLink(collectionURI string, id int) models.Link {
	link := models.Link{}
	link.URI = fmt.Sprintf("%s/%d", collectionURI, id)
	return link
}

func getNullString(value sql.NullString) *string {
	if value.Valid {
		str := value.String
		return &str
	}

	return nil
}

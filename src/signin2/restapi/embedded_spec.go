// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Sign In tracker API",
    "version": "1.0.0"
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "paths": {
    "/meetings": {
      "get": {
        "description": "Collection of meetings registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "Meetings"
        ],
        "summary": "Returns a collection of Meetings",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date",
            "description": "Meetings before this date",
            "name": "beforeDate",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date",
            "description": "Meetings after this date",
            "name": "afterDate",
            "in": "query"
          },
          {
            "type": "integer",
            "description": "Limit meetings to a specific team. Use database ID",
            "name": "teamid",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of meeting objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Meeting"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Meetings"
        ],
        "summary": "Create a new meeting",
        "parameters": [
          {
            "name": "meeting",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Meeting successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/meetings/{id}": {
      "get": {
        "tags": [
          "Meetings"
        ],
        "summary": "Get a particular meeting",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Meeting response",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Meetings"
        ],
        "summary": "Delete the meeting with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the meeting that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "Meetings"
        ],
        "summary": "Update a particular meeting",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "meeting",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated meeting response",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/people": {
      "get": {
        "description": "Collection of people registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "People"
        ],
        "summary": "Returns a collection of people",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of User objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Person"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "People"
        ],
        "summary": "Create a new Person",
        "parameters": [
          {
            "name": "person",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Person successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/people/{id}": {
      "get": {
        "tags": [
          "People"
        ],
        "summary": "Get a particular person",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Person response",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "People"
        ],
        "summary": "Delete the person with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the person that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "People"
        ],
        "summary": "Update a particular peson",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "student",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated student response",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/teams": {
      "get": {
        "description": "Collection of teams registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "Teams"
        ],
        "summary": "Returns a collection of teams",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of team objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Team"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Teams"
        ],
        "summary": "Create a new team",
        "parameters": [
          {
            "name": "mentor",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Team successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/teams/{id}": {
      "get": {
        "tags": [
          "Teams"
        ],
        "summary": "Get a particular team",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Team response",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Teams"
        ],
        "summary": "Delete the Team with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the Team that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "Teams"
        ],
        "summary": "Update a particular Team",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "mentor",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated Team response",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Day": {
      "type": "object",
      "properties": {
        "Day": {
          "type": "integer",
          "maximum": 31,
          "minimum": 1
        },
        "Month": {
          "type": "integer",
          "maximum": 12,
          "minimum": 1
        },
        "Year": {
          "type": "integer"
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "Meeting": {
      "type": "object",
      "properties": {
        "Committed": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Day": {
          "type": "string",
          "format": "date"
        },
        "EndTime": {
          "description": "End time of meeting. Using 24 hour time format. HH:MM:SS",
          "type": "string",
          "pattern": "(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)"
        },
        "SignedIn": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "SignedOut": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "StartTime": {
          "description": "Start time of meeting. Using 24 hour time format. HH:MM:SS",
          "type": "string",
          "pattern": "(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)"
        },
        "Teams": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "id": {
          "type": "string"
        }
      }
    },
    "Person": {
      "type": "object",
      "properties": {
        "PersonId": {
          "type": "integer"
        },
        "checkinid": {
          "type": "string"
        },
        "mentor": {
          "type": "object",
          "properties": {
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "name": {
          "type": "object",
          "required": [
            "first",
            "last"
          ],
          "properties": {
            "first": {
              "type": "string",
              "example": "Ryan"
            },
            "last": {
              "type": "string",
              "example": "Sjostrand"
            }
          }
        },
        "parentOf": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "relation": {
                "type": "string",
                "enum": [
                  "Father",
                  "Mother",
                  "Guardian"
                ]
              },
              "student": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "parents": {
          "type": "object",
          "properties": {
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "student": {
          "type": "object",
          "properties": {
            "graduationYear": {
              "type": "integer",
              "example": 2015
            },
            "schoolEmail": {
              "type": "string",
              "format": "email"
            },
            "schoolId": {
              "type": "string",
              "example": "sjost150"
            },
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        }
      }
    },
    "Team": {
      "type": "object",
      "properties": {
        "Competition": {
          "type": "string",
          "enum": [
            "FRC",
            "FTC",
            "FLL",
            "FLLjr"
          ],
          "example": "FRC"
        },
        "Meetings": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Mentors": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Name": {
          "type": "string",
          "example": "RoboEagles"
        },
        "Number": {
          "type": "integer",
          "example": 3081
        },
        "Students": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        }
      }
    },
    "Time": {
      "type": "object",
      "required": [
        "Hour",
        "Minute"
      ],
      "properties": {
        "Hour": {
          "type": "integer",
          "maximum": 24
        },
        "Minute": {
          "type": "integer",
          "maximum": 59
        }
      }
    },
    "count": {
      "description": "The number of items in a collection.",
      "type": "integer",
      "readOnly": true
    },
    "id": {
      "description": "The unique identifier for a resource.",
      "type": "string",
      "format": "uri-reference",
      "readOnly": true,
      "example": "foo"
    },
    "idRef": {
      "description": "A reference to a resource.",
      "type": "object",
      "properties": {
        "@meta.id": {
          "$ref": "#/definitions/id"
        }
      },
      "additionalProperties": false
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Sign In tracker API",
    "version": "1.0.0"
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "paths": {
    "/meetings": {
      "get": {
        "description": "Collection of meetings registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "Meetings"
        ],
        "summary": "Returns a collection of Meetings",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date",
            "description": "Meetings before this date",
            "name": "beforeDate",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date",
            "description": "Meetings after this date",
            "name": "afterDate",
            "in": "query"
          },
          {
            "type": "integer",
            "description": "Limit meetings to a specific team. Use database ID",
            "name": "teamid",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of meeting objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Meeting"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Meetings"
        ],
        "summary": "Create a new meeting",
        "parameters": [
          {
            "name": "meeting",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Meeting successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/meetings/{id}": {
      "get": {
        "tags": [
          "Meetings"
        ],
        "summary": "Get a particular meeting",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Meeting response",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Meetings"
        ],
        "summary": "Delete the meeting with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the meeting that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "Meetings"
        ],
        "summary": "Update a particular meeting",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "meeting",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated meeting response",
            "schema": {
              "$ref": "#/definitions/Meeting"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/people": {
      "get": {
        "description": "Collection of people registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "People"
        ],
        "summary": "Returns a collection of people",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of User objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Person"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "People"
        ],
        "summary": "Create a new Person",
        "parameters": [
          {
            "name": "person",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Person successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/people/{id}": {
      "get": {
        "tags": [
          "People"
        ],
        "summary": "Get a particular person",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Person response",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "People"
        ],
        "summary": "Delete the person with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the person that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "People"
        ],
        "summary": "Update a particular peson",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "student",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated student response",
            "schema": {
              "$ref": "#/definitions/Person"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/teams": {
      "get": {
        "description": "Collection of teams registered in the sign in tracker. By default only @meta.id's are returned in the colleciton. The links can be expanded if desired.",
        "tags": [
          "Teams"
        ],
        "summary": "Returns a collection of teams",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Expand the links",
            "name": "expand",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A array of team objects",
            "schema": {
              "type": "object",
              "required": [
                "Members",
                "@meta.id"
              ],
              "properties": {
                "@meta.id": {
                  "$ref": "#/definitions/id"
                },
                "Members": {
                  "type": "array",
                  "items": {
                    "allOf": [
                      {
                        "$ref": "#/definitions/idRef"
                      },
                      {
                        "$ref": "#/definitions/Team"
                      }
                    ]
                  }
                }
              }
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "Teams"
        ],
        "summary": "Create a new team",
        "parameters": [
          {
            "name": "mentor",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Team successfully created"
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/teams/{id}": {
      "get": {
        "tags": [
          "Teams"
        ],
        "summary": "Get a particular team",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Team response",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "delete": {
        "tags": [
          "Teams"
        ],
        "summary": "Delete the Team with the given ID",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Representation of the Team that was just deleted.",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "patch": {
        "tags": [
          "Teams"
        ],
        "summary": "Update a particular Team",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "name": "mentor",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Updated Team response",
            "schema": {
              "$ref": "#/definitions/Team"
            }
          },
          "default": {
            "description": "Unexpected error",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Day": {
      "type": "object",
      "properties": {
        "Day": {
          "type": "integer",
          "maximum": 31,
          "minimum": 1
        },
        "Month": {
          "type": "integer",
          "maximum": 12,
          "minimum": 1
        },
        "Year": {
          "type": "integer"
        }
      }
    },
    "Error": {
      "type": "object",
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "Meeting": {
      "type": "object",
      "properties": {
        "Committed": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Day": {
          "type": "string",
          "format": "date"
        },
        "EndTime": {
          "description": "End time of meeting. Using 24 hour time format. HH:MM:SS",
          "type": "string",
          "pattern": "(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)"
        },
        "SignedIn": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "SignedOut": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "StartTime": {
          "description": "Start time of meeting. Using 24 hour time format. HH:MM:SS",
          "type": "string",
          "pattern": "(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)"
        },
        "Teams": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "id": {
          "type": "string"
        }
      }
    },
    "Person": {
      "type": "object",
      "properties": {
        "PersonId": {
          "type": "integer"
        },
        "checkinid": {
          "type": "string"
        },
        "mentor": {
          "type": "object",
          "properties": {
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "name": {
          "type": "object",
          "required": [
            "first",
            "last"
          ],
          "properties": {
            "first": {
              "type": "string",
              "example": "Ryan"
            },
            "last": {
              "type": "string",
              "example": "Sjostrand"
            }
          }
        },
        "parentOf": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "relation": {
                "type": "string",
                "enum": [
                  "Father",
                  "Mother",
                  "Guardian"
                ]
              },
              "student": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "parents": {
          "type": "object",
          "properties": {
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        },
        "student": {
          "type": "object",
          "properties": {
            "graduationYear": {
              "type": "integer",
              "example": 2015
            },
            "schoolEmail": {
              "type": "string",
              "format": "email"
            },
            "schoolId": {
              "type": "string",
              "example": "sjost150"
            },
            "teams": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/idRef"
              }
            }
          }
        }
      }
    },
    "Team": {
      "type": "object",
      "properties": {
        "Competition": {
          "type": "string",
          "enum": [
            "FRC",
            "FTC",
            "FLL",
            "FLLjr"
          ],
          "example": "FRC"
        },
        "Meetings": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Mentors": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        },
        "Name": {
          "type": "string",
          "example": "RoboEagles"
        },
        "Number": {
          "type": "integer",
          "example": 3081
        },
        "Students": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/idRef"
          }
        }
      }
    },
    "Time": {
      "type": "object",
      "required": [
        "Hour",
        "Minute"
      ],
      "properties": {
        "Hour": {
          "type": "integer",
          "maximum": 24,
          "minimum": 0
        },
        "Minute": {
          "type": "integer",
          "maximum": 59,
          "minimum": 0
        }
      }
    },
    "count": {
      "description": "The number of items in a collection.",
      "type": "integer",
      "readOnly": true
    },
    "id": {
      "description": "The unique identifier for a resource.",
      "type": "string",
      "format": "uri-reference",
      "readOnly": true,
      "example": "foo"
    },
    "idRef": {
      "description": "A reference to a resource.",
      "type": "object",
      "properties": {
        "@meta.id": {
          "$ref": "#/definitions/id"
        }
      },
      "additionalProperties": false
    }
  }
}`))
}

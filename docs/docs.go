// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Nika Shamiladze",
            "email": "fbshamiladze@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/rooms": {
            "get": {
                "description": "List all rooms",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List all rooms",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.RoomDTO"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new room",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new room",
                "parameters": [
                    {
                        "description": "Room details",
                        "name": "room",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RoomDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.RoomDTO"
                        }
                    }
                }
            }
        },
        "/api/v1/rooms/{ID}/join": {
            "post": {
                "description": "Join a room",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Join a room",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.RoomDTO"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.RoomDTO": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "maxPlayers": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "private": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Game V0 API",
	Description:      "Game V0 API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

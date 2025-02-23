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
        "/room": {
            "get": {
                "description": "Get all rooms",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Room"
                ],
                "summary": "Get all rooms",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "additionalProperties": true
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a room",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Room"
                ],
                "summary": "Create a room",
                "parameters": [
                    {
                        "description": "Room",
                        "name": "room",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/presenter.RoomRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.RoomResponse"
                        }
                    }
                }
            }
        },
        "/room/join": {
            "post": {
                "description": "Join a room",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Room"
                ],
                "summary": "Join a room",
                "parameters": [
                    {
                        "description": "Room",
                        "name": "room",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/presenter.JoinRoomRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.RoomResponse"
                        }
                    }
                }
            }
        },
        "/user/me": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get current user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.MeResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user/sign-in": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign in",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/presenter.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/presenter.SignInResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user/sign-up": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign up",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/presenter.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns message and user object with excluded fields",
                        "schema": {
                            "$ref": "#/definitions/presenter.SignUpResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/presenter.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "presenter.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "presenter.JoinRoomRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "presenter.MeResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "presenter.RoomRequest": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "private": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "presenter.RoomResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "presenter.SignInRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "presenter.SignInResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "presenter.SignUpRequest": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 8
                },
                "username": {
                    "type": "string",
                    "maxLength": 32,
                    "minLength": 4
                }
            }
        },
        "presenter.SignUpResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "user": {
                    "type": "object",
                    "properties": {
                        "email": {
                            "type": "string"
                        },
                        "id": {
                            "type": "string"
                        },
                        "username": {
                            "type": "string"
                        }
                    }
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
	Title:            "Game v0 API",
	Description:      "API for the survival game",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

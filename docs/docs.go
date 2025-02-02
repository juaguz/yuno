// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cards": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create a card with card holder and PAN",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Create a new card",
                "parameters": [
                    {
                        "description": "Card Creation Request",
                        "name": "card",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CardCreation"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.Card"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cards/batch": {
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Update multiple cards in a single request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Batch update cards",
                "parameters": [
                    {
                        "description": "Batch Update Request",
                        "name": "batch",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.BatchUpdate"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.BatchUpdateStatus"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cards/{cardID}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Retrieve a card by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Get a card",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Card ID",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.Card"
                        }
                    },
                    "400": {
                        "description": "Invalid card ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Card not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Update a card's details by its ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cards"
                ],
                "summary": "Update a card",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Card ID",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Card Update Request",
                        "name": "card",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CardUpdate"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No content"
                    },
                    "400": {
                        "description": "Invalid request body or card ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Delete a card by its ID",
                "tags": [
                    "cards"
                ],
                "summary": "Delete a card",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Card ID",
                        "name": "cardID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No content"
                    },
                    "400": {
                        "description": "Invalid card ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/keys": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Generates a new public key for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "keys"
                ],
                "summary": "Create a new key",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/api.KeysResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.CardCreation": {
            "type": "object",
            "properties": {
                "card_holder": {
                    "type": "string"
                },
                "pan": {
                    "type": "string"
                }
            }
        },
        "api.CardUpdate": {
            "type": "object",
            "properties": {
                "card_holder": {
                    "type": "string"
                }
            }
        },
        "api.KeysResponse": {
            "type": "object",
            "properties": {
                "public_key": {
                    "type": "string"
                }
            }
        },
        "dtos.BatchUpdate": {
            "type": "object",
            "properties": {
                "card_holder": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "dtos.BatchUpdateStatus": {
            "type": "object",
            "properties": {
                "card": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/dtos.Status"
                }
            }
        },
        "dtos.Card": {
            "type": "object",
            "properties": {
                "card_holder": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "pan": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "dtos.Status": {
            "type": "string",
            "enum": [
                "succeeded",
                "failed"
            ],
            "x-enum-varnames": [
                "Succeeded",
                "Failed"
            ]
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/organization": {
            "get": {
                "description": "Get all organizations.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Organization"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "Create organization.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "parameters": [
                    {
                        "description": "body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/organization.CreateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.Organization"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/organization/{id}": {
            "get": {
                "description": "Get organization by ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Organization"
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "description": "Update organization.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/organization.UpdateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.Organization"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Delete organization.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Organization": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "org_id": {
                    "type": "string"
                },
                "org_key": {
                    "type": "string"
                }
            }
        },
        "organization.CreateOrganizationRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "org_key": {
                    "type": "string"
                }
            }
        },
        "organization.UpdateOrganizationRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Cronuseo API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

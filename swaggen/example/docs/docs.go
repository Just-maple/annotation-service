// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/test/add": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "This is title",
                "parameters": [
                    {
                        "type": "number",
                        "name": "float32",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "int",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "string",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/example.Ret"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/example.Ret"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test/dec": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "dec func",
                "parameters": [
                    {
                        "description": "example.Param",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/example.Param"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test/{add2}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "This is title",
                "parameters": [
                    {
                        "description": "example.Param",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/example.Param"
                        }
                    },
                    {
                        "type": "string",
                        "description": "add2",
                        "name": "add2",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/example.Ret"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/example.Ret"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test/{dec}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "dec func",
                "parameters": [
                    {
                        "type": "number",
                        "name": "float32",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "int",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "string",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dec",
                        "name": "dec",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test2/add3": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "This is title",
                "parameters": [
                    {
                        "type": "number",
                        "name": "float32",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "int",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "string",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/example.Ret"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/example.Ret"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test2/dec": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "dec func",
                "parameters": [
                    {
                        "description": "example.Param",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/example.Param"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        },
        "/test2/{dec}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "dec func",
                "parameters": [
                    {
                        "type": "number",
                        "name": "float32",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "int",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "string",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "dec",
                        "name": "dec",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/example.Ret"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "example.Param": {
            "type": "object",
            "properties": {
                "float32": {
                    "type": "number"
                },
                "int": {
                    "type": "integer"
                },
                "string": {
                    "type": "string"
                }
            }
        },
        "example.Ret": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
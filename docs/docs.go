// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/song/create": {
            "post": {
                "description": "Creating song with song name and group name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "Create song",
                "parameters": [
                    {
                        "description": "Song insert",
                        "name": "song_insert",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.SongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/entity.Song"
                        }
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/song/list": {
            "get": {
                "description": "Getting songs with pagination and filtered",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "get list songs with filters",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "With deleted",
                        "name": "d",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "group_name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "p",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "description": "Filters json body",
                        "name": "filters",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/entity.SongFilters"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/entity.SongListResponse"
                        }
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/song/{song_id}": {
            "put": {
                "description": "Put song change all fields",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "Put song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song id",
                        "name": "song_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Song with changing fields",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/entity.Song"
                        }
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deleting song if soft delete true then deleted_at in BD change on NOW() or if soft is false then delete row from bd",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "Delete song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song id",
                        "name": "song_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "Is soft delete",
                        "name": "soft",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Patch song change any fields",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "Patch song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song id",
                        "name": "song_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Song with changing fields",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/entity.Song"
                        }
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/song/{song_id}/text": {
            "get": {
                "description": "Get text with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Song"
                ],
                "summary": "get text songs with pagination",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song id",
                        "name": "song_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/entity.SongTextResponse"
                        }
                    },
                    "400": {
                        "description": "Params not valid",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Can not find ID",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "entity.Group": {
            "type": "object",
            "required": [
                "g_name"
            ],
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "description": "DeletedAt nil value if not deleted",
                    "type": "string"
                },
                "g_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "entity.Song": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "group": {
                    "$ref": "#/definitions/entity.Group"
                },
                "group_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "update_at": {
                    "type": "string"
                }
            }
        },
        "entity.SongFilters": {
            "type": "object",
            "properties": {
                "group_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "entity.SongListResponse": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "per_page": {
                    "type": "integer"
                },
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Song"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "entity.SongRequest": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "entity.SongTextResponse": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "text": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1alpha",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Music librarry API",
	Description:      "This is a simple api for music librarry",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

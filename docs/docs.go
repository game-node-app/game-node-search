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
        "contact": {
            "name": "API Support",
            "url": "gamenode.com.br/help",
            "email": "support@gamenode.com.br"
        },
        "license": {
            "name": "GNU General Public License",
            "url": "http://www.gnu.org/licenses/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/search": {
            "post": {
                "description": "Returns a parsed search response from the Manticore engine",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "search"
                ],
                "summary": "Searches using Manticore engine",
                "parameters": [
                    {
                        "description": "Account ID",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schema.GameSearchRequestDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schema.GameSearchResponseDto"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "schema.AnyMap": {
            "type": "object",
            "additionalProperties": true
        },
        "schema.GameSearchRequestDto": {
            "type": "object",
            "properties": {
                "aggs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.AnyMap"
                    }
                },
                "attrFilter": {
                    "$ref": "#/definitions/schema.AnyMap"
                },
                "expressions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.AnyMap"
                    }
                },
                "fulltextFilter": {
                    "$ref": "#/definitions/schema.AnyMap"
                },
                "highlight": {
                    "$ref": "#/definitions/schema.AnyMap"
                },
                "index": {
                    "type": "string"
                },
                "limit": {
                    "type": "integer"
                },
                "maxMatches": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "profile": {
                    "type": "boolean"
                },
                "query": {
                    "$ref": "#/definitions/schema.AnyMap"
                },
                "sort": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.AnyMap"
                    }
                },
                "source": {
                    "$ref": "#/definitions/schema.AnyMap"
                },
                "trackScores": {
                    "type": "boolean"
                }
            }
        },
        "schema.GameSearchResponseDto": {
            "type": "object",
            "properties": {
                "aggregations": {
                    "type": "object",
                    "additionalProperties": true
                },
                "hits": {
                    "$ref": "#/definitions/schema.GameSearchResponseHits"
                },
                "profile": {
                    "type": "object",
                    "additionalProperties": true
                },
                "timedOut": {
                    "type": "boolean"
                },
                "took": {
                    "type": "integer"
                },
                "warning": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "schema.GameSearchResponseHit": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "score": {
                    "type": "integer"
                },
                "source": {
                    "$ref": "#/definitions/schema.SearchGame"
                }
            }
        },
        "schema.GameSearchResponseHits": {
            "type": "object",
            "properties": {
                "hits": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schema.GameSearchResponseHit"
                    }
                },
                "maxScore": {
                    "type": "number"
                },
                "total": {
                    "type": "integer"
                },
                "totalRelation": {
                    "type": "string"
                }
            }
        },
        "schema.SearchGame": {
            "type": "object",
            "properties": {
                "aggregatedRating": {
                    "type": "number"
                },
                "aggregatedRatingCount": {
                    "type": "integer"
                },
                "category": {
                    "type": "integer"
                },
                "checksum": {
                    "type": "string"
                },
                "coverUrl": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "firstReleaseDate": {
                    "type": "string"
                },
                "genresNames": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "keywordsNames": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "numLikes": {
                    "type": "integer"
                },
                "numViews": {
                    "type": "integer"
                },
                "platformsAbbreviations": {
                    "type": "string"
                },
                "platformsNames": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "storyline": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "GameNode Search API",
	Description:      "The GameNode Search API documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

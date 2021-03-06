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
        "contact": {
            "name": "GoIn",
            "url": "localhost:8088"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/casbin/v1/application/enforce": {
            "post": {
                "description": "注册接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "创建Casbin权限",
                "parameters": [
                    {
                        "description": "创建Casbin权限参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.EnforcePost"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/casbin/v1/application/policy/all": {
            "get": {
                "description": "获取所有权限",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "获取所有权限",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/casbin/v1/application/policy/role": {
            "get": {
                "description": "获取角色所有权限",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "获取角色所有权限",
                "parameters": [
                    {
                        "type": "string",
                        "description": "角色名称",
                        "name": "role",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/casbin/v1/application/role/policy": {
            "delete": {
                "description": "批量删除角色权限",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "批量删除角色权限",
                "parameters": [
                    {
                        "description": "需要删除的权限，格式：[{",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/casbin/v1/application/user/role": {
            "get": {
                "description": "获取用户所有角色",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "获取用户所有角色",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户ID",
                        "name": "uid",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "新增用户角色",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "新增用户角色",
                "parameters": [
                    {
                        "description": "更新用户角色参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UpdateRolesForUserPut"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "post": {
                "description": "新增用户角色",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "新增用户角色",
                "parameters": [
                    {
                        "description": "新增用户角色参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.AddRolesForUserPost"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "批量删除用户角色",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "权限"
                ],
                "summary": "批量删除用户角色",
                "parameters": [
                    {
                        "description": "批量删除用户角色参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.DeleteRolesForUserPost"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/user/v1/application/info": {
            "get": {
                "description": "用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pojo.UserInfo"
                        }
                    }
                }
            }
        },
        "/user/v1/application/login": {
            "post": {
                "description": "普通登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "普通登录",
                "parameters": [
                    {
                        "description": "登陆参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.UserLoginResponse"
                        }
                    }
                }
            }
        },
        "/user/v1/application/modify/role": {
            "post": {
                "description": "修改用户角色",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "修改用户角色",
                "parameters": [
                    {
                        "description": "修改用户角色权限参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.ModifyUserRoleParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "pojo.UserInfo": {
            "type": "object",
            "properties": {
                "avatar": {
                    "description": "用户头像",
                    "type": "string"
                },
                "mailBox": {
                    "description": "邮箱",
                    "type": "string"
                },
                "mobile": {
                    "description": "手机号码",
                    "type": "string"
                },
                "nickName": {
                    "description": "用户昵称",
                    "type": "string"
                },
                "roleCode": {
                    "description": "账号角色 basic:普通员工 admin:管理员",
                    "type": "string"
                },
                "sex": {
                    "description": "性别",
                    "type": "integer"
                },
                "uid": {
                    "description": "用户唯一标识",
                    "type": "string"
                },
                "userName": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "v1.AddRolesForUserPost": {
            "type": "object",
            "required": [
                "roles",
                "uid"
            ],
            "properties": {
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "v1.DeleteRolesForUserPost": {
            "type": "object",
            "required": [
                "roles",
                "uid"
            ],
            "properties": {
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "v1.EnforcePost": {
            "type": "object",
            "required": [
                "act",
                "obj",
                "sub"
            ],
            "properties": {
                "act": {
                    "description": "用户对资源的操作(GET,POST)",
                    "type": "string"
                },
                "obj": {
                    "description": "将要被访问的资源(功能路径)",
                    "type": "string"
                },
                "sub": {
                    "description": "角色",
                    "type": "string"
                }
            }
        },
        "v1.ModifyUserRoleParam": {
            "type": "object",
            "required": [
                "roleCode",
                "uid"
            ],
            "properties": {
                "roleCode": {
                    "description": "修改为某个角色 basic:普通员工 admin:管理员",
                    "type": "string"
                },
                "uid": {
                    "description": "被修改的用户ID",
                    "type": "string"
                }
            }
        },
        "v1.UpdateRolesForUserPut": {
            "type": "object",
            "required": [
                "roles",
                "uid"
            ],
            "properties": {
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "v1.UserLogin": {
            "type": "object",
            "required": [
                "password",
                "user_name"
            ],
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "user_name": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "v1.UserLoginResponse": {
            "type": "object",
            "properties": {
                "info": {
                    "description": "用户信息",
                    "$ref": "#/definitions/pojo.UserInfo"
                },
                "token": {
                    "description": "token",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
	Version:     "2.0",
	Host:        "localhost:8088",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "GoIn",
	Description: "GoIn",
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

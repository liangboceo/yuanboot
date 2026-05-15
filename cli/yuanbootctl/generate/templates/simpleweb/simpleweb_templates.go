package simpleweb

const Main_Tel = `package main

import (
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/router"
)

func main() {
	web.CreateHttpBuilder(func(rb router.IRouterBuilder) {
		// Home page
		rb.GET("/", func(ctx *context.HttpContext) {
			ctx.HTML(200, "index.html", map[string]interface{}{
				"title": "{{.ModelName}}",
				"info":  "Welcome to Yuanboot!",
			})
		})

		// API endpoint
		rb.GET("/info", func(ctx *context.HttpContext) {
			ctx.JSON(200, context.H{
				"status": "ok",
				"app":    "{{.ModelName}}",
				"version": "1.0.0",
			})
		})

		// RESTful API example
		rb.Group("/api/v1", func(g *router.RouterGroup) {
			g.GET("/hello", func(ctx *context.HttpContext) {
				name := ctx.Query("name")
				if name == "" {
					name = "World"
				}
				ctx.JSON(200, context.H{
					"message": "Hello, " + name + "!",
				})
			})

			g.POST("/echo", func(ctx *context.HttpContext) {
				var data map[string]interface{}
				ctx.Bind(&data)
				ctx.JSON(200, context.H{
					"received": data,
				})
			})

			g.GET("/users/:id", func(ctx *context.HttpContext) {
				id := ctx.Param("id")
				ctx.JSON(200, context.H{
					"id":   id,
					"name": "User " + id,
				})
			})
		})
	}).Build().Run()
}
`

const Config_Tel = `yuanboot:
  application:
    name: {{.ModelName}}
    metadata: "develop"
    server:
      type: "fasthttp"
      address: ":8080"
      path: ""
      max_request_size: 2097152
      mvc:
        views:
          path: "./static/templates"
      static:
        patten: "/static/*"
        webroot: "./static"
      cors:
        allow_origins:
          - "*"
        allow_methods:
          - "POST"
          - "GET"
          - "PUT"
          - "DELETE"
          - "PATCH"
        allow_credentials: true
`

const Mod_Tel = `
module {{.ModelName}}

go 1.18

require (
	github.com/liangboceo/yuanboot {{.Version}}
)
`

const IndexHtml_Tel = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        .container {
            background: white;
            border-radius: 20px;
            padding: 60px;
            text-align: center;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 600px;
        }
        h1 {
            color: #667eea;
            font-size: 48px;
            margin-bottom: 20px;
        }
        .info {
            color: #666;
            font-size: 18px;
            margin-bottom: 30px;
        }
        .features {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            justify-content: center;
            margin-top: 30px;
        }
        .feature {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 10px 20px;
            border-radius: 25px;
            font-size: 14px;
        }
        .links {
            margin-top: 30px;
        }
        .links a {
            display: inline-block;
            margin: 0 10px;
            padding: 12px 30px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 25px;
            transition: transform 0.3s, box-shadow 0.3s;
        }
        .links a:hover {
            transform: translateY(-3px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Yuanboot</h1>
        <p class="info">{{.info}}</p>
        <p>Lightweight microservice framework based on Go</p>

        <div class="features">
            <span class="feature">High Performance</span>
            <span class="feature">MVC</span>
            <span class="feature">DI</span>
            <span class="feature">Middleware</span>
            <span class="feature">Service Discovery</span>
        </div>

        <div class="links">
            <a href="/info">API Info</a>
            <a href="/api/v1/hello">Hello API</a>
        </div>
    </div>
</body>
</html>
`

const Readme_Tel = `
# {{.ModelName}}

Simple web application based on Yuanboot framework

## Quick Start

` + "```bash" + `
# Run application
go run main.go

# Visit home page
open http://localhost:8080
` + "```" + `

## API List

| Method | Path | Description |
|--------|------|-------------|
| GET | / | Home page |
| GET | /info | API info |
| GET | /api/v1/hello | Hello API |
| POST | /api/v1/echo | Echo API |
| GET | /api/v1/users/:id | Get user |

## Documentation

For more documentation, please visit [Yuanboot](https://yuanboot.star2cloud.com)
`

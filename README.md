# nine

Nine is an HTTP client and server library in Go that simplifies sending HTTP
requests and processing routes. With Nine, you can easily create HTTP clients
and servers with support for middleware, response handling, and much more.

## Installation

To add the Nine library to your Go project, you can use the following command:

```sh
go get github.com/i9si-sistemas/nine
```

## Usage

### Client

The HTTP client allows you to make various types of requests (GET, POST, PUT,
PATCH, DELETE) easily.

```go
package main

import (
	"context"
	"fmt"

	"github.com/i9si-sistemas/nine"
	i9 "github.com/i9si-sistemas/nine/pkg/client"
)

func main() {
	ctx := context.Background()
	client := nine.New(ctx)
	buf, err := nine.JSON{
		"message": "Hello World",
	}.Buffer()
	if err != nil {
		fmt.Println("Erro ao gerar o buffer do payload:", err)
		return
	}
	options := &i9.Options{
		Headers: []i9.Header{
			{Data: i9.Data{Key: "Content-Type", Value: "application/json"}},
		},
		Body: buf,
	}

	res, err := client.Get("https://api.exemplo.com/endpoint", options)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return
	}
	defer res.Body.Close()
	 // use the response
}
```

### Server

To create a server, you can use the following code:

```go
package main

import (
	"fmt"
	"log"

	"github.com/i9si-sistemas/nine"
	i9 "github.com/i9si-sistemas/nine/pkg/server"
)

func main() {
	server := nine.NewServer(8080)

	server.Get("/hello", func(c *i9.Context) error {
		return c.Send([]byte("Hello World"))
	})

	server.Get("/hello/:name", func(c *i9.Context) error {
		name := c.Param("name")
		message := fmt.Sprintf("Hello %s", name)
		return c.Send([]byte(message))
	})

	log.Fatal(server.Listen())
}
```

### Route Handling

You can register routes for different HTTP methods using the Get, Post, Put,
Patch, and Delete methods.

```go
server.Post("/create", func(c *i9.Context) error {
		return c.Status(http.StatusCreated).Send([]byte("Recurso criado com sucesso"))
})

server.Route("/billing", func(router *i9.RouteGroup) {
	router.Get("/credits", func(c *i9.Context) error {
		return c.JSON(nine.JSON{"credits": 5000})
	})
})

accountGroup := server.Group("/account", func (req *i9.Request, res *i9.Response) error {
	ctx := i9.NewContext(req.Context(), req.HTTP(), res.HTTP())
	log.Printf("receiving request from ip=[%s]", ctx.IP())
	return nil
})

accountGroup.Get("/:name", func (c *i9.Context) error {
	return c.JSON(nine.JSON{"account": c.Param("name")})
})
```

### JSON Handling

The library also provides utilities for working with JSON:

```go
server.Get("/user", func(c *i9.Context) error {
	data := nine.JSON{"name": "Alice", "age": 30}
	return c.JSON(data)
})
```

## Contributing

Contributions are always welcome! Feel free to open issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for more details.

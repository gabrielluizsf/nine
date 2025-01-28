package nine

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/i9si-sistemas/assert"
)

func TestRouteGroup(t *testing.T) {
	testServer := NewServer(8080)
	type Account struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	accounts := make(GenericJSON[string, Account], 0)
	testServer.Route("/account", func(router *RouteGroup) {
		router.Post("/create", func(c *Context) error {
			var body Account
			if err := Body(c.Request, &body); err != nil {
				return c.Status(http.StatusBadRequest).JSON(JSON{
					"message": "invalid body",
				})
			}
			_, ok := accounts[body.Name]
			response := JSON{"created": !ok}
			if !ok {
				accounts[body.Name] = body
				return c.Status(http.StatusCreated).JSON(response)
			}
			return c.JSON(response)
		})
		router.Get("/:name", func(c *Context) error {
			acc, ok := accounts[c.Param("name")]
			if !ok {
				return c.SendStatus(http.StatusNotFound)
			}
			return c.JSON(JSON{
				"account": acc,
			})
		})
	})
	assertGroupEndpoints(t, testServer)
	req := httptest.NewRequest(http.MethodGet, "/account/Gabriel%20Luiz", nil)
	w := testServer.Test().Request(req)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	b := w.Body.Bytes()
	var account struct {
		Account `json:"account"`
	}
	err := DecodeJSON(b, &account)
	assert.NoError(t, err)
	assert.Equal(t, account.Name, "Gabriel Luiz")
	assert.Equal(t, account.Age, 23)
}

func TestGroup(t *testing.T) {
	testServer := NewServer(5024)
	type Account struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Money int64  `json:"money"`
	}
	accounts := make(GenericJSON[string, Account], 0)
	accountGroup := testServer.Group("/account")
	accountGroup.Post("/create", func(c *Context) error {
		var body Account
		if err := Body(c.Request, &body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(JSON{
				"message": "invalid body",
			})
		}
		_, ok := accounts[body.Name]
		response := JSON{"created": !ok}
		if !ok {
			accounts[body.Name] = body
			return c.Status(http.StatusCreated).JSON(response)
		}
		return c.JSON(response)
	})
	accountGroup.Get("/:name", func(c *Context) error {
		acc, ok := accounts[c.Param("name")]
		if !ok {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.JSON(JSON{
			"account": acc,
		})
	})
	assertGroupEndpoints(t, testServer)
	req := httptest.NewRequest(http.MethodGet, "/account/Gabriel%20Luiz", nil)
	w := testServer.Test().Request(req)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	b := w.Body.Bytes()
	var account struct {
		Account `json:"account"`
	}
	err := DecodeJSON(b, &account)
	assert.NoError(t, err)
	assert.Equal(t, account.Name, "Gabriel Luiz")
	assert.Equal(t, account.Age, 23)
	assert.Equal(t, account.Money, int64(5000))
}

func assertGroupEndpoints(t assert.T,testServer *Server) {
	var response struct {
		Created bool `json:"created"`
	}
	payload, err := JSON{
		"name": "Gabriel Luiz",
		"age":  23,
		"money": 5000,
	}.Buffer()
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/account/create", payload)
	w := testServer.Test().Request(req)
	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
	b := w.Body.Bytes()
	err = DecodeJSON(b, &response)
	assert.NoError(t, err)
	assert.True(t, response.Created)

	payload, err = JSON{
		"name": "Gabriel Luiz",
		"age":  23,
		"money": 5000000,
	}.Buffer()
	assert.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/account/create", payload)
	w = testServer.Test().Request(req)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	b = w.Body.Bytes()
	err = DecodeJSON(b, &response)
	assert.NoError(t, err)
	assert.False(t, response.Created)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	client, err := elasticsearch8.NewDefaultClient()
	if err != nil {
		e.Logger.Fatal(err)
	}

	res, err := client.Info()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer res.Body.Close()
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Printf("%+v\n", r)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"bytes"
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
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": "Brave",
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return err
		}

		res, err := client.Search(
			client.Search.WithContext(c.Request().Context()),
			client.Search.WithIndex("books"),
			client.Search.WithBody(&buf),
			client.Search.WithTrackTotalHits(true),
			client.Search.WithPretty(),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				return err
			}
			return c.JSONPretty(res.StatusCode, e, "  ")
		}

		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}

		fmt.Printf("%+v\n", r)

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

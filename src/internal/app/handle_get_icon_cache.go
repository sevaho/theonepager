package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sevaho/theonepager/src/pkg/resty"
)

type Hit struct {
	content     []byte
	contentType string
}

func handleGetIconCache() echo.HandlerFunc {
	cache := make(map[string]Hit)

	return func(ctx echo.Context) error {
		link := ctx.QueryParam("link")

		var hit Hit

		if val, ok := cache[link]; ok {
			hit = val
		} else {
			res, err := resty.Client.R().Get(link)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}

			hit = Hit{content: res.Body(), contentType: res.Header().Get("Content-Type")}
			cache[link] = hit
		}

		return ctx.Blob(http.StatusOK, hit.contentType, hit.content)
	}
}

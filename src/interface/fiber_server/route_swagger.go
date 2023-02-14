package fiber_server

import (
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (f *FiberServer) addRouteSwagger(base fiber.Router) {
	base.Get("/swagger", func(ctx *fiber.Ctx) error {
		if strings.HasSuffix(ctx.Path(), "/") {
			return ctx.Redirect("index.html", 301)
		}
		return ctx.Redirect("swagger/index.html", 301)
	})

	r := base.Group("/swagger")

	r.Get("*", swagger.New(swagger.Config{
		Layout: "StandaloneLayout",
		Plugins: []template.JS{
			"SwaggerUIBundle.plugins.DownloadUrl",
		},
		Presets: []template.JS{
			"SwaggerUIBundle.presets.apis",
			"SwaggerUIStandalonePreset",
		},
		DeepLinking:              true,
		DefaultModelsExpandDepth: 1,
		DefaultModelExpandDepth:  1,
		DefaultModelRendering:    "example",
		DocExpansion:             "list",
		SyntaxHighlight: &swagger.SyntaxHighlightConfig{
			Activate: true,
			Theme:    "monokai",
		},
		ShowMutatedRequest: true,
		URL:                "doc.json",
	}))
}

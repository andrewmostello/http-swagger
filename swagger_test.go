package httpSwagger

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

func TestWrapHandler(t *testing.T) {
	router := chi.NewRouter()

	router.Get("/*", WrapHandler)

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w2 := performRequest("GET", "/doc.json", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, "application/json; charset=utf-8", w2.Header().Get("content-type"))

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w4.Code)

	w5 := performRequest("GET", "/", router)
	assert.Equal(t, 301, w5.Code)
}

func performRequest(method, target string, h http.Handler) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	return w
}

func TestConfigURL(t *testing.T) {

	type fixture struct {
		desc  string
		cfgfn func(c *Config)
		exp   *Config
	}

	fixtures := []fixture{
		{
			desc: "configure URL",
			exp: &Config{
				URL: "https://example.org/doc.json",
			},
			cfgfn: URL("https://example.org/doc.json"),
		},
		{
			desc: "configure DeepLinking",
			exp: &Config{
				DeepLinking: true,
			},
			cfgfn: DeepLinking(true),
		},
		{
			desc: "configure DocExpansion",
			exp: &Config{
				DocExpansion: "none",
			},
			cfgfn: DocExpansion("none"),
		},
		{
			desc: "configure DomID",
			exp: &Config{
				DomID: "#swagger-ui",
			},
			cfgfn: DomID("#swagger-ui"),
		},
		{
			desc: "configure Plugins",
			exp: &Config{
				Plugins: []template.JS{
					"SomePlugin",
					"AnotherPlugin",
				},
			},
			cfgfn: Plugins([]string{
				"SomePlugin",
				"AnotherPlugin",
			}),
		},
		{
			desc: "configure UIConfig",
			exp: &Config{
				UIConfig: map[template.JS]template.JS{
					"urls": `["https://example.org/doc1.json","https://example.org/doc1.json"],`,
				},
			},
			cfgfn: UIConfig(map[string]string{
				"urls": `["https://example.org/doc1.json","https://example.org/doc1.json"],`,
			}),
		},
		{
			desc: "configure BeforeScript",
			exp: &Config{
				BeforeScript: `const SomePlugin = (system) => ({
    // Some plugin
  });`,
			},
			cfgfn: BeforeScript(`const SomePlugin = (system) => ({
    // Some plugin
  });`),
		},
		{
			desc: "configure AfterScript",
			exp: &Config{
				AfterScript: `const SomePlugin = (system) => ({
    // Some plugin
  });`,
			},
			cfgfn: AfterScript(`const SomePlugin = (system) => ({
    // Some plugin
  });`),
		},
	}

	for _, fix := range fixtures {
		t.Run(fix.desc, func(t *testing.T) {
			cfg := &Config{}
			fix.cfgfn(cfg)
			assert.Equal(t, cfg, fix.exp)
		})
	}
}

func TestUIConfigOptions(t *testing.T) {

	type fixture struct {
		desc string
		cfg  *Config
		exp  string
	}

	hdr := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link href="https://fonts.googleapis.com/css?family=Open+Sans:400,700|Source+Code+Pro:300,600|Titillium+Web:400,600,700" rel="stylesheet">
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html
    {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
    }
    *,
    *:before,
    *:after
    {
        box-sizing: inherit;
    }

    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>

<body>

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="position:absolute;width:0;height:0">
  <defs>
    <symbol viewBox="0 0 20 20" id="unlocked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V6h2v-.801C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8z"></path>
    </symbol>

    <symbol viewBox="0 0 20 20" id="locked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8zM12 8H8V5.199C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="close">
      <path d="M14.348 14.849c-.469.469-1.229.469-1.697 0L10 11.819l-2.651 3.029c-.469.469-1.229.469-1.697 0-.469-.469-.469-1.229 0-1.697l2.758-3.15-2.759-3.152c-.469-.469-.469-1.228 0-1.697.469-.469 1.228-.469 1.697 0L10 8.183l2.651-3.031c.469-.469 1.228-.469 1.697 0 .469.469.469 1.229 0 1.697l-2.758 3.152 2.758 3.15c.469.469.469 1.229 0 1.698z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow">
      <path d="M13.25 10L6.109 2.58c-.268-.27-.268-.707 0-.979.268-.27.701-.27.969 0l7.83 7.908c.268.271.268.709 0 .979l-7.83 7.908c-.268.271-.701.27-.969 0-.268-.269-.268-.707 0-.979L13.25 10z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow-down">
      <path d="M17.418 6.109c.272-.268.709-.268.979 0s.271.701 0 .969l-7.908 7.83c-.27.268-.707.268-.979 0l-7.908-7.83c-.27-.268-.27-.701 0-.969.271-.268.709-.268.979 0L10 13.25l7.418-7.141z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="jump-to">
      <path d="M19 7v4H5.83l3.58-3.59L8 6l-6 6 6 6 1.41-1.41L5.83 13H21V7z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="expand">
      <path d="M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z"/>
    </symbol>
  </defs>
</svg>

<div id="swagger-ui"></div>

<script src="./swagger-ui-bundle.js"> </script>
<script src="./swagger-ui-standalone-preset.js"> </script>
<script>
`
	ftr := `
</script>
</body>

</html>
`

	fixtures := []fixture{
		{
			desc: "default configuration",
			cfg: &Config{
				URL:          "doc.json",
				DeepLinking:  true,
				DocExpansion: "list",
				DomID:        "#swagger-ui",
			},
			exp: `window.onload = function() {
  
  const ui = SwaggerUIBundle({
    url: "doc.json",
    deepLinking:  true ,
    docExpansion: "list",
    dom_id: "#swagger-ui",
    validatorUrl: null,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  })

  window.ui = ui
}`,
		},
		{
			desc: "script configuration",
			cfg: &Config{
				URL:          "swagger.json",
				DeepLinking:  false,
				DocExpansion: "none",
				DomID:        "#swagger-ui-id",
				BeforeScript: `const SomePlugin = (system) => ({
    // Some plugin
  });
`,
				AfterScript: `const someOtherCode = function(){
    // Do something
  };
  someOtherCode();`,
				Plugins: []template.JS{
					"SomePlugin",
					"AnotherPlugin",
				},
				UIConfig: map[template.JS]template.JS{
					"showExtensions":        "true",
					"onComplete":            `() => { window.ui.setBasePath('v3'); }`,
					"defaultModelRendering": `"model"`,
				},
			},
			exp: `window.onload = function() {
  const SomePlugin = (system) => ({
    // Some plugin
  });

  
  const ui = SwaggerUIBundle({
    url: "swagger.json",
    deepLinking:  false ,
    docExpansion: "none",
    dom_id: "#swagger-ui-id",
    validatorUrl: null,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl,
      SomePlugin,
      AnotherPlugin
    ],
    defaultModelRendering: "model",
    onComplete: () => { window.ui.setBasePath('v3'); },
    showExtensions: true,
    layout: "StandaloneLayout"
  })

  window.ui = ui
  const someOtherCode = function(){
    // Do something
  };
  someOtherCode();
}`,
		},
	}

	for _, fix := range fixtures {
		t.Run(fix.desc, func(t *testing.T) {
			tmpl := template.New("swagger_index.html")
			index, err := tmpl.Parse(indexTempl)
			if err != nil {
				t.Fatal(err)
			}

			buf := bytes.NewBuffer(nil)
			if err := index.Execute(buf, fix.cfg); err != nil {
				t.Fatal(err)
			}

			exp := hdr + fix.exp + ftr

			// Compare line by line
			explns := strings.Split(exp, "\n")
			buflns := strings.Split(buf.String(), "\n")

			explen, buflen := len(explns), len(buflns)
			if explen != buflen {
				t.Errorf("expected %d lines, but got %d", explen, buflen)
			}

			printContext := func(idx int) {
				lines := 3

				firstIdx := idx - lines
				if firstIdx < 0 {
					firstIdx = 0
				}
				lastIdx := idx + lines
				if lastIdx > explen {
					lastIdx = explen
				}
				if lastIdx > buflen {
					lastIdx = buflen
				}
				t.Logf("expected:\n")
				for i := firstIdx; i < lastIdx; i++ {
					t.Logf(explns[i])
				}
				t.Logf("got:\n")
				for i := firstIdx; i < lastIdx; i++ {
					t.Logf(buflns[i])
				}
			}

			for i, expln := range explns {
				if i >= buflen {
					printContext(i)
					t.Fatalf(`first unequal line: expected "%s" but got EOF`, expln)
				}
				bufln := buflns[i]
				if bufln != expln {
					printContext(i)
					t.Fatalf(`first unequal line: expected "%s" but got "%s"`, expln, bufln)
				}
			}

			if buflen > explen {
				printContext(explen - 1)
				t.Fatalf(`first unequal line: expected EOF, but got "%s"`, buflns[explen])
			}
		})
	}
}

package docs

// ScalarHTML returns the HTML page that embeds Scalar API Reference loading OpenAPI from /openapi.json.
func ScalarHTML() string {
	return `<!DOCTYPE html>
<html lang="tr">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Gateway API – Scalar</title>
  <style> body { margin: 0; } </style>
</head>
<body>
  <div id="app"></div>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  <script>
    (function() {
      var base = window.location.origin;
      Scalar.createApiReference('#app', {
        url: base + '/openapi.json?t=' + Date.now(),
        theme: 'default',
        layout: 'modern',
        hideSidebar: false,
        hideModels: false,
        hideDownloadButton: false,
        searchHotKey: 'k'
      });
    })();
  </script>
</body>
</html>`
}

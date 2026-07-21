package docs

// ScalarHTML returns the HTML page that embeds Scalar API Reference.
func ScalarHTML() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Kuran AI Gateway – Scalar API Docs</title>
  <style> body { margin: 0; } </style>
</head>
<body>
  <div id="app"></div>
  <!-- Pinned version: modal reopen bug fix (unpinned had backdrop stuck on 2nd open) -->
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.43.15"></script>
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

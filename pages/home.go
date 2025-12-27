package pages

import (
	"net/http"
	"ti1be/handlers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	lrw, logResponse := handlers.LogRequestWithWriter(w, r)
	defer logResponse()

	lrw.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ti1be</title>
</head>
<body>
    <div class="container">
        <h1>Ti1</h1>
        <p>Burde gjøre dette bedre... men dette funker da for nå...</p>
    </div>
</body>
</html>`
	lrw.Write([]byte(html))
}

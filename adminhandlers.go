package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var metricsTemplate = template.Must(template.New("page").Parse(`<html>
  <body>
    <h1>Welcome, Admin!</h1>
    <p>Finance Calculator has been visited {{.Count}} times!</p>
  </body>
</html>
`))

func (cfg *apiConfig) handlerMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	dataMap := map[string]string{"Count": fmt.Sprint(cfg.fileserverHits.Load())}
	metricsTemplate.Execute(writer, dataMap)
}

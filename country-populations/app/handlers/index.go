package handlers

import (
	"country-populations/app/db"
	"fmt"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
	"net/http"
)

type IndexHandler struct {
	logger                  *zapray.Logger
	getPopulationForCountry func(string) (int64, bool, error)
}

func NewIndexHandler(logger *zapray.Logger, store db.DynamoDBPopulationStore) IndexHandler {
	return IndexHandler{
		logger:                  logger,
		getPopulationForCountry: store.GetPopulation,
	}
}

func (h IndexHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" || request.URL.Path != "" {
		h.logger.Warn("received request for unsupported path", zap.String("path", request.URL.Path))
		responseWriter.WriteHeader(404)
		return
	}

	country := request.URL.Query().Get("country")
	if country == "" {
		country = "united kingdom"
	}

	population, ok, err := h.getPopulationForCountry(country)
	if err != nil {
		h.logger.Error(
			"could not get population for country",
			zap.String("country", country),
			zap.Error(err),
		)
		responseWriter.WriteHeader(500)
		_, _ = responseWriter.Write([]byte("Something went wrong! " + err.Error()))
		return
	}

	var populationPtr *int64 = nil
	if ok {
		populationPtr = &population
	}

	_, err = responseWriter.Write([]byte(h.buildModel(country, populationPtr)))

	if err != nil {
		h.logger.Error("could not write response", zap.String("country", country), zap.Error(err))
		responseWriter.WriteHeader(500)
		return
	}

	h.logger.Info("responded successfully", zap.String("country", country))
}

func (h IndexHandler) buildModel(countryName string, population *int64) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <title>Country Populations</title>
  </head>
  <body>
    <main
        style="width: 400px;
        margin: 0 auto;
        text-align: center;
        border-radius: 5px;
        border: 2px solid black;
        margin-top: 3em;"
    >
  	  <h1
           style="margin-top: 0;
           margin-bottom: 0.5em;
           padding: 0.5em;
           color: white;
           background: #008fb3;
           border-radius-top-left: 5px;
           border-radius-top-right: 5px;"
      >
        %s
      </h1>
  	  <p style="font-size: 1.5em;">Population: %s</p>
    </main>
  </body>
</html>
`, countryName, h.formatPopulation(population))
}

func (h IndexHandler) formatPopulation(population *int64) string {
	if population == nil {
		return "unknown"
	}
	return h.formatWithCommaSeparators(*population)
}

func (h IndexHandler) formatWithCommaSeparators(population int64) string {
	if population < 0 {
		return "-" + h.formatWithCommaSeparators(-population)
	}
	if population < 1_000 {
		return fmt.Sprintf("%d", population)
	}
	return h.formatWithCommaSeparators(population/1_000) +
		"," +
		fmt.Sprintf("%03d", population%1000)
}

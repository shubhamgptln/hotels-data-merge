package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/shubhamgptln/hotels-data-merge/app/adapter/http/handler"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/cache/inmemory"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/extapi/acme"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/extapi/paperflies"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/extapi/patagonia"
	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/httpapi"
	"github.com/shubhamgptln/hotels-data-merge/app/usecase"
	"github.com/shubhamgptln/hotels-data-merge/bootstrap"
)

const ConfigFilePath = "./config.yml"

func main() {
	// Setup configurations
	appConf, err := bootstrap.LoadConfig(ConfigFilePath)
	if err != nil {
		panic(err)
	}

	httClient := bootstrap.NewHTTPClient(appConf)
	caller := httpapi.NewCaller(httClient)
	acmeClient := acme.NewACMEClient(appConf.HTTP.ACMEClient.EndPoint, caller)
	patagoniaClient := patagonia.NewPatagoniaClient(appConf.HTTP.PatagoniaClient.EndPoint, caller)
	paperfliesClient := paperflies.NewPaperfliesClient(appConf.HTTP.PaperfliesClient.EndPoint, caller)

	inMemoryCache := inmemory.NewInMemoryCache()
	dataProcurer := usecase.NewDataProcurer(inMemoryCache, acmeClient, patagoniaClient, paperfliesClient)
	dataMerger := usecase.NewDataMerger(appConf, dataProcurer)

	router := chi.NewRouter()
	router.Route("/hotels-data", func(ar chi.Router) {
		ar.Get("/filter", handler.FetchHotelsInformation(dataMerger))
	})
	// Start http server
	server := http.Server{
		Handler:      router,
		Addr:         appConf.App.APIPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	var httpErrChan chan error
	go func() {
		logrus.Info("Application up!")
		httpErrChan <- server.ListenAndServe()
	}()

	// Handle terminate signals
	stopSignal := make(chan os.Signal, 2)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopSignal:
		err := server.Shutdown(context.Background())
		if err != nil {
			logrus.Fatal("Cannot stop gracefully: ", err)
		}
	case err := <-httpErrChan:
		logrus.Fatal("Cannot run http server: ", err)
	}
	logrus.Info("Application stopped successfully!")
}

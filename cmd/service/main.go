package main

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/klimenkoOleg/websocket_funout/internal/handler/handle_devices"
	send_message "github.com/klimenkoOleg/websocket_funout/internal/handler/send_message"
	"github.com/klimenkoOleg/websocket_funout/internal/infra/logger"
	"github.com/klimenkoOleg/websocket_funout/internal/message"
	"github.com/klimenkoOleg/websocket_funout/internal/server"
	"github.com/klimenkoOleg/websocket_funout/internal/storage"
)

/*
	errorCause "go.avito.ru/gl/error-cause"
	"go.avito.ru/gl/json-http-protocol/v3/server"
	"go.avito.ru/gl/logger"
	loggerV3 "go.avito.ru/gl/logger/v3"
	metricsV3 "go.avito.ru/gl/metrics/v3"
	nfrMiddleware "go.avito.ru/gl/nfr/v2/middleware"
	"go.avito.ru/gl/platform"
	"go.avito.ru/gl/transport-http/v2/listener"
	"go.avito.ru/gl/transport-http/v2/zstd"

	"go.avito.ru/av/service-developments-search/internal/app"
	developments_domofond_hider_client "go.avito.ru/av/service-developments-search/internal/clients/developments_domofond_hider"
	geograph_client "go.avito.ru/av/service-developments-search/internal/clients/geograph"
	image_storage_client "go.avito.ru/av/service-developments-search/internal/clients/image_storage"
	"go.avito.ru/av/service-developments-search/internal/developer_profile"
	"go.avito.ru/av/service-developments-search/internal/development/catalog_development"
	"go.avito.ru/av/service-developments-search/internal/generated/rpc/service"
	"go.avito.ru/av/service-developments-search/internal/handler/items_delete_inactive"
	"go.avito.ru/av/service-developments-search/internal/handler/items_import"
	"go.avito.ru/av/service-developments-search/internal/infrastructure/middleware"
	"go.avito.ru/av/service-developments-search/internal/rpc/get_development_by_address_id"
	"go.avito.ru/av/service-developments-search/internal/rpc/get_development_by_key"
	"go.avito.ru/av/service-developments-search/internal/rpc/get_developments"
	"go.avito.ru/av/service-developments-search/internal/rpc/search_developments"
	"go.avito.ru/av/service-developments-search/internal/rpc/suggest_developments"
)*/

const (
	developmentsCacheName        = "developments-cache"
	developmentsCatalogCacheName = "developments-catalog-cache"
	developmentsCacheCap         = 20000

	itemsVisibilityCacheName      = "items-visibility-cache"
	itemsVisibilityCacheCap       = 100000
	itemsVisibilitiesLruCacheTTLM = 60

	itemsCountCacheName    = "items-count-cache"
	itemsCountCacheCap     = 100000
	itemsCountLruCacheTTLM = 60

	developmentsItemsCacheName = "developments-items"
	developmentsItemsCacheCap  = 100000
	developmentsItemsCacheTTLM = 60
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	var err error
	logger := logger.MustInitLogger()
	defer logger.Sync() // flush buffer

	logger.Debug()

	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Error("recover", zap.Reflect("recover error", panicErr))
			os.Exit(1)
		}

		if err != nil {
			logger.Error("error left", zap.Error(err))
			os.Exit(1)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	dispatch := make(chan message.Message)
	// quit := make(chan struct{})
	// sendMessageHandler := se.New()
	deviceStorage := storage.New(dispatch, logger) // todo quit
	deviceStorage.Start(ctx)

	sendMessageHandler := send_message.New(deviceStorage, dispatch)
	devicesHandler := handle_devices.New(deviceStorage, logger)

	mux := http.DefaultServeMux
	mux.Handle("/send", sendMessageHandler)
	// http.HandleFunc("/send", sendMessage)
	mux.Handle("/ws", devicesHandler)

	serverListener := server.New(
		server.WithLogger(logger),
	)

	logger.Fatal(serverListener.Listen(ctx, mux))

	// log.Fatal(http.ListenAndServe(":8080", nil))

	// Clients
	// developmentsDomofondHiderClient := developments_domofond_hider_client.New(metricV3)
	// geographClient := geograph_client.New(metricV3)
	// imageStorageClient := image_storage_client.New(metricV3)

	// Services
	// itemAggregator := app.MustNewItemAggregator(ctx, logV3, metricV3, grace)
	// developersService := developer_profile.New(developersClient, decoratedLogger)

	// if err != nil {
	// 	log.WithError(err).Fatal("failed to init geo service")
	// }
	// defer geoCancel()

	// developmentService := catalog_development.New(getStorage, geoService, decoratedLogger)

	// RPC
	// srv := rpcServer(observer, metricV3, parsedNfr)
	// service.SearchDevelopments(search_developments.New(searchStorage, metroTimer, developersClient, grace, decoratedLogger, tglService, developmentsSearchIndexClient, ratingsModelClient).Handle)(srv)
	// service.GetDevelopments(get_developments.New(getStorage, cachedRepo, grace, developersService, decoratedLogger).Handle)(srv)
	// service.SuggestDevelopments(suggest_developments.New(suggestStorage, suggestQuerySanitizer).Handle)(srv)
	// service.GetDevelopmentByKey(get_development_by_key.New(dcHttpClient, grace, itemStorage, goldenItemFinder, developmentsDomofondHiderClient, tglService, metroTimer).Handle)(srv)
	// service.GetDevelopmentByAddressId(get_development_by_address_id.New(getStorage).Handle)(srv)
	/*
		itemsImportHandler := items_import.New(qaasClient)
		itemsDeleteInactiveHandler := items_delete_inactive.New(qaasClient)

		// Mux handlers
		mux := http.DefaultServeMux
		mux.Handle("/items/import", itemsImportHandler)                  // handle running import all items
		mux.Handle("/items/delete_inactive", itemsDeleteInactiveHandler) // handle running deleting inactive items
	*/
	// Listen() стартует сервер с регистрацией Checkers,
	// участвующих в readinessProbe: platform.WithChecker().
	// Опциональный мультиплексор (по умолчанию http.DefaultServeMux):
	// platform.WithMux().
	/*err = app.Listen(
		ctx,
		httpListener(log, logV3, metricV3, parsedNfr),
		platform.WithReporter(decoratedLogger),
	)

	if err != nil {
		log.WithError(err).Fatal("failed to run application")
	}*/
}

// func sendMessage(w http.ResponseWriter, r *http.Request) {
//
// }
/*
func httpListener(log *logger.Logger, logV3 *loggerV3.Logger, metric *metricsV3.Metrics, parsedNfr *nfrMiddleware.Nfr) *listener.HTTPListener {
	mw := []func(handler http.Handler) http.Handler{
		middleware.ServerMetricsWithSourceMW(metric),
		// Cross-cutting контекст, доступный в Observer.
		observability.ServerMW,
		// Метрики всех HTTP запросов.
		// NOTE: Если вы используете api-composition, использовать это middleware не нужно,
		// api-composition отправляет метрики в таком же формате,
		// но учитывает динамические URL paths.
		// Если сервис использует только api-composition, нужно добавить router.WithStdMetrics() в инициализацию роутера
		platform.ServerMetricsMW(metric),
		// Поддержка Content-Encoding, Accept-Encoding: zstd
		// Всегда должен вызываться после всех записей в response
		zstd.ServerMW(),
		// Проброс логгера в контекст
		middleware.LogToContext(log),
		// Собирается метрики на основе nfr.toml, по которым строится SLA сервиса
		// Всегда должен быть последним в списке
		nfrMiddleware.Middleware(parsedNfr),
	}

	return listener.New(
		// Middleware для сервера
		listener.WithMW(mw...),
		// Logger для логирования информации о panic в обработке запросов
		listener.WithLogger(logV3),
	)
}

*/

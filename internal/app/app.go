package app

import (
	"fmt"
	"github.com/Vic07Region/avito-shop/internal/app/handlers"
	"github.com/Vic07Region/avito-shop/internal/app/mw"
	"github.com/Vic07Region/avito-shop/internal/service"
	"github.com/Vic07Region/avito-shop/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type Handlers interface {
	AuthUser(c *gin.Context)
	WalletInfo(c *gin.Context)
	SendCoin(c *gin.Context)
	BuyMerch(c *gin.Context)
}

type AvitoShop struct {
	storage  service.StorageInterface
	service  handlers.ServiceInterface
	handlers Handlers
	gin      *gin.Engine
	logger   *zap.Logger
}

func New() (*AvitoShop, error) {
	app := &AvitoShop{}

	app.gin = gin.Default()

	mode := os.Getenv("APP_MODE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSslMode := os.Getenv("DB_SSLMODE")
	dbSslRootCert := os.Getenv("DB_SSLROOTCERT")
	dbMaxOpenConns := os.Getenv("DB_MAXOPENCONNS")
	dbMsxIdleConns := os.Getenv("DB_MSXIDLECONNS")
	dbMaxLifeTime := os.Getenv("DB_MAXLIFETIME")

	connStr := fmt.Sprintf("user=%s password=%s port=%s dbname=%s",
		dbUser, dbPassword, dbPort, dbName)

	if dbHost != "" {
		connStr += fmt.Sprintf(" host=%s", dbHost)
	}
	if dbSslMode != "" {
		connStr += fmt.Sprintf(" sslmode=%s", dbSslMode)
	} else {
		connStr += " sslmode=disable"
	}

	var maxOpenConns int
	var msxIdleConns int
	var maxLifeTime time.Duration

	if dbMaxOpenConns != "" {
		moc, err := strconv.Atoi(dbMaxOpenConns)
		if err != nil {
			app.logger.Error("MaxOpenConns strconv.Atoi error", zap.Error(err))
			return nil, err
		}
		maxOpenConns = moc
	}
	if dbMsxIdleConns != "" {
		mic, err := strconv.Atoi(dbMsxIdleConns)
		if err != nil {
			app.logger.Error("MsxIdleConns strconv.Atoi error", zap.Error(err))
			return nil, err
		}
		msxIdleConns = mic
	}
	if dbMaxLifeTime != "" {
		mft, err := strconv.Atoi(dbMaxLifeTime)
		if err != nil {
			app.logger.Error("MaxLifeTime strconv.Atoi error", zap.Error(err))
			return nil, err
		}
		maxLifeTime = time.Second * time.Duration(mft)
	}

	if maxOpenConns != 0 {
		connStr += fmt.Sprintf(" max_open_conns=%d", maxOpenConns)
	}

	if dbSslRootCert != "" {
		connStr += fmt.Sprintf(" sslrootcert=%s", dbSslRootCert)
	}

	if mode != "release" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			app.logger.Error("zap.NewDevelopment error", zap.Error(err))
			return nil, err
		}
		app.logger = logger
		gin.SetMode(gin.DebugMode)
		app.logger.Info("applicatation started in debug mode")
		app.logger.Info("Set environment variable APP_MODE=release for release mode")

	} else {
		logger, err := zap.NewProduction()
		if err != nil {

			app.logger.Error("zap.NewProduction error", zap.Error(err))
			return nil, err
		}
		app.logger = logger
		gin.SetMode(gin.ReleaseMode)
	}
	app.logger.Info("connection string", zap.String("connection_string", connStr))
	dbConn, err := storage.NewDBConection(storage.ConnectionParams{
		DbDriver:         storage.DBPostgres,
		ConnectionString: connStr,
		MaxOpenConns:     maxOpenConns,
		MsxIdleConns:     msxIdleConns,
		MaxLifeTime:      maxLifeTime,
	})
	if err != nil {
		app.logger.Error("DB connection error", zap.Error(err))
		return nil, err
	}

	app.storage = storage.New(dbConn, app.logger)

	app.service = service.New(app.storage, app.logger)
	app.handlers = handlers.New(app.service, app.logger)
	middleware := mw.New(app.storage)

	app.gin.POST("/api/auth", app.handlers.AuthUser)
	mwGroupapp := app.gin.Group("/api/").Use(middleware.AuthMiddleware())
	{
		mwGroupapp.GET("/info", app.handlers.WalletInfo)
		mwGroupapp.POST("/sendCoin", app.handlers.SendCoin)
		mwGroupapp.GET("/buy/:merchName", app.handlers.BuyMerch)
	}

	return app, nil
}

func (app *AvitoShop) Run() error {
	ginAddr := os.Getenv("SERVER_ADDR")
	if ginAddr == "" {
		ginAddr = ":8080"
	}
	err := app.gin.Run(ginAddr)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	app.logger.Info("http server started", zap.String("address", ginAddr))
	return nil
}

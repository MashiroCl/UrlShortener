package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mashirocl/urlshortener/config"
	database "github.com/mashirocl/urlshortener/database/query"
	"github.com/mashirocl/urlshortener/internal/api"
	"github.com/mashirocl/urlshortener/internal/cache"
	"github.com/mashirocl/urlshortener/internal/service"
	"github.com/mashirocl/urlshortener/pkg/shortcode"
	"github.com/mashirocl/urlshortener/pkg/validator"
)

type Application struct {
	e                  *echo.Echo
	db                 *sql.DB
	redisClient        *cache.RedisCache
	urlService         *service.URLService
	urlHandler         *api.URLHandler
	cfg                *config.Config
	shortCodeGenerator *shortcode.ShortCode
}

func (a *Application) Init(filePath string) error {
	cfg, err := config.LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}
	a.cfg = cfg

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return fmt.Errorf("create database error: %w", err)
	}
	a.db = db

	redisClient, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		return err
	}
	a.redisClient = redisClient
	a.shortCodeGenerator = shortcode.NewShortCode(cfg.ShortCode.Length)

	urlService := service.NewURLService(db, a.shortCodeGenerator, cfg.App.DefaultDuration, redisClient, cfg.App.BaseURL)
	a.urlService = urlService
	a.urlHandler = api.NewUrlHandler(a.urlService)

	e := echo.New()
	e.Server.WriteTimeout = cfg.Server.WriteTimeout
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Validator = validator.NewCustomValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/api/url", a.urlHandler.CreateURL)
	e.GET("/:code", a.urlHandler.RedirectURL)
	a.e = e
	return nil
}

func (a *Application) Run() {
	go a.startServer()
	go a.cleanUp()
	a.shutdown()
}

func (a *Application) startServer() {
	if err := a.e.Start(a.cfg.Server.Addr); err != nil {
		log.Println(err)
	}
}

func (a *Application) cleanUp() {
	ticker := time.NewTicker(a.cfg.App.CleanupInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		if err := a.urlService.DeleteURL(context.Background()); err != nil {
			log.Println(err)
		}
	}
}

func (a *Application) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		if err := a.redisClient.Close(); err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.e.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

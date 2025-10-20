package main

import (
	"flag"
	"github.com/alhaos-measurement/api/internal/config"
	"github.com/alhaos-measurement/api/internal/controller"
	"github.com/alhaos-measurement/api/internal/logger"
	"github.com/alhaos-measurement/api/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

func main() {

	// Parse config filename from flags
	configFilenamePointer := flag.String("config", "config.yml", "measurement service api config file name")
	flag.Parse()
	configFilename := *configFilenamePointer

	// Init config
	cfg, err := config.New(configFilename)
	if err != nil {
		panic(err)
	}

	// Init log
	l, err := logger.New(cfg.Logger)
	if err != nil {
		panic(err)
	}

	l.Info("start logging")

	// Init connection pool
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			Database: cfg.DB.Database,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
		},
	})
	if err != nil {
		panic(err)
	}

	// Init repository
	repo := repository.New(pool)

	// Init controller
	ctrl := controller.New(repo)

	// Init router
	router := gin.Default()

	// Register routes
	ctrl.RegisterRoutes(router)

	// Run service
	err = router.Run(cfg.Address)
	if err != nil {
		panic(err)
	}
}

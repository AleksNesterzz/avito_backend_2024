package main

import (
	"avito_backend/internal/config"
	"avito_backend/internal/http-server/handlers/changeBanner"
	"avito_backend/internal/http-server/handlers/createBanner"
	"avito_backend/internal/http-server/handlers/deleteBanner"
	"avito_backend/internal/http-server/handlers/getBanner"
	"avito_backend/internal/http-server/handlers/getUserBanner"
	"avito_backend/internal/storage/postgres"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {

	//загрузка конфига
	cfg := config.MustLoad()
	//storagePath := "user=" + cfg.Db.Username + " password=hbdtkjy2012" + " dbname=" + cfg.Db.Dbname + " sslmode=disable"
	storagePath := "postgres://" + cfg.Db.Username + ":" + os.Getenv("DB_PASSWORD") + "@" + cfg.Db.Host + ":" + cfg.Db.Port + "/" + cfg.Db.Dbname + "?sslmode=disable"
	//инициализация БД
	ctx := context.Background()
	storage, err := postgres.New(ctx, storagePath)
	if err != nil {
		log.Fatal("failed to init storage\n", err)
		os.Exit(1)
	}
	//локальное кеширование с истечением 5-ти минут
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			storage.RW.Lock()
			err = storage.UpdateCache(ctx)
			if err != nil {
				fmt.Println("error updating cache")
			}
			storage.RW.Unlock()
		}
	}()
	//инициализация роутера
	router := chi.NewRouter()

	//обработчики запросов

	router.Get("/user_banner", getUserBanner.New(storage))
	router.Get("/banner", getBanner.New(storage))
	router.Post("/banner", createBanner.New(storage))
	router.Patch("/banner/{id}", changeBanner.New(storage))
	router.Delete("/banner/{id}", deleteBanner.New(storage))

	//настройка сервера

	srv := &http.Server{
		Addr:        ":" + cfg.Port,
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	//запуск сервера
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("failed to start server")
	}

}

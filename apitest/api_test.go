package apitest

import (
	"context"
	"fmt"
	"testing"

	"avito_backend/internal/config"
	"avito_backend/internal/http-server/handlers/changeBanner"
	"avito_backend/internal/http-server/handlers/createBanner"
	"avito_backend/internal/http-server/handlers/deleteBanner"
	"avito_backend/internal/http-server/handlers/getBanner"
	"avito_backend/internal/http-server/handlers/getUserBanner"
	"avito_backend/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ConfigTest struct {
	Pg string `env:"PG_STRING"`
}

type ApiSuite struct {
	suite.Suite
	router  *chi.Mux
	storage *postgres.Storage
}

func (as *ApiSuite) BeforeEach(t provider.T) {

	t.NewStep("Загрузка конфигурации окружения")
	cfg := config.MustLoad()
	storagePath := "user=" + cfg.Db.Username + " password=hbdtkjy2012" + " dbname=" + cfg.Db.Dbname + " sslmode=disable"

	fmt.Println(storagePath)
	var err error
	t.NewStep("Проверка работы базы данных окружения")

	as.storage, err = postgres.New(context.Background(), storagePath)
	if err != nil {
		t.Fatalf("error init database: %v", err)
	}

	as.router = chi.NewMux()

	// Repository
	as.router.Post("/banner", createBanner.New(as.storage))
	as.router.Get("/user_banner", getUserBanner.New(as.storage))
	as.router.Get("/banner", getBanner.New(as.storage))
	as.router.Patch("/banner/{id}", changeBanner.New(as.storage))
	as.router.Delete("/banner/{id}", deleteBanner.New(as.storage))

}

func TestRunApiTest(t *testing.T) {
	suite.RunSuite(t, new(ApiSuite))
}

package apitest

import (
	"avito_backend/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/steinfletcher/apitest"
)

const (
	activeBanner      = "activeBanner"
	inactiveBanner    = "inactiveBanner"
	cashedBanner      = "cashedBanner"
	otherCashedBanner = "otherCashedBanner"
)

func (as *ApiSuite) TestGetUserBanner(t provider.T) {
	t.Title("Тестирование апи метода GetUserBanner: GET /user_banner")
	t.NewStep("Инициализация тестовых данных")
	const path = "/user_banner"

	bannerContentList := map[string]string{
		activeBanner:      `{"content":{"title": "default_banner", "text": "some_text"}, "message":"OK"}`,
		cashedBanner:      `{"content":{"title": "cached_banner", "text": "some_text"}, "message":"OK"}`,
		otherCashedBanner: `{"content":{"title": "other_cached_banner", "text": "some_text"}, "message":"OK"}`,
		inactiveBanner:    `{"content":{"title": "disabled_banner", "text": "some_text"}, "message":"OK"`,
	}
	tim := time.Now()

	bannerList := map[string]*models.Banner{
		activeBanner: {
			Fid:        1,
			Tags:       []int{2, 4, 3},
			Is_active:  true,
			Bid:        1,
			Cnt:        models.Content{Title: "default_banner", Text: "some_text"},
			Created_at: tim,
			Updated_at: tim,
		},
		cashedBanner: {
			Fid:        3,
			Tags:       []int{5, 1, 2},
			Is_active:  true,
			Bid:        2,
			Cnt:        models.Content{Title: "cached_banner", Text: "some_text"},
			Created_at: tim,
			Updated_at: tim,
		},
		otherCashedBanner: {
			Fid:        4,
			Tags:       []int{5, 1, 2},
			Is_active:  true,
			Bid:        3,
			Cnt:        models.Content{Title: "other_cached_banner", Text: "some_text"},
			Created_at: tim,
			Updated_at: tim,
		},
		inactiveBanner: {
			Fid:        2,
			Tags:       []int{1, 4, 5},
			Is_active:  false,
			Bid:        4,
			Cnt:        models.Content{Title: "inactive_banner", Text: "some_text"},
			Created_at: tim,
			Updated_at: tim,
		},
	}

	for _, bn := range bannerList {
		for i := 0; i < len(bn.Tags); i++ {
			_, err := as.storage.CreateBanner(context.Background(), bn.Tags[i], bn.Fid, bn.Bid, bn.Cnt.Title, bn.Cnt.Text, bn.Cnt.URL, bn.Is_active, bn.Created_at)
			t.Require().NoError(err)
		}

	}
	_ = bannerList
	t.Run("Успешное получение активного баннера", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("tid", "2").Query("fid", "1").Query("last_revision", "true").
			Header("token", "user").
			Expect(t).
			Body(bannerContentList[activeBanner]).
			Status(http.StatusOK).
			End()
	})

	t.Run("Успешное получение активного баннера с обязательным запросом к базе после кэширования", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").Query("tid", "1").
			Header("token", "user").
			Expect(t).
			Body(`{"message":"Баннер не найден"}`).
			Status(http.StatusNotFound).
			End()

		err := as.storage.UpdateCache(context.Background())

		t.Require().NoError(err)

		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").Query("tid", "1").
			Header("token", "user").
			Expect(t).
			Body(bannerContentList[cashedBanner]).
			Status(http.StatusOK).
			End()
	})

	t.Run("Попытка получения неактивного баннера", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "2").Query("tid", "4").
			Header("token", "user").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Попытка получения несуществующего баннера", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").Query("tid", "4").
			Header("token", "user").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Попытка получения баннера неавторизованным пользователем", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").Query("tid", "4").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Попытка получения баннера пользователя с неверными правами", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").Query("tid", "4").
			Header("token", "admin").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Попытка получения баннера с некорректными параметрами запроса", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "mir").Query("tid", "4").
			Header("token", "user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()

		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "4").Query("tid", "mir").
			Header("token", "user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Попытка получения баннера с передачей не всех параметров запроса", func(t provider.T) {
		apitest.New().
			Handler(as.router).
			Get(path).
			Query("tid", "4").
			Header("token", "user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()

		apitest.New().
			Handler(as.router).
			Get(path).
			Query("fid", "3").
			Header("token", "user").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	err := as.storage.FlushAllTEST(context.Background(), tim)
	t.Require().NoError(err)
}

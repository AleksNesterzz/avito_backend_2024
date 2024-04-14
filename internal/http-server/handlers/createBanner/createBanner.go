package createBanner

import (
	"avito_backend/internal/http-server/responces"
	tokens "avito_backend/internal/lib"
	apierr "avito_backend/internal/lib/errors"
	"avito_backend/internal/models"
	"avito_backend/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type BannerCreator interface {
	CreateBanner(ctx context.Context, tid int, fid, bid int, title, descr, url string, is_active bool, t time.Time) (int64, error)
}

type Request struct {
	Tid       []int          `json:"tid"`
	Fid       int            `json:"fid"`
	Bid       int            `json:"bid"`
	Cnt       models.Content `json:"content"`
	Is_active bool           `json:"is_active"`
}

type Response struct {
	Bid     int    `json:"banner_id,omitempty"`
	Message string `json:"message"`
	Error   string `json:"eror,omitempty"`
}

func New(bancreator BannerCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request
		token, err := tokens.GetToken(w, r)
		if err == apierr.ErrNoAuth {
			resp := Response{Message: "Пользователь не авторизован"}
			responces.JSONResponse(w, r, 401, resp)
			return
		} else if token == tokens.AdminToken {
			err = json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp := Response{Message: "Некорректные данные",
					Error: err.Error()}
				responces.JSONResponse(w, r, 400, resp)
				return
			}

			tag_ids := req.Tid
			fid := req.Fid
			bid := req.Bid
			cnt := req.Cnt
			act := req.Is_active
			ctx := r.Context()

			var wg sync.WaitGroup
			//fmt.Println(id, name)
			t := time.Now()
			var err error
			wg.Add(len(tag_ids))
			for i := 0; i < len(tag_ids); i++ {

				go func(x int) {
					defer wg.Done()
					_, err = bancreator.CreateBanner(ctx, tag_ids[x], fid, bid, cnt.Title, cnt.Title, cnt.URL, act, t)
				}(i)
			}
			wg.Wait()
			if errors.Is(err, storage.ErrBannerExists) {

				resp := Response{Message: "Баннер уже существует"}
				responces.JSONResponse(w, r, 400, resp)
				return
			}
			if err != nil {
				resp := Response{Message: "Внутренняя ошибка сервера",
					Error: err.Error()}
				responces.JSONResponse(w, r, 500, resp)
				return
			}
			if err == nil {
				resp := Response{Bid: bid,
					Message: "Баннер создан"}
				responces.JSONResponse(w, r, 201, resp)
				return
			}
		} else {
			resp := Response{
				Message: "Пользователь не имеет доступа"}
			responces.JSONResponse(w, r, 403, resp)
			return
		}
	}
}

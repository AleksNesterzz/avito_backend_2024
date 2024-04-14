package postgres

//TODO токены админов и пользователей
//TODO е2е тестирование
import (
	_ "avito_backend/internal/lib/csv_log"
	"avito_backend/internal/models"
	"avito_backend/internal/storage"
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "errors"
)

type Storage struct {
	db    *pgxpool.Pool
	cache map[models.Pair]models.Banner //стоит ли делать так?
	RW    sync.RWMutex
}

var (
	pgInstance *Storage
	pgOnce     sync.Once
)

func New(ctx context.Context, storagePath string) (*Storage, error) {

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, storagePath)
		if err != nil {
			fmt.Println("unable to create connection pool:", err)
			return
		}
		cache := make(map[models.Pair]models.Banner)
		pgInstance = &Storage{db, cache, sync.RWMutex{}}
	})

	_, err := pgInstance.db.Exec(ctx, `CREATE TABLE IF NOT EXISTS banners(
		tag_id INTEGER NOT NULL,
		feature_id INTEGER NOT NULL,
		banner_id integer NOT NULL, 
		title text NOT NULL,
		descr TEXT NOT NULL,
		url TEXT NOT NULL,
		is_active bool NOT NULL,
		created_at timestamp,
		updated_at timestamp,
		PRIMARY KEY(tag_id,feature_id));`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement error_1: %s", err)
	}

	return pgInstance, nil
}

func (s *Storage) CreateBanner(ctx context.Context, tag_id, fid, banner_id int, title, descr, url string, is_act bool, t time.Time) (int64, error) {
	str := "INSERT INTO banners(tag_id, feature_id, banner_id, title, descr, url, is_active, created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9);"
	_, err := s.db.Exec(ctx, str, tag_id, fid, banner_id, title, descr, url, is_act, t, t)
	if err != nil {
		if strings.Contains(err.Error(), "pkey") {
			return 0, storage.ErrBannerExists
		}
		return 0, fmt.Errorf("error executing insert %w", err)
	}

	return 0, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, banner_id int) (int64, error) {
	str := "DELETE FROM banners WHERE banner_id=$1;"

	res, err := s.db.Exec(ctx, str, banner_id)
	if err != nil {
		return 0, fmt.Errorf("error executing deletion %w", err)
	}
	if res.RowsAffected() == 0 {
		return 0, storage.ErrBannerNotFound
	}

	return res.RowsAffected(), nil

}

func (s *Storage) ChangeBanner(ctx context.Context, bid int, cnt models.Content, act *bool) (string, error) {
	var str string = "UPDATE banners SET"
	title := cnt.Title
	text := cnt.Text
	url := cnt.URL
	if title != "" {
		str += " title='" + title + "',"
	}
	if text != "" {
		str += " descr='" + text + "',"
	}
	if url != "" {
		str += " url='" + url + "',"
	}
	if act != nil {
		if *act {
			str += " is_active=true,"
		} else {
			str += " is_active=false,"
		}
	}
	t := time.Now()
	str += " updated_at = $1 WHERE banner_id=$2"
	rows, err := s.db.Exec(ctx, str, t, bid)
	if err != nil {
		return "", fmt.Errorf("error preparing updating banners %w", err)
	}

	if rows.RowsAffected() == 0 {
		return "", storage.ErrBannerNotFound
	}

	return "", nil
}

func (s *Storage) GetUserBanner(ctx context.Context, tag, fid int, last bool) (*models.Content, error) {
	banner := models.Content{}
	var is_active bool
	str := "SELECT title,descr,url,is_active from banners WHERE tag_id = $1 AND feature_id = $2"
	if last {
		rows, err := s.db.Query(ctx, str, tag, fid)
		if err != nil {
			return nil, fmt.Errorf("error executing GetUserBanner %w", err)
		} else {
			for rows.Next() {
				rows.Scan(&banner.Title, &banner.Text, &banner.URL, &is_active)
			}
		}
		defer rows.Close()
		if rows.CommandTag().RowsAffected() == 0 {
			return nil, storage.ErrBannerNotFound
		}
		if is_active {
			return &banner, nil
		} else {
			return nil, storage.ErrBannerNotFound
		}
	} else {
		//запись из кэша
		s.RW.Lock()
		var temp models.Pair
		var cnt models.Content
		temp.Tag = tag
		temp.Fid = fid
		if banner, ok := s.cache[temp]; ok {
			s.RW.Unlock()
			cnt.Text = banner.Cnt.Text
			cnt.Title = banner.Cnt.Title
			cnt.URL = banner.Cnt.URL
			if banner.Is_active {
				return &cnt, nil
			} else {
				return nil, storage.ErrBannerNotFound
			}
		} else {
			s.RW.Unlock()
			return nil, storage.ErrBannerNotFound
		}

	}

}

func (s *Storage) GetBanner(ctx context.Context, tag, fid, limit, offset *int) ([]models.Banner, error) {
	//SELECT array_agg(tag_id), feature_id, banner_id, title, descr, url, created_at, updated_at, is_active from banners
	//group by feature_id,banner_id, title, descr, url, created_at, updated_at, is_active
	//having 4 = any(array_agg(tag_id))
	var banner []models.Banner
	str := "SELECT array_agg(tag_id), feature_id, banner_id, title,descr,url,created_at, updated_at, is_active from banners"

	if fid != nil {
		str_fid := strconv.Itoa(*fid)
		str += " WHERE feature_id=" + str_fid
	}
	str += " GROUP BY feature_id,banner_id, title, descr,url, created_at, updated_at, is_active"
	if tag != nil {
		str_tag := strconv.Itoa(*tag)
		str += " HAVING " + str_tag + "=ANY(array_agg(tag_id))"
	}

	if limit != nil {
		str_limit := strconv.Itoa(*limit)
		str += " LIMIT " + str_limit
	}
	if offset != nil {
		str_offset := strconv.Itoa(*offset)
		str += " OFFSET " + str_offset
	}

	rows, err := s.db.Query(ctx, str)
	if err != nil {
		return []models.Banner{}, fmt.Errorf("error executing GetBanner %w", err)
	} else {
		var temp models.Banner
		for rows.Next() {
			rows.Scan(&temp.Tags, &temp.Fid, &temp.Bid, &temp.Cnt.Title, &temp.Cnt.Text, &temp.Cnt.URL, &temp.Created_at, &temp.Updated_at, &temp.Is_active)
			banner = append(banner, temp)
		}
	}

	defer rows.Close()
	if banner == nil {
		return nil, storage.ErrBannerNotFound
	} else {
		return banner, nil
	}

}

func (s *Storage) UpdateCache(ctx context.Context) error {
	var pr models.Pair
	var cnt models.Content
	var is_active bool
	str := "SELECT tag_id,feature_id,title,descr,url,is_active from banners"
	rows, err := s.db.Query(ctx, str)
	if err != nil {
		return err
	} else {
		for rows.Next() {
			rows.Scan(&pr.Tag, &pr.Fid, &cnt.Title, &cnt.Text, &cnt.URL, &is_active)
			s.cache[pr] = models.Banner{Tags: []int{pr.Tag}, Fid: pr.Fid, Cnt: models.Content{Title: cnt.Title, Text: cnt.Text, URL: cnt.URL}, Is_active: is_active}
		}
	}
	defer rows.Close()
	return nil
}

func (s *Storage) FlushAllTEST(ctx context.Context, t time.Time) error {
	str := "DELETE FROM banners WHERE created_at=$1"
	_, err := s.db.Exec(ctx, str, t)
	if err != nil {
		return err
	} else {
		return nil
	}
}

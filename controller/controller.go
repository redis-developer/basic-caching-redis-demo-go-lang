package controller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"time"
)

type gitRepo struct {
	PublicRepos int `json:"public_repos"`
}

type Controller struct {
	r      Redis
	client *http.Client
}

func (c Controller) SetNonCachedDuration(username string, duration time.Duration) error {
	return c.r.Set(fmt.Sprintf("duration.%s", username), duration.String(), 0)
}

func (c Controller) GetNonCachedDuration(username string) (time.Duration, error) {
	value, err := c.r.Get(fmt.Sprintf("duration.%s", username))
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func (c Controller) requestRepo(username string) (*gitRepo, error) {
	URL := fmt.Sprintf("https://api.github.com/users/%s", username)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code [%d]", res.StatusCode)
	}
	repo := &gitRepo{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (c Controller) GetRepo(username string) (*Repo, error) {
	value, err := c.r.Get(username)
	if err == redis.Nil {
		repo, err := c.requestRepo(username)
		if err != nil {
			return nil, err
		}
		err = c.r.Set(username, strconv.Itoa(repo.PublicRepos), time.Hour)
		if err != nil {
			return nil, err
		}
		return &Repo{
			Username: username,
			Repos:    repo.PublicRepos,
			Cached:   false,
		}, nil
	} else if err != nil {
		return nil, err
	}
	repos, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	return &Repo{
		Username: username,
		Repos:    repos,
		Cached:   true,
	}, nil
}

var controller = &Controller{
	client: &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxConnsPerHost:     15,
			MaxIdleConns:        15,
			MaxIdleConnsPerHost: 15,
			IdleConnTimeout:     time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	},
}

func Instance() *Controller {
	return controller
}

func SetRedis(redis Redis) {
	controller.r = redis
}

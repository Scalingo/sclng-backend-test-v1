package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/google/go-github/v60/github"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)
	router.HandleFunc("/listRepo", listLastHandredPublicRepo)
	router.HandleFunc("/listAgragateRepo", agregateLastHandredPublicRepo)
	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}

}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

func listLastHandredPublicRepo(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	filter := getUrlFilter(r.URL.RawQuery)
	log.Info("The filter is: " + filter)

	client := github.NewClient(nil)

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	// Format time with RFC3339 (ISO8601)
	t := time.Now().Format(time.RFC3339)
	// Get last 100 public repos, no auth is required for public repos, vs created date
	repos, _, err := client.Search.Repositories(context.Background(), filter+"created:<"+t, opt)
	if err != nil {
		log.WithError(err).Error("Fail to list Repos")
		return err
	}
	//jsonContent, err := json.MarshalIndent(repos, "", "    ")
	err = json.NewEncoder(w).Encode(repos.Repositories)
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
		return err
	}
	return nil
}

type AgrRepo struct {
	FullName   *string         `json:"full_name,omitempty"`
	Owner      *string         `json:"owner,omitempty"`
	Repository *string         `json:"repository,omitempty"`
	Languages  json.RawMessage `json:"languages,omitempty"`
}

func agregateLastHandredPublicRepo(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	filter := getUrlFilter(r.URL.RawQuery)
	log.Info("The filter is: " + filter)

	client := github.NewClient(nil)

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	// Format time with RFC3339 (ISO8601)
	t := time.Now().Format(time.RFC3339)
	// Get last 100 public repos (no auth is required) vs created date
	repos, _, err := client.Search.Repositories(context.Background(), filter+"created:<"+t, opt)
	if err != nil {
		log.WithError(err).Error("Fail to list Repos")
		return err
	}
	var listAgrRepos []AgrRepo
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex for synchronization
	// Add to WaitGroup for each goroutine
	wg.Add(len(repos.Repositories))

	for _, repo := range repos.Repositories {
		go addRepoAgregateData(&listAgrRepos, repo, &mu, &wg)
	}
	// Wait for all goroutines to finish
	wg.Wait()

	err = json.NewEncoder(w).Encode(listAgrRepos)
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
		return err
	}
	return nil
}

func addRepoAgregateData(listAgrRepos *[]AgrRepo, repo *github.Repository, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log := logger.Get(context.Background())
	bodyStatsLang, err := getLanguagesStats(repo.GetLanguagesURL())
	if err != nil {
		log.WithError(err).Error("Fail to get languages stats")
	}
	agrRepos := &AgrRepo{
		FullName:   repo.FullName,
		Owner:      repo.Owner.Login,
		Repository: repo.Name,
		Languages:  bodyStatsLang,
	}
	// Protect access to the slice with a mutex
	mu.Lock()
	*listAgrRepos = append(*listAgrRepos, *agrRepos)
	mu.Unlock()
	log.Info(" HEre addRepoAgregateData")
}
func getUrlFilter(urlRawQuery string) string {
	log := logger.Get(context.Background())
	params, err := url.ParseQuery(urlRawQuery)
	if err != nil {
		log.WithError(err).Error("Fail to parse query params")
	}
	filter := ""
	for key, value := range params {
		// we gess that every params has only one value exp. language=Java&license=mit
		filter = filter + key + ":" + value[0] + ","
	}
	return filter
}

func getLanguagesStats(LanguagesURL string) (json.RawMessage, error) {
	log := logger.Get(context.Background())
	req, err := http.NewRequest("GET", LanguagesURL, nil)

	if err != nil {
		log.WithError(err).Error("Fail to create req to get languages")
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("Fail to get languages")
		return nil, err
	}
	defer resp.Body.Close()
	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Fail Read body languages ")
		return nil, err

	}
	var body json.RawMessage
	if err := json.Unmarshal(byteBody, &body); err != nil {
		log.WithError(err).Error("Fail to Unmarshal languages")
		return nil, err
	}
	return body, nil
}

package meetup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/agonzalezro/gotagmee/db"
	"github.com/jmcvetta/neoism"
)

type API struct {
	token, groupURLName string

	db *neoism.Database

	paginationSize int
	client         http.Client

	debug bool
}

func NewAPI(token, groupURLName string) *API {
	return &API{
		token:          token,
		groupURLName:   groupURLName,
		paginationSize: 100,
		client:         http.Client{},
		debug:          os.Getenv("DEBUG") != "",
	}
}

func (a API) endpoint(method string) string {
	u := "https://api.meetup.com"
	return u + method
}

func (a API) doRequest(endpoint string, v *url.Values) (*http.Response, error) {
	if v == nil {
		v = &url.Values{}
	}
	for key, value := range map[string]string{
		"sign":          "true",
		"key":           a.token,
		"group_urlname": a.groupURLName,
		"page":          strconv.Itoa(a.paginationSize),
	} {
		v.Set(key, value)
	}

	url := endpoint + "?" + v.Encode()
	if a.debug {
		log.Println("DEBUG:", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 429 { // Too many requests
		tresshold := 5 * time.Second

		reset := resp.Header.Get("X-RateLimit-Reset")
		if reset == "" {
			// It's possible that the previous tresshold was not enough, we are almost there, but not yet!
			time.Sleep(tresshold)
			return a.doRequest(endpoint, v)
		}

		sleep, err := time.ParseDuration(reset + "s")
		if err != nil {
			return resp, err
		}
		log.Printf("Throttled! Sleeping %v...\n", sleep+tresshold)
		time.Sleep(sleep + tresshold)
		return a.doRequest(endpoint, v)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("Expected 200, got: %d (%s)", resp.StatusCode, endpoint)
	}
	return resp, nil
}

func (a API) Members(membersChan chan db.Member) error {
	members, err := a.getMembersCount()
	if err != nil {
		return err
	}
	pages := members / a.paginationSize
	if members%a.paginationSize > 0 {
		pages += 1
	}

	var wg sync.WaitGroup
	wg.Add(pages)

	for i := 0; i < pages; i++ {
		go func(page int) {
			defer wg.Done()

			v := url.Values{}
			v.Set("offset", strconv.Itoa(page))

			resp, err := a.doRequest(a.endpoint("/2/members"), &v)
			defer resp.Body.Close()
			if err != nil {
				fmt.Println("WARNING:", err)
				if a.debug {
					fmt.Println("DEBUG:Header:", resp.Header)
					b, _ := ioutil.ReadAll(resp.Body)
					fmt.Println("DEBUG:Body:", string(b))
				}
				return
			}

			mr := MembersResponse{}
			err = json.NewDecoder(resp.Body).Decode(&mr)
			for _, m := range mr.Results {
				if a.debug {
					log.Println("DEBUG:Member:", m)
					log.Println("DEBUG:Results:", mr.Results)
				}
				var topics []string
				for _, t := range m.Topics {
					topics = append(topics, t.Name)
				}
				membersChan <- db.Member{m.Name, topics}
			}
		}(i)
	}

	wg.Wait()
	close(membersChan)
	return nil
}

func (a API) getMembersCount() (int, error) {
	gr, err := a.groupsMethod()
	if err != nil {
		return 0, err
	}
	return gr.Results[0].Members, nil
}

func (a API) groupsMethod() (*GroupsResponse, error) {
	endpoint := a.endpoint("/2/groups")

	resp, err := a.doRequest(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	gr := GroupsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		return nil, err
	}
	return &gr, nil
}

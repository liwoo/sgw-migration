package sgw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CouchbaseSgwUser struct {
	AdminChannels []string `json:"admin_channels"`
	AdminRoles    []string `json:"admin_roles"`
	Name          string   `json:"name" faker:"email"`
}

type User struct {
	Name          string   `json:"name"`
	Password      string   `json:"password"`
	AdminChannels []string `json:"admin_channels"`
	AllChannels   []string `json:"all_channels"`
	Disabled      bool     `json:"disabled"`
	AdminRoles    []string `json:"admin_roles"`
	Roles         []string `json:"roles"`
}

type Service struct {
	baseURL   string
	errorFile *os.File
	db        string
	apiKey    string
}

func NewService(baseURL string, db string, errorFile *os.File, apiKey string) *Service {
	return &Service{
		baseURL:   baseURL,
		errorFile: errorFile,
		db:        db,
		apiKey:    apiKey,
	}
}

func (s *Service) WriteToSyncGateway(users *[]CouchbaseSgwUser, results chan string, jobs chan int, workerId int, client *http.Client) {
	for job := range jobs {

		user := (*users)[job]

		sgwUsersEndpoint := s.baseURL + "/" + s.db + "/_user/"

		sgwUser := User{
			Name:          "test_" + user.Name,
			Password:      "password",
			AdminChannels: user.AdminChannels,
			AllChannels:   user.AdminChannels,
			Disabled:      false,
			AdminRoles:    user.AdminRoles,
			Roles:         user.AdminRoles,
		}

		marshalledJson, err := json.Marshal(sgwUser)
		if err != nil {
			s.writeError(err)
		}

		//create a new request
		req, err := http.NewRequest("POST", sgwUsersEndpoint, bytes.NewBuffer(marshalledJson))

		if err != nil {
			s.writeError(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", s.apiKey)

		//make http request
		resp, err := client.Do(req)
		if err != nil {
			s.writeError(err)
		}

		if resp.StatusCode != 201 {
			err := fmt.Errorf("failed to create user: %s", sgwUser.Name)
			s.writeError(err)
		} else {
			results <- fmt.Sprintf("Worker %d: finished migrating %s, left: %d", workerId, user.Name, len(results))
		}
	}
}

func (s *Service) writeError(error error) {
	_, err := s.errorFile.WriteString(error.Error() + "\n")
	if err != nil {
		fmt.Println(err)
	}
}

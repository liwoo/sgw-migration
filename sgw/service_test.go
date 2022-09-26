package sgw

import (
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"net/http"
	"os"
	"strings"
	"testing"
)

func cleanUp(t *testing.T, fakeUsers []CouchbaseSgwUser, sgwUsersEndpoint string) {
	//delete test users
	for i := 0; i < len(fakeUsers); i++ {
		req, err := http.NewRequest("DELETE", sgwUsersEndpoint+fakeUsers[i].Name, nil)
		if err != nil {
			t.Error(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestWriteToSyncGateway(t *testing.T) {

	_ = faker.SetRandomMapAndSliceSize(4)
	var users []CouchbaseSgwUser
	userNum := 10
	for i := 0; i < userNum; i++ {
		user := CouchbaseSgwUser{}
		err := faker.FakeData(&user)
		if err != nil {
			t.Error(err)
		}
		users = append(users, user)
	}

	errorFile := &os.File{}

	results := make(chan string, userNum)
	jobs := make(chan int, userNum)

	const syncGatewayURL = "http://167.99.211.75:4985"
	const db = "offline_reads"

	service := NewService(syncGatewayURL, db, errorFile, "some_api_key")

	for i := 0; i < userNum; i++ {
		go service.WriteToSyncGateway(&users, results, jobs, i)
	}

	for i := 0; i < userNum; i++ {
		jobs <- i
	}

	close(jobs)

	for i := 0; i < userNum; i++ {
		if !strings.Contains(<-results, "finished migrating") {
			t.Error("error on user ", i)
		}
	}

	const sgwUsersEndpoint = syncGatewayURL + "/" + db + "/_user/"

	//make http request
	resp, err := http.Get(sgwUsersEndpoint)
	if err != nil {
		t.Error(err)
	}

	var sgwUsers []string
	err = json.NewDecoder(resp.Body).Decode(&sgwUsers)
	if err != nil {
		t.Error(err)
	}

	fakeUser := users[0]
	found := false
	for i := 0; i < len(sgwUsers); i++ {
		if sgwUsers[i] == "test_"+fakeUser.Name {
			found = true
		}
	}

	if !found {
		t.Error("user not found")
	}

	cleanUp(t, users, sgwUsersEndpoint)

}

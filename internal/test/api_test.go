package test

import (
	"encoding/json"
	"fmt"
	"gohw/internal/api"
	"gohw/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type TestCase struct {
	Response    string
	Body        string
	CookieValue string
	StatusCode  int
}

// go test -coverprofile test/cover.out
// go tool cover -html=test/cover.out -o test/coverage.html

func TestCreateUser(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Response: `{
				"nickname": "ccc",
				"fullname": "username111",
				"about": "145454389",
				"email": "000"
			}`,
			Body: `{
				"fullname": "username111",
				"about": "145454389",
				"email": "000"
			}`,
			StatusCode: http.StatusCreated,
		},
		TestCase{
			Response: `{
				"nickname": "ddd",
				"fullname": "TestCase2",
				"about": "1454543",
				"email": "00"
				}`,
			Body: `{
				"fullname": "TestCase2",
				"email": "00",
				"about": "1454543"
			}`,
			StatusCode: http.StatusCreated,
		},
		TestCase{
			Response: `{
				"nickname": "ppp",
				"fullname": "TestCase2",
				"about": "1454543",
				"email": "0"}`,
			Body: `{
				"fullname": "TestCase2",
				"email": "0",
				"about": "1454543"
			}`,
			StatusCode: http.StatusCreated,
		},
	}

	nicknames := []string{"ccc", "ddd", "ppp"}

	API, err := api.GetHandler()
	if err != nil {
		fmt.Println("Some error happened with configuration file or database TEST" + err.Error())
		return
	}

	for caseNum, item := range cases {
		url := "/api/user/" + nicknames[caseNum] + "/create"
		req := httptest.NewRequest("POST", url, strings.NewReader(item.Body))
		req = mux.SetURLVars(req, map[string]string{"name": nicknames[caseNum]})
		w := httptest.NewRecorder()

		API.CreateUser(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d\n",
				caseNum, w.Code, item.StatusCode)
			continue
		}

		resp := w.Result()

		var (
			userTest models.User
			userReal models.User
		)
		if err := json.NewDecoder(resp.Body).Decode(&userTest); err != nil {
			t.Errorf("Test error: %s\n", err.Error())
		}
		if err := json.NewDecoder(strings.NewReader(item.Response)).Decode(&userReal); err != nil {
			t.Errorf("Test error: %s\n", err.Error())
		}

		resp.Body.Close()
		if userReal != userTest {
			t.Errorf("Nickname: Got %s  Exp %s\nAbout: Got %s  Exp %s\nFullname: Got %s  Exp %s\nEmail: Got %s  Exp %s\n\n",
				userTest.Nickname, userReal.Nickname,
				userTest.About, userReal.About,
				userTest.Fullname, userReal.Fullname,
				userTest.Email, userReal.Email)
		}

	}
}

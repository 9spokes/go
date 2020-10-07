package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/9spokes/go/api"
	"github.com/9spokes/go/logging"
	"github.com/9spokes/go/types"
)

func TestCreateConnection(t *testing.T) {

	stage := "token.CreateConnection"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			api.ErrorResponse(w, "Invalid Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		osp := r.PostForm.Get("osp")
		user := r.PostForm.Get("user")
		company := r.PostForm.Get("company")
		status := r.PostForm.Get("status")

		if osp == "10spokes" {
			api.ErrorResponse(w, "invalid OSP", http.StatusOK)
			return
		}

		conn := types.Connection{
			User:    user,
			OSP:     osp,
			Company: company,
			Status:  status,
		}

		api.SuccessResponse(w, conn, http.StatusOK)
	}))
	defer ts.Close()

	cases := []struct {
		description string
		ctx         Context
		form        map[string]string
		e           string
		response    string
	}{
		{
			description: "Missing OSP field",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			form:        map[string]string{"user": "test"},
			e:           "the field osp is required",
		},
		{
			description: "Missing User field",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			form:        map[string]string{"osp": "9spokes"},
			e:           "the field user is required",
		},
		{
			description: "Non-OK response",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			form:        map[string]string{"osp": "10spokes", "user": "test"},
			e:           "Non-OK response received from Token service: invalid OSP",
		},
		{
			description: "OK response",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			form:        map[string]string{"osp": "9spokes", "user": "test", "company": "dummy", "status": "NEW", "action": "initiate"},
			response:    `{"id":"","credentials":null,"demo":false,"token":null,"settings":null,"user":"test","config":null,"osp":"9spokes","company":"dummy","created":"0001-01-01T00:00:00Z","modified":"0001-01-01T00:00:00Z","status":"NEW"}`,
		},
	}

	for _, c := range cases {

		fmt.Printf("Testing Case: [%s]: %s\n", stage, c.description)

		if ret, err := c.ctx.CreateConnection(c.form); err != nil && !regexp.MustCompile(c.e).MatchString(err.Error()) {
			t.Fatalf("Test case failed. Expected error to contain '%s', got '%s'", c.e, err.Error())
		} else {
			if c.response != "" {
				encoded, _ := json.Marshal(ret)
				if string(encoded) != c.response {
					t.Fatalf("Test case failed. Unexpected response.  Expected '%s', received '%s'", c.response, encoded)
				}
			}
		}
	}
}


func TestManageConnection(t *testing.T) {

	stage := "token.ManageConnection"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		values := r.URL.Query()
		action := values.Get("action")

		if action == "" {
			api.ErrorResponse(w, r.URL.String(), http.StatusOK)
			return
		}

		if action == "sing" {
			api.ErrorResponse(w, "unsupported action", http.StatusOK)
			return
		}

		api.SuccessResponse(w, "Action " + action + " completed successfully", http.StatusOK)
	}))
	defer ts.Close()

	cases := []struct {
		description string
		ctx         Context
		form        map[string]string
		action      string
		e           string
		response    string
	}{
		{
			description: "Missing action",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			action:      "",
			form:        map[string]string{},
			e:           "the action must be specified",
		},
		{
			description: "Non-OK response",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			action:      "sing",
			form:        map[string]string{},
			e:           "Non-OK response received from Token service: unsupported action",
		},
		{
			description: "OK response",
			ctx:         Context{URL: ts.URL, Logger: logging.New("bogus", "INFO")},
			action:      "authorize",
			form:        map[string]string{"code": "12345"},
			response:    "Action authorize completed successfully",
		},
	}

	for _, c := range cases {

		fmt.Printf("Testing Case: [%s]: %s\n", stage, c.description)

		if err := c.ctx.ManageConnection("uuid", c.action, c.form); err != nil && !regexp.MustCompile(c.e).MatchString(err.Error()) {
			t.Fatalf("Test case failed. Expected error to contain '%s', got '%s'", c.e, err.Error())
		}
	}
}
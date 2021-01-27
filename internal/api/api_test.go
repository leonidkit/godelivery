package api

import (
	"bytes"
	"encoding/json"
	"godelivery/internal/converter/xml2json"
	"godelivery/internal/storage/redis"
	"godelivery/pkg/fakelogger"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	apiserver = &APIServer{}
)

type Case struct {
	Path     string
	Method   string
	Status   int
	Request  map[string]interface{}
	Response string
	wantErr  bool
	Error    string
}

func TestMain(m *testing.M) {
	c := xml2json.New()
	s, err := redis.New("redis://default@localhost:6379/0", time.Second*20)
	if err != nil {
		log.Fatal(err)
	}
	l := fakelogger.New()

	apiserver = New(c, s, l)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAPI(t *testing.T) {
	tcases := []Case{
		Case{
			Path:   "/create",
			Method: http.MethodPost,
			Status: http.StatusBadRequest,
			Request: map[string]interface{}{
				"request_id": "123",
			},
			wantErr: true,
			Error:   "convert request error",
		},
		Case{
			Path:   "/create",
			Method: http.MethodPost,
			Status: http.StatusBadRequest,
			Request: map[string]interface{}{
				"request_id": "",
			},
			wantErr: true,
			Error:   "request error",
		},
		Case{
			Path:   "/create",
			Method: http.MethodPost,
			Status: http.StatusOK,
			Request: map[string]interface{}{
				"request_id":  1234567,
				"format_type": "SF",
				"format":      `<hello>world</hello>`,
			},
			Response: `{"message":"OK"}` + "\n",
			wantErr:  false,
		},
		Case{
			Path:     "/read/SF/1234567",
			Method:   http.MethodGet,
			Status:   http.StatusOK,
			Response: `{"message":"OK","format":"{\"hello\": \"world\"}\n"}` + "\n",
			wantErr:  false,
		},
		Case{
			Path:     "/delete/SF/1234567",
			Method:   http.MethodDelete,
			Status:   http.StatusOK,
			Response: `{"message":"OK"}` + "\n",
			wantErr:  false,
		},
	}

	for idx, tt := range tcases {
		var req *http.Request
		if tt.Method == http.MethodPost {
			payload, _ := json.Marshal(tt.Request)
			req = httptest.NewRequest(tt.Method, tt.Path, bytes.NewBuffer(payload))

		} else {
			req = httptest.NewRequest(tt.Method, tt.Path, &bytes.Buffer{})
		}

		rr := httptest.NewRecorder()

		apiserver.Router.ServeHTTP(rr, req)

		if rr.Code != tt.Status {
			t.Fatalf("%d] want %d status code, but recieved %d", idx+1, tt.Status, rr.Code)
		}

		resp, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			t.Fatalf("%d] unexpected error: %v", idx+1, err)
		}

		if tt.wantErr {
			if !strings.Contains(string(resp), tt.Error) {
				t.Fatalf("%d] want error containing string: %s, but recieved: %s", idx+1, tt.Error, resp)
			}
			continue
		}

		if string(resp) != tt.Response {
			t.Fatalf("%d] want response: %s, but recieved: %s", idx+1, tt.Response, resp)
		}
	}
}

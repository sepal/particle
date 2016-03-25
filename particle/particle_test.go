package particle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil, "foo")
	url, _ := url.Parse(server.URL)
	client.BaseUrl = url
}

func teardown() {
	server.Close()
}

func TestNewRequest(t *testing.T) {
	token := "someRandomToken"
	c := NewClient(nil, token)

	type foo struct {
		A string
	}

	inUrl, outUrl := "/foo", apiBaseUrl+"/foo"
	inBody, outBody := &foo{"foo"}, `{"A":"foo"}`+"\n"

	req, _ := c.NewRequest("GET", inUrl, inBody)

	// Test if the inUrl expanded to the absolute url.
	if req.URL.String() != outUrl {
		t.Errorf("NewRequest with %v has URL = %v, expected %v", inUrl, req.URL, outUrl)
	}

	// Test if body was correctly JSON encoded.
	body, _ := ioutil.ReadAll(req.Body)

	if string(body) != outBody {
		t.Errorf("NewRequest has body %v, expected %v", string(body), outBody)
	}

	reqToken := req.Header.Get("Authorization")
	if reqToken != "Bearer: "+token {
		t.Errorf("NewRequest had wrong token %v, should be %v", reqToken, token)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A": "a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)

	_, err := client.Do(req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}

}

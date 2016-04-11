package particle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func TestClient_NewRequest(t *testing.T) {
	token := "someRandomTokne"
	c := NewClient(nil, token)

	type foo struct {
		A string
	}

	inURL, outURL := "/foo", apiBaseURL+"/foo"
	form := url.Values{}
	form.Add("greeting", "hello")

	req, err := c.newRequest("POST", inURL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", mediaTypeForm)

	if err != nil {
		t.Fatalf("c.NewRequest throwed an error: %v", err)
	}

	if req.URL.String() != outURL {
		t.Errorf("NewRequest with %v has URL = %v instead of %v", inURL, req.URL.String(), outURL)
	}

	err = req.ParseForm()

	if err != nil {
		t.Fatalf("Could not parse the requests body due to: %v", err)
	}

	if req.FormValue("greeting") != form.Get("greeting") {
		t.Errorf("Request has the body = %v, expected %v", req.FormValue("greeting"), form.Get("greeting"))
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request was a %v, exptected %v", m, r.Method)
		}

		fmt.Fprintf(w, `{"A": "a"}`)
	})

	req, _ := client.newRequest("GET", "/", nil)

	body := new(foo)
	_, err := client.do(req, &body)

	if err != nil {
		t.Fatalf("Error while executing request due to: %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Received wrong answer %v, expected %v", body, expected)
	}
}

func TestClient_Get(t *testing.T) {
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

	body := new(foo)
	_, err := client.get("/", &body)
	if err != nil {
		t.Fatalf("client.Get(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestClient_Post(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Answer string
	}

	form := url.Values{}
	form.Add("greeting", "hello")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "POST"; m != r.Method {
			t.Errorf("Request methood = %v, expected %v,", r.Method, m)
		}

		err := r.ParseForm()

		if err != nil {
			t.Fatalf("ParseForm error: %v", err)
		}

		if a := r.PostFormValue("greeting"); a != form.Get("greeting") {
			t.Errorf("Form value a = '%v', expected '%v'", a, form.Get("greeting"))
		}

		fmt.Fprint(w, `{"Answer": "world"}`)
	})

	body := new(foo)

	_, err := client.post("/", form, &body)

	if err != nil {
		t.Fatalf("client.Post(): %v", err)
	}

	expected := &foo{"world"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

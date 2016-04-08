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
	client.BaseURL = url
}

func teardown() {
	server.Close()
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
	_, err := client.Get("/", &body)
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

	_, err := client.Post("/", form, &body)

	if err != nil {
		t.Fatalf("client.Post(): %v", err)
	}

	expected := &foo{"world"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestNewRequest(t *testing.T) {
	token := "someRandomToken"
	c := NewClient(nil, token)

	type foo struct {
		A string
	}

	inURL, outURL := "/foo", apiBaseURL+"/foo"
	inBody, outBody := &foo{"foo"}, `{"A":"foo"}`+"\n"

	req, _ := c.NewJSONRequest("GET", inURL, inBody)

	// Test if the inUrl expanded to the absolute url.
	if req.URL.String() != outURL {
		t.Errorf("NewRequest with %v has URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// Test if body was correctly JSON encoded.
	body, _ := ioutil.ReadAll(req.Body)

	if string(body) != outBody {
		t.Errorf("NewRequest has body %v, expected %v", string(body), outBody)
	}

	reqToken := req.Header.Get("Authorization")
	if reqToken != "Bearer "+token {
		t.Errorf("NewRequest had wrong authorization header '%v', should be '%v'", reqToken, "Bearer "+token)
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

	req, _ := client.NewJSONRequest("GET", "/", nil)
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

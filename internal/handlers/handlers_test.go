package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ArmanurRahman/booking/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name   string
	url    string
	method string

	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	/*{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-07-18"},
		{key: "end", value: "2021-07-20"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-07-18"},
		{key: "end", value: "2021-07-20"},
	}, http.StatusOK},
	{"make-reservation-post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Mubeen"},
		{key: "last_name", value: "Arman"},
		{key: "email", value: "mubeenarman19@gmail.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code. Got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case reservation is not set in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non -exist room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "first_name=mubeen"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=arman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=arman@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=01012532")

	//anoder way to create form rewuest data
	//postedData := url.Values{}
	//postedData.Add("first_name", "mubeen")
	//postedData.Add("last_name", "arman")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	reservation := models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomID:    1,
	}

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler return wrong response code. Got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test from missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation = models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomID:    1,
	}

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code for missing post body. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test session data
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code for session data. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for insetion reservation failure
	reqBody = "first_name=mubeen"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=arman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=arman@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=01012532")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation = models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomID:    2,
	}

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code for insetion reservation. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for insetion restriction failure
	reqBody = "first_name=mubeen"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=arman")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=arman@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=01012532")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation = models.Reservation{
		StartDate: time.Now(),
		EndDate:   time.Now(),
		RoomID:    1000,
	}

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code for insetion restriction. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)

	}
	return ctx
}

package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ArmanurRahman/booking/internal/config"
	"github.com/ArmanurRahman/booking/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//what am i put in session
	gob.Register(models.Reservation{})
	//change this value to true in production
	testApp.IsProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp
	os.Exit(m.Run())

}

type myWtiter struct{}

func (tw *myWtiter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWtiter) WriteHeader(i int) {

}

func (tw *myWtiter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

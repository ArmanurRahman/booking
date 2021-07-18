package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ArmanurRahman/booking/internal/config"
	"github.com/ArmanurRahman/booking/internal/handlers"
	"github.com/ArmanurRahman/booking/internal/models"
	"github.com/ArmanurRahman/booking/internal/render"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting listining to port ", port)
	//_ = http.ListenAndServe(port, nil)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//what am i put in session
	gob.Register(models.Reservation{})
	//change this value to true in production
	app.IsProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc

	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)
	return nil
}

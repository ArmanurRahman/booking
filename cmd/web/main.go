package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ArmanurRahman/booking/internal/config"
	"github.com/ArmanurRahman/booking/internal/drivers"
	"github.com/ArmanurRahman/booking/internal/handlers"
	"github.com/ArmanurRahman/booking/internal/helpers"
	"github.com/ArmanurRahman/booking/internal/models"
	"github.com/ArmanurRahman/booking/internal/render"

	"github.com/alexedwards/scs/v2"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	fmt.Println("stating mail listener")

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

func run() (*drivers.DB, error) {
	//what am i put in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	//initiate mail chan
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change this value to true in production
	app.IsProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	//connect to database

	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=booking user=postgres password=mubeen")
	if err != nil {
		log.Fatal("Cannot connect to DB. Dying...")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	app.UseCache = false
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}

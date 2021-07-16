package handlers

import (
	"net/http"

	"github.com/ArmanurRahman/booking/pkg/config"
	"github.com/ArmanurRahman/booking/pkg/models"
	"github.com/ArmanurRahman/booking/pkg/render"
)

//Repo the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIp := r.RemoteAddr

	//save ip in session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	strngMap := make(map[string]string)

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	strngMap["remote_ip"] = remoteIp
	strngMap["test"] = "Hello world"
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{StringMap: strngMap})

}

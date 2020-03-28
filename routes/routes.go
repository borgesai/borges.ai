package routes

import (
	"borges.ai/data"
	"borges.ai/models"
	"borges.ai/text"
	"borges.ai/utils"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
)

var r *render.Render = render.New(render.Options{
	IsDevelopment: os.Getenv("ENV") != "prod",
	Directory:     os.Getenv("TEMPLATES_DIR"),
	Layout:        "layout",
	Extensions:    []string{".html"},
	Funcs: []template.FuncMap{template.FuncMap{
		"loop": func(n int) []struct{} {
			return make([]struct{}, n)
		},
		"slug": func(textStr string) string {
			return text.Slug(textStr)
		},
		"formatFloat32": func(f float32) string { return fmt.Sprintf("%.1f", f) },
	},
	},
})

func GetUser(req *http.Request) models.User {
	defer utils.Duration(utils.Track("GetUser"))
	session, err := Store.Get(req, sessionName)
	if err != nil {
		log.WithError(err).Error("cannot get session")
	}
	var user = models.User{}
	userID := session.Values[sessionUserIDKey]
	if userID == nil {
		return user
	}
	db, err := models.NewDB()
	if err != nil {
		return user
	}
	defer db.Close()
	user, err = data.FindUserByID(db, userID.(uint64))
	if err != nil {
		log.WithField("user_id", userID).WithError(err).Error("cannot get user")
	}
	return user
}

func PingHandler(w http.ResponseWriter, req *http.Request) {
	r.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

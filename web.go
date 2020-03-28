package main

import (
	goodreadsOAuth "borges.ai/goodreads/oauth"
	"borges.ai/models"
	"borges.ai/routes"
	"borges.ai/twitter"
	"borges.ai/workers"
	"github.com/bufferapp/sqs-worker-go/worker"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/oauth2"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Book{})
	db.AutoMigrate(&models.Author{})
	db.AutoMigrate(&models.Review{})
	db.AutoMigrate(&models.Reading{})
	db.AutoMigrate(&models.BookEdition{})
	db.AutoMigrate(&models.AuditEvent{})
	db.AutoMigrate(&models.Job{})

	db.Exec("CREATE TRIGGER books_tsvectorupdate BEFORE INSERT OR UPDATE ON books FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(searchable_text, 'pg_catalog.english', title, subtitle, authors_cache_str)")
	// to drop DROP TRIGGER books_tsvectorupdate ON books;

	db.Exec("CREATE UNIQUE INDEX author_goodreads_unique_idx ON authors(goodreads_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX author_slug_unique_idx ON authors(slug) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX author_goodreads_unique_idx ON authors(goodreads_id) WHERE deleted_at IS NULL IS NULL AND goodreads_id <> ''")

	db.Exec("CREATE INDEX user_review_status_idx ON reviews(creator_id,status) WHERE deleted_at IS NULL")
	db.Exec("CREATE INDEX user_review_idx ON reviews(creator_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX user_book_review_unique_idx ON reviews(creator_id,book_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX user_review_goodreads_unique_idx ON reviews(goodreads_id) WHERE deleted_at IS NULL  IS NULL AND goodreads_id <> ''")


	db.Exec("CREATE UNIQUE INDEX book_goodreads_unique_idx ON books(goodreads_id) WHERE deleted_at IS NULL  IS NULL AND goodreads_id <> ''")
	db.Exec("CREATE UNIQUE INDEX book_best_edition_unique_idx ON books(best_edition_goodreads_id) WHERE deleted_at IS NULL  IS NULL AND best_edition_goodreads_id <> ''")
	db.Exec("CREATE UNIQUE INDEX book_edition_goodreads_unique_idx ON book_editions(goodreads_id) WHERE deleted_at IS NULL  IS NULL AND goodreads_id <> ''")
	db.Exec("CREATE UNIQUE INDEX book_edition_goodreads_book_unique_idx ON book_editions(goodreads_book_id) WHERE deleted_at IS NULL  IS NULL AND goodreads_book_id <> ''")


	db.Exec("CREATE UNIQUE INDEX user_twitter_id_unique_idx ON users(twitter_id) WHERE deleted_at IS NULL IS NULL AND twitter_id <> 0")
	db.Exec("CREATE UNIQUE INDEX user_goodreads_id_unique_idx ON users(goodreads_id) WHERE deleted_at IS NULL IS NULL AND goodreads_id <> 0")
	db.Exec("CREATE UNIQUE INDEX user_dropbox_id_unique_idx ON users(dropbox_id) WHERE deleted_at IS NULL IS NULL AND dropbox_id <> 0")
	db.Exec("CREATE UNIQUE INDEX user_username_unique_idx ON users(username) WHERE deleted_at IS NULL IS NULL AND username <> ''")
	db.Exec("CREATE UNIQUE INDEX user_email_unique_idx ON users(email) WHERE deleted_at IS NULL IS NULL AND email <> ''")

	db.Exec("CREATE INDEX user_audit_event_idx ON audit_events(creator_id,model_id,model_type,type) WHERE deleted_at IS NULL")

	db.Exec("CREATE INDEX user_readings_idx ON readings(creator_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE INDEX user_book_reading_idx ON readings(creator_id,book_id) WHERE deleted_at IS NULL")
	db.Exec("CREATE UNIQUE INDEX user_book_goodreads_reading_unique_idx ON readings(creator_id,book_id,goodreads_id) WHERE deleted_at IS NULL")

	log.Info("migrated")
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func main() {
	db, err := models.NewDB()
	if err != nil {
		log.WithError(err).Error("failed to connect")
		panic("failed to connect database")
	}
	defer db.Close()


	//migrate(db)

	log.Info("started")

	router := mux.NewRouter()

	// auth
	router.Handle("/twitter/callback", twitter.CallbackHandler(routes.GetTwitterOAuth1Config(), routes.IssueSession(), routes.HandleAuthFailure()))
	router.Handle("/twitter/login", twitter.LoginHandler(routes.GetTwitterOAuth1Config(), routes.HandleAuthFailure()))
	router.HandleFunc("/logout", routes.LogoutHandler)

	// goodreads
	router.Handle("/goodreads/callback", routes.GoodreadsCallbackHandler(goodreadsOAuth.GetGoodreadsOAuth1Config(), routes.HandleAuthFailure()))
	router.Handle("/goodreads/login", routes.GoodreadsLoginHandler(goodreadsOAuth.GetGoodreadsOAuth1Config(), routes.HandleAuthFailure()))

	// dropbox
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	router.Handle("/dropbox/callback", oauth2.StateHandler(stateConfig, oauth2.CallbackHandler(routes.GetDropboxOauth2Config(), routes.DropboxCallbackHandler(), routes.HandleAuthFailure())))
	router.Handle("/dropbox/login", oauth2.StateHandler(stateConfig, oauth2.LoginHandler(routes.GetDropboxOauth2Config(), routes.HandleAuthFailure())))

	// web routes
	router.HandleFunc("/", routes.IndexHandler).Methods("GET")
	router.HandleFunc("/changelog", routes.ChangelogHandler).Methods("GET")
	router.HandleFunc("/about", routes.AboutHandler).Methods("GET")
	router.HandleFunc("/ping", routes.PingHandler).Methods("GET")
	router.HandleFunc("/settings", routes.UserSettingsHandler).Methods("GET")
	router.HandleFunc("/@{username}", routes.UserHandler).Methods("GET")
	router.HandleFunc("/@{username}/to-read", routes.UserWantToReadHandler).Methods("GET")
	router.HandleFunc("/@{username}/reading", routes.UserReadingHandler).Methods("GET")
	router.HandleFunc("/@{username}/best", routes.UserBestHandler).Methods("GET")
	router.HandleFunc("/@{username}/charts", routes.UserChartsHandler).Methods("GET")


	router.HandleFunc("/to-read", routes.UserWantToReadHandler).Methods("GET")
	router.HandleFunc("/reading", routes.UserReadingHandler).Methods("GET")
	router.HandleFunc("/best", routes.UserBestHandler).Methods("GET")
	router.HandleFunc("/charts", routes.UserChartsHandler).Methods("GET")


	router.HandleFunc("/@{username}/b/{book-name-slug}-{book-short-id}", routes.UserBookHandler).Methods("GET")
	// this are routes that modify user resources. we probably can create middleware that verifies permissions here
	router.HandleFunc("/@{username}/s/goodreads/sync", routes.UserSyncGoodreadsHandler).Methods("POST")
	router.HandleFunc("/@{username}/s/goodreads/sync-setting", routes.UserSyncGoodreadsSettingHandler).Methods("POST")
	router.HandleFunc("/@{username}/b/{book-name-slug}-{book-short-id}/status", routes.UserChangeBookStatusHandler).Methods("POST")
	router.HandleFunc("/@{username}/b/{book-name-slug}-{book-short-id}/rating", routes.UserChangeBookRatingHandler).Methods("POST")
	router.HandleFunc("/api/@{username}/b/{book-name-slug}-{book-short-id}/review", routes.UserUpdateReviewAPIHandler).Methods("PATCH")
	router.HandleFunc("/api/@{username}/b/{book-name-slug}-{book-short-id}/readings/{reading-id}/note", routes.UserReadingNoteAPIHandler).Methods("PATCH")
	router.HandleFunc("/api/@{username}/b/{book-name-slug}-{book-short-id}/readings/{reading-id}/start-date", routes.UserReadingStartDateAPIHandler).Methods("PATCH")
	router.HandleFunc("/api/@{username}/b/{book-name-slug}-{book-short-id}/readings/{reading-id}/finish-date", routes.UserReadingFinishDateAPIHandler).Methods("PATCH")

	router.HandleFunc("/a/{author-name-slug}", routes.AuthorHandler).Methods("GET")
	router.HandleFunc("/b/{book-name-slug}-{book-short-id}", routes.BookHandler).Methods("GET")
	router.HandleFunc("/search", routes.SearchHandler).Methods("GET")

	router.NotFoundHandler = routes.Borges404Handler()

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("STATICS_DIR"))))



	worker.MaxNumberOfMessage = 5
	wSyncGoodreadsUser, err := worker.NewService("borges-sync-goodreads-user-queue.fifo")
	if err != nil {
		log.WithError(err).Error("error creating worker service")
	}
	go wSyncGoodreadsUser.Start(worker.HandlerFunc(workers.ProcessSyncGoodreadsUser))

	wSyncGoodreadsBook, err := worker.NewService("borges-sync-goodreads-book-queue.fifo")
	if err != nil {
		log.WithError(err).Error("error creating worker service")
	}
	go wSyncGoodreadsBook.Start(worker.HandlerFunc(workers.ProcessSyncGoodreadsBook))

	wSyncUserReviewToGoodreads, err := worker.NewService("borges-sync-user-review-to-goodreads-queue.fifo")
	if err != nil {
		log.WithError(err).Error("error creating worker service")
	}
	go wSyncUserReviewToGoodreads.Start(worker.HandlerFunc(workers.ProcessSyncUserReviewToGoodreads))

	wSyncUserReviewToDropbox, err := worker.NewService("borges-sync-user-review-to-dropbox-queue.fifo")
	if err != nil {
		log.WithError(err).Error("error creating worker service")
	}
	go wSyncUserReviewToDropbox.Start(worker.HandlerFunc(workers.ProcessSyncUserReviewToDropbox))

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:" + os.Getenv("PORT"),
		// Enforce timeouts
		IdleTimeout: 120 * time.Second,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  60 * time.Second,
		ReadHeaderTimeout:  60 * time.Second,
	}

	log.WithField("port", os.Getenv("PORT")).Info("started web server")
	log.Fatal(srv.ListenAndServe())
}

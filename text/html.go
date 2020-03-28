package text

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/repetitive/html-to-markdown"
	"github.com/russross/blackfriday"
	log "github.com/sirupsen/logrus"
)

func CreateContentHTML(content string) string {
	contentBytes := []byte(content)
	unsafe := blackfriday.MarkdownCommon(contentBytes)
	html := string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	log.WithField("content_html", html).Info("html")
	return html
}

func ConvertHTMLToMarkdown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.WithError(err).WithField("html", html).Error("failed to convert html to markdown")
	}
	return markdown
}

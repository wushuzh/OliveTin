package httpservers

import (
	"google.golang.org/grpc/metadata"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	localUserSessions = make(map[string]string) // sid -> username, used for local user sessions
)

func parseLocalUserCookie(req *http.Request) (string, string, string) {
	cookie, err := req.Cookie("olivetin-sid-local")

	if err != nil {
		return "", "", ""
	}

	cookieValue := cookie.Value

	username, ok := localUserSessions[cookieValue]

	if !ok {
		log.WithFields(log.Fields{
			"sid":      cookieValue,
			"provider": "local",
		}).Warnf("Stale session")
		return "", "", ""
	}

	return username, "", cookie.Value
}

func forwardResponseHandlerLoginLocalUser(md metadata.MD, w http.ResponseWriter) error {
	setUser := getMetadataKeyOrEmpty(md, "set-user")

	if setUser != "" {
		sid := uuid.NewString()
		localUserSessions[sid] = setUser

		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "olivetin-sid-local",
				Value:    sid,
				MaxAge:   31556952, // 1 year
				HttpOnly: true,
				Path:     "/",
			},
		)
	}

	return nil
}

package livesession

import (
	"fmt"
	"os"

	"github.com/darkphotonKN/collabradoc/internal/customerrors"
	"github.com/darkphotonKN/collabradoc/internal/document"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/user"
	"github.com/google/uuid"
)

func CreateLiveSessionService(userId uint, req CreateLiveSessionReq) (model.LiveSession, error) {

	// validate user exists
	user, err := user.FindUserById(userId)

	if err != nil {
		return model.LiveSession{}, err
	}

	// validate document exists and belongs to specific user
	doc, err := document.GetDocumentById(req.DocumentID, userId)
	fmt.Println("err finding doc", err)

	if err != nil {
		return model.LiveSession{}, err
	}

	sessionId := GenerateSessionID()

	return CreateLiveSession(user, sessionId, doc)

}

func GenerateSessionID() string {
	return uuid.NewString()
}

func GetLiveSessionService(userId uint, documentId uint) (LiveSessionLink, error) {

	liveSession, err := QueryLiveSession(userId, documentId)

	if err != nil {
		return "", customerrors.LiveSessionUnauthorized
	}

	domain := os.Getenv("SITE_DOMAIN")

	// construct link with liveSession's sessionId which only allows authenticated users
	// who own both the doc and the session to access
	return LiveSessionLink(fmt.Sprintf("%s/%s", domain, liveSession.SessionID)), nil
}

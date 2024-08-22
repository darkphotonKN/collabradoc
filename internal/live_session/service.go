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

func CreateLiveSessionService(userId uint, documentId uint) (model.LiveSession, error) {
	// validate user exists
	user, err := user.FindUserById(userId)

	if err != nil {
		return model.LiveSession{}, err
	}

	// validate document exists and belongs to specific user
	doc, err := document.GetDocumentById(documentId, userId)
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
	// validates live session belongs to the user, and retreives it
	liveSession, err := QueryLiveSession(userId, documentId)

	if err != nil {
		return "", customerrors.LiveSessionUnauthorized
	}

	return GenerateLiveSessionURL(liveSession.SessionID, documentId), nil
}

func AuthorizeLiveSessionService(userId uint, sessionId string) (bool, error) {
	// get relation of user with the sessionId
	_, err := user.FindUserById(userId)

	if err != nil {
		return false, fmt.Errorf("User does not exist.")
	}

	err = QueryLiveSessionForUser(userId, sessionId)

	if err != nil {
		return false, fmt.Errorf("This user is not authorized to access this session.")
	}

	return true, nil
}

func InviteToliveSessionService(userId uint, email string, documentId uint) (model.LiveSession, error) {
	sendingUser, err := user.FindUserById(userId)
	if err != nil {
		return model.LiveSession{}, fmt.Errorf("Error when retrieving user: %s", err)
	}
	fmt.Println(sendingUser)

	targetUser, err := user.FindUserByEmail(email)
	if err != nil {
		return model.LiveSession{}, fmt.Errorf("Error when retrieving target user: %s", err)
	}

	fmt.Println("targetUser:", targetUser)

	return model.LiveSession{}, nil
}

/**
* Constructs Live Session URL
**/
func GenerateLiveSessionURL(sessionId string, documentId uint) LiveSessionLink {

	domain := os.Getenv("SITE_DOMAIN")

	// construct link with liveSession's sessionId which only allows authenticated users
	// who own both the doc and the session to access
	return LiveSessionLink(fmt.Sprintf("%s/docs-live?sessionId=%s&documentId=%d", domain, sessionId, documentId))
}

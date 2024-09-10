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

	fmt.Printf("attempting to query with live session userId %d and documentId %d", user.ID, documentId)

	// check if live session already exists, including one that user has been invited to
	existingLiveSession, err := QueryLiveSession(user.ID, documentId)

	fmt.Printf("existingLiveSession: %+v", existingLiveSession)

	if err != nil {
		return model.LiveSession{}, err
	}

	// check existing sessionID exists
	if existingLiveSession.SessionID != "" {

		// live session exist, just return that one
		fmt.Println("found existing live session")
		return existingLiveSession, nil
	}

	// live session does not exist, generate unique session id and create as the owner
	sessionId := GenerateSessionID()
	// confirm user is document owner and retrieve the document
	doc, err := document.GetDocumentById(documentId, user.ID)

	if err != nil {
		fmt.Println("err finding doc", err)
		return model.LiveSession{}, err
	}

	fmt.Println("Live session does not exist, creating...")
	return CreateLiveSession(user, sessionId, doc)
}

func GenerateSessionID() string {
	return uuid.NewString()
}

func GetLiveSessionService(userId uint, documentId uint) (model.LiveSession, error) {
	fmt.Printf("attempting to query with live session userId %d and documentId %d", userId, documentId)

	// validates live session belongs to the user, and retreives it
	liveSession, err := QueryLiveSession(userId, documentId)

	if err != nil {
		return liveSession, customerrors.LiveSessionUnauthorized
	}

	return liveSession, nil
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

func InviteToLiveSessionService(userId uint, email string, sessionId string) (model.LiveSession, error) {
	// validate the user sending the invite and user being invited both exist
	sendingUser, err := user.FindUserById(userId)
	if err != nil {
		return model.LiveSession{}, fmt.Errorf("Error when retrieving user: %s", err)
	}

	targetUser, err := user.FindUserByEmail(email)
	if err != nil {
		return model.LiveSession{}, fmt.Errorf("Error when retrieving target user: %s", err)
	}

	// add target to existing live session
	liveSession, err := InsertUserToLiveSession(sendingUser, targetUser, sessionId)

	if err != nil {
		return model.LiveSession{}, err
	}

	return liveSession, nil
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

func GetInvitedLiveSessionsService(userId uint) ([]LiveSessionInvites, error) {
	_, err := user.FindUserById(userId)

	if err != nil {
		return nil, err
	}

	return QueryAllNonOwnedLiveSessions(userId)
}

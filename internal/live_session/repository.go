package livesession

import (
	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/db"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateLiveSession(user model.User, sessionId string, doc model.Document) (model.LiveSession, error) {
	db := db.DBCon

	newLiveSession := model.LiveSession{
		SessionID:  sessionId,
		DocumentID: doc.ID,
		IsActive:   true,
		Users:      []model.User{user}, // add creator to live session
	}

	if err := db.Create(&newLiveSession).Error; err != nil {
		return model.LiveSession{}, err
	}

	return newLiveSession, nil
}

func QueryLiveSession(userId uint, documentId uint) (model.LiveSession, error) {
	db := db.DBCon

	var liveSession model.LiveSession

	result := db.Where("document_id =?", documentId).First(&liveSession)

	if result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		return model.LiveSession{}, fmt.Errorf("Document does not belong to user you are attempting to create a live session with.")
	}

	return liveSession, nil
}

func QueryLiveSessionForUser(userId uint, sessionId string) error {
	db := db.DBCon
	db.Debug()

	var liveSession model.LiveSession

	fmt.Printf("sessionId: %s\n", sessionId)

	result := db.Joins("JOIN live_session_users ON live_session_users.live_session_id = live_sessions.id").Where("live_sessions.session_id = ? AND live_session_users.user_id = ?", sessionId, userId).First(&liveSession)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

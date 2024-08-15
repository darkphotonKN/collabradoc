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

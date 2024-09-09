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

/**
* Queries for live sessions that the user is part of.
**/
func QueryLiveSession(userId uint, documentId uint) (model.LiveSession, error) {
	db := db.DBCon

	// join users and live sessions
	var existingLiveSession ExistingLiveSession

	r := db.Model(model.LiveSession{}).Select("live_sessions.session_id, live_sessions.document_id").Joins("JOIN live_session_users ON live_session_users.live_session_id = live_sessions.id").Joins("JOIN users ON users.id = live_session_users.user_id").Where("users.id = ? AND live_sessions.document_id = ?", userId, documentId).Scan(&existingLiveSession)

	if r.Error != nil {
		fmt.Println("Problem querying all live session and user relations:", r.Error)
	}

	fmt.Printf("\nFound relations %+v\n\n", existingLiveSession)

	return model.LiveSession{
		SessionID:  existingLiveSession.SessionID,
		DocumentID: existingLiveSession.DocumentID,
	}, nil
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

func InsertUserToLiveSession(user model.User, newUser model.User, sessionId string) (model.LiveSession, error) {
	db := db.DBCon

	// get old record
	var existingLiveSession model.LiveSession

	result := db.Preload("Users").Where("session_id =?", sessionId).First(&existingLiveSession)

	if result.Error != nil {
		return model.LiveSession{}, result.Error
	}

	// insert new user into live session
	existingLiveSession.Users = append(existingLiveSession.Users, newUser)

	saveResult := db.Save(&existingLiveSession)
	if saveResult.Error != nil {
		return model.LiveSession{}, saveResult.Error
	}

	// fmt.Printf("Added user in live session %+v\n", existingLiveSession)
	return existingLiveSession, nil
}

func QueryAllNonOwnedLiveSessions(userId uint) ([]LiveSessionInvites, error) {
	db := db.DBCon

	var liveSessionsRes []LiveSessionInvites

	// queries for live sessions that the user is part of but of which is NOT the owner of
	// the original document, effectively making this query find all live sessions in which
	// the user was *invited*.
	result := db.Model(&model.LiveSession{}).
		Select("live_sessions.id, live_sessions.created_at, documents.title, live_sessions.document_id, live_sessions.session_id, live_sessions.is_active").
		Joins("JOIN live_session_users ON live_session_users.live_session_id = live_sessions.id").
		Joins("JOIN documents ON documents.id = live_sessions.document_id").
		Where("live_session_users.user_id = ? AND documents.user_id <> ?", userId, userId).
		Scan(&liveSessionsRes)

	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Println("querying all live sessions:", liveSessionsRes)

	return liveSessionsRes, nil
}

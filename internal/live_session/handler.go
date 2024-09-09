package livesession

import (
	"encoding/json"

	"fmt"
	"net/http"
	"strconv"

	"github.com/darkphotonKN/collabradoc/internal/document"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
	"github.com/darkphotonKN/collabradoc/internal/utils/request"
	"github.com/go-playground/validator/v10"
)

/**
* Creats a live session based on user's id from JWT token and the documentId
* of the document they are creating a live session for.
**/
func CreateLiveSessionHandler(w http.ResponseWriter, r *http.Request) {

	userId, _ := request.ExtractUserID(r.Context())

	var req CreateLiveSessionReq
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validation
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create live session for the specific document and user
	docIdInt64, err := req.DocumentID.Int64()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	liveSession, err := CreateLiveSessionService(userId, uint(docIdInt64))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	liveSessionLink := GenerateLiveSessionURL(liveSession.SessionID, liveSession.DocumentID)

	newLiveSessionRes := model.Response[LiveSessionLink]{
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Successfully created new live session %s for user %d", liveSession.SessionID, userId),
		Data:    liveSessionLink,
	}

	out, err := json.Marshal(newLiveSessionRes)

	if err != nil {
		fmt.Printf("Error when encoding created live session response: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

/**
* Invites a user to join a live session based on invitee's user id, the target user's
* user's email, and the documentId in which they are sending the invite from.
**/
func InviteToLiveSessionHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	sessionId := r.URL.Query().Get("sessionId")
	email := r.URL.Query().Get("email")

	inviteLiveSessionRes, err := InviteToLiveSessionService(userId, email, sessionId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("user invited to live session, %+v\n\n", inviteLiveSessionRes)

	liveSessionRes := model.Response[model.LiveSession]{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Successfully added %s to live session.", email),
		Data:    inviteLiveSessionRes,
	}

	out, err := json.Marshal(liveSessionRes)

	if err != nil {
		fmt.Printf("Error when encoding created live session response: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

/**
* Gets an exiting live session based on userId and documentId.
**/
func GetLiveSessionHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	documentIdQuery := r.URL.Query().Get("documentId")

	// convert to uint to conform to the actual documentId column
	documentIdUint64, err := strconv.ParseUint(documentIdQuery, 10, 32)
	documentId := uint(documentIdUint64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	liveSession, err := GetLiveSessionService(userId, documentId)

	liveSessionRes := model.Response[model.LiveSession]{
		Status:  http.StatusOK,
		Message: "Current live session.",
		Data:    liveSession,
	}

	out, err := json.Marshal(liveSessionRes)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

/**
* Authorizes a user's request to access a live session via userId and sessionId.
**/
func AuthorizeLiveSessionHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	sessionId := r.URL.Query().Get("sessionId")

	fmt.Printf("sessionId: %s", sessionId)

	sessionAuthorized, err := AuthorizeLiveSessionService(userId, sessionId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	// check final session authorized boolean to dictate response package
	if sessionAuthorized {
		successRes := model.Response[bool]{
			Status:  http.StatusOK,
			Message: "Live session was authorized.",
			Data:    true,
		}

		out, err := json.Marshal(successRes)

		if err != nil {
			fmt.Printf("Error when encoding authorize live session handler response: %s\n", err)
		}
		w.Write(out)
	} else {

		rejectRes := model.Response[bool]{
			Status:  http.StatusUnauthorized,
			Message: "Live session is not authorized for this user / session.",
			Data:    false,
		}

		out, err := json.Marshal(rejectRes)

		if err != nil {
			fmt.Printf("Error when encoding authorize live session handler response: %s\n", err)
		}
		w.Write(out)
	}
}

/**
* TODO: Update to emailing the user the live session or adding it to their personal docs.
* Retrieves a user's accessible live sessions which docs don't directly belong to them.
*
* NOTE: (Temp) This is the temporary version of invited live sessions.
**/
func GetInvitedLiveSessionsHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := request.ExtractUserID(r.Context())

	liveSessionInvites, err := GetInvitedLiveSessionsService(userId)

	fmt.Printf("liveSessionInvites: %+v", liveSessionInvites)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var liveSessionInvitesRes []LiveSessionInvitesRes

	for _, liveSessionInvite := range liveSessionInvites {
		liveSessionInvitesRes = append(liveSessionInvitesRes, LiveSessionInvitesRes{ID: liveSessionInvite.DocumentID,
			CreatedAt:       liveSessionInvite.CreatedAt,
			Title:           liveSessionInvite.Title,
			IsActive:        liveSessionInvite.IsActive,
			LiveSessionInfo: document.LiveSessionInfo{SessionID: liveSessionInvite.SessionID},
		})

	}

	invitedLiveSessionsRes := model.Response[[]LiveSessionInvitesRes]{
		Status:  http.StatusOK,
		Message: "Successfully retrieved all currently invited live sessions.",
		Data:    liveSessionInvitesRes,
	}

	out, err := json.Marshal(invitedLiveSessionsRes)

	if err != nil {
		fmt.Printf("Error when encoding authorize live session handler response: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

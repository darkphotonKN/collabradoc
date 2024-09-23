package document

import (
	// "fmt"

	"fmt"

	"github.com/darkphotonKN/collabradoc/internal/rating"
	model "github.com/darkphotonKN/collabradoc/internal/shared"
)

func CreateDocumentService(doc CreateDocumentReq, userId uint) (model.Document, error) {
	return CreateDocument(doc, userId)
}

func GetDocuments(userId uint) ([]model.Document, error) {
	return QueryDocuments(userId)
}

/**
* Gets all community (public) documents.
**/
func GetCommunityDocsService() ([]DocumentRes, error) {
	docs, err := QueryPublicDocuments()

	if err != nil {
		return []DocumentRes{}, err
	}

	docsRes := make([]DocumentRes, len(docs))

	// get ratings and count average
	for index, doc := range docs {
		avgRating, err := rating.CountRatingsAvg(doc.ID)

		fmt.Printf("avgRating %f\n", avgRating)

		if err != nil {
			return []DocumentRes{}, err
		}

		docsRes[index] = DocumentRes{
			ID:        doc.ID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			Title:     doc.Title,
			Content:   doc.Content,
			UserId:    doc.UserId,
			LiveSessionInfo: LiveSessionInfo{
				SessionID: doc.LiveSession.SessionID,
			},
			Comment:       doc.Comment,
			Privacy:       doc.Privacy,
			AverageRating: avgRating,
		}
	}

	fmt.Printf("docRes %+v", docsRes)
	return docsRes, nil
}

func GetDocumentById(id uint, userId uint) (model.Document, error) {
	return QueryDocumentById(id, userId)
}

/**
* Toggles document's privacy, if it was private make it public and vice versa.
**/
func ToggleDocPrivacyService(userId uint, documentId uint) (model.Document, error) {
	// confirm document ownership and acquire document
	existingDoc, err := QueryDocumentById(documentId, userId)

	if err != nil {
		return existingDoc, err
	}

	// fmt.Printf("\n\nexistingDoc: \n%+v\n\n", existingDoc)

	// update doc
	if existingDoc.Privacy == private {
		existingDoc.Privacy = public
	} else if existingDoc.Privacy == public {
		existingDoc.Privacy = private
	}

	// update doc with new privacy setting
	updatedDoc, err := UpdateDocument(existingDoc)
	if err != nil {
		return updatedDoc, err
	}

	return updatedDoc, nil
}

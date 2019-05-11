package auth

import (
	"time"

	"github.com/dueckminor/mypi-api/go/rand"
)

type AuthRequest struct {
	Id          string
	TTL         time.Time
	RedirectURI string
	Path        string
}

var authRequests []*AuthRequest

func init() {
	authRequests = make([]*AuthRequest, 0, 0)
}

// NewRequest creates a new AuthRequest
func NewRequest() (request *AuthRequest, err error) {
	request = &AuthRequest{}
	request.Id, err = rand.GetString(32)
	if err != nil {
		return nil, err
	}

	request.TTL = time.Now().Add(time.Minute * 5)

	authRequests = append(authRequests, request)

	return request, nil
}

// GetRequest gets an existing AuthRequest by id
func GetRequest(id string) (request *AuthRequest) {
	for _, request = range authRequests {
		if request.Id == id {
			if request.TTL.After(time.Now()) {
				return request
			}
			return nil
		}
	}
	return nil
}

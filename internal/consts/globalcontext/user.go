package globalcontext

import (
	"fmt"
	"sync"

	"github.com/hantbk/vbackup/internal/model"
)

var (
	mu          sync.Mutex
	currentUser *model.Userinfo
)

// SetCurrentUser sets the current user information globally
func SetCurrentUser(user *model.Userinfo) {
	mu.Lock()
	defer mu.Unlock()
	currentUser = user
}

// GetCurrentUser retrieves the current user information from global context
func GetCurrentUser() (*model.Userinfo, error) {
	mu.Lock()
	defer mu.Unlock()

	if currentUser == nil {
		return nil, fmt.Errorf("no user information available")
	}
	return currentUser, nil
}

package store

import (
	"accounter/backend/core"
)

type mainStore struct {
	baseStore
	user *core.CurrentUser
}

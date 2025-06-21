package testutils

import (
	"testing"

	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

var client *ent.Client

// SetupTestDB returns an in-memory ent client
// remember to close the client after setting up db
func SetupTestDB(t *testing.T) *ent.Client {
	if client == nil {
		client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	}
	return client
}

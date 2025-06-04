// handler_test tests the core biz logic of the handler package
// handler package is in charge of doing biz logic inbetween the transport/server layer and the core IO logic
package handler_test

import (
	"fmt"
	"testing"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
	"github.com/jbarzegar/ron-mod-manager/testutils"
	"golang.org/x/net/context"
)

var testCfg = appconfig.AppConfig{
	GameDir: "/tmp",
	ModDir:  "/tmp",
}

func initTestHandler(t *testing.T) (handler.Handler, *ent.Client) {

	client := testutils.SetupTestDB(t)
	ioHandler := handlerio.NewMockIOHandler()
	h := handler.NewHandler(client, &testCfg, &ioHandler)

	return h, client
}

// TestShouldAddMod tests that a mod can be added via an archive
// This test specifically should return a mod with no choices
func TestShouldAddMod(t *testing.T) {
	h, client := initTestHandler(t)

	expectedName := "test-mod"
	r, err := h.AddMod("/path/to/archive.zip", expectedName)
	if err != nil {
		t.Fatal(err)
	}

	// check if mod name is the expected mod name
	if r.Archive.Name != expectedName {
		fmt.Println("Created mod does not match expected name")
		fmt.Println("Got: " + r.Archive.Name)
		fmt.Println("Expected: " + expectedName)
		t.Fail()
	}

	exists, err := client.Archive.
		Query().
		Where(archive.NameEQ(expectedName)).
		Exist(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		fmt.Println("could not find saved archive")
		t.Fail()
	}
}

// TestShouldModWithChoices tests that a mod can be installed after an archive is added
// Stipulation being that the added mod returns a list of choices that must be provided
// Choices are specific to optional files that are detected
// func TestShouldAddModWithChoices() {

// }

// TestShouldInstallMod tests that a mod can be installed once an archive is added
// func TestShouldInstallMod() {

// }

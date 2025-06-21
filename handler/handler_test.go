// handler_test tests the core biz logic of the handler package handler package
// is in charge of doing biz logic inbetween the transport/server layer and
// calling core IO logic (supplied via a IOHandler)
package handler_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	a "github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
	"github.com/jbarzegar/ron-mod-manager/testutils"
	"golang.org/x/net/context"
)

var testCfg = appconfig.AppConfig{
	GameDir: "/tmp",
	ModDir:  "/tmp",
}

var noChoices = []a.Choice{}

// var client = ent.

func initTestHandler(t *testing.T, choices []a.Choice) (handler.Handler, *ent.Client) {
	client := testutils.SetupTestDB(t)
	ioHandler := handlerio.NewMockIOHandler(choices)
	h := handler.NewHandler(client, &testCfg, &ioHandler)

	return h, client
}

// TestShouldAddMod tests that a mod can be added via an archive
// This test specifically should return a mod with no choices
func TestShouldAddMod(t *testing.T) {
	h, client := initTestHandler(t, noChoices)

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
// Choices are specific to optional files that are detected that can be added
// TODO: figure out how to manage conflicts
// currently choices has no concept of dependencies meaning someone can easily add multiple conflciting dependencies
// there's no easy way to determine if two choices conflict with one another atm
func TestShouldAddModWithChoices(t *testing.T) {
	choices := []a.Choice{
		{Name: "choice-a", FullPath: "/a/path/to/choice-a"},
		{Name: "choice-b", FullPath: "/a/path/to/choice-b"},
	}

	expectedName := "test-mod"
	h, _ := initTestHandler(t, choices)

	r, err := h.AddMod("/path/to/archive.zip", expectedName)
	if err != nil {
		t.Fatal(err)
	}

	if len(r.Choices) != len(choices) {
		fmt.Println("Did not get amount of expected choices")
		t.Fail()
	}
}

type multiMod struct {
	Name           string
	ExpectedPakLen int
	Path           string
	Choices        []a.Choice
}

// TestShouldAddMultipleMods tests that multiple mods can be added to the db
// the test will also verify that paks will be queried from the correct mods
func TestShouldAddMultipleMods(t *testing.T) {
	testMods := []multiMod{
		{
			Name:           "first-mod",
			ExpectedPakLen: 1,
			Path:           "/path/to/archive.zip",
			Choices:        noChoices,
		},
		{
			Name:           "second-mod",
			ExpectedPakLen: 3,
			Path:           "/apth/to/second/mod",
			Choices: []a.Choice{
				{Name: "choice-b", FullPath: "/apth/to/choice-b"},
				{Name: "choice-c", FullPath: "/apth/to/choice-c"},
			},
		},
	}

	// add each mod
	for i, m := range testMods {
		h, _ := initTestHandler(t, m.Choices)
		_, err := h.AddMod(m.Path, m.Name)
		if err != nil {
			fmt.Printf("failed to create mod on index: %v\n", i)
			t.Fatal(err)
		}
	}

	// handler is a noop since we just want the db
	_, db := initTestHandler(t, noChoices)
	for _, testMod := range testMods {
		// get the newly created mod using the names provided
		modQueryResult, err := db.
			Mod.
			Query().
			Where(mod.NameEQ(testMod.Name)).
			Only(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if modQueryResult == nil {
			t.Fatal(errors.New(fmt.Sprintf("mod %v not found", testMod.Name)))
		}

		// get all the mod Versions (though we only expect a single version for this test)
		mvList, err := modQueryResult.QueryVersions().All(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(mvList) > 1 {
			t.Fatal(errors.New("only 1 mod version expected"))
		}

		mv := mvList[0]
		paks, err := mv.QueryPaks().All(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		// ensure the expected amount of paks is correct for each mod/modversion
		if len(paks) != testMod.ExpectedPakLen {
			t.Fatalf("paks len was %v. Expected %v", len(paks), testMod.ExpectedPakLen)
		}
	}
}

// TestShouldInstallMod tests that a mod can be installed once an archive is added
// func TestShouldInstallMod(t *testing.T) {
// 	h, _ := initTestHandler(t, noChoices)

// 	h.InstallMod()

// }

// TestShouldInstallModWithChoices tests that a mod can be installed with choices
// func TestShouldInstallModWithChoices(t *testing.T) {

// }

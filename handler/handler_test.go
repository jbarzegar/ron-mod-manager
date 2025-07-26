// handler_test tests the core biz logic of the handler package handler package
// is in charge of doing biz logic inbetween the transport/server layer and
// calling core IO logic (supplied via a IOHandler)
package handler_test

import (
	"fmt"
	"testing"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	a "github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/ent/modversion"
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

func initTestHandler(t *testing.T, choices []a.Choice) (*handler.Handler, *handlerio.MockIOHandler, *ent.Client) {
	client := testutils.SetupTestDB(t)
	// ioHandler := handlerio.NewMockIOHandler(choices)
	ioHandler := &handlerio.MockIOHandler{
		MockedChoices: choices,
		Installed:     map[string]handlerio.MockMod{},
	}
	h := handler.Handler{Db: client, Config: &testCfg, Io: ioHandler}

	return &h, ioHandler, client
}

// TestShouldAddMod tests that a mod can be added via an archive
// This test specifically should return a mod with no choices
func TestShouldAddMod(t *testing.T) {
	h, _, client := initTestHandler(t, noChoices)

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
	h, _, _ := initTestHandler(t, choices)

	r, err := h.AddMod("/path/to/archive.zip", expectedName)
	if err != nil {
		t.Fatal(err)
	}

	if len(r.Choices) != len(choices)+1 {
		fmt.Println(len(r.Choices), len(choices))
		fmt.Println("Did not get amount of expected choices")
		t.Fail()
	}
}

type modToInstall struct {
	Name           string
	ExpectedPakLen int
	Path           string
	Choices        []a.Choice
}

// TestShouldAddMultipleMods tests that multiple mods can be added to the db
// the test will also verify that paks will be queried from the correct mods
// func TestShouldAddMultipleMods(t *testing.T) {
// 	testMods := []modToInstall{
// 		{
// 			Name:           "first-mod",
// 			ExpectedPakLen: 1,
// 			Path:           "/path/to/archive.zip",
// 			Choices:        noChoices,
// 		},
// 		{
// 			Name:           "second-mod",
// 			ExpectedPakLen: 3,
// 			Path:           "/apth/to/second/mod",
// 			Choices: []a.Choice{
// 				{Name: "choice-b", FullPath: "/apth/to/choice-b"},
// 				{Name: "choice-c", FullPath: "/apth/to/choice-c"},
// 			},
// 		},
// 	}

// 	var db *ent.Client = nil
// 	// add each mod
// 	for i, m := range testMods {
// 		h, _, d := initTestHandler(t, m.Choices)
// 		if db == nil {
// 			db = d
// 		}
// 		_, err := h.AddMod(m.Path, m.Name)
// 		if err != nil {
// 			fmt.Printf("failed to create mod on index: %v\n", i)
// 			t.Fatal(err)
// 		}
// 	}

// 	for _, testMod := range testMods {
// 		// get the newly created mod using the names provided
// 		modQueryResult, err := db.
// 			Mod.
// 			Query().
// 			Where(mod.NameEQ(testMod.Name)).
// 			Only(context.TODO())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if modQueryResult == nil {
// 			t.Fatal(fmt.Errorf("mod %v not found", testMod.Name))
// 		}

// 		// get all the mod Versions (though we only expect a single version for this test)
// 		mvList, err := modQueryResult.QueryVersions().All(context.TODO())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if len(mvList) > 1 {
// 			t.Fatal(errors.New("only 1 mod version expected"))
// 		}

// 		mv := mvList[0]
// 		paks, err := mv.QueryPaks().All(context.TODO())
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		// ensure the expected amount of paks is correct for each mod/modversion
// 		if len(paks) != testMod.ExpectedPakLen {
// 			t.Fatalf("paks len was %v. Expected %v", len(paks), testMod.ExpectedPakLen)
// 		}
// 	}
// }

// TestShouldInstallMod tests that a mod can be installed once an archive is added
func TestShouldInstallModWithNoChoices(t *testing.T) {
	// add a mod w no choices (for testing purposes)
	testMod := modToInstall{
		Name:           "first-mod",
		ExpectedPakLen: 1,
		Path:           "/path/to/archive.zip",
	}
	h, hio, db := initTestHandler(t, noChoices)

	// do an addMod
	addModResp, err := h.AddMod(testMod.Path, testMod.Name)
	if err != nil {
		t.Fatal(err)
	}

	// determine it worked
	// test if db state has used updated modversion id we expect
	mv, err := db.ModVersion.Query().Where(modversion.UUIDEQ(addModResp.ModVersion.UUID)).Only(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if mv.UUID != addModResp.ModVersion.UUID {
		t.Fatal("UUIDS do not match")
	}

	mod, err := mv.QueryMod().Only(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	paks, err := mod.QueryVersions().QueryPaks().All(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	// prep deps to install mod
	installable := handlerio.Installable{
		Mod:         mod,
		ArchivePath: "idk",
		Assets:      handlerio.InstallableAssets{Paks: paks},
		OutPath:     "idk",
	}
	// test if mock io has installed the mod
	if err := hio.InstallMod(installable); err != nil {
		t.Fatal(err)
	}
	if hio.Installed[installable.Mod.Name].Installed != true {
		t.Fatal("Mod was not installed")
	}
}

package plex

import (
	"os"
	"testing"

	"github.com/kjbreil/go-plex/pkg/library"
)

func getTestConnection(t *testing.T) *Plex {
	t.Helper()
	plexHost := os.Getenv("PLEX_HOST")
	plexToken := os.Getenv("PLEX_TOKEN")

	if plexHost == "" {
		t.Skip("PLEX_HOST not set")
	}

	conn, err := New(plexHost, plexToken)
	if err != nil {
		t.Fatalf("failed to create plex connection: %v", err)
	}
	return conn
}

func TestPlex_GetLibraries(t *testing.T) {
	conn := getTestConnection(t)
	libraries, err := conn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
}

func TestPlex_GetLibraryShows(t *testing.T) {
	conn := getTestConnection(t)
	libraries, err := conn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
	libs := libraries.Type(library.TypeShow)
	if len(libs) == 0 {
		t.Skip("no show libraries found")
	}
	lib := libs[0]

	err = conn.GetLibraryShows(lib, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.Shows) == 0 {
		t.Fatal("no shows found")
	}
}

func TestPlex_GetShowEpisodes(t *testing.T) {
	conn := getTestConnection(t)
	libraries, err := conn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
	libs := libraries.Type(library.TypeShow)
	if len(libs) == 0 {
		t.Skip("no show libraries found")
	}
	lib := libs[0]

	err = conn.GetLibraryShows(lib, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.Shows) == 0 {
		t.Fatal("no shows found")
	}

	sh := lib.Shows.FindTitle("Bluey")
	if sh == nil {
		t.Skip("Bluey not found")
	}

	err = conn.GetShowEpisodes(sh)
	if err != nil {
		t.Fatal(err)
	}
	if len(sh.Seasons) == 0 {
		t.Fatal("no seasons found")
	}
	for _, season := range sh.Seasons {
		if len(season.Episodes) == 0 {
			t.Fatal("no episodes found")
		}
	}
}

func TestPlex_PopulateLibraries(t *testing.T) {
	conn := getTestConnection(t)
	err := conn.InitLibraries()
	if err != nil {
		t.Fatal(err)
	}
	done := conn.PopulateLibraries()
	done()
	if len(conn.Libraries) == 0 {
		t.Fatal("no libraries found")
	}
}

func TestPlex_Scrobble(t *testing.T) {
	conn := getTestConnection(t)
	err := conn.Scrobble("9744")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlex_Scan(t *testing.T) {
	conn := getTestConnection(t)
	libraries, err := conn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
	libs := libraries.Type(library.TypeShow)
	if len(libs) == 0 {
		t.Skip("no show libraries found")
	}
	lib := libs[0]
	err = conn.ScanLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
}

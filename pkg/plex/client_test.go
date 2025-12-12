package plex

import (
	"fmt"
	"os"
	"testing"

	"github.com/kjbreil/go-plex/pkg/library"
)

var (
	plexHost  string
	plexToken string
	plexConn  *Plex
)

func init() {
	plexHost = os.Getenv("PLEX_HOST")
	plexToken = os.Getenv("PLEX_TOKEN")

	if plexHost != "" {
		var err error
		if plexConn, err = New(plexHost, plexToken); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func TestPlex_GetLibraries(t *testing.T) {
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	libraries, err := plexConn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
}

func TestPlex_GetLibraryShows(t *testing.T) {
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	libraries, err := plexConn.GetLibraries()
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

	err = plexConn.GetLibraryShows(lib, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.Shows) == 0 {
		t.Fatal("no shows found")
	}
}

func TestPlex_GetShowEpisodes(t *testing.T) {
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	libraries, err := plexConn.GetLibraries()
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

	err = plexConn.GetLibraryShows(lib, "")
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

	err = plexConn.GetShowEpisodes(sh)
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
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	err := plexConn.InitLibraries()
	if err != nil {
		t.Fatal(err)
	}
	done := plexConn.PopulateLibraries()
	done()
	if len(plexConn.Libraries) == 0 {
		t.Fatal("no libraries found")
	}
}

func TestPlex_Scrobble(t *testing.T) {
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	err := plexConn.Scrobble("9744")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlex_Scan(t *testing.T) {
	if plexConn == nil {
		t.Skip("PLEX_HOST not set")
	}
	libraries, err := plexConn.GetLibraries()
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
	err = plexConn.ScanLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
}

package plex

import (
	"fmt"
	"github.com/kjbreil/go-plex/library"
	"os"
	"testing"
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
	libraries, err := plexConn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
}

func TestPlex_GetLibraryShows(t *testing.T) {
	libraries, err := plexConn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
	lib := libraries.Type(library.TypeShow)[0]

	err = plexConn.GetLibraryShows(lib, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.Shows) == 0 {
		t.Fatal("no shows found")
	}
}

func TestPlex_GetShowEpisodes(t *testing.T) {
	libraries, err := plexConn.GetLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}
	lib := libraries.Type(library.TypeShow)[0]

	err = plexConn.GetLibraryShows(lib, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lib.Shows) == 0 {
		t.Fatal("no shows found")
	}

	sh := lib.Shows.Title("Bluey")

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
	libraries, err := plexConn.PopulateLibraries()
	if err != nil {
		t.Fatal(err)
	}
	if len(libraries) == 0 {
		t.Fatal("no libraries found")
	}

}

func TestPlex_Scrobble(t *testing.T) {
	// 9744
	err := plexConn.Scrobble("9744")
	if err != nil {
		t.Fatal(err)
	}
}

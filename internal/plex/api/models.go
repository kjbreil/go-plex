package api

import (
	"encoding/json"
	"strconv"

	"github.com/kjbreil/go-plex/pkg/library"
)

// LibrarySections metadata of your library contents.
type LibrarySections struct {
	MediaContainer struct {
		Directory library.Libraries `json:"Directory"`
	} `json:"MediaContainer"`
}

// SearchMediaContainer ...
type SearchMediaContainer struct {
	MediaContainer
	Provider []Provider
}

// MediaContainer contains media info.
type MediaContainer struct {
	Metadata            []Metadata `json:"Metadata"`
	AllowSync           bool       `json:"allowSync"`
	Identifier          string     `json:"identifier"`
	LibrarySectionID    int        `json:"librarySectionID"`
	LibrarySectionTitle string     `json:"librarySectionTitle"`
	LibrarySectionUUID  string     `json:"librarySectionUUID"`
	MediaTagPrefix      string     `json:"mediaTagPrefix"`
	MediaTagVersion     int        `json:"mediaTagVersion"`
	Size                int        `json:"size"`
}

// Provider ...
type Provider struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// Metadata ...
type Metadata struct {
	Player                Player       `json:"Player"`
	Session               Session      `json:"Session"`
	User                  User         `json:"User"`
	AddedAt               int          `json:"addedAt"`
	Art                   string       `json:"art"`
	ContentRating         string       `json:"contentRating"`
	Duration              int          `json:"duration"`
	GrandparentArt        string       `json:"grandparentArt"`
	GrandparentKey        string       `json:"grandparentKey"`
	GrandparentRatingKey  string       `json:"grandparentRatingKey"`
	GrandparentTheme      string       `json:"grandparentTheme"`
	GrandparentThumb      string       `json:"grandparentThumb"`
	GrandparentTitle      string       `json:"grandparentTitle"`
	GUID                  string       `json:"guid"`
	AltGUIDs              AltGUIDs     `json:"Guid"`
	Index                 int64        `json:"index"`
	Key                   string       `json:"key"`
	LastViewedAt          int          `json:"lastViewedAt"`
	LibrarySectionID      json.Number  `json:"librarySectionID"`
	LibrarySectionKey     string       `json:"librarySectionKey"`
	LibrarySectionTitle   string       `json:"librarySectionTitle"`
	OriginallyAvailableAt string       `json:"originallyAvailableAt"`
	ParentIndex           int64        `json:"parentIndex"`
	ParentKey             string       `json:"parentKey"`
	ParentRatingKey       string       `json:"parentRatingKey"`
	ParentThumb           string       `json:"parentThumb"`
	ParentTitle           string       `json:"parentTitle"`
	RatingCount           int          `json:"ratingCount"`
	AudienceRating        float64      `json:"audienceRating"`
	UserRating            float64      `json:"userRating"`
	Rating                Ratings      `json:"Rating"`
	RatingKey             string       `json:"ratingKey"`
	SessionKey            string       `json:"sessionKey"`
	Summary               string       `json:"summary"`
	Thumb                 string       `json:"thumb"`
	Media                 []Media      `json:"Media"`
	Title                 string       `json:"title"`
	TitleSort             string       `json:"titleSort"`
	Type                  string       `json:"type"`
	UpdatedAt             int          `json:"updatedAt"`
	ViewCount             int          `json:"viewCount"`
	ViewOffset            int          `json:"viewOffset"`
	Year                  int          `json:"year"`
	Director              []TaggedData `json:"Director"`
	Writer                []TaggedData `json:"Writer"`
}

// User plex server user. only difference is id is a string.
type User struct {
	// ID is an int when signing in to Plex.tv but a string when access own server
	ID                  string `json:"id"`
	UUID                string `json:"uuid"`
	Email               string `json:"email"`
	JoinedAt            string `json:"joined_at"`
	Username            string `json:"username"`
	Thumb               string `json:"thumb"`
	HasPassword         bool   `json:"hasPassword"`
	AuthToken           string `json:"authToken"`
	AuthenticationToken string `json:"authenticationToken"`
	Subscription        struct {
		Active   bool     `json:"active"`
		Status   string   `json:"Active"`
		Plan     string   `json:"lifetime"`
		Features []string `json:"features"`
	} `json:"subscription"`
	Roles struct {
		Roles []string `json:"roles"`
	} `json:"roles"`
	Entitlements []string `json:"entitlements"`
	ConfirmedAt  string   `json:"confirmedAt"`
	ForumID      string   `json:"forumId"`
	RememberMe   bool     `json:"rememberMe"`
	Title        string   `json:"title"`
}

// TaggedData ...
type TaggedData struct {
	Tag    string      `json:"tag"`
	Filter string      `json:"filter"`
	ID     json.Number `json:"id"`
}

// Player ...
type Player struct {
	Address             string `json:"address"`
	Device              string `json:"device"`
	Local               bool   `json:"local"`
	MachineIdentifier   string `json:"machineIdentifier"`
	Model               string `json:"model"`
	Platform            string `json:"platform"`
	PlatformVersion     string `json:"platformVersion"`
	Product             string `json:"product"`
	Profile             string `json:"profile"`
	RemotePublicAddress string `json:"remotePublicAddress"`
	State               string `json:"state"`
	Title               string `json:"title"`
	UserID              int    `json:"userID"`
	Vendor              string `json:"vendor"`
	Version             string `json:"version"`
}

// Session ...
type Session struct {
	Bandwidth int    `json:"bandwidth"`
	ID        string `json:"id"`
	Location  string `json:"location"`
}

// CurrentSessions metadata of users consuming media.
type CurrentSessions struct {
	MediaContainer struct {
		Metadata []Metadata `json:"Metadata"`
		Size     int        `json:"size"`
	} `json:"MediaContainer"`
}

// Media media info.
type Media struct {
	AspectRatio           json.Number `json:"aspectRatio"`
	AudioChannels         int         `json:"audioChannels"`
	AudioCodec            string      `json:"audioCodec"`
	AudioProfile          string      `json:"audioProfile"`
	Bitrate               int         `json:"bitrate"`
	Container             string      `json:"container"`
	Duration              int         `json:"duration"`
	Has64bitOffsets       bool        `json:"has64bitOffsets"`
	Height                int         `json:"height"`
	ID                    json.Number `json:"id"`
	OptimizedForStreaming BoolOrInt   `json:"optimizedForStreaming"` // plex can return int (GetMetadata(), GetPlaylist()) or boolean (GetSessions()): 0 or 1; true or false

	Selected        bool   `json:"selected"`
	VideoCodec      string `json:"videoCodec"`
	VideoFrameRate  string `json:"videoFrameRate"`
	VideoProfile    string `json:"videoProfile"`
	VideoResolution string `json:"videoResolution"`
	Width           int    `json:"width"`
	Part            []Part `json:"Part"`
}

// Part ...
type Part struct {
	AudioProfile          string      `json:"audioProfile"`
	Container             string      `json:"container"`
	Decision              string      `json:"decision"`
	Duration              int         `json:"duration"`
	File                  string      `json:"file"`
	Has64bitOffsets       bool        `json:"has64bitOffsets"`
	HasThumbnail          string      `json:"hasThumbnail"`
	ID                    json.Number `json:"id"`
	Key                   string      `json:"key"`
	OptimizedForStreaming BoolOrInt   `json:"optimizedForStreaming"`
	Selected              bool        `json:"selected"`
	Size                  int         `json:"size"`
	Stream                []Stream    `json:"Stream"`
	VideoProfile          string      `json:"videoProfile"`
}

// Stream ...
type Stream struct {
	AlbumGain          string      `json:"albumGain"`
	AlbumPeak          string      `json:"albumPeak"`
	AlbumRange         string      `json:"albumRange"`
	Anamorphic         bool        `json:"anamorphic"`
	AudioChannelLayout string      `json:"audioChannelLayout"`
	BitDepth           int         `json:"bitDepth"`
	Bitrate            int         `json:"bitrate"`
	BitrateMode        string      `json:"bitrateMode"`
	Cabac              string      `json:"cabac"`
	Channels           int         `json:"channels"`
	ChromaLocation     string      `json:"chromaLocation"`
	ChromaSubsampling  string      `json:"chromaSubsampling"`
	Codec              string      `json:"codec"`
	CodecID            string      `json:"codecID"`
	ColorRange         string      `json:"colorRange"`
	ColorSpace         string      `json:"colorSpace"`
	Default            bool        `json:"default"`
	DisplayTitle       string      `json:"displayTitle"`
	Duration           string      `json:"duration"`
	FrameRate          float64     `json:"frameRate"`
	FrameRateMode      string      `json:"frameRateMode"`
	Gain               string      `json:"gain"`
	HasScalingMatrix   bool        `json:"hasScalingMatrix"`
	Height             int         `json:"height"`
	ID                 json.Number `json:"id"`
	Index              int         `json:"index"`
	Language           string      `json:"language"`
	LanguageCode       string      `json:"languageCode"`
	Level              int         `json:"level"`
	Location           string      `json:"location"`
	Loudness           string      `json:"loudness"`
	Lra                string      `json:"lra"`
	Peak               string      `json:"peak"`
	PixelAspectRatio   string      `json:"pixelAspectRatio"`
	PixelFormat        string      `json:"pixelFormat"`
	Profile            string      `json:"profile"`
	RefFrames          int         `json:"refFrames"`
	SamplingRate       int         `json:"samplingRate"`
	ScanType           string      `json:"scanType"`
	Selected           bool        `json:"selected"`
	StreamIdentifier   string      `json:"streamIdentifier"`
	StreamType         int         `json:"streamType"`
	Width              int         `json:"width"`
}

// AltGUIDs represents a list of Globally Unique Identifier for a metadata provider that is not.
type AltGUIDs []AltGUID

func (ag AltGUIDs) TVDB() int {
	for _, alt := range ag {
		if len(alt.ID) >= 4 && alt.ID[:4] == "tvdb" {
			id, err := strconv.Atoi(alt.ID[7:])
			if err != nil {
				return 0
			}
			return id
		}
	}

	return 0
}

func (ag AltGUIDs) TMDB() int {
	for _, alt := range ag {
		if len(alt.ID) >= 4 && alt.ID[:4] == "tmdb" {
			id, err := strconv.Atoi(alt.ID[7:])
			if err != nil {
				return 0
			}
			return id
		}
	}

	return 0
}

// AltGUID represents a Globally Unique Identifier for a metadata provider that is not actively being used.
type AltGUID struct {
	ID string `json:"id"`
}

// SearchResults ...
type SearchResults struct {
	MediaContainer SearchMediaContainer `json:"MediaContainer"`
}

// SearchResultsEpisode contains metadata about an episode.
type SearchResultsEpisode struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

// MediaMetadata wraps MediaContainer for metadata responses.
type MediaMetadata struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

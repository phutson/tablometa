package tablometadata

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TABLODATEPATTERN = `"[0-9]+-(0?[1-9]|[1][0-2])-[0-9]+T(0?[0-9]|1[0-9]|2[0-3]):[0-9]+Z"`
	RFC3339PATTERN   = `"[0-9]+-(0?[1-9]|[1][0-2])-[0-9]+T(0?[0-9]|1[0-9]|2[0-3]):[0-9]+:[0-9]+\.[0-9]+Z"`
)

type TabloType interface {
	GetTabloType() string
}

type TabloDate struct {
	StoredTime time.Time
}

func (tt *TabloDate) UnmarshalJSON(data []byte) error {
	var rfcPattern = regexp.MustCompile(RFC3339PATTERN)
	var tabloPattern = regexp.MustCompile(TABLODATEPATTERN)
	if rfcPattern.Match(data) {
		return tt.StoredTime.UnmarshalJSON(data)
	} else if tabloPattern.Match(data) {
		workingDateString := string(data[:])
		workingDateString = strings.TrimSuffix(workingDateString, `Z"`)
		workingDateString = workingDateString + `:00.00Z"`
		return tt.StoredTime.UnmarshalJSON([]byte(workingDateString))
	} else {
		return errors.New("unknown date pattern")
	}
}

func (tt TabloDate) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(tt.StoredTime)
	if err != nil {
		return nil, err
	}
	jsonString := strings.Replace(string(jsonData[:]), `:00.00Z"`, `Z"`, 1)
	return []byte(jsonString), nil
}

type Relationships struct {
	RecMovie   int   `json:"recMovie"`
	RecChannel int   `json:"recChannel"`
	Genres     []int `json:"genres"`
	RecSeason  int   `json:"recSeason"`
	RecSeries  int   `json:"recSeries"`
}

func (tr Relationships) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	jsonData = append(jsonData, []byte("{ ")...)
	needsComma := false

	if tr.RecMovie >= 0 {
		jsonData = append(jsonData, []byte(`"recMovie": `)...)
		jsonData = append(jsonData, []byte(strconv.Itoa(tr.RecMovie))...)
		needsComma = true
	}

	if tr.RecSeason >= 0 {
		if needsComma {
			jsonData = append(jsonData, []byte(`, `)...)
		}
		jsonData = append(jsonData, []byte(`"recSeason": `)...)
		jsonData = append(jsonData, []byte(strconv.Itoa(tr.RecSeason))...)
		needsComma = true
	}
	if tr.RecSeries >= 0 {
		if needsComma {
			jsonData = append(jsonData, []byte(`, `)...)
		}
		jsonData = append(jsonData, []byte(`"recSeries": `)...)
		jsonData = append(jsonData, []byte(strconv.Itoa(tr.RecSeries))...)
		needsComma = true
	}
	if tr.RecChannel >= 0 {
		if needsComma {
			jsonData = append(jsonData, []byte(`, `)...)
		}
		jsonData = append(jsonData, []byte(`"recChannel": `)...)
		jsonData = append(jsonData, []byte(strconv.Itoa(tr.RecChannel))...)
		needsComma = true
	}
	if tr.Genres != nil && len(tr.Genres) > 0 {
		genreJSONData, err := json.Marshal(tr.Genres)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, []byte(`"genres": `)...)
		jsonData = append(jsonData, genreJSONData...)
		needsComma = true

	}

	jsonData = append(jsonData, []byte(" }")...)
	return jsonData, nil
}

type VideoInfo struct {
	State               string  `json:"state"`
	Size                uint64  `json:"size"`
	Width               int     `json:"width"`
	Height              int     `json:"height"`
	Duration            float32 `json:"duration"`
	ScheduleOffsetStart float32 `json:"scheduleOffsetStart"`
	ScheduleOffsetEnd   float32 `json:"scheduleOffsetEnd"`
}

type UserInfo struct {
	UserType  string  `json:"type"`
	Watched   bool    `json:"watched"`
	Protected bool    `json:"protected"`
	Position  float32 `json:"position"`
}

type ImageData struct {
	Type       string `json:"type"`
	ImageID    int    `json:"imageID"`
	ImageType  string `json:"imageType"`
	ImageStyle string `json:"imageStyle"`
}

type ImageJSONData struct {
	Images []ImageData `json:"images"`
}

type ClientJSON struct {
	Title            string        `json:"title"`
	Plot             string        `json:"plot"`
	Runtime          int           `json:"runtime"`
	MPAARating       string        `json:"mpaaRating"`
	ReleaseYear      int           `json:"releaseYear"`
	Cast             []string      `json:"cast"`
	Directors        []string      `json:"directors"`
	QualityRating    float32       `json:"qualityRating"`
	Relationships    Relationships `json:"relationships"`
	Type             string        `json:"type"`
	ObjectId         int           `json:"objectID"`
	AirDate          TabloDate     `json:"airDate"`
	ScheduleDuration float32       `json:"scheduleDuration"`
	Video            VideoInfo     `json:"video"`
	User             UserInfo      `json:"user"`
	Description      string        `json:"description"`
	EpisodeNumber    int           `json:"episodeNumber"`
	SeasonNumber     int           `json:"seasonNumber"`
	OriginalAirDate  string        `json:"originalAirDate"`
	Qualifiers       []string      `json:"qualifiers"`
	Duration         int           `json:"duration"`
}

func (tr ClientJSON) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	jsonData = append(jsonData, []byte("{ ")...)

	jsonData = append(jsonData, []byte(`"type": "`)...)
	jsonData = append(jsonData, []byte(tr.Type)...)
	jsonData = append(jsonData, []byte(`", `)...)

	if tr.Type == "recMovieAiring" {
		jsonData = append(jsonData, []byte(`"objectID": `)...)
		jsonData = append(jsonData, []byte(strconv.Itoa(tr.ObjectId))...)
		jsonData = append(jsonData, []byte(`, `)...)

		jsonData = append(jsonData, []byte(`"airDate": "`)...)
		airDateData, err := tr.AirDate.MarshalJSON()
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, airDateData...)
		jsonData = append(jsonData, []byte(`", `)...)

		jsonData = append(jsonData, []byte(`"scheduleDuration": `)...)
		jsonData = append(jsonData, []byte(strconv.FormatFloat(float64(tr.ScheduleDuration), 'f', -1, 32))...)
		jsonData = append(jsonData, []byte(`, `)...)

		jsonData = append(jsonData, []byte(`"relationships": `)...)
		relationshipJSONData, err := tr.Relationships.MarshalJSON()
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, relationshipJSONData...)
		jsonData = append(jsonData, []byte(`, `)...)

		jsonData = append(jsonData, []byte(`"video": `)...)
		videoJSONData, err := json.Marshal(tr.Video)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, videoJSONData...)
		jsonData = append(jsonData, []byte(`, `)...)

		jsonData = append(jsonData, []byte(`"user": `)...)
		userJSONData, err := json.Marshal(tr.User)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, userJSONData...)

	}
	jsonData = append(jsonData, []byte(" }")...)
	return jsonData, nil
}

type MovieAiring struct {
	JSONForClient ClientJSON    `json:"jsonForClient"`
	ImageJSON     ImageJSONData `json:"imageJson"`
}

func (ma *MovieAiring) GetTabloType() string {
	return ma.JSONForClient.Type
}

type RecMovie struct {
	JSONForClient ClientJSON    `json:"jsonForClient"`
	ImageJSON     ImageJSONData `json:"imageJson"`
}

func (rm *RecMovie) GetTabloType() string {
	return rm.JSONForClient.Type
}

type RecSeries struct {
	JSONForClient ClientJSON    `json:"jsonForClient"`
	ImageJSON     ImageJSONData `json:"imageJson"`
}

func (rs *RecSeries) GetTabloType() string {
	return rs.JSONForClient.Type
}

type RecSeason struct {
	JSONForClient ClientJSON `json:"jsonForClient"`
}

func (rs *RecSeason) GetTabloType() string {
	return rs.JSONForClient.Type
}

type RecEpisode struct {
	JSONForClient ClientJSON    `json:"jsonForClient"`
	ImageJSON     ImageJSONData `json:"imageJson"`
}

func (re *RecEpisode) GetTabloType() string {
	return re.JSONForClient.Type
}

type Recording struct {
	RecordedEpisode RecEpisode  `json:"recEpisode"`
	RecordedSeries  RecSeries   `json:"recSeries"`
	RecordedSeason  RecSeason   `json:"recSeason"`
	Airing          MovieAiring `json:"recMovieAiring"`
	RecordedMovie   RecMovie    `json:"recMovie"`
}

func buildJSONObject(objectToMarshal interface{}, objectKey string, needsComma bool, totalJSONData []byte) ([]byte, error) {
	objectJSONData, err := json.Marshal(objectToMarshal)
	if err != nil {
		return nil, err
	}
	if needsComma {
		totalJSONData = append(totalJSONData, ',')
	}
	totalJSONData = append(totalJSONData, []byte(` "`)...)

	totalJSONData = append(totalJSONData, []byte(objectKey)...)
	totalJSONData = append(totalJSONData, []byte(`": `)...)
	totalJSONData = append(totalJSONData, objectJSONData...)

	return totalJSONData, nil
}

func (tr Recording) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	jsonData = append(jsonData, []byte("{ ")...)
	var needsComma = false
	var err error

	if len(tr.RecordedEpisode.GetTabloType()) > 0 {
		episodeField, wasFound := reflect.TypeOf(tr).FieldByName("RecordedEpisode")
		if wasFound {
			objectKey := episodeField.Tag.Get("json")
			jsonData, err = buildJSONObject(tr.RecordedEpisode, objectKey, needsComma, jsonData)
			if err != nil {
				return nil, err
			}
			needsComma = true
		}
	}

	if len(tr.RecordedSeries.GetTabloType()) > 0 {
		episodeField, wasFound := reflect.TypeOf(tr).FieldByName("RecordedSeries")
		if wasFound {
			objectKey := episodeField.Tag.Get("json")
			jsonData, err = buildJSONObject(tr.RecordedSeries, objectKey, needsComma, jsonData)
			if err != nil {
				return nil, err
			}
			needsComma = true
		}
	}

	if len(tr.RecordedSeason.GetTabloType()) > 0 {
		episodeField, wasFound := reflect.TypeOf(tr).FieldByName("RecordedSeason")
		if wasFound {
			objectKey := episodeField.Tag.Get("json")
			jsonData, err = buildJSONObject(tr.RecordedSeason, objectKey, needsComma, jsonData)
			if err != nil {
				return nil, err
			}
			needsComma = true
		}
	}
	if len(tr.Airing.GetTabloType()) > 0 {
		episodeField, wasFound := reflect.TypeOf(tr).FieldByName("Airing")
		if wasFound {
			objectKey := episodeField.Tag.Get("json")
			jsonData, err = buildJSONObject(tr.Airing, objectKey, needsComma, jsonData)
			if err != nil {
				return nil, err
			}
			needsComma = true
		}
	}
	if len(tr.RecordedMovie.GetTabloType()) > 0 {
		episodeField, wasFound := reflect.TypeOf(tr).FieldByName("RecordedMovie")
		if wasFound {
			objectKey := episodeField.Tag.Get("json")
			jsonData, err = buildJSONObject(tr.RecordedMovie, objectKey, needsComma, jsonData)
			if err != nil {
				return nil, err
			}
			needsComma = true
		}
	}

	jsonData = append(jsonData, []byte(" }")...)
	return jsonData, nil
}

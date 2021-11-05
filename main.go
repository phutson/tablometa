package tablometadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	TABLODATEPATTERN         = `"[0-9]+-(0?[1-9]|[1][0-2])-[0-9]+T(0?[0-9]|1[0-9]|2[0-3]):[0-9]+Z"`
	RFC3339PATTERN           = `"[0-9]+-(0?[1-9]|[1][0-2])-[0-9]+T(0?[0-9]|1[0-9]|2[0-3]):[0-9]+:[0-9]+\.[0-9]+Z"`
	CLIENTRECAIRINGMOVIEFMT  = `{"%s":"%s","%s":%d,"%s":%s,"%s":%.1f,"%s":%s,"%s":%s,"%s":%s} `
	CLIENTRECMOVIE           = `{"%s":"%s","%s":"%s","%s":%d,"%s":"%s","%s":%d,"%s":%s,"%s":%s,"%s":%.3f,"%s":%s,"%s":"%s","%s":%d}`
	RECMOVIEFMT              = `{"%s":%s,"%s":%s}`
	RELATIONSHIPSMOVIESFMT   = `{"%s":%d,"%s":%d}`
	RELATIONSHIPSEPISODESFMT = `{"%s":%d,"%s":%d,"%s":%d}`
	RELATIONSHIPSGENRESFMT   = `{"%s":%s}`
	RECORDINGMOVIEFMT        = `{"%s":%s,"%s":%s}`
	VIDEOINFOFMT             = `{"%s":"%s","%s":%d,"%s":%d,"%s":%d,"%s":%.1f,"%s":%.1f,"%s":%.1f}`
	USERINFOFMT              = `{"%s":"%s","%s":%t,"%s":%t,"%s":%.1f}`
)

func getJSONFieldNameByName(structureOfInterest interface{}, fieldName string) (string, error) {
	episodeField, wasFound := reflect.TypeOf(structureOfInterest).FieldByName(fieldName)
	if wasFound {
		jsonFieldName := episodeField.Tag.Get("json")
		if len(jsonFieldName) < 1 {
			return "", errors.New("json tag not found")
		} else {
			return jsonFieldName, nil
		}
	} else {
		return "", errors.New("field not found")
	}
}

type TabloType interface {
	GetTabloType() string
}

type TabloDate struct {
	StoredTime time.Time
}

func (tt TabloDate) Format(layout string) string {
	return tt.StoredTime.Format(layout)
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
	jsonString := fmt.Sprintf(`"%s"`, tt.StoredTime.Format("2006-01-02T15:04Z"))
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
	recChannelFieldName, err := getJSONFieldNameByName(tr, "RecChannel")
	if err != nil {
		return nil, err
	}

	if tr.RecMovie > 0 {
		recMovieFieldName, err := getJSONFieldNameByName(tr, "RecMovie")
		if err != nil {
			return nil, err
		}

		jsonString := fmt.Sprintf(RELATIONSHIPSMOVIESFMT, recMovieFieldName, tr.RecMovie, recChannelFieldName, tr.RecChannel)
		jsonData = append(jsonData, []byte(jsonString)...)
	}

	if tr.RecSeason > 0 {

		recSeasonFieldName, err := getJSONFieldNameByName(tr, "RecSeason")
		if err != nil {
			return nil, err
		}
		recSeriesFieldName, err := getJSONFieldNameByName(tr, "RecSeries")
		if err != nil {
			return nil, err
		}

		jsonString := fmt.Sprintf(RELATIONSHIPSEPISODESFMT, recSeasonFieldName, tr.RecSeason, recSeriesFieldName, tr.RecSeries, recChannelFieldName, tr.RecChannel)
		jsonData = append(jsonData, []byte(jsonString)...)
	}
	if tr.Genres != nil && len(tr.Genres) > 0 {
		genreJSONData, err := json.Marshal(tr.Genres)
		if err != nil {
			return nil, err
		}
		genresFieldName, err := getJSONFieldNameByName(tr, "Genres")
		if err != nil {
			return nil, err
		}
		jsonString := fmt.Sprintf(RELATIONSHIPSGENRESFMT, genresFieldName, string(genreJSONData[:]))
		jsonData = append(jsonData, []byte(jsonString)...)
	}

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

func (vr VideoInfo) MarshalJSON() ([]byte, error) {

	var jsonData []byte
	stateFieldName, err := getJSONFieldNameByName(vr, "State")
	if err != nil {
		return nil, err
	}

	sizeFieldName, err := getJSONFieldNameByName(vr, "Size")
	if err != nil {
		return nil, err
	}
	widthFieldName, err := getJSONFieldNameByName(vr, "Width")
	if err != nil {
		return nil, err
	}
	heightFieldName, err := getJSONFieldNameByName(vr, "Height")
	if err != nil {
		return nil, err
	}
	durationFieldName, err := getJSONFieldNameByName(vr, "Duration")
	if err != nil {
		return nil, err
	}

	scheduleOffsetStartFieldName, err := getJSONFieldNameByName(vr, "ScheduleOffsetStart")
	if err != nil {
		return nil, err
	}

	scheduleOffsetEndFieldName, err := getJSONFieldNameByName(vr, "ScheduleOffsetEnd")
	if err != nil {
		return nil, err
	}

	jsonString := fmt.Sprintf(VIDEOINFOFMT, stateFieldName, vr.State, sizeFieldName, vr.Size, widthFieldName, vr.Width, heightFieldName, vr.Height,
		durationFieldName, vr.Duration, scheduleOffsetStartFieldName, vr.ScheduleOffsetStart, scheduleOffsetEndFieldName, vr.ScheduleOffsetEnd)
	jsonData = append(jsonData, []byte(jsonString)...)

	return jsonData, nil
}

type UserInfo struct {
	UserType  string  `json:"type"`
	Watched   bool    `json:"watched"`
	Protected bool    `json:"protected"`
	Position  float32 `json:"position"`
}

func (ur UserInfo) MarshalJSON() ([]byte, error) {

	var jsonData []byte
	usertypeFieldName, err := getJSONFieldNameByName(ur, "UserType")
	if err != nil {
		return nil, err
	}

	watchedFieldName, err := getJSONFieldNameByName(ur, "Watched")
	if err != nil {
		return nil, err
	}
	protectedFieldName, err := getJSONFieldNameByName(ur, "Protected")
	if err != nil {
		return nil, err
	}
	positionFieldName, err := getJSONFieldNameByName(ur, "Position")
	if err != nil {
		return nil, err
	}
	jsonString := fmt.Sprintf(USERINFOFMT, usertypeFieldName, ur.UserType, watchedFieldName, ur.Watched, protectedFieldName, ur.Protected, positionFieldName, ur.Position)
	jsonData = append(jsonData, []byte(jsonString)...)

	return jsonData, nil
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
	ObjectID         int           `json:"objectID"`
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
	typeFieldName, err := getJSONFieldNameByName(tr, "Type")
	if err != nil {
		return nil, err
	}

	relationshipsFieldName, err := getJSONFieldNameByName(tr, "Relationships")
	if err != nil {
		return nil, err
	}

	relationshipJSONData, err := tr.Relationships.MarshalJSON()
	if err != nil {
		return nil, err
	}

	objectIDFieldName, err := getJSONFieldNameByName(tr, "ObjectID")
	if err != nil {
		return nil, err
	}

	if tr.Type == "recMovieAiring" {
		//type objecid airdate

		airDateFieldName, err := getJSONFieldNameByName(tr, "AirDate")
		if err != nil {
			return nil, err
		}

		airDateData, err := tr.AirDate.MarshalJSON()
		if err != nil {
			return nil, err
		}

		scheduleDurationFieldName, err := getJSONFieldNameByName(tr, "ScheduleDuration")
		if err != nil {
			return nil, err
		}

		videoFieldName, err := getJSONFieldNameByName(tr, "Video")
		if err != nil {
			return nil, err
		}

		videoJSONData, err := json.Marshal(tr.Video)
		if err != nil {
			return nil, err
		}

		userFieldName, err := getJSONFieldNameByName(tr, "User")
		if err != nil {
			return nil, err
		}

		userJSONData, err := json.Marshal(tr.User)
		if err != nil {
			return nil, err
		}

		jsonString := fmt.Sprintf(CLIENTRECAIRINGMOVIEFMT, typeFieldName, tr.Type, objectIDFieldName, tr.ObjectID, airDateFieldName, string(airDateData[:]), scheduleDurationFieldName,
			tr.ScheduleDuration, relationshipsFieldName, string(relationshipJSONData[:]), videoFieldName, string(videoJSONData[:]), userFieldName, userJSONData)
		jsonData = append(jsonData, []byte(jsonString)...)
	} else if tr.Type == "recMovie" {
		titleFieldName, err := getJSONFieldNameByName(tr, "Title")
		if err != nil {
			return nil, err
		}
		plotFieldName, err := getJSONFieldNameByName(tr, "Plot")
		if err != nil {
			return nil, err
		}
		runtimeFieldName, err := getJSONFieldNameByName(tr, "Runtime")
		if err != nil {
			return nil, err
		}
		mpaaRatingFieldName, err := getJSONFieldNameByName(tr, "MPAARating")
		if err != nil {
			return nil, err
		}
		releaseYearFieldName, err := getJSONFieldNameByName(tr, "ReleaseYear")
		if err != nil {
			return nil, err
		}
		castFieldName, err := getJSONFieldNameByName(tr, "Cast")
		if err != nil {
			return nil, err
		}

		castJSONData, err := json.Marshal(tr.Cast)
		if err != nil {
			return nil, err
		}

		directorsFieldName, err := getJSONFieldNameByName(tr, "Directors")
		if err != nil {
			return nil, err
		}
		directorsJSONData, err := json.Marshal(tr.Directors)
		if err != nil {
			return nil, err
		}
		qualityRationFieldName, err := getJSONFieldNameByName(tr, "QualityRating")
		if err != nil {
			return nil, err
		}

		jsonString := fmt.Sprintf(CLIENTRECMOVIE, titleFieldName, tr.Title, plotFieldName, tr.Plot, runtimeFieldName, tr.Runtime,
			mpaaRatingFieldName, tr.MPAARating, releaseYearFieldName, tr.ReleaseYear, castFieldName, string(castJSONData[:]), directorsFieldName, directorsJSONData,
			qualityRationFieldName, tr.QualityRating, relationshipsFieldName, relationshipJSONData, typeFieldName, tr.Type, objectIDFieldName, tr.ObjectID)
		jsonData = append(jsonData, []byte(jsonString)...)

	}
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

func (tr Recording) MarshalJSON() ([]byte, error) {
	var jsonData []byte

	if len(tr.Airing.GetTabloType()) > 0 {
		airingFieldName, err := getJSONFieldNameByName(tr, "Airing")
		if err != nil {
			return nil, err
		}

		airingJSONData, err := json.Marshal(tr.Airing)
		if err != nil {
			return nil, err
		}

		recMovieFieldName, err := getJSONFieldNameByName(tr, "RecordedMovie")
		if err != nil {
			return nil, err
		}

		recMovieJSONData, err := json.Marshal(tr.RecordedMovie)
		if err != nil {
			return nil, err
		}
		jsonString := fmt.Sprintf(RECORDINGMOVIEFMT, airingFieldName, string(airingJSONData[:]), recMovieFieldName, string(recMovieJSONData[:]))

		jsonData = append(jsonData, []byte(jsonString)...)
	}

	return jsonData, nil
}

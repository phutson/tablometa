package tablometadata_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	tablometadata "github.com/phutson/tablometa"
)

func TestParseMovie(t *testing.T) {
	var correctMovieJSON = `{"recMovieAiring":{"jsonForClient":{"type":"recMovieAiring","objectID":117665,"airDate":"2016-11-06T23:00Z","scheduleDuration":7200.0,"relationships":{"recMovie":117666,"recChannel":5465},"video":{"state":"finished","size":3415293952,"width":1280,"height":720,"duration":7520.0,"scheduleOffsetStart":-15.0,"scheduleOffsetEnd":304.0},"user":{"type":"recordingUserInfo","watched":false,"protected":false,"position":0.0}},"imageJson":{"images":[{"type":"image","imageID":123122,"imageType":"snapshot","imageStyle":"snapshot"}]}},"recMovie":{"jsonForClient":{"title":"Buying the Cow","plot":"A man hits the dating scene when his girlfriend gives him two months to decide whether or not he wants to marry her. Uncertain of commitment he spots another woman and instantly falls for her, but when she disappears he decides the only way to be sure of the relationship is to track the mysterious girl down.","runtime":5160,"mpaaRating":"r","releaseYear":2001,"cast":["Jerry O'Connell","Bridgette L. Wilson","Ryan Reynolds","Alyssa Milano","Annabeth Gish","Bill Bellamy","Brian Beacock","C.C. Boyce","Bix Barnaba","Erinn Bartlett","Adam Bitterman","Sonya Eddy","Nipper Knapp","Ron Livingston","Nina Petronzio"],"directors":["Walt Becker"],"qualityRating":0.250,"relationships":{"genres":[1063]},"type":"recMovie","objectID":117666},"imageJson":{"images":[{"type":"image","imageID":114189,"imageType":"movie_2x3_small","imageStyle":"thumbnail"},{"type":"image","imageID":114190,"imageType":"iconic_4x3_large","imageStyle":"cover"},{"type":"image","imageID":114191,"imageType":"iconic_4x3_large","imageStyle":"background"}]}}}`
	var recording tablometadata.Recording

	err := json.Unmarshal([]byte(correctMovieJSON), &recording)
	if err != nil {
		t.Fatal(`Failed to unmarshal string`)
	}

	jsonData, err := json.Marshal(recording)
	if err != nil {
		fmt.Print(err)
		t.Fatal("Failed to Marshal string")
	}
	if !(strings.Contains(string(jsonData[:]), correctMovieJSON)) {
		t.Fail()
	}
}

func TestParseEpisode(t *testing.T) {
	var correctEpisodeJSON = `{"recEpisode":{"jsonForClient":{"type":"recEpisode","title":"The Virgin Sacrifice","description":"Manfred leads the Midnighters to take back the town from the evil forces occupying it; while Bobo tries to save Fiji, Olivia and Creek confront the wraiths; Manfred, Lem, Joe and the Rev work to kill the demon and close the veil.","episodeNumber":10,"seasonNumber":1,"airDate":"2017-09-19T05:00Z","originalAirDate":"2017-09-18","scheduleDuration":3600,"qualifiers":["cc"],"relationships":{"recSeason":301535,"recSeries":301534,"recChannel":185238},"video":{"state":"finished","size":5302616064,"width":1920,"height":1080,"duration":5417.0,"scheduleOffsetStart":-15.0,"scheduleOffsetEnd":1805.0},"user":{"type":"recordingUserInfo","watched":false,"protected":false,"position":0.0},"objectID":343176},"imageJson":{"images":[{"type":"image","imageID":353557,"imageType":"snapshot","imageStyle":"snapshot"}]}},"recSeries":{"jsonForClient":{"title":"Midnight, Texas","description":"Based on Charlaine Harris' book series by the same name, \"Midnight, Texas\" follows the lives of the inhabitants of a small town where the concept of normal is relative. A haven for vampires, witches, psychics, hit men and others with extraordinary backgrounds, Midnight gives outsiders a place to belong. The town members form a strong and unlikely family as they work together to fend off the pressures of unruly biker gangs, questioning police officers and shades of their own dangerous pasts.","originalAirDate":"2017-07-24","duration":3600,"cast":["François Arnaud","Dylan Bruce","Parisa Fitz-Henley","Arielle Kebbel","Sarah Ramos","Peter Mensah","Yul Vazquez","Jason Lewis","Sean Bridgers"],"relationships":{"genres":[108,335,100019]},"objectID":301534,"type":"recSeries"},"imageJson":{"images":[{"type":"image","imageID":290612,"imageType":"series_3x4_small","imageStyle":"thumbnail"},{"type":"image","imageID":290613,"imageType":"series_4x3_large","imageStyle":"cover"},{"type":"image","imageID":290614,"imageType":"iconic_4x3_large","imageStyle":"background"}]}},"recSeason":{"jsonForClient":{"seasonNumber":1,"relationships":{"recSeries":301534},"objectID":301535,"type":"recSeason"}}}`

	var recording tablometadata.Recording

	err := json.Unmarshal([]byte(correctEpisodeJSON), &recording)
	if err != nil {
		t.Fatal(`Failed to unmarshal string`)
	}

	jsonData, err := json.Marshal(recording)
	if err != nil {
		fmt.Print(err)
		t.Fatal("Failed to Marshal string")
	}

	if !(strings.Contains(string(jsonData[:]), correctEpisodeJSON)) {
		t.Fail()
	}

}

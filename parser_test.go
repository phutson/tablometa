package tablometadata_test

import (
	"encoding/json"
	"fmt"
	"plh9.com/tablometadata"
	"strings"
	"testing"
)

func TestParseMovie(t *testing.T) {
	var correctMovieJson = `
	{
		"recMovieAiring": {
			"jsonForClient": {
				"type": "recMovieAiring",
				"objectID": 117665,
				"airDate": "2016-11-06T23:00Z",
				"scheduleDuration": 7200.0,
				"relationships": {
					"recMovie": 117666,
					"recChannel": 5465
				},
				"video": {
					"state": "finished",
					"size": 3415293952,
					"width": 1280,
					"height": 720,
					"duration": 7520.0,
					"scheduleOffsetStart": -15.0,
					"scheduleOffsetEnd": 304.0
				},
				"user": {
					"type": "recordingUserInfo",
					"watched": false,
					"protected": false,
					"position": 0.0
				}
			},
			"imageJson": {
				"images": [
					{
						"type": "image",
						"imageID": 123122,
						"imageType": "snapshot",
						"imageStyle": "snapshot"
					}
				]
			}
		},
		"recMovie": {
			"jsonForClient": {
				"title": "Buying the Cow",
				"plot": "A man hits the dating scene when his girlfriend gives him two months to decide whether or not he wants to marry her. Uncertain of commitment he spots another woman and instantly falls for her, but when she disappears he decides the only way to be sure of the relationship is to track the mysterious girl down.",
				"runtime": 5160,
				"mpaaRating": "r",
				"releaseYear": 2001,
				"cast": [
					"Jerry O'Connell",
					"Bridgette L. Wilson",
					"Ryan Reynolds",
					"Alyssa Milano",
					"Annabeth Gish",
					"Bill Bellamy",
					"Brian Beacock",
					"C.C. Boyce",
					"Bix Barnaba",
					"Erinn Bartlett",
					"Adam Bitterman",
					"Sonya Eddy",
					"Nipper Knapp",
					"Ron Livingston",
					"Nina Petronzio"
				],
				"directors": [
					"Walt Becker"
				],
				"qualityRating": 0.250,
				"relationships": {
					"genres": [
						1063
					]
				},
				"type": "recMovie",
				"objectID": 117666
			},
			"imageJson": {
				"images": [
					{
						"type": "image",
						"imageID": 114189,
						"imageType": "movie_2x3_small",
						"imageStyle": "thumbnail"
					},
					{
						"type": "image",
						"imageID": 114190,
						"imageType": "iconic_4x3_large",
						"imageStyle": "cover"
					},
					{
						"type": "image",
						"imageID": 114191,
						"imageType": "iconic_4x3_large",
						"imageStyle": "background"
					}
				]
			}
		}
	}
	`
	var recording tablometadata.Recording
	cleanJSON := strings.ReplaceAll(correctMovieJson, "\n", "")
	err := json.Unmarshal([]byte(correctMovieJson), &recording)
	if err != nil {
		t.Fatal(`Failed to unmarshal string`)
	}

	jsonData, err := json.Marshal(recording)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(string(jsonData[:]))
	if strings.Contains(string(jsonData[:]), cleanJSON) {
		fmt.Println("Was there")
	} else {
		fmt.Println("Not there")
	}

}

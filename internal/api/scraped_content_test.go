package api

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestMarshalScrapedSingleEmptyContentReturnsNil(t *testing.T) {
	if got, err := marshalScrapedScene(nil); err != nil || got != nil {
		t.Fatalf("scene empty content = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedPerformer(nil); err != nil || got != nil {
		t.Fatalf("performer empty content = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedGallery(nil); err != nil || got != nil {
		t.Fatalf("gallery empty content = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedImage(nil); err != nil || got != nil {
		t.Fatalf("image empty content = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedMovie(nil); err != nil || got != nil {
		t.Fatalf("movie empty content = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedGroup(nil); err != nil || got != nil {
		t.Fatalf("group empty content = (%#v, %v), want (nil, nil)", got, err)
	}
}

func TestMarshalScrapedSingleTypedNilContentReturnsNil(t *testing.T) {
	var scene *models.ScrapedScene
	var performer *models.ScrapedPerformer
	var gallery *models.ScrapedGallery
	var image *models.ScrapedImage
	var movie *models.ScrapedMovie
	var group *models.ScrapedGroup

	checks := []struct {
		name string
		got  any
		err  error
	}{
		{"scene", func() (any, error) { v, e := marshalScrapedScene(scene); return v, e }()},
		{"performer", func() (any, error) { v, e := marshalScrapedPerformer(performer); return v, e }()},
		{"gallery", func() (any, error) { v, e := marshalScrapedGallery(gallery); return v, e }()},
		{"image", func() (any, error) { v, e := marshalScrapedImage(image); return v, e }()},
		{"movie", func() (any, error) { v, e := marshalScrapedMovie(movie); return v, e }()},
		{"group", func() (any, error) { v, e := marshalScrapedGroup(group); return v, e }()},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			if check.err != nil || check.got != nil {
				t.Fatalf("typed nil content = (%#v, %v), want (nil, nil)", check.got, check.err)
			}
		})
	}
}

func TestMarshalScrapedSinglePreservesValidContent(t *testing.T) {
	scene := &models.ScrapedScene{Title: "scene"}
	got, err := marshalScrapedScene(scene)
	if err != nil || got != scene {
		t.Fatalf("scene = (%#v, %v), want original pointer", got, err)
	}
}

func TestMarshalScrapedSceneRejectsWrongContentType(t *testing.T) {
	if got, err := marshalScrapedScene(&models.ScrapedPerformer{}); err == nil || got != nil {
		t.Fatalf("wrong content = (%#v, %v), want conversion error", got, err)
	}
}

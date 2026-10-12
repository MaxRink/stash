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

	if got, err := marshalScrapedScene(scene); err != nil || got != nil {
		t.Fatalf("typed nil scene = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedPerformer(performer); err != nil || got != nil {
		t.Fatalf("typed nil performer = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedGallery(gallery); err != nil || got != nil {
		t.Fatalf("typed nil gallery = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedImage(image); err != nil || got != nil {
		t.Fatalf("typed nil image = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedMovie(movie); err != nil || got != nil {
		t.Fatalf("typed nil movie = (%#v, %v), want (nil, nil)", got, err)
	}
	if got, err := marshalScrapedGroup(group); err != nil || got != nil {
		t.Fatalf("typed nil group = (%#v, %v), want (nil, nil)", got, err)
	}
}

func TestMarshalScrapedSinglePreservesValidContent(t *testing.T) {
	scene := &models.ScrapedScene{}
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

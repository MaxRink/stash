package api

import (
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper"
)

func TestMarshalScrapedSingleNoResult(t *testing.T) {
	// scrapers return an untyped nil when nothing was found, and the script
	// scraper returns a typed nil pointer when the script prints null
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedPerformer)(nil)} {
		if got, err := marshalScrapedPerformer(content); err != nil || got != nil {
			t.Errorf("marshalScrapedPerformer(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedScene)(nil)} {
		if got, err := marshalScrapedScene(content); err != nil || got != nil {
			t.Errorf("marshalScrapedScene(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedGallery)(nil)} {
		if got, err := marshalScrapedGallery(content); err != nil || got != nil {
			t.Errorf("marshalScrapedGallery(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedImage)(nil)} {
		if got, err := marshalScrapedImage(content); err != nil || got != nil {
			t.Errorf("marshalScrapedImage(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedMovie)(nil)} {
		if got, err := marshalScrapedMovie(content); err != nil || got != nil {
			t.Errorf("marshalScrapedMovie(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
	// the script scraper returns a *models.ScrapedMovie for group content
	for _, content := range []scraper.ScrapedContent{nil, (*models.ScrapedGroup)(nil), (*models.ScrapedMovie)(nil)} {
		if got, err := marshalScrapedGroup(content); err != nil || got != nil {
			t.Errorf("marshalScrapedGroup(%#v) = (%v, %v), want (nil, nil)", content, got, err)
		}
	}
}

func TestMarshalScrapedScenesSkipsNil(t *testing.T) {
	scene := &models.ScrapedScene{}
	got, err := marshalScrapedScenes([]scraper.ScrapedContent{nil, (*models.ScrapedScene)(nil), scene})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != scene {
		t.Errorf("marshalScrapedScenes = %v, want only the non-nil scene", got)
	}
}

func TestMarshalScrapedSceneValidAndWrongType(t *testing.T) {
	scene := &models.ScrapedScene{}
	if got, err := marshalScrapedScene(scene); err != nil || got != scene {
		t.Errorf("marshalScrapedScene(scene) = (%v, %v), want original pointer", got, err)
	}
	if got, err := marshalScrapedScene(&models.ScrapedPerformer{}); err == nil || got != nil {
		t.Errorf("marshalScrapedScene(performer) = (%v, %v), want conversion error", got, err)
	}
}

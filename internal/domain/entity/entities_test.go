package entity

import (
	"testing"
)

func TestEntity_DeletedAtNull(t *testing.T) {
	blog := Blog{
		ID:                  1,
		Title:               "Test",
		TitleEng:            "Test",
		DescriptionShort:    "Short",
		Description:         "Long",
		DescriptionShortEng: "Short EN",
		DescriptionEng:      "Long EN",
	}

	if blog.DeletedAt != nil {
		t.Error("DeletedAt should be nil by default")
	}

	experience := Experience{
		ID:        1,
		YearStart: "2020-01-01",
		Company:   "Test Corp",
	}

	if experience.DeletedAt != nil {
		t.Error("Experience.DeletedAt should be nil by default")
	}

	customSection := CustomSection{
		ID:       1,
		Title:    "Test",
		TitleEng: "Test",
		Visible:  false,
	}

	if customSection.DeletedAt != nil {
		t.Error("CustomSection.DeletedAt should be nil by default")
	}
}

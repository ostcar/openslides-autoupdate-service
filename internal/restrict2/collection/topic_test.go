package collection_test

import (
	"testing"

	"github.com/OpenSlides/openslides-autoupdate-service/internal/restrict/perm"
	"github.com/OpenSlides/openslides-autoupdate-service/internal/restrict2/collection"
)

func TestTopicModeA(t *testing.T) {
	t.Parallel()
	f := collection.Topic{}.Modes("A")

	testCase(
		"no perm",
		t,
		f,
		false,
		`---
		topic/1:
			meeting_id: 2
			agenda_item_id: 7
		agenda_item/7/meeting_id: 30
		`,
	)

	testCase(
		"see agenda item",
		t,
		f,
		true,
		`---
		topic/1:
			meeting_id: 30
			agenda_item_id: 3

		agenda_item/3:
			is_internal: true
			meeting_id: 30
		`,
		withPerms(30, perm.AgendaItemCanSeeInternal),
	)

	testCase(
		"can not see agenda_item",
		t,
		f,
		false,
		`---
		topic/1:
			meeting_id: 30
			agenda_item_id: 3

		agenda_item/3:
			is_internal: true
			meeting_id: 30
		`,
		withPerms(30, perm.AgendaItemCanSee),
	)
}
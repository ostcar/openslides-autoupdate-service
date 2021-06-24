package slide

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenSlides/openslides-autoupdate-service/internal/projector"
	"github.com/OpenSlides/openslides-autoupdate-service/pkg/datastore"
)

type dbAgendaItem struct {
	ID              int    `json:"id"`
	ItemNumber      string `json:"item_number"`
	ContentObjectID string `json:"content_object_id"`
	MeetingID       int    `json:"meeting_id"`
	IsHidden        bool   `json:"is_hidden"`
	IsInternal      bool   `json:"is_internal"`
	Depth           int    `json:"level"`
}

func agendaItemFromMap(in map[string]json.RawMessage) (*dbAgendaItem, error) {
	bs, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encoding agenda item data: %w", err)
	}

	var ai dbAgendaItem
	if err := json.Unmarshal(bs, &ai); err != nil {
		return nil, fmt.Errorf("decoding agenda item data: %w", err)
	}
	return &ai, nil
}

type dbAgendaItemList struct {
	AgendaItemIds      []int `json:"agenda_item_ids"`
	AgendaShowInternal bool  `json:"agenda_show_internal_items_on_projector"`
}

func agendaItemListFromMap(in map[string]json.RawMessage) (*dbAgendaItemList, error) {
	bs, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encoding agenda item list data: %w", err)
	}

	var ail dbAgendaItemList
	if err := json.Unmarshal(bs, &ail); err != nil {
		return nil, fmt.Errorf("decoding agenda item list data: %w", err)
	}
	return &ail, nil
}

type outAgendaItem struct {
	TitleInformation json.RawMessage `json:"title_information"`
	Depth            int             `json:"depth"`
}

// AgendaItem renders the agenda_item slide.
func AgendaItem(store *projector.SlideStore) {
	store.RegisterSliderFunc("agenda_item", func(ctx context.Context, ds projector.Datastore, p7on *projector.Projection) (encoded []byte, keys []string, err error) {
		fetch := datastore.NewFetcher(ds)
		defer func() {
			if err == nil {
				err = fetch.Error()
			}
		}()

		data := fetch.Object(
			ctx,
			[]string{
				"id",
				"item_number",
				"content_object_id",
				"meeting_id",
				"is_hidden",
				"is_internal",
				"level",
			},
			p7on.ContentObjectID,
		)

		agendaItem, err := agendaItemFromMap(data)
		if err != nil {
			return nil, nil, fmt.Errorf("get agenda item: %w", err)
		}

		collection := strings.Split(agendaItem.ContentObjectID, "/")[0]
		titler := store.GetTitleInformationFunc(collection)
		if titler == nil {
			return nil, nil, fmt.Errorf("no titler function registered for %s", collection)
		}

		titleInfo, err := titler.GetTitleInformation(ctx, fetch, agendaItem.ContentObjectID, agendaItem.ItemNumber)
		if err != nil {
			return nil, nil, fmt.Errorf("get title func: %w", err)
		}

		out := outAgendaItem{
			TitleInformation: titleInfo,
			Depth:            agendaItem.Depth,
		}

		responseValue, err := json.Marshal(out)
		if err != nil {
			return nil, nil, fmt.Errorf("encoding response slide agenda item: %w", err)
		}
		return responseValue, fetch.Keys(), err
	})
}

// AgendaItemList renders the agenda_item_list slide.
func AgendaItemList(store *projector.SlideStore) {
	store.RegisterSliderFunc("agenda_item_list", func(ctx context.Context, ds projector.Datastore, p7on *projector.Projection) (encoded []byte, keys []string, err error) {
		fetch := datastore.NewFetcher(ds)
		defer func() {
			if err == nil {
				err = fetch.Error()
			}
		}()

		data := fetch.Object(
			ctx,
			[]string{
				"agenda_item_ids",
				"agenda_show_internal_items_on_projector",
			},
			p7on.ContentObjectID,
		)
		agendaItemList, err := agendaItemListFromMap(data)
		if err != nil {
			return nil, nil, fmt.Errorf("get agenda item list: %w", err)
		}

		var options struct {
			OnlyMainItems bool `json:"only_main_items"`
		}
		if p7on.Options != nil {
			if err := json.Unmarshal(p7on.Options, &options); err != nil {
				return nil, nil, fmt.Errorf("decoding projection options: %w", err)
			}
		}
		var allAgendaItems []outAgendaItem
		for _, aiID := range agendaItemList.AgendaItemIds {
			data = fetch.Object(
				ctx,
				[]string{
					"id",
					"item_number",
					"content_object_id",
					"meeting_id",
					"is_hidden",
					"is_internal",
					"level",
				},
				"agenda_item/%d",
				aiID,
			)
			agendaItem, err := agendaItemFromMap(data)
			if err != nil {
				return nil, nil, fmt.Errorf("get agenda item: %w", err)
			}

			if agendaItem.IsHidden || (agendaItem.IsInternal && !agendaItemList.AgendaShowInternal) {
				continue
			}

			if options.OnlyMainItems && agendaItem.Depth > 0 {
				continue
			}

			collection := strings.Split(agendaItem.ContentObjectID, "/")[0]
			titler := store.GetTitleInformationFunc(collection)
			if titler == nil {
				return nil, nil, fmt.Errorf("no titler function registered for %s", collection)
			}

			titleInfo, err := titler.GetTitleInformation(ctx, fetch, agendaItem.ContentObjectID, agendaItem.ItemNumber)
			if err != nil {
				return nil, nil, fmt.Errorf("get title func: %w", err)
			}

			allAgendaItems = append(
				allAgendaItems,
				outAgendaItem{
					TitleInformation: titleInfo,
					Depth:            agendaItem.Depth,
				},
			)
		}

		out := struct {
			Items []outAgendaItem `json:"items"`
		}{
			allAgendaItems,
		}

		responseValue, err := json.Marshal(out)
		if err != nil {
			return nil, nil, fmt.Errorf("encoding response for slide agenda item list: %w", err)
		}
		return responseValue, fetch.Keys(), err
	})
}
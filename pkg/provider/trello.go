package provider

import (
	"os"

	"github.com/Shikugawa/pingoo/pkg/task"
	"github.com/adlio/trello"
)

const ProviderName = "trello"

type TrelloProvider struct {
	client *trello.Client
}

func NewTrelloProvider() *TrelloProvider {
	return &TrelloProvider{
		client: trello.NewClient(os.Getenv("TRELLO_APP_ID"), os.Getenv("TRELLO_TOKEN")),
	}
}

func (c *TrelloProvider) Get() ([]task.Task, error) {
	board, err := c.client.GetBoard("jYfamf8A", trello.Defaults())
	if err != nil {
		return nil, err
	}

	lists, err := board.GetLists()
	if err != nil {
		return nil, err
	}

	var tasks []task.Task

	for _, l := range lists {
		cards, err := l.GetCards()
		if err != nil {
			continue
		}

		for _, card := range cards {
			labels := card.Labels

			if len(labels) < 1 {
				continue
			}

			priority := task.LabelToPriority(labels[0].Name)

			if priority == task.None {
				continue
			}

			tasks = append(tasks, task.Task{
				ProviderName: ProviderName,
				TaskID:       card.ID,
				Group:        l.Name,
				Title:        card.Name,
				Priority:     priority,
				Deadline:     card.Badges.Due,
				Url:          card.ShortURL,
			})
		}
	}

	return tasks, nil
}

func (c *TrelloProvider) OnDelete(id string) error {
	board, err := c.client.GetBoard(os.Getenv("TRELLO_BOARD_ID"), trello.Defaults())
	if err != nil {
		return err
	}

	lists, err := board.GetLists()
	if err != nil {
		return err
	}

	for _, l := range lists {
		cards, err := l.GetCards()
		if err != nil {
			continue
		}

		for _, card := range cards {
			if card.ID == id {
				if err := card.Delete(); err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

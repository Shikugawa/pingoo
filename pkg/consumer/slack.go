package consumer

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Shikugawa/pingoo/pkg/task"
	"github.com/slack-go/slack"
	"go.uber.org/multierr"
)

type SlackClient struct {
	client *slack.Client
}

func NewSlackClient() *SlackClient {
	return &SlackClient{
		client: slack.New(os.Getenv("SLACK_TOKEN")),
	}
}

func (s *SlackClient) Post(tasks []task.Task) (errs error) {
	attachments := make(map[string][]slack.Attachment)

	for _, task := range tasks {
		chanName, err := s.getTargetChannel(task)

		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		if chanName == "" {
			continue
		}

		if _, ok := attachments[chanName]; !ok {
			attachments[chanName] = []slack.Attachment{}
		}

		attachments[chanName] = append(attachments[chanName], slack.Attachment{
			Text:       task.Title,
			CallbackID: task.TaskID,
			Color:      s.priorityToColor(task),
			Actions: []slack.AttachmentAction{
				{
					Type: slack.ActionType(slack.METButton),
					Name: "Done",
					Text: "Done",
				},
			},
		})
	}

	for chanName, atts := range attachments {
		sort.Slice(atts, func(i, j int) bool {
			leftHex, _ := strconv.ParseInt(strings.TrimPrefix(atts[i].Color, "#"), 16, 32)
			rightHex, _ := strconv.ParseInt(strings.TrimPrefix(atts[j].Color, "#"), 16, 32)

			return leftHex < rightHex
		})

		_, _, err := s.client.PostMessage(
			"#"+chanName,
			slack.MsgOptionText("tasks!!", false),
			slack.MsgOptionAttachments(atts...),
		)

		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
	}

	return errs
}

func (s *SlackClient) getTargetChannel(task task.Task) (string, error) {
	if !strings.HasPrefix(task.Group, "TODO - ") {
		return "", nil
	}
	seps := strings.Split(task.Group, " - ")
	requiredChanName := seps[1] + "-" + task.ProviderName

	channels, _, err := s.client.GetConversations(&slack.GetConversationsParameters{})
	if err != nil {
		return "", err
	}

	for _, ch := range channels {
		chanName := ch.Name
		if chanName == requiredChanName {
			return requiredChanName, nil
		}
	}

	return "", nil
}

func (s *SlackClient) priorityToColor(t task.Task) string {
	if t.Priority == task.High {
		return "#B00B13"
	} else if t.Priority == task.Mid {
		return "#FFFF00"
	} else if t.Priority == task.Low {
		return "#00FFFF"
	}
	return "#FFFFFF"
}

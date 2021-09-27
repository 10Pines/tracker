package reporter

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"

	"github.com/10Pines/tracker/v2/internal/report"
)

type SlackReporter struct {
	api     *slack.Client
	channel string
}

func NewSlackReporter(token, channel string) SlackReporter {
	api := slack.New(token)
	return SlackReporter{
		api:     api,
		channel: channel,
	}
}

func (s SlackReporter) SendReport(report report.Report) error {
	blocks := []slack.Block{
		slack.NewHeaderBlock(slack.NewTextBlockObject(slack.PlainTextType, ":newspaper:  Daily Report  :newspaper:", false, false)),
		slack.NewContextBlock("date",
			slack.NewTextBlockObject(slack.MarkdownType, with(report.Timestamp()), false, false),
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock("xxx",
			slack.NewTextBlockObject(slack.MarkdownType, ":calendar: |   *PROCESSED TASKS*  | :calendar:", false, false),
		),
	}
	blocks = append(blocks, taskSections(report)...)
	blocks = append(blocks,
		slack.NewDividerBlock(),
		slack.NewContextBlock("footer",
			slack.NewTextBlockObject(slack.MarkdownType, ":pushpin: Do you have something to include in the newsletter? Here's *how to submit content*.", false, false),
		),
	)

	_, _, _, err := s.api.SendMessage(s.channel, slack.MsgOptionBlocks(blocks...))
	return err
}

func taskSections(report report.Report) []slack.Block {
	var blocks []slack.Block
	for _, taskStatus := range report.Statuses() {
		msg := fmt.Sprintf("`11/20-11/22` *%s*", taskStatus.Task.Name)
		taskBlock := slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, msg, false, false), nil, nil)
		blocks = append(blocks, taskBlock)
	}
	return blocks
}

func with(timestamp time.Time) string {
	return fmt.Sprintf("*%s*  |  Sales Team Announcements", timestamp.Format(time.ANSIC))
}

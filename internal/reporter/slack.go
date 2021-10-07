package reporter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/slack-go/slack"

	"github.com/10Pines/tracker/v2/internal/report"
)

const (
	shortISO = "2006-01-02 15:04"
	danger   = "#fa352a"
	rickroll = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

type slackReporter struct {
	api     *slack.Client
	channel string
}

// NewSlackReporter create a Reporter instance that sends messages via Slack
func NewSlackReporter(token, channel string) Reporter {
	api := slack.New(token)
	return slackReporter{
		api:     api,
		channel: channel,
	}
}

func (s slackReporter) Name() string {
	return "slack"
}

func (s slackReporter) Process(report report.Report) error {
	var blocks []slack.Block
	if h := header(report); h != nil {
		blocks = append(blocks, h)
	}
	if f := footer(report); f != nil {
		blocks = append(blocks, f)
	}
	failedBackups := failedTasksAttachments(report)
	content := slack.MsgOptionBlocks(blocks...)
	_, _, _, err := s.api.SendMessage(s.channel, content, slack.MsgOptionAttachments(failedBackups...))
	return err
}

func header(r report.Report) slack.Block {
	if r.IsOK() {
		msg := slack.NewTextBlockObject(slack.MarkdownType, "- _Pío-fui-pío, todo bien por aquí._", false, false)
		return slack.NewSectionBlock(msg, nil, nil, slack.SectionBlockOptionBlockID("header"))
	}
	return nil
}

func failedTasksAttachments(report report.Report) []slack.Attachment {
	var failedTasks []slack.Attachment
	for i, taskStatus := range report.Statuses() {
		if !taskStatus.IsOK() {
			failedCount := taskStatus.Expected - taskStatus.BackupCount
			msg := fmt.Sprintf(" *%s*\n Falló %d veces en los últimos %d reportes. Tolerancia %d", taskStatus.Task.Name, failedCount, taskStatus.Task.Datapoints, taskStatus.Task.Tolerance)
			viewMoreBtn := slack.NewButtonBlockElement(fmt.Sprintf("view-more-%d", i), "asd", slack.NewTextBlockObject(slack.PlainTextType, "Mas info", false, false))
			viewMoreBtn.URL = rickroll
			viewMoreBtn.WithStyle(slack.StylePrimary)
			moreInfo := slack.NewAccessory(viewMoreBtn)
			taskBlock := slack.NewSectionBlock(slack.NewTextBlockObject(slack.MarkdownType, msg, false, false), nil, moreInfo, slack.SectionBlockOptionBlockID(strconv.Itoa(i)))
			failedTask := slack.Attachment{
				ID:     i,
				Color:  danger,
				Blocks: slack.Blocks{BlockSet: []slack.Block{taskBlock}},
			}
			failedTasks = append(failedTasks, failedTask)
		}
	}
	return failedTasks
}

func footer(r report.Report) slack.Block {
	var msg strings.Builder
	msg.WriteString("Tareas reportadas\n")
	for _, taskStatus := range r.Statuses() {
		ts := taskStatus.Task.CreatedAt.Format(shortISO)
		msg.WriteString(ts)
		msg.WriteString(" ")
		msg.WriteString(fmt.Sprintf("*%s*\n", taskStatus.Task.Name))
	}
	tasks := slack.NewTextBlockObject(slack.MarkdownType, msg.String(), false, false)
	return slack.NewSectionBlock(tasks, nil, nil, slack.SectionBlockOptionBlockID("footer"))
}

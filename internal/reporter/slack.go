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
	danger   = "#FF5248"
	rickroll = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

type statusEmoji string

const (
	ok      statusEmoji = ":zorzal-tarea-ok:"
	nok                 = ":zorzal-tarea-error:"
	pending             = ":zorzal-tarea-noactiva:"
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
			msg := fmt.Sprintf(" *%s*\n Falló %d veces en los últimos %d reportes.", taskStatus.Task.Name, failedCount, taskStatus.Task.Datapoints)
			viewMoreBtn := slack.NewButtonBlockElement(fmt.Sprintf("view-more-%d", i), "", slack.NewTextBlockObject(slack.PlainTextType, "Mas info", false, false))
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
	msg.WriteString("Tareas observadas:\n")
	for _, taskStatus := range r.Statuses() {
		if taskStatus.TaskHasBackups() {
			em := taskEmoji(taskStatus)
			ts := taskStatus.LastBackup.Format(shortISO)
			msg.WriteString(fmt.Sprintf("%s %s *%s*\n", em, ts, taskStatus.Task.Name))
		} else {
			msg.WriteString(fmt.Sprintf("%s Sin datos: *%s*\n", pending, taskStatus.Task.Name))
		}
	}
	tasks := slack.NewTextBlockObject(slack.MarkdownType, msg.String(), false, false)
	return slack.NewSectionBlock(tasks, nil, nil, slack.SectionBlockOptionBlockID("footer"))
}

func taskEmoji(taskStatus report.TaskStatus) statusEmoji {
	if taskStatus.IsOK() {
		return ok
	}
	return nok
}

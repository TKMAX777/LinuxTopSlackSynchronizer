package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TKMAX777/LinuxTopSlackSynchronizer/slack_webhook"
	"github.com/TKMAX777/LinuxTopSlackSynchronizer/top"
	"github.com/pkg/errors"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

type SlackHandler struct {
	token string
	hook  *slack_webhook.Handler

	hostname      string
	processNumber int

	regularMessageText string

	lastTS string
}

func NewSlackHandler(token string, processNum int) *SlackHandler {
	name, _ := os.Hostname()
	var handler = &SlackHandler{
		token:         token,
		processNumber: processNum,
		hook:          slack_webhook.New(token),
		hostname:      strings.ToUpper(name),
	}

	handler.regularMessageText = fmt.Sprintf("ProcessRegularMessage,%s", handler.hook.Identity.User)

	return handler
}

func (h *SlackHandler) RemoveLastRegularMessage(channelID string) error {
	messages, err := h.hook.GetMessages(channelID, "", 100)
	if err != nil {
		return errors.Wrap(err, "GetMessages")
	}

	for _, m := range messages {
		if m.Text == h.regularMessageText {
			h.hook.Remove(channelID, m.TS)
			break
		}
	}
	return nil
}

func (h *SlackHandler) RegularSend(ps []top.Process, channelID string) error {
	var blocks = h.buildBlocks(ps)

	if h.lastTS != "" {
		ts, err := h.hook.Update(slack_webhook.Message{
			Username: fmt.Sprintf("[%s] Process", h.hostname),
			TS:       h.lastTS,
			Text:     h.regularMessageText,
			Blocks:   blocks,
			Channel:  channelID,
		})

		h.lastTS = ts
		return err
	}

	ts, err := h.Send(ps, channelID)
	h.lastTS = ts

	return err
}

func (h SlackHandler) Send(ps []top.Process, channelID string) (string, error) {
	return h.hook.Send(slack_webhook.Message{
		Username: fmt.Sprintf("[%s] Process", h.hostname),
		Text:     h.regularMessageText,
		Blocks:   h.buildBlocks(ps),
		Channel:  channelID,
	})
}

func (h SlackHandler) buildBlocks(ps []top.Process) []slack_webhook.BlockBase {
	var blocks = []slack_webhook.BlockBase{
		slack_webhook.ContextBlock(
			slack_webhook.MrkdwnElement("CPU", true),
			slack_webhook.MrkdwnElement("Memory", true),
			slack_webhook.MrkdwnElement("User", true),
			slack_webhook.MrkdwnElement("Command", true),
			slack_webhook.MrkdwnElement("CPU", true),
		),
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"CPU", "MEMORY", "USER", "COMMAND", "CPU BAR"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:  "CPU BAR",
			Align: text.AlignLeft,
		},
		{
			Name:     "COMMAND",
			Align:    text.AlignLeft,
			WidthMax: 30,
		},
		{
			Name:  "USER",
			Align: text.AlignLeft,
		},
	})
	t.SetAutoIndex(true)

	for i, p := range ps {
		if i >= h.processNumber {
			break
		}

		var CPUbar string

		for j := 0; j < int(p.CPU/10); j++ {
			CPUbar += "|"
		}

		t.AppendRow(table.Row{p.CPU, p.Memory, p.User, p.Command, CPUbar})
		fmt.Println(p.Command)
	}

	t.SetStyle(table.StyleLight)

	var section = slack_webhook.SectionBlock()
	section.Text = slack_webhook.MrkdwnElement("```\n"+t.Render()+"\n```", false)
	blocks = []slack_webhook.BlockBase{
		section,
	}

	return blocks
}

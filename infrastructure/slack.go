package infrastructure

import (
	"context"
	"net/url"

	"github.com/slack-go/slack"
)

type notificationLevel int

const (
	notificationLevelInfo = iota
	notificationLevelError
)

func (l notificationLevel) color() string {
	switch l {
	case notificationLevelError:
		return "danger"
	case notificationLevelInfo:
		fallthrough
	default:
		return "good"
	}
}

type SlackConfig struct {
	Channel  string
	UserName *string
	IconURL  *url.URL
}

type SlackNotifier struct {
	config SlackConfig
	client *slack.Client
}

func NewSlackNotifier(token string, config SlackConfig) *SlackNotifier {
	return &SlackNotifier{
		config: config,
		client: slack.New(token),
	}
}

func (n *SlackNotifier) Info(ctx context.Context, title string, txt string) error {
	return n.notify(ctx, notificationLevelInfo, title, txt)
}

func (n *SlackNotifier) Error(ctx context.Context, title string, txt string) error {
	return n.notify(ctx, notificationLevelError, title, txt)
}

func (n *SlackNotifier) notify(ctx context.Context, level notificationLevel, title string, txt string) error {
	var msgOpts []slack.MsgOption
	if n.config.UserName != nil {
		msgOpts = append(msgOpts, slack.MsgOptionUsername(*n.config.UserName))
	}
	if n.config.IconURL != nil {
		msgOpts = append(msgOpts, slack.MsgOptionIconURL(n.config.IconURL.String()))
	}
	msgOpts = append(msgOpts, slack.MsgOptionAttachments(
		slack.Attachment{
			Title: title,
			Color: level.color(),
			Text:  txt,
		},
	))

	_, _, err := n.client.PostMessageContext(
		ctx, n.config.Channel,
		msgOpts...,
	)
	return err
}

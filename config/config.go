package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	toml "github.com/pelletier/go-toml"
)

var (
	userHomeDir = os.UserHomeDir
)

type Config struct {
	SlackURL       string
	Token          string
	PrimaryChannel string
	Channel        string
	SnippetChannel string
	Username       string
	IconEmoji      string
	Duration       time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadEnv() error {
	if c.SlackURL == "" {
		c.SlackURL = os.Getenv("NOTIFY_SLACK_WEBHOOK_URL")
	}

	if c.Token == "" {
		c.Token = os.Getenv("NOTIFY_SLACK_TOKEN")
	}

	if c.Channel == "" {
		c.Channel = os.Getenv("NOTIFY_SLACK_CHANNEL")
	}

	if c.SnippetChannel == "" {
		c.SnippetChannel = os.Getenv("NOTIFY_SLACK_SNIPPET_CHANNEL")
	}

	if c.Username == "" {
		c.Username = os.Getenv("NOTIFY_SLACK_USERNAME")
	}

	if c.IconEmoji == "" {
		c.IconEmoji = os.Getenv("NOTIFY_SLACK_ICON_EMOJI")
	}

	durationStr := os.Getenv("NOTIFY_SLACK_INTERVAL")
	if durationStr != "" {
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("incorrect value to inteval option from NOTIFY_SLACK_INTERVAL: %s: %w", durationStr, err)
		}
		c.Duration = duration
	}

	return nil
}

func (c *Config) LoadTOML(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	config, err := toml.LoadBytes(b)
	if err != nil {
		return err
	}

	slackConfig := config.Get("slack").(*toml.Tree)

	if c.SlackURL == "" {
		slackURL, ok := slackConfig.Get("url").(string)
		if ok {
			c.SlackURL = slackURL
		}
	}
	if c.Token == "" {
		token, ok := slackConfig.Get("token").(string)
		if ok {
			c.Token = token
		}
	}
	if c.Channel == "" {
		channel, ok := slackConfig.Get("channel").(string)
		if ok {
			c.Channel = channel
		}
	}
	if c.SnippetChannel == "" {
		snippetChannel, ok := slackConfig.Get("snippet_channel").(string)
		if ok {
			c.SnippetChannel = snippetChannel
		}
	}
	if c.Username == "" {
		username, ok := slackConfig.Get("username").(string)
		if ok {
			c.Username = username
		}
	}
	if c.IconEmoji == "" {
		iconEmoji, ok := slackConfig.Get("icon_emoji").(string)
		if ok {
			c.IconEmoji = iconEmoji
		}
	}

	durationStr, ok := slackConfig.Get("interval").(string)
	if ok {
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("incorrect value to inteval option: %s: %w", durationStr, err)
		}
		c.Duration = duration
	}

	return nil
}

func LoadTOMLFilename(filename string) string {
	if filename != "" {
		return filename
	}

	homeDir, err := userHomeDir()
	if err == nil {
		tomlFile := filepath.Join(homeDir, ".notify_slack.toml")
		if fileExists(tomlFile) {
			return tomlFile
		}

		tomlFile = filepath.Join(homeDir, "/etc/notify_slack.toml")
		if fileExists(tomlFile) {
			return tomlFile
		}
	}

	tomlFile := "/etc/notify_slack.toml"
	if fileExists(tomlFile) {
		return tomlFile
	}

	return ""
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}

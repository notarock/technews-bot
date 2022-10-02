package articles

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	HN_PROFILE_URL = "https://news.ycombinator.com/user?id=%s"
	LB_PROFILE_URL = "https://lobste.rs/u/%s"
)

type Article struct {
	ID     string
	Title  string
	Link   string
	Tags   []string
	Author string
	Source string
}

func LinkToID(link string) string {
	hmd5 := md5.Sum([]byte(link))
	return fmt.Sprintf("%x", hmd5)
}

func authorEmbedFromSource(source, author string) *discordgo.MessageEmbedAuthor {
	var profileUrl string
	if source == "HACKERNEWS" {
		profileUrl = fmt.Sprintf(HN_PROFILE_URL, author)
	} else if source == "LOBSTERS" {
		profileUrl = fmt.Sprintf(LB_PROFILE_URL, author)
	} else {
		profileUrl = ""
	}

	var iconUrl string
	if source == "HACKERNEWS" {
		iconUrl = "https://news.ycombinator.com/y18.gif"
	} else if source == "LOBSTERS" {
		iconUrl = "https://upload.wikimedia.org/wikipedia/commons/e/e7/Lobsters_logo.png"
	} else {
		iconUrl = ""
	}

	embed := discordgo.MessageEmbedAuthor{
		URL:          profileUrl,
		Name:         author,
		IconURL:      iconUrl,
		ProxyIconURL: "",
	}
	return &embed
}

func (a Article) ToDiscordEmbed() *discordgo.MessageEmbed {
	authorEmbed := authorEmbedFromSource(a.Source, a.Author)

	embed := discordgo.MessageEmbed{
		URL:         a.Link,
		Type:        "link",
		Title:       a.Title,
		Description: a.Description(),
		Timestamp:   "",
		Color:       0,
		Author:      authorEmbed,
		Fields:      []*discordgo.MessageEmbedField{},
	}

	return &embed
}

func (a Article) RelatesTo(subject string) bool {
	return a.titleContains(subject) || a.tagsContains(subject)
}

func (a Article) titleContains(subject string) bool {
	return strings.Contains(strings.ToLower(a.Title), strings.ToLower(subject))
}

func (a Article) tagsContains(subject string) bool {
	for _, tag := range a.Tags {
		if strings.Contains(strings.ToLower(tag), strings.ToLower(subject)) {
			return true
		}
	}
	return false
}

func (a Article) Description() string {
	desc := fmt.Sprintf("Article posted to %s by user %s.", a.Source, a.Author)
	if len(a.Tags) > 0 {
		desc += fmt.Sprintf("\nArticle was tagged using these tags: %s", strings.Join(a.Tags, ", "))
	}

	return desc
}

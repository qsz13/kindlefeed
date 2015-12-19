package main

type Config struct {
	Database Database
	FeedURL  []string
	Cron     string
}

type Database struct {
	Address  string
	Username string
	Password string
	Protocol string
	DBname   string
}

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	LastBuildDate string `xml:"lastBuildDate"`
	Language      string `xml:"language"`
	Image         Image  `xml:"image"`
	Items         []Item `xml:"item"`
}

type Image struct {
	Url   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	Creator     string `xml:"creator"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Content     string `xml:"encoded"`
}

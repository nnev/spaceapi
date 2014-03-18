package main

import ()

type Tristate int

const (
	Undefined Tristate = iota
	True
	False
)

type Endpoint struct {
	Api                 string   `json:"api"`
	Space               string   `json:"space"`
	Logo                string   `json:"logo"`
	Url                 string   `json:"url"`
	Location            Location `json:"location"`
	State               State    `json:"state"`
	Contact             Contact  `json:"contact"`
	IssueReportChannels []string `json:"issue_report_channels"`
	Feeds               Feeds
}

type Location struct {
	Address string
	Lat     float32 `json:"lat"`
	Lon     float32 `json:"lon"`
}
type State struct {
	Open Tristate `json:"open"`
}

type Contact struct {
	Irc       string `json:"irc"`
	List      string `json:"ml"`
	IssueMail string `json:"issue_mail"`
}

type Feeds struct {
	Calendar Feed `json:"calendar"`
}

type Feed struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

func NewEndpoint() *Endpoint {
	ep := &Endpoint{}

	ep.Api = "0.13"
	ep.Space = "Chaostreff Heidelberg"
	ep.Logo = "https://www.noname-ev.de//img/noname.svg"
	ep.Url = "https://www.noname-ev.de/"

	ep.Location = Location{
		Address: "Im Neuenheimer Feld 368, 69120 Heidelberg",
		Lat:     49.41759,
		Lon:     8.66834,
	}

	ep.Contact.Irc = "ircs://irc.twice-irc.de/chaos-hd"
	ep.Contact.List = "ccchd@ccchd.de"
	ep.Contact.IssueMail = "mero@merovius.de"

	ep.IssueReportChannels = []string{"issue_mail"}

	ep.Feeds = Feeds{Calendar: Feed{
		Type: "ical",
		Url:  "https://www.noname-ev.de/c14h.ics",
	}}

	return ep
}

func (s Tristate) MarshalJSON() ([]byte, error) {
	switch s {
	case Undefined:
		return []byte("null"), nil
	case True:
		return []byte("true"), nil
	case False:
		return []byte("false"), nil
	default:
		panic("Undefined value for tristate")
	}
}

func (s Tristate) String() string {
	b, _ := s.MarshalJSON()
	return string(b)
}

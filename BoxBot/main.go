package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"time"
	"strconv"
	"regexp"
	//"github.com/valyala/fastjson"
	//"encoding/json"
)

var (
	Token  string
	prefix string
)

/*
type GlobalStats struct {
	Success string
	Data    GlobalStats_Data
}
type GlobalStats_Data struct {
	Sessions int
	Vpn      int
	Machines int
	Latency  string
}
*/

type Embed struct {
	*discordgo.MessageEmbed
}

// Constants for message embed character limits
const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
	EmbedLimitFieldName   = 256
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
	EmbedLimit            = 4000
)

//NewEmbed returns a new embed object
func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

//SetTitle ...
func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

//SetDescription [desc]
func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
	return e
}

//AddField [name] [value]
func (e *Embed) AddField(name, value string) *Embed {
	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(name) > 1024 {
		name = name[:1024]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:  name,
		Value: value,
	})

	return e

}

//SetFooter [Text] [iconURL]
func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	return e
}

//SetImage ...
func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

//SetThumbnail ...
func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

//SetAuthor ...
func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

//SetURL ...
func (e *Embed) SetURL(URL string) *Embed {
	e.URL = URL
	return e
}

//SetColor ...
func (e *Embed) SetColor(clr int) *Embed {
	e.Color = clr
	return e
}

// InlineAllFields sets all fields in the embed to be inline
func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}

// Truncate truncates any embed value over the character limit.
func (e *Embed) Truncate() *Embed {
	e.TruncateDescription()
	e.TruncateFields()
	e.TruncateFooter()
	e.TruncateTitle()
	return e
}

// TruncateFields truncates fields that are too long
func (e *Embed) TruncateFields() *Embed {
	if len(e.Fields) > 25 {
		e.Fields = e.Fields[:EmbedLimitField]
	}

	for _, v := range e.Fields {

		if len(v.Name) > EmbedLimitFieldName {
			v.Name = v.Name[:EmbedLimitFieldName]
		}

		if len(v.Value) > EmbedLimitFieldValue {
			v.Value = v.Value[:EmbedLimitFieldValue]
		}

	}
	return e
}

// TruncateDescription ...
func (e *Embed) TruncateDescription() *Embed {
	if len(e.Description) > EmbedLimitDescription {
		e.Description = e.Description[:EmbedLimitDescription]
	}
	return e
}

// TruncateTitle ...
func (e *Embed) TruncateTitle() *Embed {
	if len(e.Title) > EmbedLimitTitle {
		e.Title = e.Title[:EmbedLimitTitle]
	}
	return e
}

// TruncateFooter ...
func (e *Embed) TruncateFooter() *Embed {
	if e.Footer != nil && len(e.Footer.Text) > EmbedLimitFooter {
		e.Footer.Text = e.Footer.Text[:EmbedLimitFooter]
	}
	return e
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	prefix = "~#"
}

func main() {
	//Create Session
	dg, err :=  discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Ctrl-C to terminate.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == (prefix + "status.get()") {
		v := url.Values{}
		ss := v.Encode()
		resp, err := http.Post("https://www.hackthebox.eu/api/stats/global", "application/json", strings.NewReader(ss))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ioutil.ReadAll() error: %v\n", err)
			return
		}
		//fmt.Println(string(data))
		mach := gjson.Get(string(data), "data.machines")
		vpn := gjson.Get(string(data), "data.vpn")
		sessions := gjson.Get(string(data), "data.sessions")
		latency := gjson.Get(string(data), "data.latency")
		//s.ChannelMessageSend(m.ChannelID, value.String())
		currentTime := time.Now()
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color: 0x00ff00,
			Description: "HackTheBox Status",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "Total Machines:",
					Value: mach.String(),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name: "Total VPN Connections:",
					Value: vpn.String(),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name: "Total Active Users:",
					Value: sessions.String(),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name: "Average Response Time:",
					Value: latency.String()+"ms",
					Inline: true,
				},
			},
			/*Image: &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
			},*/
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://cdn.discordapp.com/avatars/608327221723267072/667b7cd8ea038286f86163dd978112c3.webp?size=128",
			},
			Timestamp: currentTime.Format(time.RFC3339),
			Title: "~#status.get()",
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		
	}

	if m.Content == (prefix + "top5.get(img)") {

		//v := url.Values{}
		//ss := v.Encode()
		resp, err := http.Get("https://www.hackthebox.eu/api/charts/users/scores/?api_token=")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ioutil.ReadAll() error: %v\n", err)
			return
		}

		//vpn := gjson.Get(string(data), "data.vpn")


		var names[5]string
		var points[5]string
		for i := 0; i < 5; i++ {
			
			stri := strconv.Itoa(i)
			fmt.Println(string("chartData."+stri+".name"))
			tempstr := gjson.Get(string(data), string("chartData."+stri+".name"))
			names[i] = tempstr.String()
			fmt.Println(names[i])
			fmt.Println(string("chartData."+stri+".timeline."+stri+".y"))
			tempint := gjson.Get(string(data), string("chartData."+stri+".timeline."+stri+".y"))
			points[i] = strconv.FormatInt(int64(tempint.Int()), 10)
			fmt.Println(points[i])

			v := url.Values{}
			v.Set("api_token", "")
			v.Add("username", string(names[i]))
			ss := v.Encode()
			//fmt.Printf(ss)
			resp2, err2 := http.Post("https://www.hackthebox.eu/api/user/id", "application/x-www-form-urlencoded", strings.NewReader(ss))
			if err2 != nil {
				log.Fatalln(err2)
			}
			defer resp2.Body.Close()
			
			
			data2, err2 := ioutil.ReadAll(resp2.Body)
			if err2 != nil {
				fmt.Printf("ioutil.ReadAll() error: %v\n", err2)
				return
			}

			tempUserID := gjson.Get(string(data2), "id")
			embed := NewEmbed().
			SetImage("https://www.hackthebox.eu/badge/image/"+tempUserID.String()).
			SetColor(0x00ff00).MessageEmbed
			s.ChannelMessageSendEmbed(m.ChannelID, embed)	
			
		}

	}

	if m.Content == (prefix + "top5.get()") {

		//v := url.Values{}
		//ss := v.Encode()
		resp, err := http.Get("https://www.hackthebox.eu/api/charts/users/scores/?api_token=S7vfMOiKiKza8TyBQKI3U3HnjyQjbtvrtTqbPnznfhwbuCZ2hoG8Deq7lHcj")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ioutil.ReadAll() error: %v\n", err)
			return
		}

		//vpn := gjson.Get(string(data), "data.vpn")


		var names[5]string
		var points[5]string
		for i := 0; i < 5; i++ {
			
			stri := strconv.Itoa(i)
			//fmt.Println(string("chartData."+stri+".name"))
			tempstr := gjson.Get(string(data), string("chartData."+stri+".name"))
			names[i] = tempstr.String()
			//fmt.Println(names[i])
			//fmt.Println(string("chartData."+stri+".timeline."+stri+".y"))
			tempint := gjson.Get(string(data), string("chartData."+stri+".timeline."+stri+".y"))
			points[i] = strconv.FormatInt(int64(tempint.Int()), 10)
			//fmt.Println(points[i])
		}
		embed := NewEmbed().
		SetTitle("~#top5.get()").
		SetDescription("HackTheBox Top 5").
		SetThumbnail("https://cdn.discordapp.com/avatars/608327221723267072/667b7cd8ea038286f86163dd978112c3.webp?size=128").
		AddField(string(names[0]), string(points[0])).
		AddField(string(names[1]), string(points[1])).
		AddField(string(names[2]), string(points[2])).
		AddField(string(names[3]), string(points[3])).
		AddField(string(names[4]), string(points[4])).
		SetColor(0x00ff00).MessageEmbed

		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	}

	if strings.HasPrefix(m.Content, (prefix + "badge.get")) == true {
		rgx := regexp.MustCompile(`\((.*?)\)`)
		tempUser := rgx.FindStringSubmatch(m.Content)
		fmt.Printf(tempUser[1])

		v := url.Values{}
		v.Set("api_token", "")
		v.Add("username", string(tempUser[1]))
		ss := v.Encode()
		//fmt.Printf(ss)
		resp2, err2 := http.Post("https://www.hackthebox.eu/api/user/id", "application/x-www-form-urlencoded", strings.NewReader(ss))
		if err2 != nil {
			log.Fatalln(err2)
		}
		defer resp2.Body.Close()
		
		data2, err2 := ioutil.ReadAll(resp2.Body)
		if err2 != nil {
			fmt.Printf("ioutil.ReadAll() error: %v\n", err2)
			return
		}
		tempUserID := gjson.Get(string(data2), "id")
		embed := NewEmbed().
		SetImage("https://www.hackthebox.eu/badge/image/"+tempUserID.String()).
		SetColor(0x00ff00).MessageEmbed
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	}

}

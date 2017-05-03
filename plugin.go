package main

import (
	"crypto/tls"

	"github.com/Sirupsen/logrus"
	"github.com/rendom/ircFormat"
	"github.com/thoj/go-ircevent"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag        string
		Event      string
		Number     int
		Commit     string
		Ref        string
		Branch     string
		Author     string
		Message    string
		Status     string
		Link       string
		CommitLink string
		Started    int64
		Created    int64
	}

	Config struct {
		Server  string
		Channel string
		Nick    string
		Tls     bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {
	con := irc.IRC(p.Config.Nick, "Drone IRC")
	con.UseTLS = p.Config.Tls
	con.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	con.Debug = false

	logrus.WithFields(logrus.Fields{
		"server":  p.Config.Server,
		"nick":    p.Config.Nick,
		"tls":     p.Config.Tls,
		"channel": p.Config.Channel,
	}).Info("Connection to irc")

	err := con.Connect(p.Config.Server)

	if err != nil {
		logrus.Fatal(err)
	}

	go func() {
		if err := <-con.ErrorChan(); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	con.AddCallback("001", func(event *irc.Event) {
		con.Join(p.Config.Channel)
		con.Privmsgf(p.Config.Channel,
			"*%s* <%s/%s#%s> %s (%s) by %s\n%s",
			status(p.Build),
			p.Repo.Owner,
			p.Repo.Name,
			p.Build.Commit[:8],
			p.Build.Message,
			p.Build.Branch,
			p.Build.Author,
			p.Repo.CommitLink,
		)
		con.Quit()
	})
	con.Loop()
	return nil
}

func status(build Build) string {
	var fg int
	bg := ircFormat.None

	switch build.Status {
	case "success":
		fg = ircFormat.Green
	case "failure", "error", "killed":
		fg = ircFormat.Red
	default:
		fg = ircFormat.Yellow
	}

	return ircFormat.New(build.Status).
		SetBold().
		SetFg(fg).
		SetBg(bg).
		String()
}

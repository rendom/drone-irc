package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	CHAN   string
	SERVER string
	NICK   string
	TLS    bool
	DEBUG  bool
)

func main() {
	app := cli.NewApp()
	app.Name = "irc plugin"
	app.Usage = "irc plugin"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Usage:  "irc server",
			EnvVar: "PLUGIN_SERVER,IRC_SERVER",
		},
		cli.StringFlag{
			Name:   "channel",
			Usage:  "irc channel",
			EnvVar: "PLUGIN_CHANNEL,IRC_CHANNEL",
		},
		cli.StringFlag{
			Name:   "nick",
			Usage:  "IRC_NICK",
			EnvVar: "PLUGIN_NICK,IRC_NICK",
		},
		cli.BoolFlag{
			Name:   "tls",
			Usage:  "Use tls",
			EnvVar: "PLUGIN_TLS,IRC_TLS",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:        c.String("build.tag"),
			Number:     c.Int("build.number"),
			Event:      c.String("build.event"),
			Status:     c.String("build.status"),
			Commit:     c.String("commit.sha"),
			Ref:        c.String("commit.ref"),
			Branch:     c.String("commit.branch"),
			Author:     c.String("commit.author"),
			Message:    c.String("commit.message"),
			Link:       c.String("build.link"),
			Started:    c.Int64("build.started"),
			Created:    c.Int64("build.created"),
			CommitLink: c.String("commit.link"),
		},
		Config: Config{
			Server:  c.String("server"),
			Channel: c.String("channel"),
			Nick:    c.String("nick"),
			Tls:     c.Bool("tls"),
		},
	}

	return plugin.Exec()
}

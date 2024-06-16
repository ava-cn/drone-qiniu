package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	c := cli.NewApp()
	c.Name = "qiniu drone plugin"
	c.Usage = ""
	c.Version = ""

	c.Action = func(ctx *cli.Context) (err error) {
		p := Plugin{
			AccessKey: ctx.String("access-key"),
			SercetKey: ctx.String("secret-key"),
			Bucket:    ctx.String("bucket"),
			Zone:      ctx.String("zone"),
			Prefix:    ctx.String("prefix"),
			Dir:       ctx.String("dir"),
			Parallel:  ctx.Int("parallel"),
		}

		return p.Exec()
	}

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "access-key",
			Usage:   "qiniu access key",
			EnvVars: []string{"ACCESS_KEY", "PLUGIN_ACCESS_KEY"},
		},
		&cli.StringFlag{
			Name:    "secret-key",
			Usage:   "qiniu secret key",
			EnvVars: []string{"SECRET_KEY", "PLUGIN_SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "zone",
			Usage:   "bucket zone",
			EnvVars: []string{"ZONE", "PLUGIN_ZONE"},
		},
		&cli.StringFlag{
			Name:    "bucket",
			Usage:   "bucket name",
			EnvVars: []string{"BUCKET", "PLUGIN_BUCKET"},
		},
		&cli.StringFlag{
			Name:    "prefix",
			Usage:   "upload key prefix",
			EnvVars: []string{"PREFIX", "PLUGIN_PREFIX"},
		},
		&cli.StringFlag{
			Name:    "dir",
			Usage:   "local dir",
			EnvVars: []string{"DIR", "PLUGIN_DIR"},
		},
		&cli.StringFlag{
			Name:    "parallel",
			Usage:   "parallel num",
			EnvVars: []string{"PARALLEL", "PLUGIN_PARALLEL"},
		},
	}

	err := c.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

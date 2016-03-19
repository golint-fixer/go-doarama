package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/twpayne/go-doarama"
)

func newDoaramaClient(c *cli.Context) *doarama.Client {
	return doarama.NewClient(c.GlobalString("apiurl"), c.GlobalString("apiname"), c.GlobalString("apikey"))
}

func newAuthenticatedDoaramaClient(c *cli.Context) (*doarama.Client, error) {
	client := newDoaramaClient(c)
	userID := c.GlobalString("userid")
	userKey := c.GlobalString("userkey")
	switch {
	case userID != "" && userKey == "":
		return client.Anonymous(userID), nil
	case userID == "" && userKey != "":
		return client.Delegate(userKey), nil
	default:
		return nil, errors.New("exactly one of -userid and -userkey must be specified")
	}
}

func newVisualisationURLOptions(c *cli.Context) *doarama.VisualisationURLOptions {
	var vuo doarama.VisualisationURLOptions
	if c.StringSlice("name") != nil {
		vuo.Names = c.StringSlice("name")
	}
	if c.StringSlice("avatar") != nil {
		vuo.Avatars = c.StringSlice("avatar")
	}
	if c.String("avatarbaseurl") != "" {
		vuo.AvatarBaseURL = c.String("avatarbaseurl")
	}
	if c.Bool("fixedaspect") {
		vuo.FixedAspect = c.Bool("fixedaspect")
	}
	if c.Bool("minimalview") {
		vuo.MinimalView = c.Bool("minimalview")
	}
	if c.String("dzml") != "" {
		vuo.DZML = c.String("dzml")
	}
	return &vuo
}

func activityCreateOne(client *doarama.Client, filename string) (*doarama.Activity, error) {
	gpsTrack, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer gpsTrack.Close()
	return client.CreateActivity(filepath.Base(filename), gpsTrack)
}

func activityCreate(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	typeID := c.Int("typeid")
	for _, arg := range c.Args() {
		a, err := activityCreateOne(client, arg)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("ActivityId: %d\n", a.ID)
		if err := a.SetInfo(&doarama.ActivityInfo{
			TypeID: typeID,
		}); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func activityDelete(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	var ids []int
	for _, arg := range c.Args() {
		id64, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		ids = append(ids, int(id64))
	}
	for _, id := range ids {
		a := client.Activity(id)
		if err := a.Delete(); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func create(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	typeID := c.Int("typeid")
	var as []*doarama.Activity
	for _, arg := range c.Args() {
		a, err := activityCreateOne(client, arg)
		if err != nil {
			break
		}
		err = a.SetInfo(&doarama.ActivityInfo{
			TypeID: typeID,
		})
		if err != nil {
			break
		}
		fmt.Printf("ActivityId: %d\n", a.ID)
		as = append(as, a)
	}
	if err != nil {
		for _, a := range as {
			a.Delete()
		}
		return err
	}
	if len(as) == 0 {
		return errors.New("no activitiess specified")
	}
	v, err := client.CreateVisualisation(as)
	if err != nil {
		return err
	}
	fmt.Printf("VisualisationKey: %s\n", v.Key)
	vuo := newVisualisationURLOptions(c)
	fmt.Printf("VisualisationURL: %s\n", v.URL(vuo))
	return nil
}

type ByName doarama.ActivityTypes

func (ats ByName) Len() int           { return len(ats) }
func (ats ByName) Less(i, j int) bool { return ats[i].Name < ats[j].Name }
func (ats ByName) Swap(i, j int)      { ats[i], ats[j] = ats[j], ats[i] }

func queryActivityTypes(c *cli.Context) error {
	client := newDoaramaClient(c)
	ats, err := client.ActivityTypes()
	if err != nil {
		return err
	}
	sort.Sort(ByName(ats))
	for _, at := range ats {
		fmt.Printf("%s: %d\n", at.Name, at.ID)
	}
	return nil
}

func visualisationCreate(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	var as []*doarama.Activity
	for _, arg := range c.Args() {
		id64, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		a := client.Activity(int(id64))
		as = append(as, a)
	}
	v, err := client.CreateVisualisation(as)
	if err != nil {
		return err
	}
	fmt.Printf("VisualisationKey: %s\n", v.Key)
	return nil
}

func visualisationURL(c *cli.Context) error {
	client := newDoaramaClient(c)
	vuo := newVisualisationURLOptions(c)
	for _, arg := range c.Args() {
		v := client.Visualisation(arg)
		fmt.Printf("VisualisationURL: %s\n", v.URL(vuo))
	}
	return nil
}

func logError(f func(*cli.Context) error) func(*cli.Context) {
	return func(c *cli.Context) {
		if err := f(c); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "doarama"
	app.Usage = "A command line interface to doarama.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apiurl",
			Value:  doarama.APIURL,
			Usage:  "Doarama API URL",
			EnvVar: "DOARAMA_API_URL",
		},
		cli.StringFlag{
			Name:   "apikey",
			Usage:  "Doarama API key",
			EnvVar: "DOARAMA_API_KEY",
		},
		cli.StringFlag{
			Name:   "apiname",
			Usage:  "Doarama API name",
			EnvVar: "DOARAMA_API_NAME",
		},
		cli.StringFlag{
			Name:   "userid",
			Usage:  "Doarama user ID",
			EnvVar: "DOARAMA_USER_ID",
		},
		cli.StringFlag{
			Name:   "userkey",
			Usage:  "Doarama user key",
			EnvVar: "DOARAMA_USER_KEY",
		},
	}
	typeIDFlag := cli.IntFlag{
		Name:  "typeid",
		Usage: "type id",
	}
	nameFlag := cli.StringSliceFlag{
		Name:  "name",
		Usage: "name",
	}
	avatarFlag := cli.StringSliceFlag{
		Name:  "avatar",
		Usage: "avatar",
	}
	avatarBaseUrlFlag := cli.StringFlag{
		Name:  "avatarbaseurl",
		Usage: "avatar base URL",
	}
	fixedAspectFlag := cli.BoolFlag{
		Name:  "fixedaspect",
		Usage: "fixed aspect",
	}
	minimalViewFlag := cli.BoolFlag{
		Name:  "minimalview",
		Usage: "minimal view",
	}
	dzmlFlag := cli.StringFlag{
		Name:  "dzml",
		Usage: "DZML",
	}
	app.Commands = []cli.Command{
		{
			Name:    "activity",
			Aliases: []string{"a"},
			Usage:   "Manages activities",
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Creates an activity from one or more tracklogs",
					Action:  logError(activityCreate),
					Flags: []cli.Flag{
						typeIDFlag,
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Deletes one or more activities by id",
					Action:  logError(activityDelete),
				},
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Creates a visualisation URL from one or more tracklogs",
			Action:  logError(create),
			Flags: []cli.Flag{
				typeIDFlag,
				nameFlag,
				avatarFlag,
				avatarBaseUrlFlag,
				fixedAspectFlag,
				minimalViewFlag,
				dzmlFlag,
			},
		},
		{
			Name:    "query-activity-types",
			Aliases: []string{"qat"},
			Usage:   "Queries activity types",
			Action:  logError(queryActivityTypes),
		},
		{
			Name:    "visualisation",
			Aliases: []string{"v"},
			Usage:   "Manages visualisations",
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Creates a visualisation from a list of activities",
					Action:  logError(visualisationCreate),
				},
				{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "Creates a visualisation URL from a visualisation key",
					Action:  logError(visualisationURL),
					Flags: []cli.Flag{
						nameFlag,
						avatarFlag,
						avatarBaseUrlFlag,
						fixedAspectFlag,
						minimalViewFlag,
						dzmlFlag,
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

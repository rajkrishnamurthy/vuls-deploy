package commands

import (
	"context"
	"encoding/xml"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/subcommands"
	"github.com/inconshreveable/log15"
	c "github.com/kotakanbe/goval-dictionary/config"
	"github.com/kotakanbe/goval-dictionary/db"
	"github.com/kotakanbe/goval-dictionary/fetcher"
	"github.com/kotakanbe/goval-dictionary/models"
	"github.com/kotakanbe/goval-dictionary/util"
	"github.com/ymomoi/goval-parser/oval"
)

// FetchRedHatCmd is Subcommand for fetch RedHat OVAL
type FetchRedHatCmd struct {
	Debug     bool
	DebugSQL  bool
	Quiet     bool
	NoDetails bool
	LogDir    string
	LogJSON   bool
	DBPath    string
	DBType    string
	HTTPProxy string
}

// Name return subcommand name
func (*FetchRedHatCmd) Name() string { return "fetch-redhat" }

// Synopsis return synopsis
func (*FetchRedHatCmd) Synopsis() string { return "Fetch Vulnerability dictionary from RedHat" }

// Usage return usage
func (*FetchRedHatCmd) Usage() string {
	return `fetch-redhat:
	fetch-redhat
		[-dbtype=sqlite3|mysql|postgres|redis]
		[-dbpath=$PWD/oval.sqlite3 or connection string]
		[-http-proxy=http://192.168.0.1:8080]
		[-debug]
		[-debug-sql]
		[-quiet]
		[-no-details]
		[-log-dir=/path/to/log]
		[-log-json]


For details, see https://github.com/kotakanbe/goval-dictionary#usage-fetch-oval-data-from-redhat
	$ goval-dictionary fetch-redhat 5 6 7 8
    	or
	$ for i in {5..8}; do goval-dictionary fetch-redhat $i; done

`
}

// SetFlags set flag
func (p *FetchRedHatCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.Debug, "debug", false, "debug mode")
	f.BoolVar(&p.DebugSQL, "debug-sql", false, "SQL debug mode")
	f.BoolVar(&p.Quiet, "quiet", false, "quiet mode (no output)")
	f.BoolVar(&p.NoDetails, "no-details", false, "without vulnerability details")

	defaultLogDir := util.GetDefaultLogDir()
	f.StringVar(&p.LogDir, "log-dir", defaultLogDir, "/path/to/log")
	f.BoolVar(&p.LogJSON, "log-json", false, "output log as JSON")

	pwd := os.Getenv("PWD")
	f.StringVar(&p.DBPath, "dbpath", pwd+"/oval.sqlite3",
		"/path/to/sqlite3 or SQL connection string")

	f.StringVar(&p.DBType, "dbtype", "sqlite3",
		"Database type to store data in (sqlite3, mysql, postgres or redis supported)")

	f.StringVar(
		&p.HTTPProxy,
		"http-proxy",
		"",
		"http://proxy-url:port (default: empty)",
	)
}

// Execute execute
func (p *FetchRedHatCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c.Conf.Quiet = p.Quiet
	c.Conf.DebugSQL = p.DebugSQL
	c.Conf.Debug = p.Debug
	c.Conf.DBPath = p.DBPath
	c.Conf.DBType = p.DBType
	c.Conf.HTTPProxy = p.HTTPProxy
	c.Conf.NoDetails = p.NoDetails

	util.SetLogger(p.LogDir, c.Conf.Quiet, c.Conf.Debug, p.LogJSON)
	if !c.Conf.Validate() {
		return subcommands.ExitUsageError
	}

	if len(f.Args()) == 0 {
		log15.Error("Specify versions to fetch")
		return subcommands.ExitUsageError
	}

	driver, locked, err := db.NewDB(c.RedHat, c.Conf.DBType, c.Conf.DBPath, c.Conf.DebugSQL)
	if err != nil {
		if locked {
			log15.Error("Failed to open DB. Close DB connection before fetching", "err", err)
			return subcommands.ExitFailure
		}
		log15.Error("Failed to open DB", "err", err)
		return subcommands.ExitFailure
	}

	// Distinct
	vers := []string{}
	v := map[string]bool{}
	for _, arg := range f.Args() {
		ver, err := strconv.Atoi(arg)
		if err != nil || ver < 5 {
			log15.Error("Specify version to fetch (from 5 to latest RHEL version)", "arg", arg)
			return subcommands.ExitUsageError
		}
		v[arg] = true
	}
	for k := range v {
		vers = append(vers, k)
	}

	results, err := fetcher.FetchRedHatFiles(vers)
	if err != nil {
		log15.Error("Failed to fetch files", "err", err)
		return subcommands.ExitFailure
	}

	for _, r := range results {
		ovalroot := oval.Root{}
		if err = xml.Unmarshal(r.Body, &ovalroot); err != nil {
			log15.Error("Failed to unmarshal", "url", r.URL, "err", err)
			return subcommands.ExitUsageError
		}
		log15.Info("Fetched", "URL", r.URL, "OVAL definitions", len(ovalroot.Definitions.Definitions))
		defs := models.ConvertRedHatToModel(&ovalroot)

		var timeformat = "2006-01-02T15:04:05"
		t, err := time.Parse(timeformat, ovalroot.Generator.Timestamp)
		if err != nil {
			panic(err)
		}

		root := models.Root{
			Family:      c.RedHat,
			OSVersion:   r.Target,
			Definitions: defs,
			Timestamp:   time.Now(),
		}

		ss := strings.Split(r.URL, "/")
		fmeta := models.FetchMeta{
			Timestamp: t,
			FileName:  ss[len(ss)-1],
		}

		if err := driver.InsertOval(c.RedHat, &root, fmeta); err != nil {
			log15.Error("Failed to insert oval", "err", err)
			return subcommands.ExitFailure
		}
		if err := driver.InsertFetchMeta(fmeta); err != nil {
			log15.Error("Failed to insert meta", "err", err)
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/maskshell/zkcli/core"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/namsral/flag"
)

var gitCommit = "unknown"
var built = "unknown"

const version = "0.5.1"

func main() {
	servers := flag.String("s", "127.0.0.1:2181", "Servers")
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	showVersion := flag.Bool("version", false, "Show version info")
	verboseLog := flag.Bool("v", false, "Set to true if want to enable zk log, usefull for diagnose zk problems")
	homePath, _ := homedir.Dir()
	defaultConf := filepath.Join(homePath, ".config/zkcli.conf")
	if _, err := os.Stat(defaultConf); err != nil {
		defaultConf = ""
	}
	flag.String(flag.DefaultConfigFlagname, defaultConf, "path to config file")
	flag.Parse()
	args := flag.Args()

	if *showVersion {
		fmt.Printf("Version:\t%s\nGit commit:\t%s\nBuilt: %s\n",
			version, gitCommit, built)
		os.Exit(0)
	}

	config := core.NewConfig(strings.Split(*servers, ","), !*verboseLog)
	if *username != "" && *password != "" {
		auth := core.NewAuth(
			"digest", fmt.Sprintf("%s:%s", *username, *password),
		)
		config.Auth = auth
	}
	conn, err := config.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Simplified process for command line mode, it will use ParseCmd4Cli to process args.
	name, options := core.ParseCmd4Cli(args)
	cmd := core.NewCmd(name, options, conn, config)
	if len(args) > 0 {
		cmd.ExitWhenErr = true
		cmd.Run()
		return
	}

	// When no args, it will enter interactive mode.
	// This will use ParseCmd to process input internally.
	p := prompt.New(
		core.GetExecutor(cmd),
		core.GetCompleter(cmd),
		prompt.OptionTitle("zkcli: A interactive Zookeeper client"),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
}

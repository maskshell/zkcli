package core

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"os"
	"strings"
)

const flag int32 = 0

var acl = zk.WorldACL(zk.PermAll)
var ErrUnknownCmd = errors.New("unknown command")

type Cmd struct {
	Name        string
	Options     []string
	ExitWhenErr bool
	Conn        *zk.Conn
	Config      *Config
}

func NewCmd(name string, options []string, conn *zk.Conn, config *Config) *Cmd {
	return &Cmd{
		Name:    name,
		Options: options,
		Conn:    conn,
		Config:  config,
	}
}

func ParseCmd(input string) (string, []string) {
	var part1 string
	var part2 []string

	var temp string
	var inQuotes bool

	// Parse the string character by character
	for i, r := range input {
		c := string(r)

		if c == "\"" {
			if i > 0 && string(input[i-1]) == "\\" {
				temp = temp[:len(temp)-1] + "\""
			} else {
				inQuotes = !inQuotes
			}
		} else if c == " " && !inQuotes {
			if part1 == "" {
				part1 = temp
			} else {
				part2 = append(part2, temp)
			}
			temp = ""
		} else {
			temp += c
		}
	}

	if temp != "" {
		if part1 == "" {
			part1 = temp
		} else {
			part2 = append(part2, temp)
		}
	}

	// Process part1:
	// If there are paired double quotes at the beginning and end of part1, remove them,
	// and then, if there are double quotes escaped by backslashes, remove the backslashes.
	if part1 != "" {
		if part1[0] == '"' && part1[len(part1)-1] == '"' {
			part1 = part1[1 : len(part1)-1]
		}
		part1 = strings.Replace(part1, `\"`, `"`, -1)
	}

	// Process all string elements in part2:
	// If there are paired double quotes at the beginning and end of the element, remove them,
	// and then, if there are double quotes escaped by backslashes, remove the backslashes.
	for i, str := range part2 {
		if str[0] == '"' && str[len(str)-1] == '"' {
			part2[i] = str[1 : len(str)-1]
		}
		part2[i] = strings.Replace(part2[i], `\"`, `"`, -1)
	}

	return part1, part2
}

// ParseCmd4Cli
// zk command parsing, cli simplified version
func ParseCmd4Cli(cmd []string) (name string, options []string) {
	if len(cmd) < 2 {
		return
	}

	return cmd[0], cmd[1:]
}

// CombineArgs
// The input parameter is an array of string ([]string),
// and the function combines the array elements into a string.
// Each element string should be enclosed in a pair of double quotes, and the elements are separated by spaces.
// In particular, if the string content of the array element contains double quotes,
// then the double quotes should be replaced with the escape character "\"".
func CombineArgs(args []string) string {
	for i, str := range args {
		str = strings.Replace(str, "\"", "\\\"", -1)
		args[i] = "\"" + str + "\""
	}
	return strings.TrimSpace(strings.Join(args, " "))
}

func (c *Cmd) ls() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	children, _, err := c.Conn.Children(p)
	if err != nil {
		return
	}
	fmt.Printf("[%s]\n", strings.Join(children, ", "))
	return
}

func (c *Cmd) get() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	value, stat, err := c.Conn.Get(p)
	if err != nil {
		return
	}
	fmt.Printf("%+v\n%s\n", string(value), fmtStat(stat))
	return
}

func (c *Cmd) create() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	p = cleanPath(p)
	_, err = c.Conn.Create(p, []byte(data), flag, acl)
	if err != nil {
		return
	}
	fmt.Printf("Created %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

func (c *Cmd) set() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	data := ""
	options := c.Options
	if len(options) > 0 {
		p = options[0]
		if len(options) > 1 {
			data = options[1]
		}
	}
	p = cleanPath(p)
	stat, err := c.Conn.Set(p, []byte(data), -1)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", fmtStat(stat))
	return
}

func (c *Cmd) delete() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)
	err = c.Conn.Delete(p, -1)
	if err != nil {
		return
	}
	fmt.Printf("Deleted %s\n", p)
	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

func (c *Cmd) deleteAll() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	p := "/"
	options := c.Options
	if len(options) > 0 {
		p = options[0]
	}
	p = cleanPath(p)

	err = c.deleteRecursive(p)
	if err != nil {
		return err
	}
	return
}

func (c *Cmd) deleteRecursive(p string) (err error) {
	p = cleanPath(p)

	children, _, err := c.Conn.Children(p)
	if err != nil {
		return
	}

	for _, child := range children {
		path := fmt.Sprintf("%s/%s", p, child)
		grandChildren, _, serr := c.Conn.Children(path)
		if serr != nil {
			return serr
		}

		if len(grandChildren) > 0 {
			for _, grandChild := range grandChildren {
				grandPath := fmt.Sprintf("%s/%s", path, grandChild)
				err := c.deleteRecursive(grandPath)
				if err != nil {
					return err
				}
				if serr != nil {
					return serr
				}
			}
		}

		serr = c.Conn.Delete(path, -1)
		if serr != nil {
			return serr
		}
		fmt.Printf("Deleted %s\n", path)
	}

	err = c.Conn.Delete(p, -1)
	if err != nil {
		return
	}
	fmt.Printf("Deleted %s\n", p)

	root, _ := splitPath(p)
	suggestCache.del(root)
	return
}

func (c *Cmd) close() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	c.Conn.Close()
	if !c.connected() {
		fmt.Println("Closed")
	}
	return
}

func (c *Cmd) connect() (err error) {
	options := c.Options
	var conn *zk.Conn
	if len(options) > 0 {
		cf := NewConfig(strings.Split(options[0], ","), false)
		conn, err = cf.Connect()
		if err != nil {
			return err
		}
	} else {
		conn, err = c.Config.Connect()
		if err != nil {
			return err
		}
	}
	if c.connected() {
		c.Conn.Close()
	}
	c.Conn = conn
	fmt.Println("Connected")
	return err
}

func (c *Cmd) addAuth() (err error) {
	err = c.checkConn()
	if err != nil {
		return
	}

	options := c.Options
	if len(options) < 2 {
		return errors.New("addauth <scheme> <auth>")
	}
	scheme := options[0]
	auth := options[1]
	err = c.Conn.AddAuth(scheme, []byte(auth))
	if err != nil {
		return
	}
	fmt.Println("Added")
	return err
}

func (c *Cmd) connected() bool {
	state := c.Conn.State()
	return state == zk.StateConnected || state == zk.StateHasSession
}

func (c *Cmd) checkConn() (err error) {
	if !c.connected() {
		err = errors.New("connection is disconnected")
	}
	return
}

func (c *Cmd) run() (err error) {
	switch c.Name {
	case "ls":
		return c.ls()
	case "get":
		return c.get()
	case "create":
		return c.create()
	case "set":
		return c.set()
	case "delete":
		return c.delete()
	case "deleteall":
		return c.deleteAll()
	case "close":
		return c.close()
	case "connect":
		return c.connect()
	case "addauth":
		return c.addAuth()
	default:
		return ErrUnknownCmd
	}
}

func (c *Cmd) Run() {
	err := c.run()
	if err != nil {
		if err == ErrUnknownCmd {
			printHelp()
			if c.ExitWhenErr {
				os.Exit(2)
			}
		} else {
			printRunError(err)
			if c.ExitWhenErr {
				os.Exit(3)
			}
		}
	}
}

func printHelp() {
	fmt.Println(`get <path>
ls <path>
create <path> [<data>]
set <path> [<data>]
delete <path>
deleteall <path>
connect <host:port>
addauth <scheme> <auth>
close
exit`)
}

func printRunError(err error) {
	fmt.Println(err)
}

func cleanPath(p string) string {
	if p == "/" {
		return p
	}
	return strings.TrimRight(p, "/")
}

func GetExecutor(cmd *Cmd) func(s string) {
	return func(s string) {
		name, options := ParseCmd(s)
		cmd.Name = name
		cmd.Options = options
		if name == "exit" {
			os.Exit(0)
		}
		cmd.Run()
	}
}

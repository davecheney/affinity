package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/cmars/affinity"
	"github.com/cmars/affinity/client"
)

type groupCmd struct {
	subCmd
	url     string
	group   string
	homeDir string
	client  *client.Client
}

func groupFlags(cmd *groupCmd) {
	cmd.flags.StringVar(&cmd.url, "url", "", "Affinity server URL")
	cmd.flags.StringVar(&cmd.group, "group", "", "Affinity group")
	cmd.flags.StringVar(&cmd.homeDir, "homedir", "", "Affinity client home (default: ~/.affinity)")
}

func (c *groupCmd) Main(h cmdHandler) {
	if c.url == "" {
		Usage(h, "--url is required")
	}
	if c.group == "" {
		Usage(h, "--group is required")
	}
	if c.homeDir == "" {
		c.homeDir = path.Join(os.Getenv("HOME"), ".affinity")
	}
	authFile := path.Join(c.homeDir, "auth")
	authContents, err := ioutil.ReadFile(authFile)
	if err != nil {
		die(err)
	}
	auth := strings.TrimSpace(string(authContents))
	c.client = &client.Client{Auth: auth, Url: c.url}
}

type userCmd struct {
	groupCmd
	user string
	User User
}

func userFlags(cmd *userCmd) {
	groupFlags(&cmd.groupCmd)
	cmd.flags.StringVar(&cmd.user, "user", "", "Affinity user")
}

func (c *userCmd) Main(h cmdHandler) {
	c.groupCmd.Main(h)
	if c.user == "" {
		Usage(h, "--user is required")
	}
	var err error
	c.User, err = ParseUser(c.user)
	if err != nil {
		die(err)
	}
}

type addGroupCmd struct {
	groupCmd
}

func newAddGroupCmd() *addGroupCmd {
	cmd := &addGroupCmd{}
	groupFlags(&cmd.groupCmd)
	return cmd
}

func (c *addGroupCmd) Name() string { return "add-group" }

func (c *addGroupCmd) Desc() string { return "Add affinity group" }

func (c *addGroupCmd) Main() {
	c.groupCmd.Main(c)
	err := c.client.AddGroup(c.group)
	die(err)
}

type removeGroupCmd struct {
	groupCmd
}

func newRemoveGroupCmd() *removeGroupCmd {
	cmd := &removeGroupCmd{}
	groupFlags(&cmd.groupCmd)
	return cmd
}

func (c *removeGroupCmd) Name() string { return "remove-group" }

func (c *removeGroupCmd) Desc() string { return "Remove affinity group" }

func (c *removeGroupCmd) Main() {
	c.groupCmd.Main(c)
	err := c.client.DeleteGroup(c.group)
	die(err)
}

type showGroupCmd struct {
	groupCmd
}

func newShowGroupCmd() *showGroupCmd {
	cmd := &showGroupCmd{}
	groupFlags(&cmd.groupCmd)
	return cmd
}

func (c *showGroupCmd) Name() string { return "show-group" }

func (c *showGroupCmd) Desc() string { return "Show affinity group" }

func (c *showGroupCmd) Main() {
	c.groupCmd.Main(c)
	g, err := c.client.GetGroup(c.group)
	if err != nil {
		die(err)
	}
	out, err := json.MarshalIndent(g, "", "\t")
	if err != nil {
		die(err)
	}
	os.Stdout.Write(out)
}

type addUserCmd struct {
	userCmd
}

func newAddUserCmd() *addUserCmd {
	cmd := &addUserCmd{}
	userFlags(&cmd.userCmd)
	return cmd
}

func (c *addUserCmd) Name() string { return "add-user" }

func (c *addUserCmd) Desc() string { return "Add user to affinity group" }

func (c *addUserCmd) Main() {
	c.userCmd.Main(c)
	err := c.client.AddUser(c.group, c.User)
	die(err)
}

type removeUserCmd struct {
	userCmd
}

func newRemoveUserCmd() *removeUserCmd {
	cmd := &removeUserCmd{}
	userFlags(&cmd.userCmd)
	return cmd
}

func (c *removeUserCmd) Name() string { return "remove-user" }

func (c *removeUserCmd) Desc() string { return "Remove user from affinity group" }

func (c *removeUserCmd) Main() {
	c.userCmd.Main(c)
	err := c.client.DeleteUser(c.group, c.User)
	die(err)
}

type checkUserCmd struct {
	userCmd
}

func newCheckUserCmd() *checkUserCmd {
	cmd := &checkUserCmd{}
	userFlags(&cmd.userCmd)
	return cmd
}

func (c *checkUserCmd) Name() string { return "check-user" }

func (c *checkUserCmd) Desc() string { return "Check user membership in affinity group" }

func (c *checkUserCmd) Main() {
	c.userCmd.Main(c)
	err := c.client.CheckUser(c.group, c.User)
	die(err)
}
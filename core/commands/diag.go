package commands

import (
	cmds "github.com/stateless-minds/go-ipfs-cmds"
)

var DiagCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Generate diagnostic reports.",
	},

	Subcommands: map[string]*cmds.Command{
		"sys":     sysDiagCmd,
		"cmds":    ActiveReqsCmd,
		"profile": sysProfileCmd,
	},
}

// Package sessioncmd holds `semiclaw session` command tree (list / view /
// delete / ask / continue-stream / stop) for chat history and agent invocation.
//
// Package name `sessioncmd` (not `session`) so callers can `import sdk
// "github.com/vagawind/semiclaw/client"` and use `sdk.Session` without
// shadowing - same hygiene as `profilecmd`.
package sessioncmd

import (
	"github.com/spf13/cobra"

	"github.com/vagawind/semiclaw/cli/internal/cmdutil"
)

// NewCmd builds the `semiclaw session` parent command.
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage chat sessions",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdView(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdAsk(f))
	cmd.AddCommand(NewCmdContinueStream(f))
	cmd.AddCommand(NewCmdStop(f))
	return cmd
}

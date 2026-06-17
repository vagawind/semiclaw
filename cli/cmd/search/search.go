// Package search implements the `semiclaw search` command tree:
// chunks / kb / docs / sessions.
package search

import (
	"github.com/spf13/cobra"

	"github.com/vagawind/semiclaw/cli/internal/cmdutil"
)

// NewCmdSearch builds the `semiclaw search` parent. Pure dispatcher to the
// four subcommands - users must pick a verb (chunks / kb / docs / sessions).
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search across chunks, knowledge bases, documents, or sessions",
		Long: `Verb-noun search tree:

  search chunks   "<q>" --kb X     hybrid retrieval (RAG search)
  search kb       "<q>"            find KBs by name / description
  search docs     "<q>" --kb X     find documents inside a KB
  search sessions "<q>"            find chat sessions by title / description`,
		Example: `  semiclaw search chunks "what is RAG?" --kb engineering
  semiclaw search kb     "marketing"
  semiclaw search docs   "Q3 forecast" --kb finance
  semiclaw search sessions "onboarding"`,
	}

	cmd.AddCommand(NewCmdChunks(f))
	cmd.AddCommand(NewCmdKB(f))
	cmd.AddCommand(NewCmdDocs(f))
	cmd.AddCommand(NewCmdSessions(f))
	return cmd
}

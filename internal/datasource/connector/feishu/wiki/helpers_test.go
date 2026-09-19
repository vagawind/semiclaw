package wiki

import "github.com/vagawind/semiclaw/internal/datasource/connector/feishu/core"

func txt(s string) *core.BlockText {
	return &core.BlockText{Elements: []core.TextElement{{TextRun: &core.TextRun{Content: s}}}}
}

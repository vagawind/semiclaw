package format

import sdk "github.com/vagawind/semiclaw/client"

	sdk "github.com/vagawind/semiclaw/client"
)

// IndexReferences projects full SDK search results into stable lookup keys.
// It never mutates the SDK events, which keeps the raw NDJSON path lossless.
// fallbackKBID is used by `chat`, whose single KB is known by the CLI even
// when an older server omits knowledge_base_id from a reference.
func IndexReferences(refs []*sdk.SearchResult, fallbackKBID string) []ReferenceIndex {
	indexes := make([]ReferenceIndex, 0, len(refs))
	for _, r := range refs {
		if r == nil || r.ID == "" {
			continue
		}
		kbID := r.KnowledgeBaseID
		if kbID == "" {
			kbID = fallbackKBID
		}
		indexes = append(indexes, ReferenceIndex{
			KBID:          kbID,
			ChunkID:       r.ID,
			ParentChunkID: r.ParentChunkID,
		})
	}
	return indexes
}

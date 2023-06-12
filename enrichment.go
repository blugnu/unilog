package unilog

import "context"

type EnrichmentFunc func(context.Context, Enricher) Entry

var enrichmentFuncs []EnrichmentFunc

// RegisterEnrichment adds a new enrichment function.  All enrichment
// functions are called whenever a new unilog.Entry is initialised.
func RegisterEnrichment(d EnrichmentFunc) {
	enrichmentFuncs = append(enrichmentFuncs, d)
}

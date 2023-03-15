package unilog

import "context"

type EnrichmentFunc func(context.Context, Enricher) Entry

var enrichmentFuncs []EnrichmentFunc

// EnrichWith adds a new enrichment function to the registered
// functions.  All enrichment functions are called whenever a new
// unilog.Entry is initialised.
func EnrichWith(d EnrichmentFunc) {
	enrichmentFuncs = append(enrichmentFuncs, d)
}

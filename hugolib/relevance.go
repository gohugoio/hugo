package hugolib

import (
	"math"

	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// More information about Tf-Idf algorithm 
// http://en.wikipedia.org/wiki/Tf%E2%80%93idf

// TokenRelevance represents the score for a given token (Value)
type TokenRelevance struct {
	Value string
	Freq float64
	Tf float64
	Idf float64
	TfIdf float64
}

// Relevance stores all scores for tokens indexed by value
type Relevance map[string]TokenRelevance

// appendIfMissing appends p element in slice if missing
func appendIfMissing(slice []*Page, p *Page) []*Page {
	for _, el := range slice {
		if el == p {
			return slice
		}
	}
	return append(slice, p)
}

// BuildRelatedGraph is used to create page relations weighted with Tf-Idf score
func BuildRelatedGraph(s *Site) {

	taxonomies := viper.GetStringMapString("Taxonomies")

	for _, plural := range taxonomies {
		// Extract all tags as vocabulary
		m := s.Taxonomies[plural]
		vocabulary := make([]string, len(m))
		i := 0
	    for k, _ := range m {
	        vocabulary[i] = k
	        i++
	    }

	    // For all pages
		for _, p := range s.Pages {
			// Get taxononomy values
			vals := p.GetParam(plural)

			// If pages has identified taxonomy
			if vals != nil {
				v, ok := vals.([]string)
				if ok {
					if ( p.relevance == nil) {
						p.relevance = make(Relevance)	
					}
					
					// Generate the relevance score according to all token present in vocabuary
					for _, voc := range vocabulary {
						// Default token is not used in this page
						freq := 0.
						if( helpers.InStringArray(v, voc) ) {
							// Token is used in this page
							freq = 1.
						}
						// Token frequency
						tf := freq / float64(len(v))
						// Inverted document frequency
						idf := math.Log(float64(len(vocabulary)) / float64(1 + s.Taxonomies[plural][voc].Len()))
						// Score
						tfidf := tf * idf

						// Store scores on the page
						x := TokenRelevance{voc,freq,tf,idf,tfidf}
						p.relevance[voc] = x
					}

					// for all tagged values of the current page
					for _, idx := range v {
						// Set potential related documents to all token tagged pages
						potential_related := s.Taxonomies[plural][idx].Pages()
						// Remove unrelated documents
						for _, rpage := range potential_related {
							if(rpage.relevance[idx].TfIdf > 0. && rpage != p) {
								// Store related page to current page
								p.RelatedPages = appendIfMissing(p.RelatedPages, rpage)
								// Store related page to reverse association
								rpage.RelatedPages = appendIfMissing(rpage.RelatedPages, p)
							}
						}
					}
				} else {
					jww.ERROR.Printf("Invalid %s in %s\n", plural, p.File.LogicalName())
				}
			}
		}

	}

}


package segmentsconv

import "github.com/m1al04949/avito-tech-service/internal/model"

func SegmentsConv(s []model.Segment) []string {

	var segments []string

	for i := 0; i < len(s); i++ {
		segments = append(segments, s[i].Slug)
	}

	return segments
}

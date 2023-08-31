package segment



type Segment struct {
	Name        string `json:"name" ksql:"name"`
	AudienceCvg int    `json:"audience_cvg" ksql:"audience_cvg"`
}


func NewSegment(name string, audience_cvg int) Segment {
	return Segment{
		Name:        name,
		AudienceCvg: audience_cvg,
	}
}


func (s Segment) GetAudienceCvg() int {
	return s.AudienceCvg
}


func (s Segment) GetName() string {
	return s.Name
}



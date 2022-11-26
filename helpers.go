package opensky

import "time"

func GetStartAndEndOfDay(t time.Time) (time.Time, time.Time) {
	begin := t.UTC().Truncate(24 * time.Hour)
	end := begin.Add(24 * time.Hour)
	return begin, end
}

func (a *Altitude) ConvertToFeet() float64 {
	return a.Value * metersToFeet
}

func (s *Speed) ConverToKnots() float64 {
	return s.Value * mPerSToKnots
}

func (s Speed) ConvertToFtPerMin() float64 {
	return s.Value * mPerSToFtPerMin
}

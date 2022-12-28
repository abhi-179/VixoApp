package config

import "time"

func TimeConvert(loc string, currTime time.Time) time.Time {
	//init the loc
	location, _ := time.LoadLocation(loc)
	//set timezone,
	now := currTime.In(location)
	return now
}

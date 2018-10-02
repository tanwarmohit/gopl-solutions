package xkcd

// Comic data structure for holding json response
type Comic struct {
	Day, Month, Year string
	Num              int
	Title            string
	Transcript       string
	Img              string
	Alt              string
}

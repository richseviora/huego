package resources

type Reference struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

type Dimming struct {
	Brightness float64 `json:"brightness"`
}

type XYCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ResourceError struct {
	Description string `json:"description"`
}

type ResourceList[T any] struct {
	Data   []T             `json:"data"`
	Errors []ResourceError `json:"errors"`
}

type ResourceUpdateResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Reference `json:"data"`
}

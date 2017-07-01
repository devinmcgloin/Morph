package request

type PatchImageRequest struct {
	Tags         []string
	Aperature    string
	ExposureTime string
	FocalLength  string
	ISO          int
	Make         string
	Model        string
	LensMake     string
	LensModel    string
	CaptureTime  string
}

package internal

type Path struct {
	SrcPath       string            `json:"srcPath"`
	Status        Status            `json:"status"`
	TgtPath       *string           `json:"tgtPath"`
	ResolvedTgt   *DetectedProfile  `json:"resolvedTgt"`
	DetectedProfs []DetectedProfile `json:"detectedProfs"`
}

type Status string

const (
	StatusOk                   Status = "ok"
	StatusErrorSrcNotFound     Status = "error_src_not_found"
	StatusErrorTgtNotFound     Status = "error_tgt_not_found"
	StatusErrorSrcNotSymlink   Status = "error_tgt_not_prof"
	StatusErrorTgtUnresolvable Status = "error_tgt_not_resolvable"
)

type DetectedProfile struct {
	Name string
	Path string
}

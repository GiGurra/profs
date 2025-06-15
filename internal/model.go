package internal

import (
	"fmt"
	"path/filepath"
)

type Path struct {
	SrcPath       string            `json:"srcPath"`
	Status        Status            `json:"status"`
	TgtPath       *string           `json:"tgtPath"`
	ResolvedTgt   *DetectedProfile  `json:"resolvedTgt"`
	DetectedProfs []DetectedProfile `json:"detectedProfs"`
}

func (path *Path) ProfsDir() (string, error) {
	if path.TgtPath == nil {
		return "", fmt.Errorf("target path is nil for source path: %s", path.SrcPath)
	}
	if *path.TgtPath == "" {
		return "", fmt.Errorf("target path is empty for source path: %s", path.SrcPath)
	}
	res := filepath.Dir(*path.TgtPath)
	if fileOrDirExists(res) {
		return res, nil
	} else {
		return "", fmt.Errorf("target profs directory does not exist: %s", res)
	}
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

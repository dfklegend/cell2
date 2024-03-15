package consts

import (
	"errors"
)

var (
	ErrNotImportant    = errors.New("something error, not important")
	ErrToDefine        = errors.New("something error, to be define")
	ErrCanNotFindScene = errors.New("can not find scene")
)

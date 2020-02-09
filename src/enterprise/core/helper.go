package core

import (
	"github.com/labstack/gommon/log"
)

// Try は成功するか上限回数まで処理を繰り返し行います
func Try(f func() error, limit int) error {
	count := 0
	for {
		err := f()
		if err == nil {
			return nil
		}
		count++
		if count >= limit {
			return err
		}
		log.Warn(err)
	}
}

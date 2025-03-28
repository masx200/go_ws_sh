package go_ws_sh

import (
	"time"

	"github.com/golang-module/carbon/v2"
)

func Formattimewithstdtime(tt time.Time) string {
	return FormatTimeWithCarbon(carbon.CreateFromStdTime(tt))
}

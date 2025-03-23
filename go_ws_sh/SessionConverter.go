package go_ws_sh

import (
	"github.com/golang-module/carbon/v2"
)

// SessionsToMapSlice 将 []Session 转换为 []map[string]any
func SessionsToMapSlice(sessions []Session) []map[string]any {
	result := make([]map[string]any, len(sessions))
	for i, session := range sessions {
		result[i] = map[string]any{
			"name":       session.Name,
			"cmd":        session.Cmd,
			"args":       session.Args,
			"dir":        session.Dir,
			"created_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(session.CreatedAt)),
			"updated_at": FormatTimeWithCarbon(carbon.CreateFromStdTime(session.UpdatedAt)),
		}
	}
	return result
}

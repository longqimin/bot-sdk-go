package video_player

import (
	"github.com/longqimin/bot-sdk-go/bot/directive"
)

type StopDirective struct {
	directive.BaseDirective
}

func NewStopDirective() *StopDirective {
	stop := &StopDirective{}
	stop.Type = "VideoPlayer.Stop"
	return stop
}

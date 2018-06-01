package model

import (
	"encoding/json"
	"regexp"

	"github.com/dueros/bot-sdk-go/bot/data"
)

type Response struct {
	session *Session
	request interface{}
	data    map[string]interface{}
}

func NewResponse(session *Session, request interface{}) *Response {
	d := make(map[string]interface{})
	return &Response{
		data:    d,
		session: session,
		request: request,
	}
}

/**
 * @desc 询问用户时，返回的speech.
 *
 * 此时设备的麦克风会进入收音状态，比如设备灯光亮起
 * TIP: 一般技能要完成一项任务，还缺少一些信息，主动发起对用户的询问的时候使用
 */
func (this *Response) Ask(speech string) *Response {
	this.Tell(speech)
	this.HoldOn()
	return this
}

func (this *Response) AskSlot(speech string, slot string) *Response {
	this.Ask(speech)

	request, ok := this.request.(IntentRequest)
	if ok {
		request.Dialog.ElicitSlot(slot)
	}
	return this
}

/**
 * @desc 回复用户，返回的speech
 */
func (this *Response) Tell(speech string) *Response {
	this.data["outputSpeech"] = this.formatSpeech(speech)
	return this
}

/**
 * @desc 返回卡片.
 *
 * 针对有屏幕的设备，比如: 电视、show，可以呈现更多丰富的信息给用户
 * 卡片协议参考：TODO
 */
func (this *Response) DisplayCard(card interface{}) *Response {
	this.data["card"] = card

	return this
}

/**
 * @desc 返回指令. 比如，返回音频播放指令，使设备开始播放音频
 *
 * TIP: 可以同时返回多个指令，设备按返回顺序执行这些指令，指令协议参考TODO
 */
func (this *Response) Command(directive interface{}) *Response {
	_, ok := this.data["directives"]
	if !ok {
		this.data["directives"] = make([]interface{}, 0)
	}

	directives, ok := this.data["directives"].([]interface{})
	directives = append(directives, directive)

	this.data["directives"] = directives

	return this
}

/**
 * @desc 保持会话.
 *
 * 此时设备的麦克风会自动开启监听用户说话
 */
func (this *Response) HoldOn() *Response {
	this.data["shouldEndSession"] = false
	return this
}

func (this *Response) Build() string {
	//session
	attributes := this.session.GetData().Attributes

	ret := map[string]interface{}{
		"version":  "2.0",
		"session":  data.SessionResponse{Attributes: attributes},
		"response": this.data,
	}

	//intent request
	request, ok := this.request.(IntentRequest)
	if ok {
		ret["context"] = data.ContextResponse{Intent: request.Dialog.Intents[0].GetData()}

		directive := request.Dialog.GetDirective()
		if directive != nil {
			this.Command(directive)
		}
	}

	response, _ := json.Marshal(ret)

	return string(response)
}

func (this *Response) formatSpeech(speech string) data.Speech {
	match, _ := regexp.MatchString("^<speak>", speech)

	if match {
		return data.Speech{
			Type: "SSML",
			Ssml: speech,
		}
	}

	return data.Speech{
		Type: "PlainText",
		Text: speech,
	}
}
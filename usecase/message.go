package usecase

import "fmt"

var (
	replyText                = fmt.Sprintf("今から調べるから数分待っててね%s", EMOJI_NAGIGATE_MAN)
	littleConditionErrorText = fmt.Sprintf("%s入力エラー%s\n検索条件が足りないよ%s%s\n「エリア名 キーワード」で検索してね！", EMOJI_WARNING, EMOJI_WARNING, EMOJI_TIRED_FACE, EMOJI_TIRED_FACE)
	manyConditionErrorText   = fmt.Sprintf("%s入力エラー%s\n検索条件が多すぎるよ%s%s\n「エリア名 キーワード」で検索してね！", EMOJI_WARNING, EMOJI_WARNING, EMOJI_TIRED_FACE, EMOJI_TIRED_FACE)
)

const (
	EMOJI_TIRED_FACE   = "😫"
	EMOJI_WARNING      = "🚫"
	EMOJI_NAGIGATE_MAN = "💁"
)

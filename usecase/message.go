package usecase

import "fmt"

var (
	replyText                = fmt.Sprintf("今から調べるから数分待っててね%c", EMOJI_NAGIGATE_MAN)
	littleConditionErrorText = fmt.Sprintf("%c入力エラー%c\n検索条件が足りないよ%c%c\n「エリア名 キーワード」で検索してね！", EMOJI_WARNING, EMOJI_WARNING, EMOJI_TIRED_FACE, EMOJI_TIRED_FACE)
	manyConditionErrorText   = fmt.Sprintf("%c入力エラー%c\n検索条件が多すぎるよ%c%c\n「エリア名 キーワード」で検索してね！", EMOJI_WARNING, EMOJI_WARNING, EMOJI_TIRED_FACE, EMOJI_TIRED_FACE)
)

const (
	EMOJI_TIRED_FACE   = 0x1F62B
	EMOJI_WARNING      = 0x1F6AB
	EMOJI_NAGIGATE_MAN = 0x1F481
)

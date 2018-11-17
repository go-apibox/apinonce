// 错误定义

package apinonce

import (
	"github.com/go-apibox/api"
)

// error type
const (
	errorMissingNonce = iota
	errorInvalidNonce
	errorNonceExist
	errorNonceCountExceed
)

var ErrorDefines = map[api.ErrorType]*api.ErrorDefine{
	errorMissingNonce: api.NewErrorDefine(
		"MissingNonce",
		[]int{0},
		map[string]map[int]string{
			"en_us": {
				0: "Missing nonce!",
			},
			"zh_cn": {
				0: "缺少随机串！",
			},
		},
	),
	errorInvalidNonce: api.NewErrorDefine(
		"InvalidNonce",
		[]int{0},
		map[string]map[int]string{
			"en_us": {
				0: "Invalid nonce!",
			},
			"zh_cn": {
				0: "无效的随机串！",
			},
		},
	),
	errorNonceExist: api.NewErrorDefine(
		"NonceExist",
		[]int{0},
		map[string]map[int]string{
			"en_us": {
				0: "Nonce already exists!",
			},
			"zh_cn": {
				0: "随机串已存在！",
			},
		},
	),
	errorNonceCountExceed: api.NewErrorDefine(
		"NonceCountExceed",
		[]int{0},
		map[string]map[int]string{
			"en_us": {
				0: "Nonce count exceed!",
			},
			"zh_cn": {
				0: "随机串数量超出限制！",
			},
		},
	),
}

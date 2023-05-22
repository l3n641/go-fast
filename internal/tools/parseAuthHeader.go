package tools

import (
	"errors"
	"strings"
)

func ParseAuthHeader(auth string) (string, error) {
	if auth == "" {
		return "", errors.New("请求头中的auth为空")
	}
	state := strings.HasPrefix(auth, "Bearer ")
	if state == false {
		return "", errors.New("请求头中的auth格式错误")
	}
	result := strings.Split(auth, " ")
	if !(len(result) == 2 && result[0] == "Bearer") {
		return "", errors.New("请求头中的auth格式错误")
	}
	return result[1], nil
}

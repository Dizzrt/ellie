package http

import "strings"

const _BASE_CONTENT_TYPE = "application"

func contentType(subType string) string {
	return _BASE_CONTENT_TYPE + "/" + subType
}

func contentSubType(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}

	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}

	if right <= left {
		return ""
	}

	return contentType[left+1 : right]
}

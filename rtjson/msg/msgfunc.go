package msg

import "sync"

var msgList = make(map[string]map[int]string, 0) // 错误列表
var rLock sync.RWMutex

// SetError 设置错误
func SetError(business string, lang string, info map[int]string) {
	rLock.Lock()
	defer rLock.Unlock()

	// 判断业务类型
	businessList, ok := codeRange[business]
	if !ok {
		panic(BuiltMessage[UnknownBusinessTypeErr])
	}

	// 与内置message合并并校验传递的错误码
	codeList := make(map[int]string, 0)
	for k, v := range info {
		if k < businessList[0] || k > businessList[1] {
			panic(BuiltMessage[CodeRangeErr])
		}
		if _, ok := codeList[k]; ok {
			panic(BuiltMessage[CodeDuplicationErr])
		}
		codeList[k] = v
	}
	for k, v := range BuiltMessage {
		codeList[k] = v
	}

	msgList[lang] = codeList
}

// GetError 返回错误
func GetError(code int) string {
	rLock.RLock()
	defer rLock.RUnlock()

	// 查找语言
	msgInfo, ok := msgList["en"]
	if !ok {
		return BuiltMessage[UnknownErr]
	}

	// 查找错误
	message, ok := msgInfo[code]
	if !ok {
		return BuiltMessage[UnsetErr]
	}

	return message
}

// GetLangError 根据语言返回错误
func GetLangError(lang string, code int) string {
	rLock.RLock()
	defer rLock.RUnlock()

	if len(lang) == 0 {
		lang = "en"
	}
	// 查找语言
	msgInfo, ok := msgList[lang]
	if !ok {
		return BuiltMessage[UnknownErr]
	}

	// 查找错误
	message, ok := msgInfo[code]
	if !ok {
		return BuiltMessage[UnsetErr]
	}

	return message
}

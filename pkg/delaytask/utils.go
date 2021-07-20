package delaytask

import (
	"reflect"
	"strings"
)

//判断当前有效事件
func isUsableKey(key, namespace string) bool {
	return strings.HasPrefix(key, ProcessQueue + "-" + namespace)
}

func getQueuePrefix(namespace string) string {
	return ProcessQueue + "-" + namespace + "_content_"
}

func getOriginKey(key string) string {
	keyList := strings.Split(key, "_")
	if len(keyList) < 3 {
		return ""
	}
	return strings.Join(keyList[2:], "_")
}

func generateKey(key, namespace string) string {
	return ProcessQueue + "-" + namespace + "_key_" + key
}

func generateContentKey(key, namespace string) string {
	return ProcessQueue + "-" + namespace + "_content_" + key
}

func transformKey2Content (key string) string {
	keyList := strings.Split(key, "_")
	if len(keyList) < 3 {
		return ""
	}
	keyList[1] = "content"
	return strings.Join(keyList, "_")
}

func transformContent2Key (content string) string {
	keyList := strings.Split(content, "_")
	if len(keyList) < 3 {
		return ""
	}
	keyList[1] = "key"
	return strings.Join(keyList, "_")
}

func getStructTagDefaultValue(data interface{}, name string) string {
	field, _ := reflect.TypeOf(data).FieldByName(name)
	return field.Tag.Get("default")
}

func getConsumeQueueName(namespace string) string {
	return ConsumeQueue + "_" + namespace
}
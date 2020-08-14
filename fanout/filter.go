package fanout

func AllowAll(interface{}) bool {
	return true
}

func DenyAll(interface{}) bool {
	return false
}

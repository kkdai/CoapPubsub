package CoapPubsub

func parseToString(in interface{}) string {
	val, ok := in.([]uint8)
	if ok {
		return string(val)
	} else {
		return ""
	}
}

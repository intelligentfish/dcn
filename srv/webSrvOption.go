package srv

// WebSrvOption web service option
type WebSrvOption func(object *WebSrv)

// WebSrvListenAddrOption web service listen addr option
func WebSrvListenAddrOption(lAddr string) WebSrvOption {
	return func(object *WebSrv) {
		object.lAddr = lAddr
	}
}

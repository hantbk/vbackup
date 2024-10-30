package consts

const (
	appname = "backup"
)

// Key retrieves the cache key.
func Key(funcn string, parms ...string) (res string) {
	res = appname + "_" + funcn
	for _, parm := range parms {
		res = res + "_" + parm
	}
	return
}

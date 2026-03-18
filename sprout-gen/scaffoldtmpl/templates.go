package scaffoldtmpl

// SproutTmpl 是默认的 .sprout 文件模板。
var SproutTmpl = `syntax = "v1"

server {
    prefix "/api"
}

type PingResp {
    Message string ` + "`" + `json:"message"` + "`" + `
}

service {{.Name}} {
    public {
        GET "/ping" Ping() -> PingResp
    }
}
`

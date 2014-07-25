package templates

func init() {

	registerTemplate("proxy.pac.tpl", `
function FindProxyForURL(url, host) {
    var DEFAULT_PROXY = "DIRECT";
    
    var RUNNING_PROXY = "{{with .Proxy }}{{range .}}{{.Category}} {{.Address}};{{end}}{{end}};";

    {{with .Role}}{{range .}}
    if({{.Name|html}}.test(url)){ return {{if myeq .Category  "0" }}RUNNING_PROXY{{else}}DEFAULT_PROXY{{end}} }
    {{end}}{{end}}


    {{with .GFW}}{{range .}}
    if({{.|html}}.test(url)){ return  RUNNING_PROXY+DEFAULT_PROXY;  }
    {{end}}{{end}}
    return DEFAULT_PROXY;
    }
}`)
}

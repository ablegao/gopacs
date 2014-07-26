package templates

func init() {

	registerTemplate("proxy.pac.tpl", `
function FindProxyForURL(url, host) {
    var DEFAULT_PROXY = "DIRECT";


    if (host == "127.0.0.1" ||
        isInNet(host, "10.0.0.0", "255.0.0.0") ||
        isInNet(host, "192.168.0.0", "255.255.0.0")) {
        return DEFAULT_PROXY;
    }
    
    var RUNNING_PROXY = "{{with .Proxy }}{{range .}} {{.Category}} {{.Address}};{{end}}DIRECT{{end}}";

    {{with .Role}}{{range .}}
    if({{.Name|html}}.test(url)){ return {{if myeq .Category  "a" }}RUNNING_PROXY{{else}}DEFAULT_PROXY{{end}} }
    {{end}}{{end}}


   {{with .GFW}}{{range .}}
   if({{.|html}}.test(url)){ return  RUNNING_PROXY;  }
   {{end}}{{end}}

    return DEFAULT_PROXY;
    
}`)
}

/*
  {{with .GFW}}{{range .}}
   if({{.|html}}.test(url)){ return  RUNNING_PROXY;  }
   {{end}}{{end}}
*/

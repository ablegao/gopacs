package templates

func tplLayout(tpl string) string {
	text := `
<!doctype html>
<html lang="zh-cn" ng-app="goProxy">
<head>
<meta http-equiv="Content-Type" content="text/html;charset=utf-8">​
<title>GoPacProxy</title>
<link rel="stylesheet" href="/static/lib/onsen/css/onsenui.css">  
<link rel="stylesheet" href="/static/lib/onsen/css/onsen-css-components.min.css">  
<link rel="stylesheet" href="/static/style/app.css">  

  <script src="/static/lib/onsen/js/angular/angular.js"></script>    
<script src="/static/lib/onsen/js/angular/angular-touch.js"></script>
<script src="/static/lib/onsen/js/angular/angular-animate.js"></script>
<script src="/static/lib/onsen/js/onsenui.js"></script>
<script src="/static/js/app.js"></script>
</head>
<body>
	` + tpl + `

</body>
</html>
`
	return text
}

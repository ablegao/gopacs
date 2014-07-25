package templates

func init() {

	registerTemplate("home.tpl", tplLayout(`
	<ons-split-view
    var="app.splitView"
    secondary-page="/static/html/menu.html"
    main-page="/static/html/page1.html"
    main-page-width="70%"
    collapse="width 500px">
  </ons-split-view>
`))
}

gopacs
======

这一个pac 代理文件服务。



使用golang 语言编写服务部分。 

angularjs + onsenui 实现的管理界面部分。 

引用了autoproxy 的gfwlist 列表。 


##使用：

基础的golang 环境。 

	git clone git@github.com:ablegao/gopacs.git


	cd  gopacs 
	go build . 
	./gopacs


浏览器访问 http://localhost:8888 可以添加和规则和代理服务器地址。 

pac文件地址: http://localhost:8888/proxy.pac 

默认局域网内其他机器可以直接访问， http://服务IP:8888/proxy.pac


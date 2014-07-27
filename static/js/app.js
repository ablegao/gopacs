(function() {
    'use strict';
    var module = angular.module('goProxy', ['onsen']);

    module.controller('ProxyListCtrl', ["$scope", 'DbApi',"Params",
        function($scope, DbApi,Params) {
            console.log("ProxyListCtrl");

            $scope.proxy = new DbApi("/proxy_ctrl" , {
                        name: "",
                        category: "SOCKS5",
                        address: "0.0.0.0:8080"
                    });
            $scope.proxy.loadData();

            ons.ready(function() {
                myNavigator.on("prepop",function(){ $scope.proxy.loadData(); });
            });


            $scope.gotoEdit = function(id) {
                console.log(id);
                
                Params.ProxyListIndex= id;
                myNavigator.pushPage('/static/html/proxy_plus.html', {
                    animation: 'slide'
                });
            }
            //console.log($scope.proxy.index);
        }
    ]);

    module.controller('ProxyEditCtrl', ['$scope', 'DbApi',"Params",
        function($scope, DbApi , Params) {
            $scope.proxys =  new DbApi("/proxy_ctrl" , {
                        name: "",
                        category: "SOCKS5",
                        address: "0.0.0.0:8080"
                    });

            $scope.proxys.loadData(function(){

                $scope.proxy = $scope.proxys.items[Params.ProxyListIndex];
                $scope.index =Params.ProxyListIndex;

                $scope.$watch('proxy', function(to, form) {
                    console.log("proxy .. watch.save")
                    $scope.proxys.save();
                }); 
            });
            
        }
    ]);
    module.controller('ParamEditCtrl', ['$scope', 'DbApi',"Params",
        function($scope, DbApi,Params) {

            console.log(Params.paramsEditIndex);
            $scope.modes = new DbApi("/params_ctrl", {
                name: "",
                categray: "a",
            });
            $scope.modes.catename = {"a":"连接可用代理", "b":"跳过代理"};

            $scope.modes.loadData(function() {
                $scope.mode = $scope.modes.items[Params.paramsEditIndex];
                $scope.index = Params.paramsEditIndex;
                $scope.$watch('mode', function(to, form) {
                    $scope.modes.save();
                });
            });



        }
    ]);



    module.controller('ParamsCtrl', ["$scope", 'DbApi',"Params",
        function($scope, DbApi,Params) {
            $scope.mode = new DbApi("/params_ctrl", {
                name: "",
                category: "a",
            });

            ons.ready(function() {
                paramNavigator.on("prepop",function(){ $scope.mode.loadData(); });
            });
            
            //console.log("=====");
            $scope.mode.catename = {"a":"连接可用代理", "b":"跳过代理"};
            $scope.mode.loadData();

            $scope.gotoEdit = function(id) {
                Params.paramsEditIndex = id;
                paramNavigator.pushPage('/static/html/role_plus.html', {
                    animation: 'slide'
                });

            }

        }
    ]);




    module.controller('SSHCtrl', ['$scope', 'DbApi',"Params",
        function($scope, DbApi,Params ) {
           
            $scope.modes = newSSHModel(DbApi);

            $scope.modes.loadData(function() {
                $scope.mode = $scope.modes.items[Params.sshEditIndex];
                $scope.index = Params.sshEditIndex;
                $scope.$watch('mode', function(to, form) {
                    $scope.modes.save();
                });
            });
        }
    ]);



    module.controller('SSHListCtrl', ["$scope","$http", 'DbApi',"Params","UpdateStatus",
        function($scope,$http, DbApi,Params, UpdateStatus) {
            $scope.mode = newSSHModel(DbApi);

             $scope.state = new UpdateStatus();

            setInterval(function(){
                $scope.state.update();
                console.log($scope.state)
            },3000);

            ons.ready(function() {
                myNavigator.on("prepop",function(){ $scope.mode.loadData(); });
            });
            
            //console.log("=====");
            $scope.mode.stateName = Params.ssh_state
            $scope.mode.loadData();

            $scope.gotoEdit = function(id) {
                Params.sshEditIndex = id;
                myNavigator.pushPage('/static/html/ssh_plus.html', {
                    animation: 'slide'
                });
            }
            $scope.serverStop= function(id){
                $http({
                    method:"POST",
                    url:"/ssh_stop" , 
                    data:$scope.mode.items[id].name
                }).success=function(response, status, headers, config){
                    console.log(response , status)
                }
            }
            $scope.serverStart = function(id){
                var p  = $http({
                    method:"POST",
                    url:"/ssh_start" , 
                    data:$scope.mode.items[id]
                });
                p.success=function(response, status, headers, config){
                    console.log(response , status)
                }
            }
        }
    ]);

    function newSSHModel(DbApi){
        return new DbApi("/ssh_ctrl", {
                name: "ssh名称",
                address: "root@www.com",
                server_port:"22",
                local_port:"8080",
                passwd:"",
                state: 0
            });
    }


    module.factory('Params', [ function(){
        var params ={
            ssh_state:["关闭" , "连接中" , "正常","无法启动"]
        }

        return params
    }])
    module.factory('UpdateStatus', ['$http', function($http){
        var mode = function(){
            this.state = new Array();
        }
        mode.prototype.update=function(){
            var self = this;
            $http({
                method:"GET",
                url:"/ssh_state",
            }).success(function(response, status, headers, config){
                self.state = response;
            });

        }
        return mode
    }]);

    module.factory('DbApi', ["$http",
        function($http) {

            var mode = function(url, fields) {
                this.url = url;
                this.fields = fields;
                //this.items={};
            }

            mode.prototype = {
                index: 0,
                items: new Array(),
                loadData: function(callback) {
                    var p = $http({
                        method: "GET",
                        url: this.url
                    });
                    var self = this;

                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            self.items = response;
                        }else{
                            self.items= new Array();
                        }
                        if (callback) {
                            callback();
                        }
                    });
                },
                get: function(id) {
                    return this.items[id];
                },
                add: function() {
                    //this.items.push(this.fields);
                    //this.items[]= this.fields;
                    this.items.push(angular.copy(this.fields));
                    console.log(this.items);
                    this.save();

                },
                delete: function(index) {
                    console.log("delete " + index)
                    this.items.splice(index, 1);
                    this.save();

                },
                save: function(callback) {
                    var p = $http({
                        method: "POST",
                        data: this.items,
                        url: this.url
                    });
                    var self = this;
                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            self.items = response;
                        }
                        if (callback) {
                            callback();
                        }
                    });
                }

            }
            return mode;
        }
    ]);



})();

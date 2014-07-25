(function() {
    'use strict';
    var module = angular.module('goProxy', ['onsen']);

    module.controller('ProxyListCtrl', ["$scope", 'ProxyList',
        function($scope, ProxyList) {
            console.log("ProxyListCtrl");

            $scope.proxy = ProxyList;
            
            $scope.proxy.loadData(function(){
              $scope.$watch('proxy', function(to, form) {
                console.log("proxy .. watch")
                $scope.proxy.save();
              });
            });
            //console.log("=====");

            //$scope.proxy.load();

            $scope.gotoEdit = function(id){
                console.log(id);
                $scope.proxy.cate=1;
                $scope.proxy.index = id;
                myNavigator.pushPage('/static/html/proxy_plus.html', {animation: 'slide'});
            }
            //console.log($scope.proxy.index);
        }
    ]);

    module.controller('ProxyEditCtrl', ['$scope' , 'ProxyList', function($scope,ProxyList){
        $scope.proxys= ProxyList;
        $scope.proxy = ProxyList.items[ProxyList.index];
        $scope.index = ProxyList.index ;

        $scope.$watch('proxy', function(to, form) {
                console.log("proxy .. watch.save")
                $scope.proxys.save();
        });
    }]);


    module.factory('ProxyList', ["$http",
        function($http) {
            var proxys = {
                loadData: function(callback) {
                    var p = $http({
                        method: "GET",
                        url: "/proxy_ctrl"
                    });

                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            proxys.items = response;
                        }
                        if (callback){callback();}
                    });


                },
                add:function(){
                  this.items.unshift({name:"点击填写名称",category:"SOCKET5",address:"0.0.0.0:8080",state:0})
                },
                delete: function(index) {
                    this.items.splice(index, 1);
                    this.save();
                },
                save: function(callback) {
                    var p = $http({
                        method: "POST",
                        data: proxys.items,
                        url: "/proxy_ctrl"
                    });

                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            p.items = response;
                        }
                        if (callback){callback();}
                    });

                },
                items: [ ],
                cate:1, //#编辑 , 2添加
                index:0 
            }

            return proxys;
        }
    ]);

    module.controller('ParamEditCtrl', ['$scope' , 'ParamsList', function($scope,Params){
        $scope.modes= Params;
        $scope.mode = Params.items[Params.index];
        $scope.index = Params.index ;

        $scope.$watch('mode', function(to, form) {
            $scope.modes.save();
        });
    }]);



    module.controller('ParamsCtrl', ["$scope", 'ParamsList',
        function($scope, List) {
            console.log("ParamsCtrl");

            $scope.mode = List;
            //console.log("=====");

            $scope.mode.loadData(function(){
              $scope.$watch('mode', function(to, form) {
                $scope.mode.save();
              });
            });

            $scope.gotoEdit = function(id){

                $scope.mode.cate=1;
                $scope.mode.index = id;
                paramNavigator.pushPage('/static/html/role_plus.html', {animation: 'slide'});

            }

        }
    ]);


    module.factory('ParamsList', ["$http",
        function($http) {
            var proxys = {
                loadData: function(callback) {
                    var p = $http({
                        method: "GET",
                        url: "/params_ctrl"
                    });
                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            proxys.items = response;
                        }
                        if (callback){callback();}
                    });

                },
                add:function(){
                  this.items.unshift({name:"点击填写名称",categray:0,state:0})
                },
                delete: function(index) {
                    this.items.splice(index, 1);
                    this.save();

                },
                save: function(callback) {
                    var p = $http({
                        method: "POST",
                        data: proxys.items,
                        url: "/params_ctrl"
                    });

                    p.success(function(response, status, headers, config) {
                        if (response.length > 0) {
                            p.items = response;
                        }
                        if (callback){callback();}
                    });

                },
                items: [],
                catename:["连接可用代理" , "跳过代理"],
                cate:1, //#编辑 , 2添加
                index:0 
            }

            return proxys;
        }
    ]);


})();

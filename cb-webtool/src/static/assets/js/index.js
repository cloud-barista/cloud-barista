
        var wmts =  new ol.layer.Tile({
            source: new ol.source.XYZ({
                url: 'http://api.vworld.kr/req/wmts/1.0.0/E4A59B05-0CF4-3654-BD0C-A169F70CCB34/Base/{z}/{y}/{x}.png'
            })
        })
        var map = new ol.Map({
          target: 'map',
          layers: [wmts],
          view: new ol.View({
            center: ol.proj.transform([126.9380517322744,37.16792263658907], 'EPSG:4326', 'EPSG:900913'),
            zoom: 7
          })
        }); 
        var features = new Array();
        var styleCache = new Array();
        var search = function(){
            $.ajax({
                type: "get",
                url: "http://api.vworld.kr/req/search",
                data : $('#searchForm').serialize(),
                dataType: 'jsonp',
                async: false,
                success: function(data) {
                    for(var o in data.response.result.items){ 
                        if(o==0){
                            move(data.response.result.items[o].point.x*1,data.response.result.items[o].point.y*1);
                        }
                        //Feature 객체에 저장하여 활용 
                        features[o] = new ol.Feature({
                            geometry: new ol.geom.Point(ol.proj.transform([ data.response.result.items[o].point.x*1,data.response.result.items[o].point.y*1],'EPSG:4326', "EPSG:900913")),
                            title: data.response.result.items[o].title,
                            parcel: data.response.result.items[o].address.parcel,
                            road: data.response.result.items[o].address.road,
                            category: data.response.result.items[o].category,
                            point: data.response.result.items[o].point
                        });
                        features[o].set("id",data.response.result.items[o].id);
                          
                        var iconStyle = new ol.style.Style({
                            image: new ol.style.Icon(/** @type {olx.style.IconOptions} */ ({
                                anchor: [0.5, 10],
                                anchorXUnits: 'fraction',
                                anchorYUnits: 'pixels',
                                src: 'http://map.vworld.kr/images/ol3/marker_blue.png'
                            }))
                        });
                        features[o].setStyle(iconStyle);
                          
                    }
                      
                    var vectorSource = new ol.source.Vector({
                          features: features
                    });
                    /*
                        기존검색결과를 제거하기 위해 키 값 생성
                    */
                    var vectorLayer = new ol.layer.Vector({
                        source: vectorSource
                    });
                      
                    /*
                        기존검색결과를 제거하기 위해 키 값 생성
                    */
                    vectorLayer.set("vectorLayer","search_vector")
                      
                    map.getLayers().forEach(function(layer){
                        if(layer.get("vectorLayer")=="search_vector"){
                            map.removeLayer(layer);
                        }
                    });
                      
                    map.addLayer(vectorLayer);
                },
                error: function(xhr, stat, err) {}
            });
        }
          
        var move = function(x,y){//127.10153, 37.402566
            map.getView().setCenter(ol.proj.transform([ x, y ],'EPSG:4326', "EPSG:900913")); // 지도 이동
            map.getView().setZoom(12);
        }
          
        /* 클릭 이벤트 제어 */
        map.on("click", function(evt) {
            var coordinate = evt.coordinate //좌표정보
            var pixel = evt.pixel
            var cluster_features = [];
            var features = [];
              
            //선택한 픽셀정보로  feature 체크 
            map.forEachFeatureAtPixel(pixel, function(feature, layer) {
                var title = feature.get("title");
                if(title.length>0){
                      
                    var overlayElement= document.createElement("div"); // 오버레이 팝업설정 
                      
                    overlayElement.setAttribute("class", "overlayElement");
                    overlayElement.setAttribute("style", "background-color: #3399CC; border: 2px solid white; color:white");
                    overlayElement.setAttribute("onclick", "deleteOverlay('"+feature.get("id")+"')");
                    overlayElement.innerHTML="<p>"+title+"</p>";
                      
                    var overlayInfo = new ol.Overlay({
                        id:feature.get("id"),
                        element:overlayElement,
                        offset: [0, -70],
                        position: ol.proj.transform([feature.get("point").x*1, feature.get("point").y*1],'EPSG:4326', "EPSG:900913")
                    });
                      
                    if(feature.get("id") != null){
                        map.removeOverlay(map.getOverlayById(feature.get("id")));
                    }
                      
                    map.addOverlay(overlayInfo);
                }
            });
        });
          
        /**
            오버레이 삭제
        */
        var deleteOverlay = function(id){
            map.removeOverlay(map.getOverlayById(id));
        }
          
        init();
        function init(){
          var wms=new OpenLayers.Layer.WMS('USA Population','https://demo.geo-solutions.it/geoserver/wms',{layers:'topp:states'},{attribution:'<a target="_blank" href="https://demo.geo-solutions.it/geoserver" title="USA Population">USA Population</a>'}),
            vector=new OpenLayers.Layer.Vector(),
            map=new OpenLayers.Map('map',{
              center:[-98,38],
              zoom:5,
              layers:[wms,vector]
            });
          map.events.register('click',wms,GetFeatureInfoWMS);
          function GetFeatureInfoWMS(e){
            var me=this,
              xy=e.xy,
              map=me.map,
              loc=map.getLonLatFromPixel(xy),			
              bounds=map.getExtent(),
              size=map.getSize(),
              obj=OpenLayers.Util.getParameters(me.getURL(bounds)),			
              params ={
                REQUEST:'GetFeatureInfo',
                //VERSION:obj.VERSION,//comes from getMapUrl
                QUERY_LAYERS:obj.LAYERS,
                INFO_FORMAT:'application/json',
                FEATURE_COUNT:5,
                //EXCEPTIONS:'application/vnd.ogc.se_xml',//deafault
                //BBOX:bounds.toBBOX(),//comes from getMapUrl
                WIDTH:size.w,//WIDTH, comes from getMapUrl isn't working
                HEIGHT:size.h //HEIGHT, comes from getMapUrl isn't working
              };
            if(parseFloat(obj.VERSION)>=1.3){
              params.I=xy.x;params.J=xy.y;
            }else{
              params.X=xy.x;params.Y=xy.y;
            }
            OpenLayers.Util.applyDefaults(params,obj);
            OpenLayers.Request.GET({
              url:me.url,
              params:params,
              success:function(res){
                var dif=new OpenLayers.Format.GeoJSON(),
                  features=dif.read(res.responseText),
                  html='',
                  popup;
                vector.removeAllFeatures();
                if(features.length){
                  vector.addFeatures(features);
                  for(var k in features){
                    var feature=features[k],
                      attributes=feature.attributes;
                    html+='<br/><b>State Code: </b>'+attributes.STATE_ABBR+'<br/><b>State Name: </b>'+attributes.STATE_NAME+'<br/>';
                  }
                  popup=new OpenLayers.Popup.FramedCloud('popup',loc,size,html,null,true,
                    onPopupClose=function(e){
                      var me=this;
                      me.destroy();
                      vector.removeAllFeatures();
                      OpenLayers.Event.stop(e);
                    }
                  );
                }else{
                  html='No features found...';
                  popup=new OpenLayers.Popup.FramedCloud('popup',loc,size,html,null,true);
                }
                map.addPopup(popup,true);
              },
              failure:function(res){
                vector.removeAllFeatures();
                html='Unable to complete the request...:'+res.responseText;
                popup=new OpenLayers.Popup.FramedCloud('popup',loc,size,html,null,true);
                map.addPopup(popup,true);
              }
            });
          }
        }
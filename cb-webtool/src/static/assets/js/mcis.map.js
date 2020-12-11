//서버에서 처리 필요 없다.ㅜㅡ
function getIPStackRegion(ip_address){
    var apiUrl = "http://api.ipstack.com/"
    var access_key = "86c895286435070c0369a53d2d0b03d1"
    var url = apiUrl+ip_address+"?access_key="+access_key

    console.log("api get region url:",url);
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then((result)=>{
        console.log("api get result : ",result);
        var data = result.data
        var lat = data.latitude
        var long = data.longitude
        
    })
}
function viewMap(){
    var mcis_id = $("#mcis_id").val();
    $("#map_detail").show();
    $("#map2").empty();
    var map = map_init_target('map2')
    fnMove('map_detail');
    getGeoLocationInfo(mcis_id,map);
}
function getGeoLocationInfo(mcis_id,map){
  var JZMap = map;
  $.ajax({
    type:'GET',
    url: '/map/geo/'+mcis_id,
   // async:false,
    }).done(function(result){
        console.log("region Info : ",result)
        var polyArr = new Array();
  
        for(var i in result){
            console.log("region lat long info : ",result[i])
            // var json_parse = JSON.parse(result[i])
            // console.log("json_parse : ",json_parse.longitude)
            var long = result[i].longitude
            var lat = result[i].latitude
            var fromLonLat = long+" "+lat;
            polyArr.push(fromLonLat)
            drawMap(JZMap,long,lat,result[i])
        }
        var polygon = "";
        if(polyArr.length > 1){
          polygon = polyArr.join(", ")
          polygon = "POLYGON(("+polygon+"))";
        }else{
          polygon = "POLYGON(("+polyArr[0]+"))";
        }
       
        if(polyArr.length >1){
          drawPoligon(JZMap,polygon);
        }
        //drawPoligon(map,wkt);
       
    })
}

function map_init_target(target){
 
  const osmLayer = new ol.layer.Tile({
    source: new ol.source.OSM(),
  });
  

var m = new ol.Map({
    target: target,
    layers: [
      osmLayer
    ],
    view: new ol.View({
      center: [0,0],
      zoom: 1
    })
  });
   var JZMap = m;
;
  JZMap.on('click',function(evt){
   // var element = document.getElementById('map_pop2');
    var element = document.createElement('div');
    var feature = JZMap.forEachFeatureAtPixel(evt.pixel,function(feature){
      return feature;
    })
    console.log("feature click info : ",feature.get("vm_id"));
   
    

   
    if(feature){
      var coordinates = feature.getGeometry().getCoordinates();
      
      // $(element).html('<div class="popover" role="tooltip"><div class="arrow"></div><h3 class="popover-header"></h3><div class="popover-body"></div></div>');
     
     
      element.setAttribute("class", "popover");
      
      
      //element.setAttribute("onclick", "deleteOverlay('"+feature.get("vm_id")+"')");
      element.setAttribute("onclick", "$(this).hide()");
     // element.innerHTML="<p>"+feature.get("title")+"</p>";
      // element.innerHTML='<div tabindex="0" class="btn btn-lg btn-danger" role="button" data-toggle="popover" data-trigger="focus" title="Dismissible popover" data-content="And here\'s some amazing content. It\'s very engaging. Right?">Dismissible popover</a>';
     
      
      // $(element).empty()
      // $(element).show()    

      $(element).popover({
        placement: 'auto',
        html: true,
        content: "<div onclick='alert(\"Hello\")'><p>ID : "+feature.get('vm_id')+"</p>"+"Status :"+feature.get('vm_status')+"</div>",
        title: feature.get('title'),
        trigger:'click',
      });

      var popup = new ol.Overlay({
        element: element,
        id:feature.get("id"),
        positioning: 'bottom-center',
        stopEvent: false,
        offset: [0, -50]
      });
      popup.setPosition(coordinates);
      // var popup = new ol.Overlay({
      //   element: element,
      //   id:feature.get("vm_id"),
      //   positioning: 'bottom-center',
      //   position: coordinates,
      //   offset: [0, -70],
      //   stopEvent: false,
      //   offset: [0, -50]
      // });
      
       
      JZMap.addOverlay(popup);
      
      $(element).popover('toggle');
    }else{
      $(element).popover('hide');
    }
  });
 
return m;
}
function map_init(){
 
    const osmLayer = new ol.layer.Tile({
      source: new ol.source.OSM(),
    });

  var m = new ol.Map({
      target: 'map',
      layers: [
        osmLayer
      ],
      view: new ol.View({
        center: ol.proj.fromLonLat([37.41, 8.82]),
        zoom: 0
      })
    });
    var JZMap = m;

    
    

  //   var deleteOverlay = function(id){
  //     JZMap.removeOverlay(JZMap.getOverlayById(id));
  // }
 
  // JZMap.on('pointermove', function(e) {
  //   var feature = JZMap.forEachFeatureAtPixel(e.pixel,function(feature){
  //     return feature;
  //   })
  //   if (e.dragging) {
  //     var element = feature.get("element");
  //     var id = feature.get("id");
  //     $(element).popover('hide');
  //     JZMap.removeOverlay(JZMap.getOverlayById(id));
  //     return;
  //   }
  //   var pixel = JZMap.getEventPixel(e.originalEvent);
  //   var hit = JZMap.hasFeatureAtPixel(pixel);
  
  // });
    JZMap.on('click',function(evt){
     
      var feature = JZMap.forEachFeatureAtPixel(evt.pixel,function(feature){
        return feature;
      })
      
      console.log("feature click info : ",feature.get("id"));
      var id = feature.get("id")
      if(feature.get("id") != null){
        JZMap.removeOverlay(JZMap.getOverlayById(id));
      }
      if(feature){
        var element = document.createElement('div');
      element.setAttribute("class", "popover");
      element.setAttribute("onclick", "$(this).hide()");
      element.innerHTML="<div data-toggle='popover' style='width:100%;'>"+feature.get("title")+"</div>"
      
      

      
      
      feature.set("element");
  
      var popup = new ol.Overlay({
        element: element,
        positioning: 'auto',
        stopEvent: false,
        offset: [0, 0],
        id:feature.get("id")
      });
        var coordinates = feature.getGeometry().getCoordinates();
        popup.setPosition(coordinates);
        // $(element).html('<div class="popover" role="tooltip"><div class="arrow"></div><h3 class="popover-header"></h3><div class="popover-body"></div></div>');
       
       
        // element.setAttribute("class", "overlayElement");
        // element.setAttribute("id", feature.get('vm_id'));
        // element.setAttribute("onclick", "deleteOverlay('"+feature.get("vm_id")+"')");
       
        
        // $(element).empty()
        // $(element).show()    

        // $(element).popover({
        //   placement: 'top',
        //   html: true,
        //   content: "ID : "+feature.get('vm_id')+"\n"+"Status :"+feature.get('vm_status'),
        //   title: feature.get('title'),
        // });

        
        // var popup = new ol.Overlay({
        //   element: element,
        //   id:feature.get("vm_id"),
        //   positioning: 'bottom-center',
        //   position: coordinates,
        //   offset: [0, -70],
        //   stopEvent: false,
        //   offset: [0, -50]
        // });
        
         
        JZMap.addOverlay(popup);
        
        $(element).popover('show');
      }else{
        $(element).popover('hide');
      }
    });
  return m;
}
function drawMap(map,long,lat,info){
  var JZMap = map;
  console.log("JZMap : ",JZMap);
  
  var icon = new ol.style.Style({
    image: new ol.style.Icon({
        src:'/assets/img/marker/purple.png', // pin Image
        anchor: [0.5, 1],
        scale: 0.5
    
    })
})
  var map_center = ol.proj.fromLonLat([long, lat]);
  var point_gem = new ol.geom.Point(map_center);
  var point_feature = new ol.Feature(point_gem);
  point_feature.setStyle([icon])
  //feature 에 set info
  console.log("info : ",info)
  point_feature.set('title',info.name)
  point_feature.set('vm_status',info.status)
  point_feature.set('vm_id',info.id)
  point_feature.set('id',info.id)

  var stackVectorMap = new ol.source.Vector({
    features : [point_feature]
  })

  var stackLayer = new ol.layer.Vector({
    source: stackVectorMap
  })
  JZMap.addLayer(stackLayer)
 
}

function drawPoligon(JZMap,polygon){
  var wkt = polygon;
  console.log("polygon : ",wkt)
  var format = new ol.format.WKT();

  var feature = format.readFeature(wkt, {
    dataProjection: "EPSG:4326",
    featureProjection: "EPSG:3857"
  });
  var stackVectorMap = new ol.source.Vector({
    features : [feature]
  })

  var stackLayer = new ol.layer.Vector({
    source: stackVectorMap
  })
  JZMap.addLayer(stackLayer);
  
}



function escapeXml(string) {
    return string.replace(/[<>]/g, function (c) {
      switch (c) {
        case '<': return '\u003c';
        case '>': return '\u003e';
      }
    });
  }


//서버에서 처리 필요 없다.ㅜㅡ
function getIPStackRegion(ip_address) {
    var apiUrl = "http://api.ipstack.com/"
    var access_key = "86c895286435070c0369a53d2d0b03d1"
    var url = apiUrl + ip_address + "?access_key=" + access_key

    console.log("api get region url:", url);
    var apiInfo = ApiInfo
    axios.get(url, {
        headers: {
            'Authorization': apiInfo
        }
    }).then((result) => {
        console.log("api get result : ", result);
        var data = result.data
        var lat = data.latitude
        var long = data.longitude

    })
}
function viewMap() {
    var mcis_id = $("#mcis_id").val();
    $("#map_detail").show();
    $("#map2").empty();
    var map = map_init_target('map2')
    fnMove('map_detail');
    getGeoLocationInfo(mcis_id, map);
}
function getGeoLocationInfo(mcis_id, map) {
    var JZMap = map;
    $.ajax({
        type: 'GET',
        url: '/map/geo/' + mcis_id,
        // async:false,
    }).done(function (result) {
        console.log("region Info : ", result)
        var polyArr = new Array();

        for (var i in result) {
            console.log("region lat long info : ", result[i])
            // var json_parse = JSON.parse(result[i])
            // console.log("json_parse : ",json_parse.longitude)
            var long = result[i].longitude
            var lat = result[i].latitude
            var fromLonLat = long + " " + lat;
            polyArr.push(fromLonLat)
            drawMap(JZMap, long, lat, result[i])
        }
        var polygon = "";
        if (polyArr.length > 1) {
            polygon = polyArr.join(", ")
            polygon = "POLYGON((" + polygon + "))";
        } else {
            polygon = "POLYGON((" + polyArr[0] + "))";
        }

        if (polyArr.length > 1) {
            drawPoligon(JZMap, polygon);
        }
        //drawPoligon(map,wkt);

    })
}

function map_init_target(target) {

    const osmLayer = new ol.layer.Tile({
        source: new ol.source.OSM(),
    });


    var m = new ol.Map({
        target: target,
        logo: false,
        layers: [
            osmLayer
        ],
        view: new ol.View({
            center: [0, 0],
            zoom: 1
        })
    });
    var JZMap = m;
    ;

    JZMap.on('click', function (evt) {
        // var element = document.getElementById('map_pop2');
        var element = document.createElement('div');
        var feature = JZMap.forEachFeatureAtPixel(evt.pixel, function (feature) {
            return feature;
        })
        //console.log("feature click info : ", feature.get("vm_id"));

        if (feature) {
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
                content: "<div onclick='alert(\"Hello\")'><p>ID : " + feature.get('vm_id') + "</p>" + "Status :" + feature.get('vm_status') + "</div>",
                title: feature.get('title'),
                trigger: 'click',
            });

            var popup = new ol.Overlay({
                element: element,
                id: feature.get("id"),
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
        } else {
            $(element).popover('hide');
        }
    });

    return m;
}

function map_init() {

    const osmLayer = new ol.layer.Tile({
        source: new ol.source.OSM(),
    });
    var control = new ol.control.FullScreen();
    var m = new ol.Map({
        target: 'map',
        logo: false,
        // controls: ol.control.defaults().extend([
        //     new ol.control.FullScreen()
        // ]),
        controls: [control],
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
    JZMap.on('click', function (evt) {
        var selectedFeature
        var selectedFeature = JZMap.forEachFeatureAtPixel(evt.pixel, function (feature) {
            return feature;
        })
        if (selectedFeature) {
            console.log(selectedFeature)
            console.log("feature click info : ", selectedFeature.get("id"));//Cannot read property 'get' of undefined at e.<anonymous> (mcis.map.js:196)
            var overlayId = selectedFeature.get("id")
            if (selectedFeature.get("id") == undefined) { return; }
            //     JZMap.removeOverlay(JZMap.getOverlayById(overlayId));
            // }
            JZMap.getOverlays().forEach(function (overlay) {
                JZMap.removeOverlay(overlay);
            });

            var element = document.createElement('div');
            element.setAttribute("class", "popover");
            element.setAttribute("onclick", "$(this).hide()");
            element.innerHTML = "<div data-toggle='popover' style='width:100%;min-width: 100px;'>" + selectedFeature.get("title") + "</div>"

            selectedFeature.set("element");

            var popup = new ol.Overlay({
                element: element,
                positioning: 'auto',
                stopEvent: false,
                offset: [0, 0],
                id: selectedFeature.get("id")
            });
            var coordinates = selectedFeature.getGeometry().getCoordinates();
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
        } else {
            $(element).popover('hide');
        }
    });
    return m;
}

// 주어진 index에 맞는 marker표시.
function getMarkerSrc(markerIndex) {
    var markerSrc = ""
    console.log("markerIndex " + markerIndex)
    if (markerIndex == undefined) {
        markerIndex = 1
    }
    var remainder = markerIndex % 10
    console.log("markerIndex " + markerIndex + " : " + remainder)
    switch (remainder) {
        case 0:
            markerSrc = "/assets/img/marker/purple.png"
            break;
        case 1:
            markerSrc = "/assets/img/marker/blue.png"
            break;
        case 2:
            markerSrc = "/assets/img/marker/green.png"
            break;
        case 3:
            markerSrc = "/assets/img/marker/grey.png"
            break;
        case 4:
            markerSrc = "/assets/img/marker/orange.png"
            break;
        case 5:
            markerSrc = "/assets/img/marker/black.png"
            break;
        case 6:
            markerSrc = "/assets/img/marker/red.png"
            break;
        case 7:
            markerSrc = "/assets/img/marker/white.png"
            break;
        case 8:
            markerSrc = "/assets/img/marker/yellow.png"
            break;
        default:
            markerSrc = "/assets/img/marker/black.png"
            break;
    }
    return markerSrc
}
function drawMap(map, long, lat, info) {
    console.log("in drawMap")

    var JZMap = map;
    console.log("JZMap : ", JZMap);

    var icon = new ol.style.Style({
        image: new ol.style.Icon({
            // src:'/assets/img/marker/purple.png', // pin Image
            src: getMarkerSrc(info.markerIndex),
            anchor: [0.5, 1],
            scale: 0.5
        })
    })

    var map_center = ol.proj.fromLonLat([long, lat]);
    var point_gem = new ol.geom.Point(map_center);
    var point_feature = new ol.Feature(point_gem);
    point_feature.setStyle([icon])

    //feature 에 set info
    console.log("info : ", info)
    point_feature.set('title', info.name)
    point_feature.set('vm_status', info.status)
    point_feature.set('vm_id', info.id)
    point_feature.set('id', info.id)

    var stackVectorMap = new ol.source.Vector({
        features: [point_feature]
    })

    var stackLayer = new ol.layer.Vector({
        source: stackVectorMap
    })
    JZMap.addLayer(stackLayer)
    layersMap.set("pin_" + info.pinIndex, stackLayer)

}

function drawPoligon(JZMap, polygon, polygonId, colorIndex) {
    var wkt = polygon;
    console.log("polygon : ", wkt)
    var format = new ol.format.WKT();

    var polyColor = 'rgba(100, 100, 100, 0.1)';
    var polyLineColor = 'fc8d10';
    if (colorIndex) {
        switch (colorIndex) {
            case 0:
                polyColor = 'rgba(100, 0, 0, 0.1)';
                polyLineColor = '#cb00f5';//purple
                break;
            case 1:
                polyColor = 'rgba(0, 255, 0, 0.1)';;
                polyLineColor = '#1e2b67';//blue
                break;
            case 2:
                polyColor = 'rgba(0, 0, 1, 0.1)';
                polyLineColor = '#86b049';//green
                break;
            case 3:
                polyColor = 'rgba(255, 1, 0, 0.1)';
                polyLineColor = '#545b62';//grey
                break;
            case 4:
                polyColor = 'rgba(255, 1, 1, 0.1)';
                polyLineColor = '#ff6700';//orange
                break;
            case 5:
                polyColor = 'rgba(0, 255, 0, 0.1)';
                polyLineColor = '#f1dddf';//black
                break;
            case 6://OK
                polyColor = 'rgba(255, 2, 1, 0.1)';
                polyLineColor = '#bf3c41';//red
                break;
            case 7:
                polyColor = 'rgba(255, 3, 2, 0.1)';
                polyLineColor = '#754000';//white
                break;
            case 8:
                polyColor = 'rgba(255, 4, 3, 0.1)';
                polyLineColor = '#ffcb17';//yellow
                break;

            default:
                polyColor = 'rgba(255, 0, 0, 0.1)';
                polyLineColor = '#fc8d10';
        }
    }

    var feature = format.readFeature(wkt, {
        dataProjection: "EPSG:4326",
        featureProjection: "EPSG:3857"
    });

    var stackVectorMap = new ol.source.Vector({
        features: [feature]
    })

    console.log("colorIndex = " + colorIndex + ", polyColor = " + polyColor)
    var styles = new ol.style.Style({
        fill: new ol.style.Fill({
            color: polyColor,
            weight: 1
        }),
        stroke: new ol.style.Stroke({
            color: polyLineColor,
            width: 1
        })
    });
    //   var styles = [ new ol.style.Style({ stroke: new ol.style.Stroke({ color: polyColor, width: 3, }), }) ];
    // var stackLayer = new ol.layer.Vector({
    //   source: stackVectorMap,
    //     style: styles
    // })

    var stackLayer = new ol.layer.Vector({
        source: stackVectorMap,
        style: styles
    })
    JZMap.addLayer(stackLayer);
    layersMap.set("polygon_" + polygonId, stackLayer)
}

// 지도에 그려진 polygon 초기화, exceptionIndex가 -1 이면 전체 clear, 아니면 해당 index 배고 clear
var layersMap = new Map();//polygon layer들이 생성되면 넣는 map
//function clearPolygon(JZMap){
function clearLayers(JZMap) {
    console.log("clearPolygon");

    layersMap.forEach((value, key) => {
        console.log("clearPolygon = " + key)
        console.log(value)
        JZMap.removeLayer(value);
    })
    // console.log("clearPolygon ")
    // console.log(mapLayerMap);
    // for(let key in mapLayerMap){
    //     layerIndex = mapLayerMap[key];
    //     console.log("layerIndex " + layerIndex)
    //     if( exceptionIndex > -1){
    //         if( layerIndex == exceptionIndex )continue;
    //
    //         JZMap.getLayers().getArray()[layerIndex].getSource().clear();
    //     }else{
    //         JZMap.getLayers().getArray()[layerIndex].getSource().clear();
    //     }
    // }

}


function escapeXml(string) {
    return string.replace(/[<>]/g, function (c) {
        switch (c) {
            case '<': return '\u003c';
            case '>': return '\u003e';
        }
    });
}


function addPin(map, long, lat) {
    var JZMap = map;
    console.log("JZMap : ", JZMap);

    var icon = new ol.style.Style({
        image: new ol.style.Icon({
            // src:'/assets/img/marker/purple.png', // pin Image
            src: getMarkerSrc(0), // 0은 임시,  markerIndex에 따라 다른 색의 pin image표시됨
            anchor: [0.5, 1],
            scale: 0.5
        })
    })

    var map_center = ol.proj.fromLonLat([long, lat]);
    var point_gem = new ol.geom.Point(map_center);
    var point_feature = new ol.Feature(point_gem);
    point_feature.setStyle([icon])

    var stackVectorMap = new ol.source.Vector({
        features: [point_feature]
    })

    var stackLayer = new ol.layer.Vector({
        source: stackVectorMap
    })
    JZMap.addLayer(stackLayer)
}


function addClickPin(map) {
    map.on('click', function (evt) {
        var lonlat = ol.proj.transform(evt.coordinate, 'EPSG:3857', 'EPSG:4326');
        addPin(map, lonlat[0], lonlat[1])
        $("#longitude").val(lonlat[0])
        $("#latitude").val(lonlat[1])
    });
}


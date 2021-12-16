var config = {
    type: 'line',
    data: {
        labels:[] ,// 시간을 배열로 받아서 처리
        datasets: [{
            label : "",//cpu 관련 내용들 
            //backgroundColor:window.chartColors.red,
           // borderColor:window.chartColors.red,
            data:[],//
        }]

    }
}

function showMonitoring(mcis_id, vm_id){
    $("#cpu").empty()
    $("#memory").empty()
    $("#disk").empty()
    $("#network").empty()
     var arr = ["cpu","memory","disk","network"];
     var periodType = "m";
     var duration = "10m";
     var statisticsCriteria = "last";
     
     for (var i in arr){
       var chart_target = "canvas_"+i;
        getMetric(chart_target,arr[i],mcis_id,vm_id,arr[i],periodType,statisticsCriteria,duration);
  
     }
}

function show_monitoring(){
    var mcis_id = $("#mcis_id").val();
    var vm_id = $("#current_vmid").val();
    var checkDragonValue = $("#check_dragonFly").val();
    var public_ip = $("#current_publicIP").val();

    if(checkDragonValue == "200"){
        showMonitoring(mcis_id,vm_id);
    }else{
        agentSetup(mcis_id,vm_id,public_ip);
    }
    
}

var vmChart;
function showMonitoring(mcis_id, vm_id, metric, periodType, duration){
	// $("#cpu").empty()
	// $("#memory").empty()
	// $("#disk").empty()
	// $("#network").empty()
	$("#canvas_vm").empty();
	var statisticsCriteria = "last";
    
	getVmMetric(vmChart,"canvas_vm",metric,mcis_id,vm_id,metric,periodType,statisticsCriteria,duration);
}
function genChartFmt(chart_target){

	var ctx = document.getElementById(chart_target).getContext('2d')
	var chart = new Chart(ctx,{
		type:"line",
		data:{},
		options:{
			responsive: true,
			title: {
				display: true,
				text: target
			},
			tooltips: {
				mode: 'index',
				intersect: false,
			},
			hover: {
				mode: 'nearest',
				intersect: true
			},
			scales: {
				x: {
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Time'
					}
				},
				y: {
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Value'
					}
				}
			}
		}
	});	

	return chart;
}

function getMetric(chart_target,target, mcis_id, vm_id, metric, periodType,statisticsCriteria, duration){
	console.log("====== Start GetMetric ====== ")
	var color = "";
    var metric_size ="";
    
	
	
    var nsid = NAMESPACE;
    console.log("get metric namespace : ",nsid);
	var url = DragonFlyURL+"/ns/"+nsid+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
    console.log("Request URL : ",url)
    
    var ctx = document.getElementById(chart_target).getContext('2d')
    var chart = new Chart(ctx,{
        type:"line",
        data:{},
        options:{
            responsive: true,
            title: {
                display: true,
                text: target
            },
            tooltips: {
                mode: 'index',
                intersect: false,
            },
            hover: {
                mode: 'nearest',
                intersect: true
            },
            scales: {
                x: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Time'
                    }
                },
                y: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Value'
                    }
                }
            }
        }
    });
	chart.clear()
	$.ajax({
		url: url,
		async:false,
		type:'GET',
		success : function(result){
			var data = result
			console.log("Get Monitoring Data : ",data)
			console.log("info items : ", target);
            console.log("======== start mapping data ======");
            $("#"+chart_target).empty();
           
    
            //data sets
            var key =[]
            var values = data.values[0]
            for(var i in values){                
                key.push(i)
            }
            console.log("Key values time except:",key);
	
            var labels = key;
            var datasets = data.values;
            // 각 값의 배열 데이터들
            //console.log("info labels : ",labels);
            console.log("info datasets : ",datasets);

            var obj = {}
            obj.columns = labels
	        obj.values = datasets

			var time_obj = time_arr(obj,target);
			console.log("chart_target :",chart_target);
			console.log("info datasets : ", time_obj);
			
			// var myChart = new Chart(ctx, time_obj);
			chart.data = time_obj;
			chart.update();
            
		},
		error : function(request,status, error){
			console.log(request.status, request.responseText,error)
		}
		
	})
}

// function getMetric(chart_target,target, mcis_id, vm_id, metric, periodType,statisticsCriteria, duration){
//     console.log("====== Start GetMetric ====== ")
  
//    var ctx = document.getElementById(chart_target).getContext('2d')
//    var chart = new Chart(ctx,{
//        type:"line",
//        data:{},
//        options:{
//         responsive: true,
//         title: {
//             display: true,
//             text: target
//         },
//         tooltips: {
//             mode: 'index',
//             intersect: false,
//         },
//         hover: {
//             mode: 'nearest',
//             intersect: true
//         },
//         scales: {
//             x: {
//                 display: true,
//                 scaleLabel: {
//                     display: true,
//                     labelString: 'Time'
//                 }
//             },
//             y: {
//                 display: true,
//                 scaleLabel: {
//                     display: true,
//                     labelString: 'Value'
//                 }
//             }
//         }
//     }
//    });
   
   
//    var url = DragonFlyURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
 
//    console.log("Request URL : ",url)
//    var html = "";
//    var apiInfo = ApiInfo;
//    $.ajax({
//    url: url,
//    async:false,
//    type:'GET',
//    beforeSend : function(xhr){
//     xhr.setRequestHeader("Authorization", apiInfo);
//     xhr.setRequestHeader("Content-type","application/json");
// },
//    success : function(result){
//        var data = result
//          console.log("Get Monitoring Data : ",data)
//          console.log("======== start mapping data ======");
//          var time_obj = time_arr(data,target);
//          console.log("chart_target :",chart_target);
        
//         // var myChart = new Chart(ctx, time_obj);
//         chart.data = time_obj;
//         chart.update();
//         $("#chart_detail").show();
//         fnMove('chart_detail');

//    },
//    error : function(request,status, error){
//        console.log("ERROR request status at DragonFly : ",status);
     
       
       
//    }
   
// })
//    // var apiInfo = ApiInfo
//     // axios.get(url,{
//     //     headers:{
//     //         'Authorization': apiInfo
//     //     }
//     // })then(result=>{
//    //       var data = result.data
//    //       console.log("Get Monitoring Data : ",data)
//    //       console.log("======== start mapping data ======");
//    //       var time_obj = time_arr(data,target);
//    //       console.log("chart_target :",chart_target);
//    //       var ctx = document.getElementById(chart_target).getContext('2d')
//    //       var myChart = new Chart(ctx, time_obj);
//    //       myChart.update();
//    //      // Chart.Line('canvas1',time_obj);
//    //       console.log("==time series==",time_obj);
//    //      // var metricObject = mappingMetric(data);
//    //      // var m_len = metricObject.length-1

//    //     //   if(target == "cpu"){
//    //     //     color += '<div class="icon icon-shape bg-success text-white rounded-circle shadow">'
//    //     //           + '<i class="fas fa-microchip"></i>';
//    //     //     var num = parseFloat(metricObject[m_len].cpu_utilization)
//    //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>'                  
//    //     //   }else if(target == "memory"){
//    //     //     color += '<div class="icon icon-shape bg-warning text-white rounded-circle shadow">'
//    //     //           + '<i class="fas fa-memory"></i>'
//    //     //     var num = parseFloat(metricObject[m_len].mem_utilization)
//    //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>' 
//    //     //   }else if(target == "disk"){
//    //     //     color += '<div class="icon icon-shape bg-danger text-white rounded-circle shadow">'
//    //     //           + '<i class="far fa-save"></i>'
//    //     //     var num = parseFloat(metricObject[m_len].used_percent)
//    //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>' 
//    //     //   }else if(target == "network"){
//    //     //     color += '<div class="icon icon-shape bg-primary text-white rounded-circle shadow">'
//    //     //           + '<i class="fas fa-network-wired"></i>'
//    //     //     var num = parseFloat(metricObject[m_len].bytes_in)
//    //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(1)+'byte</span>' 
//    //     //   }
//    //     //         html += '<div class="card card-stats mb-4 mb-xl-0">'
//    //     //              +'<div class="card-body">'
//    //     //              +'<div class="row">'
//    //     //              +'<div class="col">'
//    //     //              +'<h5 class="card-title text-uppercase text-muted mb-0">'+metric+'</h5>'
//    //     //              //+'<span class="h2 font-weight-bold mb-0">2,356</span>'
//    //     //              +metric_size
//    //     //              +'</div>'
//    //     //              +'<div class="col-auto">'
                  
//    //     //              +color
                    
//    //     //              //+'<i class="fas fa-chart-pie"></i>'
//    //     //              +'</div>'
//    //     //              +'</div>'
//    //     //              +'</div>'
//    //     //              +'<p class="mt-3 mb-0 text-muted text-sm">'
//    //     //              +'<span class="text-danger mr-2"> 3.48%</span>'
//    //     //              +'<span class="text-nowrap">'+metricObject[0].time+'</span>'
//    //     //              +'</p>'
//    //     //              +'</div>'
//    //     //              +'</div>';
         
//    //     //   $("#"+target+"").empty()
//    //     //   $("#"+target+"").append(html)
//    //   })
// }

function checkDragonFly(mcis_id, vm_id){
   console.log("====== Start Check DragonFly ====== ")
   var periodType = "m";
   var duration = "10m";
   var statisticsCriteria = "last";
   var metric = "cpu" 
   var apiInfo = ApiInfo;
   var url = DragonFlyURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
   
   console.log("Request URL : ",url)
   
   $.ajax({
        url: url,
        async:false,
        type:'GET',
        beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
        success : function(result){
            console.log("check dragon fly : ",result)
          //  $("#check_dragonFly").val("200");
            $("#mcis_detail_info_check_monitoring").prop("checked",true)
            $("#mcis_detail_info_check_monitoring").attr("disabled",true)
            $("#Monitoring_tab").show();
            var duration = "5m"
            var period_type = "m"
            var metric_arr = ["cpu","memory","disk","network"];
            var statisticsCriteria = "last";
            for(var i in metric_arr){
                getMetric("canvas_"+i,metric_arr[i],mcis_id,vm_id,metric_arr[i],period_type,statisticsCriteria,duration);
            }
        },
        error : function(request,status, error){
            console.log("check dragon fly : ",status)
           // $("#check_dragonFly").val("400");
            $("#mcis_detail_info_check_monitoring").prop("checked",false)
            $("#mcis_detail_info_check_monitoring").attr("disabled",false)
            $("#Monitoring_tab").hide();
            
        }          
    })
   
}

function checkDragonFly2(mcis_id, vm_id){
    console.log("====== Start Check DragonFly ====== ")
    var periodType = "m";
    var duration = "10m";
    var statisticsCriteria = "last";
    var metric = "cpu" 
    var apiInfo = ApiInfo;
    var url = DragonFlyURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
    
    console.log("Request URL : ",url)
    
    $.ajax({
         url: url,
         async:false,
         type:'GET',
         beforeSend : function(xhr){
             xhr.setRequestHeader("Authorization", apiInfo);
             xhr.setRequestHeader("Content-type","application/json");
         },
         success : function(result){
             
           //  $("#check_dragonFly").val("200");
           var input_duration = $("#input_duration").val();
           var duration_type = $("#duration_type").val();
           var duration = input_duration+duration_type
           var period_type = $("#vm_period").val();
           var metric = $("#select_metric").val();
           showMonitoring(mcis_id,vm_id,metric,period_type,duration);
         },
         error : function(request,status, error){
         
            // $("#check_dragonFly").val("400");
          alert("It is Not installed Monitoring Agent!!");
             
         }          
     })
    
 }



function time_arr(obj, title){
    //data sets
    console.log("labels:",obj)
    console.log("")
   var labels = obj.columns;
   var datasets = obj.values;
   
    // 각 값의 배열 데이터들
   var series_label = new Array();
   var data_set = new Array();
   for(var i in labels){
       var ky = labels[i]
       var series_data = new Array(); 
       if(ky == "time"){
        for(var k in datasets){
            for(var o in datasets[k]){
                if(o == ky){
                    series_label.push(datasets[k][o])
                }
            }
          }

       }else{
        
        var dt = {}
        
        dt.label = ky
        var color1 = Math.floor(Math.random() * 256);
        var color2 = Math.floor(Math.random() * 256);
        var color3 = Math.floor(Math.random() * 256);
        var color = 'rgb('+color1+","+color2+","+color3+")"
        dt.borderColor = color
        dt.backgroundColor = color;

      
      
       dt.fill= false;
           for(var k in datasets){
             for(var o in datasets[k]){
                 if(o == ky){
                   series_data.push(datasets[k][o])
                 }
             }
           }
        dt.data = series_data
        data_set.push(dt)
       }
       
    }
  var new_obj = {};
   console.log("data set : ",data_set);
   console.log("time series : ",series_label);
   new_obj.labels = series_label //시간만 담김 배열
   new_obj.datasets =  data_set//각 데이터 셋의 배열
   console.log("Chart Object : ",new_obj);
   config.type = 'line',
   config.data = new_obj
   config.options = {
    responsive: true,
    title: {
        display: true,
        text: title
    },
    tooltips: {
        mode: 'index',
        intersect: false,
    },
    hover: {
        mode: 'nearest',
        intersect: true
    },
    scales: {
        x: {
            display: true,
            scaleLabel: {
                display: true,
                labelString: 'Time'
            }
        },
        y: {
            display: true,
            scaleLabel: {
                display: true,
                labelString: 'Value'
            }
        }
    }
}
   return new_obj;
}

window.chartColors = {
	red: 'rgb(255, 99, 132)',
	orange: 'rgb(255, 159, 64)',
	yellow: 'rgb(255, 205, 86)',
	green: 'rgb(75, 192, 192)',
	blue: 'rgb(54, 162, 235)',
	purple: 'rgb(153, 102, 255)',
    grey: 'rgb(201, 203, 207)',
    mint: 'rgb(99, 255, 243)'
};


//vm 의 통계조회
function getVmMetric(vmChart, chartTarget,target, mcisID, vmID, metric, periodType,statisticsCriteria, duration){
	console.log("====== Start GetMetric ====== ")
	var color = "";
    var metric_size ="";

    if( vmChart){
        vmChart.destroy();
    }
    vmChart = setVmChartInit(vmChart, chartTarget, target);
    // var vmChart = setVmChart(chartTarget,target);
	// vmChart.clear()
    
	var url = "/operation/manages/mcismng/proc/vmmonitoring"    
    console.log("Request URL : ",url)
    axios.post(url,{
        headers: { },
        mcisID:mcisID,
        vmID:vmID,
        metric:metric,
        periodType:periodType,
        statisticsCriteria:statisticsCriteria,
        duration:duration
    }).then(result=>{    
        console.log(result)    

        var statusCode = result.data.status;
        var message = result.data.message;
        
        if( statusCode != 200 && statusCode != 201) {
            commonAlert(message +"(" + statusCode + ")");
            return;
        }

        var data = result.data.VMMonitoringInfo
        console.log("Get Monitoring Data : ",data)
        console.log("info items : ", target);
        console.log("======== start mapping data ======");
        $("#"+chartTarget).empty();       

        //data sets
        var key =[]
        var values = data.values[0]
        for(var i in values){                
            key.push(i)
        }
        console.log("Key values time except:",key);

        var labels = key;
        var datasets = data.values;
        // 각 값의 배열 데이터들
        //console.log("info labels : ",labels);
        console.log("info datasets : ",datasets);

        var obj = {}
        obj.columns = labels
        obj.values = datasets

        var timeObj = xAxisSet(obj,target);
        console.log("chart_target :",chartTarget);
        console.log("info datasets : ", timeObj);			
        
        vmChart.data = timeObj;
        vmChart.update();
    // }).catch(function(error){
    //     var statusCode = error.response.data.status;
    //     var message = error.response.data.message;
    //     commonErrorAlert(statusCode, message)        
    // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)

        try{
            var statusCode = error.response.data.status;
            var errorMessage = error.response.data.error;
            commonErrorAlert(statusCode, errorMessage + " " + metric + " 조회실패") 
        }catch(e){
            var statusCode1 = error.response.status;
            var errorMessage1 = error.response.statusText;
            commonErrorAlert(statusCode1, errorMessage1 + " " + metric + " 조회실패") 
        }
    });
	
}

function setVmChartInit(vmChart, chartTarget,target){
    var ctx = document.getElementById(chartTarget).getContext('2d')
    vmChart = new Chart(ctx,{
        type:"line",
        data:{},
        options:{
            responsive: true,
            title: {
                display: true,
                text: target
            },
            tooltips: {
                mode: 'index',
                intersect: false,
            },
            hover: {
                mode: 'nearest',
                intersect: true
            },
            scales: {
                x: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Time'
                    }
                },
                y: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Value'
                    }
                }
            }
        }
    });
    return vmChart;
}

// 이전버전 : vmChart를 매번 생성하여 return
function setVmChart(chartTarget,target){
    var ctx = document.getElementById(chartTarget).getContext('2d')
    var vmChart = new Chart(ctx,{
        type:"line",
        data:{},
        options:{
            responsive: true,
            title: {
                display: true,
                text: target
            },
            tooltips: {
                mode: 'index',
                intersect: false,
            },
            hover: {
                mode: 'nearest',
                intersect: true
            },
            scales: {
                x: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Time'
                    }
                },
                y: {
                    display: true,
                    scaleLabel: {
                        display: true,
                        labelString: 'Value'
                    }
                }
            }
        }
    });
    return vmChart;
}

// x축 설정
function xAxisSet(obj, title){
    //data sets
    console.log("labels:",obj)
    console.log("")
    var labels = obj.columns;
    var datasets = obj.values;

    // 각 값의 배열 데이터들
    var series_label = new Array();
    var data_set = new Array();
    for(var i in labels){
        var ky = labels[i]
        var series_data = new Array(); 
        if(ky == "time"){
            for(var k in datasets){
                for(var o in datasets[k]){
                    if(o == ky){
                        series_label.push(datasets[k][o])
                    }
                }
             }

        }else{
        
            var dt = {}

            dt.label = ky
            var color1 = Math.floor(Math.random() * 256);
            var color2 = Math.floor(Math.random() * 256);
            var color3 = Math.floor(Math.random() * 256);
            var color = 'rgb('+color1+","+color2+","+color3+")"
            dt.borderColor = color
            dt.backgroundColor = color;      
      
            dt.fill= false;
            for(var k in datasets){
                for(var o in datasets[k]){
                    if(o == ky){
                       series_data.push(datasets[k][o])
                    }
                }
            }
            dt.data = series_data
            data_set.push(dt)
        }       
    }// end of for
    
    var newObj = {};
    console.log("data set : ",data_set);
    console.log("time series : ",series_label);
    newObj.labels = series_label //시간만 담김 배열
    newObj.datasets =  data_set//각 데이터 셋의 배열
    console.log("Chart Object : ",newObj);
    config.type = 'line',
    config.data = newObj
    config.options = {
        responsive: true,
        title: {
            display: true,
            text: title
        },
        tooltips: {
            mode: 'index',
            intersect: false,
        },
        hover: {
            mode: 'nearest',
            intersect: true
        },
        scales: {
            x: {
                display: true,
                scaleLabel: {
                    display: true,
                    labelString: 'Time'
                }
            },
            y: {
                display: true,
                scaleLabel: {
                    display: true,
                    labelString: 'Value'
                }
            }
        }
    }// end of config.options
   return newObj;
}

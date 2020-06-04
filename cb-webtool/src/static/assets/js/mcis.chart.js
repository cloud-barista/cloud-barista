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

function getMetric(chart_target,target, mcis_id, vm_id, metric, periodType,statisticsCriteria, duration){
    console.log("====== Start GetMetric ====== ")
  
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
   
   
   var url = DragonFlyURL+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
  // url = 'http://182.252.135.42:9090/dragonfly/mcis/mzc-aws-montest-01-mcis/vm/aws-mon-test-east-01/metric/'+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
   console.log("Request URL : ",url)
   var html = "";
   $.ajax({
   url: url,
   async:false,
   type:'GET',
   success : function(result){
       var data = result
         console.log("Get Monitoring Data : ",data)
         console.log("======== start mapping data ======");
         var time_obj = time_arr(data,target);
         console.log("chart_target :",chart_target);
        
        // var myChart = new Chart(ctx, time_obj);
        chart.data = time_obj;
        chart.update();
        $("#chart_detail").show();
        fnMove('chart_detail');

   },
   error : function(request,status, error){
       console.log("ERROR request status at DragonFly : ",status);
     
       
       
   }
   
})
   // axios.get(url).then(result=>{
   //       var data = result.data
   //       console.log("Get Monitoring Data : ",data)
   //       console.log("======== start mapping data ======");
   //       var time_obj = time_arr(data,target);
   //       console.log("chart_target :",chart_target);
   //       var ctx = document.getElementById(chart_target).getContext('2d')
   //       var myChart = new Chart(ctx, time_obj);
   //       myChart.update();
   //      // Chart.Line('canvas1',time_obj);
   //       console.log("==time series==",time_obj);
   //      // var metricObject = mappingMetric(data);
   //      // var m_len = metricObject.length-1

   //     //   if(target == "cpu"){
   //     //     color += '<div class="icon icon-shape bg-success text-white rounded-circle shadow">'
   //     //           + '<i class="fas fa-microchip"></i>';
   //     //     var num = parseFloat(metricObject[m_len].cpu_utilization)
   //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>'                  
   //     //   }else if(target == "memory"){
   //     //     color += '<div class="icon icon-shape bg-warning text-white rounded-circle shadow">'
   //     //           + '<i class="fas fa-memory"></i>'
   //     //     var num = parseFloat(metricObject[m_len].mem_utilization)
   //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>' 
   //     //   }else if(target == "disk"){
   //     //     color += '<div class="icon icon-shape bg-danger text-white rounded-circle shadow">'
   //     //           + '<i class="far fa-save"></i>'
   //     //     var num = parseFloat(metricObject[m_len].used_percent)
   //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(3)+'%</span>' 
   //     //   }else if(target == "network"){
   //     //     color += '<div class="icon icon-shape bg-primary text-white rounded-circle shadow">'
   //     //           + '<i class="fas fa-network-wired"></i>'
   //     //     var num = parseFloat(metricObject[m_len].bytes_in)
   //     //     metric_size +='<span class="h2 font-weight-bold mb-0">'+num.toFixed(1)+'byte</span>' 
   //     //   }
   //     //         html += '<div class="card card-stats mb-4 mb-xl-0">'
   //     //              +'<div class="card-body">'
   //     //              +'<div class="row">'
   //     //              +'<div class="col">'
   //     //              +'<h5 class="card-title text-uppercase text-muted mb-0">'+metric+'</h5>'
   //     //              //+'<span class="h2 font-weight-bold mb-0">2,356</span>'
   //     //              +metric_size
   //     //              +'</div>'
   //     //              +'<div class="col-auto">'
                  
   //     //              +color
                    
   //     //              //+'<i class="fas fa-chart-pie"></i>'
   //     //              +'</div>'
   //     //              +'</div>'
   //     //              +'</div>'
   //     //              +'<p class="mt-3 mb-0 text-muted text-sm">'
   //     //              +'<span class="text-danger mr-2"> 3.48%</span>'
   //     //              +'<span class="text-nowrap">'+metricObject[0].time+'</span>'
   //     //              +'</p>'
   //     //              +'</div>'
   //     //              +'</div>';
         
   //     //   $("#"+target+"").empty()
   //     //   $("#"+target+"").append(html)
   //   })
}

function checkDragonFly(mcis_id, vm_id){
   console.log("====== Start Check DragonFly ====== ")
   var periodType = "m";
   var duration = "10m";
   var statisticsCriteria = "last";
   var metric = "cpu" 
   var url = DragonFlyURL+"/mcis/"+mcis_id+"/vm/"+vm_id+"/metric/"+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
   //url = 'http://182.252.135.42:9090/dragonfly/mcis/mzc-aws-montest-01-mcis/vm/aws-mon-test-east-01/metric/'+metric+"/info?periodType="+periodType+"&statisticsCriteria="+statisticsCriteria+"&duration="+duration;
   console.log("Request URL : ",url)
   
   $.ajax({
        url: url,
        async:false,
        type:'GET',
        success : function(result){
            
            $("#check_dragonFly").val("200");
        },
        error : function(request,status, error){
        
            $("#check_dragonFly").val("400");
            
        }          
    })
   
}



function time_arr(obj, title){
    //data sets
   var labels = obj.columns;
   var datasets = obj.values;
    // 각 값의 배열 데이터들
   
   var series_label = new Array();
   var data_set = new Array();
   // 최종 객체 data
   var new_obj = {}
   var color_arr = ['rgb(255, 99, 132)','rgb(255, 159, 64)', 'rgb(255, 205, 86)','rgb(75, 192, 192)','rgb(54, 162, 235)','rgb(153, 102, 255)','rgb(201, 203, 207)','rgb(99, 255, 243)']   

   for(var i in labels){
    var dt = {}  
    var series_data = new Array();  
    for(var k in datasets){
        if(i == 0){
            series_label.push(datasets[k][i]) //이건 시간만 담는다.
        }else{
            dt.label = labels[i];
            series_data.push(datasets[k][i]) //그외 나머지 데이터만 담는다.
            dt.borderColor = color_arr[i];
            dt.backgroundColor = color_arr[i];
            dt.fill= false;
           // dt.data
        }  
    }
    if(i > 0){
       dt.data = series_data
       data_set.push(dt)
    }
   
    
   
   }
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
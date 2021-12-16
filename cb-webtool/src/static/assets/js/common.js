// funtcion requestAjax(url, method, data){
//     console.log("Request URL : ",url)
//     var met = method.toLowerCase
//     $.ajax({
//         url : url,
//         type: method,
//         data: data

//     }).then(function(result){
//         console.log(result)
//     })
// }

//폼의 Validation을 체크함.
//<input> tag에 "required" 옵션이 추가된 항목의 값이 공백인 경우 false  그렇지 않으면 true 리턴
function chkFormValidate(formObj) {
	try {
		var objs = formObj.find("[required]");
		//alert(objs.length)

		// required 옵션이 체크된 필드 들의 값을 조회 함.(현재는 Text 필드만 가능)
		for (var i = 0; i < objs.length; i++) {
			if (objs.eq(i).val() == '') {
				alert("Please enter a value.");
				objs.eq(i).focus();
				return false;
			}
		}
		return true;
	} catch (e) {
		alert(e);
		return false;
	}
}

function getOSType(image_id) {
	var url = CommonURL + "/ns/" + NAMESPACE + "/resources/image/" + image_id
	console.log("api Info : ", ApiInfo);
	return axios.get(url, {
		headers: {
			'Authorization': apiInfo
		}

	}).then(result => {
		var data = result.data
		var osType = data.guestOS
		console.log("Image Data : ", data);
		return osType;
	})
}
function checkNS() {
	var url = CommonURL + "/ns";
	var apiInfo = ApiInfo
	axios.get(url, {
		headers: {
			'Authorization': apiInfo
		}
	}).then(result => {
		var data = result.data.ns
		if (!data) {
			commonAlert("NameSpace가 등록되어 있지 않습니다.\n등록페이지로 이동합니다.")
			location.href = "/NS/reg";
			return;
		}
	})

}
function getNameSpace() {
	var url = CommonURL + "/ns"
	var apiInfo = ApiInfo
	axios.get(url, {
		headers: {
			'Authorization': apiInfo
		}
	}).then(result => {
		var data = result.data.ns
		var namespace = ""
		for (var i in data) {
			if (i == 0) {
				namespace = data[i].id
			}
		}
		$("#namespace1").val(namespace);

	})
}
function cancel_btn() {
	if (confirm("Cancel it?")) {
		history.back();
	} else {
		return;
	}
}
function close_btn() {
	if (confirm("close it?")) {
		$("#transDiv").hide();
	} else {
		return;
	}
}
function fnMove(target) {
	var offset = $("#" + target + "").offset();
	console.log("FnMove offset : ", offset)
	$('html, body').animate({ scrollTop: offset.top }, 400);
}

function goFocus(target) {
	console.log(event)
	event.preventDefault()
	$("#" + target).focus();
	fnMove(target)
}

// MCIS 상태값 중 일부만 사용 
// ex) Partial-Suspended-1(2/2)  : 가운데값만 사용
// todo : 일부정지인데 stop으로 표시하고 있는데....
function getMcisStatusDisp(mcisFullStatus) {
	console.log("getMcisStatus " + mcisFullStatus);
	var statusArr = mcisFullStatus.split("-");
	returnStatus = statusArr[0].toLowerCase();

	// const MCIS_STATUS_RUNNING = "running"
	// const MCIS_STATUS_INCLUDE = "include"
	// const MCIS_STATUS_SUSPENDED = "suspended"
	// const MCIS_STATUS_TERMINATED = "terminated"
	// const MCIS_STATUS_PARTIAL = "partial"
	// const MCIS_STATUS_ETC = "etc"
	// console.log("before status " + returnStatus)
	// if (returnStatus == MCIS_STATUS_RUNNING) {
	// 	returnStatus = "running"
	// } else if (returnStatus == MCIS_STATUS_INCLUDE) {
	// 	returnStatus = "stop"
	// } else if (returnStatus == MCIS_STATUS_SUSPENDED) {
	// 	returnStatus = "stop"
	// } else if (returnStatus == MCIS_STATUS_TERMINATED) {
	// 	returnStatus = "terminate"
	// } else if (returnStatus == MCIS_STATUS_PARTIAL) {
	// 	returnStatus = "stop"
	// } else if (returnStatus == MCIS_STATUS_ETC) {
	// 	returnStatus = "stop"
	// } else {
	// 	returnStatus = "stop"
	// }

	if (mcisFullStatus.toLowerCase().indexOf("running") > -1) {
		returnStatus = "running"
	} else if (mcisFullStatus.toLowerCase().indexOf("suspend") > -1) {
		returnStatus = "stop"
	} else if (mcisFullStatus.toLowerCase().indexOf("terminate") > -1) {
		returnStatus = "terminate"
		// TODO : partial도 있는데... 처리를 어떻게 하지??
	} else {
		returnStatus = "terminate"
	}
	console.log("after status " + returnStatus)
	return returnStatus
}

function getMcisStatusIcon(mcisDispStatus){
	var mcisStatusIcon = "";
	if(mcisDispStatus == "running"){ mcisStatusIcon = "icon_running_db.png"
	}else if(mcisDispStatus == "include" ){ mcisStatusIcon = "icon_stop_db.png"
	}else if(mcisDispStatus == "suspended"){mcisStatusIcon = "icon_stop_db.png"
	}else if(mcisDispStatus == "terminate"){mcisStatusIcon = "icon_terminate_db.png"
	}else{
		mcisStatusIcon = "icon_stop_db.png"
	}
	return mcisStatusIcon
}
// VM 상태를 UI에서 표현하는 방식으로 변경
function getVmStatusDisp(vmFullStatus) {
	console.log("getVmStatusDisp " + vmFullStatus);
	var returnVmStatus = vmFullStatus.toLowerCase() // 소문자로 변환

	const VM_STATUS_RUNNING = "running"
	const VM_STATUS_STOPPED = "stop"
	const VM_STATUS_RESUMING = "resuming";
	const VM_STATUS_INCLUDE = "include"
	const VM_STATUS_SUSPENDED = "suspended"
	const VM_STATUS_TERMINATED = "terminated"
	const VM_STATUS_FAILED = "failed"

	if (returnVmStatus == VM_STATUS_RUNNING) {
		returnVmStatus = "running"
	} else if (returnVmStatus == VM_STATUS_TERMINATED) {
		returnVmStatus = "terminate"
	} else if (returnVmStatus == VM_STATUS_FAILED) {
		returnVmStatus = "terminate"
	} else {
		returnVmStatus = "stop"
	}
	return returnVmStatus
}

function getVmStatus(vm_name, connection_name) {
	var url = "/vmstatus/" + vm_name + "?connection_name=" + connection_name
	var apiInfo = ApiInfo;
	$.ajax({
		url: url,
		async: false,
		type: 'GET',
		beforeSend: function (xhr) {
			xhr.setRequestHeader("Authorization", apiInfo);
			xhr.setRequestHeader("Content-type", "application/json");
		},
		success: function (res) {
			var vm_status = res.Status

		}
	})
}

function getVmStatusClass(vmDispStatus){
	var vmStatusClass = "bgbox_g";
	if (vmDispStatus == "running") {
		vmStatusClass = "bgbox_b"
	} else if (vmDispStatus == "include") {
		vmStatusClass = "bgbox_g"
	} else if (vmDispStatus == "suspended") {
		vmStatusClass = "bgbox_g"
	} else if (vmDispStatus == "terminated") {
		vmStatusClass = "bgbox_r"
	} else {
		vmStatusClass = "bgbox_r"
	}
	return vmStatusClass;
}
function getVmStatusIcon(vmDispStatus){
	var vmStatusIcon = "icon_running_db.png";
	if(vmDispStatus == "running"){
	    vmStatusIcon = "icon_running_db.png";
	}else if(vmDispStatus == "stop"){
	    vmStatusIcon = "icon_stop_db.png";
	}else if(vmDispStatus == "suspended"){
	    vmStatusIcon = "icon_stop_db.png";
	}else if(vmDispStatus == "terminate"){
	    vmStatusIcon = "icon_terminate_db.png";
	}else{
	    vmStatusIcon = "icon_stop_db.png";
	}
	return vmStatusIcon;
}
// 좌측메뉴 선택 표시
// 경로를 split하여 첫번째 : Operation / Setting, 두번째 선택, 세번째 선택하도록 
// //http://localhost:1234/setting/connections/cloudconnectionconfig/mngform
// //http://localhost:1234/setting/resources/network/mngform
// 이 때 
// menu 2 의 id 는 menu_level2_connections, menu_level2_resources
// menu 3 의 id 는 menu_level3_cloudconnectionconfig, menu_level3_network

function lnb_on() {
	var url = new URL(location.href)
	var path = url.pathname
	path = path.split("/")
	var target1 = path[1]
	var target2 = path[2]
	var target3 = path[3]
	// console.log('lnb_on path : ' + path)
	// console.log('target1=' + target1)
	// console.log('target2=' + target2)
	console.log('target3=' + target3)
	if (target1 == undefined || target1 == "main") {
		target1 = "operation";
	}

	$("#tab_" + target1).addClass("active")

	// menu의 첫번째 단계인 operation, setting 은 common.css 에 id로 style이 적용되어있어 변경이 어려움.
	$("#" + target1).addClass("active")
	$("#" + target1).addClass("on")
	$("#" + target1).addClass("show")
	//show active

	// $("#"+target1).addClass("active")

	// $("#" + target1) // Setting
	$("#menu_level2_" + target2).addClass("active")
	$("#menu_level3_" + target3).addClass("on")

	$(".leftmenu .tab-content ul > li").each(function () {

	})

}
//webmoa common
$(function () {
	//body scrollbar
	jQuery('.scrollbar-dynamic').scrollbar();
	//Server List scrollbar
	jQuery('.ds_cont .listbox.scrollbar-inner').scrollbar();
	//selectbox
	//jQuery('.selectbox').niceSelect();
	//menu_level3_cloudconnectionconfig
	/* lnb s */


	var $menu_li = $('.menu > li'),
		$ul_sub = $('.menu > li ul'),
		$lnb = $('#lnb'),
		$mobileCate = $('#mobileCate'),
		$contents = $('#contents'),
		$menubg = $('#lnb.on .bg'),
		$topmenu = $contents.find('.topmenu'),
		$btn_menu = $('#btn_menu'),
		$btn_top = $('#btn_top');
	console.log(" $menu_li ")
	console.log($menu_li)
	//left menu upDwon
	$menu_li.children('a').not('.link').click(function () {
		console.log("left menu updownd clicked ")
		if ($(this).next().css('display') === 'none') {
			console.log("left menu display none ")
			$menu_li.removeClass('on');
			$ul_sub.slideUp(300);
			$(this).parent().addClass('on');
			$(this).next().slideDown(300);
		} else {
			console.log("left menu display else ")
			$(this).parent().removeClass('on');
			$(this).next().slideUp(300);
		}
		return false;
	});


	//mobile on(open)
	$btn_menu.click(function () {
		console.log(" $btn_menu " + btn_menu)
		$menubg.stop(true, true).fadeIn(300);
		$lnb.animate({ right: 0 }, 300);
		$lnb.addClass('on');
		$('html, body').addClass('body_hidden');
	});
	//mobile topmenu copy
	$lnb.find('.bottom').append($topmenu.clone());

	//mobile off(close)
	$('#m_close, #lnb .bg').click(function () {
		$menubg.stop(true, true).fadeOut(300);
		$lnb.animate({ right: -350 }, 300);
		$lnb.removeClass('on');
		$('html, body').removeClass('body_hidden');
	});

	//left Name Space mouse over
	$("#lnb .topbox .txt_2").each(function () {
		var $btn = $(this);
		var list = $btn.find('.infobox');
		var menuTime;
		$btn.mouseenter(function () {
			list.fadeIn(300);
			clearTimeout(menuTime);
		}).mouseleave(function () {
			clearTimeout(menuTime);
			menuTime = setTimeout(mTime, 200);
		});
		function mTime() {
			list.stop().fadeOut(200);
		}
	});
	/* lnb e */

	//header menu mouse over
	$("#lnb .topmenu > ul > li").each(function () {
		var $btn = $(this);
		var list = $btn.find('.infobox');
		var menuTime;
		$btn.mouseenter(function () {
			list.fadeIn(300);
			clearTimeout(menuTime);
		}).mouseleave(function () {
			clearTimeout(menuTime);
			menuTime = setTimeout(mTime, 200);
		});
		function mTime() {
			list.stop().fadeOut(200);
		}
	});

	//header menu click(toggle)
	$(".header .topmenu > ul > li").each(function () {
		var $btn = $(this);
		var list = $btn.find('.infobox');
		var badge = $btn.find('.badge');
		$btn.click(function () {
			list.fadeToggle(300, function () {
				badge.innerHTML = "0";
				badge.hide();
			});
		});
	});

	//Action menu mouse over
	$(".dashboard .top_info > ul > li").each(function () {
		var $btn = $(this);
		var list = $btn.find('.infobox');
		var menuTime;
		$btn.mouseenter(function () {
			list.fadeIn(300);
			clearTimeout(menuTime);
		}).mouseleave(function () {
			clearTimeout(menuTime);
			menuTime = setTimeout(mTime, 200);
		});
		function mTime() {
			list.stop().fadeOut(200);
		}
	});

	//common table on/off
	$(".dashboard .status_list tbody tr").each(function () {
		var $td_list = $(this),
			$status = $(".server_status"),
			$detail = $(".server_info");
		$td_list.off("click").click(function () {
			console.log("common td list click add on")
			$td_list.addClass("on");
			$td_list.siblings().removeClass("on");
			$status.addClass("view");
			$status.siblings().removeClass("on");
			$(".dashboard.register_cont").removeClass("active");
			$td_list.off("click").click(function () {
				if ($(this).hasClass("on")) {
					$td_list.removeClass("on");
					$status.removeClass("view");
					$detail.removeClass("active");
				} else {
					$td_list.addClass("on");
					$td_list.siblings().removeClass("on");
					$status.addClass("view");
					$status.siblings().removeClass("view");
					$(".dashboard.register_cont").removeClass("active");
				}
			});
		});
	});

	//RuleSet(s) mouse over
	$(".bubble_box .box").each(function () {
		var $list = $(this);
		var bubble = $list.find('.bb_info');
		var menuTime;
		$list.mouseenter(function () {
			bubble.fadeIn(300);
			clearTimeout(menuTime);
		}).mouseleave(function () {
			clearTimeout(menuTime);
			menuTime = setTimeout(mTime, 100);
		});
		function mTime() {
			bubble.stop().fadeOut(100);
		}
	});

	//Manage MCIS Server List on/off
	$(".dashboard .ds_cont .area_cont .listbox li.sel_cr").each(function () {
		console.log("sel_cr");
		var $sel_list = $(this),
			$detail = $(".server_info");
		$sel_list.off("click").click(function () {
			$sel_list.addClass("active");
			$sel_list.siblings().removeClass("active");
			$detail.addClass("active");
			$detail.siblings().removeClass("active");
			$sel_list.off("click").click(function () {
				if ($(this).hasClass("active")) {
					$sel_list.removeClass("active");
					$detail.removeClass("active");
				} else {
					$sel_list.addClass("active");
					$sel_list.siblings().removeClass("active");
					$detail.addClass("active");
					$detail.siblings().removeClass("active");
				}
			});
		});
	});


	//Monitoring MCIS Server List on/off
	$(".ds_cont_mbox .mtbox .g_list .listbox li.sel_cr").each(function () {
		var $sel_list = $(this),
			$detail_view = $(".monitoring_view");
		$sel_list.off("click").click(function () {
			console.log("sel_list click add active")
			$sel_list.addClass("active");
			$sel_list.siblings().removeClass("active");
			$detail_view.addClass("active");
			$detail_view.siblings().removeClass("active");
			$sel_list.off("click").click(function () {
				if ($(this).hasClass("active")) {
					console.log("sel_list click remove active")
					$sel_list.removeClass("active");
					$detail_view.removeClass("active");
				} else {
					console.log("sel_list click remove active, add sibling")
					$sel_list.addClass("active");
					$sel_list.siblings().removeClass("active");
					$detail_view.addClass("active");
					$detail_view.siblings().removeClass("active");
				}
			});
		});
	});

	/*
	$(".graph_list .glist .gbox").each(function(){
		var $glist = $(this),
				$detail_view = $(".g_detail_view");
		$glist.off("click").click(function(){
			$glist.addClass("active");
			$glist.siblings().removeClass("active");
			$detail_view.addClass("active");
			$detail_view.siblings().removeClass("active");
			  $glist.off("click").click(function(){
				if( $(this).hasClass("active") ) {
					$glist.removeClass("active");
					$detail_view.removeClass("active");
			} else {
					$glist.addClass("active");
					$glist.siblings().removeClass("active");
					$detail_view.addClass("active");
					$detail_view.siblings().removeClass("active");
			}
			});
		});
	});
	*/

	$(".dashboard.dashboard_cont .ds_cont .dbinfo").each(function () {
		var $list = $(this);
		$list.on('click', function () {
			if ($(this).hasClass("active")) {
				$list.removeClass("active");
			} else {
				$list.addClass("active");
				$list.siblings().removeClass("active");
			}
		});
	});

	// btn_top
	$("#footer .btn_top").click(function () {
		$("html,body,#wrap").stop().animate({
			scrollTop: 0
		});
	});

	$(".pop_setting_chbox input:checkbox").on('click', function () {
		if ($(this).prop('checked')) {
			$(this).parent().addClass("selected");
		} else {
			$(this).parent().removeClass("selected");
		}
	});


});

// mobile table
$(function () {

	$(".dataTable tr span.ov").each(function () {
		$(this).on('click', function () {
			$(this).parent().parent().find(".btn_mtd").toggleClass("over");
			$(this).parent().parent().find(".overlay").toggleClass("hidden");
		});
	});
});

/*
$('.graph_list .glist a[href*="#"]').click(function(event) {
  if (location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '') && location.hostname == this.hostname) {
	var target = $(this.hash);
	target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
	if (target.length) {
	  event.preventDefault();
	  $('html, body ,#wrap').animate({
		scrollTop: target.offset().top
	  }, 500);
	  return false;
	}
  }
});
*/

/*
	지원하는 cloud driver 목록
	target : target object = id (name아님)
	getCloudOS(apiInfo, target)
	
	ex)
	var spiderURL =  "{{ .comURL.SpiderURL}}"
	var apiInfo = "{{ .apiInfo}}";
	getCloudOS(spiderURL,apiInfo,'ProviderName')
 */
// function getCloudOS(apiInfo, target){
// 	var url = SpiderURL;

//     var req_url = SpiderURL+"/cloudos"
// 	console.log("getCloudOS ::: " + " : " + req_url );
//     var initCSP = ""
//     axios.get(req_url,{
//     headers:{
//             'Authorization': apiInfo
//         }
//     }).then(result=>{
//         var data = result.data.cloudos
//         var html =""
//         if(data){
//             html += '<option>Select Provider</option>'
//             data.map(csp=>(html += '<option value="'+csp+'">'+csp+'</option>'))
//         }
//         html += '<option value="MOCK">MOCK</option>'
//         $("#"+target).empty()
//         $("#"+target).append(html)

//         initCSP = data[0]

//         // changeProvider(url,initCSP)// 이게 필요한가? 바뀔때 Event가 이미 있을 텐데??
//     }) 
// }







// namespace 목록에서 한 개 선택. 해당값을 임시로 저장하고 confirm 창 띄우기
// 실제 set은  setDefaultNameSpace function에서  ajax호출로
// set과 select 혼돈하지 말 것.
function selectDefaultNameSpace(callerLocation, nameSpaceID) {
	// 변경할 namespaceId를 임시로 
	console.log("selectDefaultNameSpace " + callerLocation + ", " + nameSpaceID)
	if (callerLocation == "TobBox") {
		$("#tempSelectedNameSpaceID").val(nameSpaceID);
		commonConfirmOpen("ChangeNameSpace");
	} else if (callerLocation == "LNBPopup") {
		console.log("selectDefaultNameSpace " + callerLocation + ", " + nameSpaceID + " set!!")
		// Modal 내 namespace 값을 hidden으로 set
		$("#tempSelectedNameSpaceID").val(nameSpaceID);
		// 선택했고 OK버튼이 나타난다. OK버튼 클릭시 저장 됨
		console.log("선택했음. Set을 해야 실제로 저장 됨")
	} else if (callerLocation == "Main") {
		console.log("selectDefaultNameSpace " + callerLocation + ", " + nameSpaceID + " set!!")
		// Modal 내 namespace 값을 hidden으로 set
		$("#tempSelectedNameSpaceID").val(nameSpaceID);
		console.log("선택했음. Set을 해야 실제로 저장 됨")
	}
}

// namespace 선택 후 OK 버튼 클릭시(modal, main)에서 
function nameSpaceSet(callerLocation) {
	var nameSpaceID = $("#tempSelectedNameSpaceID").val();
	console.log("nameSpaceSet OK " + nameSpaceID)
	setDefaultNameSpace(nameSpaceID, callerLocation)
}

// store에 defaultnamespace 변경. namespace가 등록되어 있지 않으면 ns 설정 page로 이동
function setDefaultNameSpace(nsid, callerLocation) {
	console.log("setNameSpace : " + nsid)
	if (nsid) {
		//reqUrl = "/SET/NS/"+nsid;
		var url = "/setting/namespaces/namespace/set/" + nsid;
		console.log(url);
		axios.get(url, {
			// headers:{
			//     'Authorization': apiInfo
			// }
		}).then(result => {
			var data = result.data.LoginInfo
			console.log(data);
			// 성공했으면 해당 namespace 선택 또는 조회
			console.log(" defaultNameSpaceID : " + data.DefaultNameSpaceID)
			$('#topboxDefaultNameSpaceID').val(data.DefaultNameSpaceID)
			$('#topboxDefaultNameSpaceName').text(data.DefaultNameSpaceName)

			if (callerLocation == "Main") {
				// $('#loadingContainer').show();// page 이동 전 loading bar를 보여준다.
				// location.href = "/operation/dashboards/dashboardnamespace/mngform"
				var targetUrl = "/operation/dashboards/dashboardnamespace/mngform"
				changePage(targetUrl)
			} else if (callerLocation == "NameSpace") {
				// commonAlert(data.DefaultNameSpaceID + "가 기본 NameSpace로 변경되었습니다.")
				commonAlert("기본 NameSpace로 변경되었습니다")
				location.reload(); // TODO : 호출한 곳에서 reload를 할 것인지 redirect를 할 것인지
			} else {
				location.reload(); // TODO : 호출한 곳에서 reload를 할 것인지 redirect를 할 것인지
			}
			// 
			// }).catch(function(error){
			// 	console.log("setNameSpace error : ",error);        
			// });
		}).catch((error) => {
			console.warn(error);
			console.log(error.response)
		});
	} else {
		commonAlert("NameSpace가 선택되어 있지 않습니다.\n등록되어 있지 않은 경우 등록하세요.")
		//location.href ="/NS/reg";
	}
}

// this.value -> 특정 obj 에 넣을 때 사용
function copyValue(targetValue, targetObjId) {
	$("#" + targetObjId).val(targetValue);
}


// 이름 Validation : 소문자, 숫자, 하이프(-)만 가능   [a-z]([-a-z0-9]*[a-z0-9])?
function validateCloudbaristaKeyName(elementValue, maxLength) {
	var returnStr = "first letter = small letter \n middle letter = small letter, number, hyphen(-) only \n last letter = small letter";
	//var charsPattern = /^[a-zA-Z0-9-]*$/;
	//var charsPattern = /^[a-z0-9-]*$/;
	//var charsPattern = /^[a-z]([-a-z0-9]*[a-z0-9])$/;
	//var regex = new RegExp('^[0-9]*\\.[0-9]{'+b+'}$') ;

	// min = 3 이므로 4자이상. maxlength + 1 이하 ex( 3, 12) 면 4자~13자 까지 허용
	var regex = new RegExp('^[a-z]([-a-z0-9]*[a-z0-9])$');
	console.log("validation " + elementValue + " : " + maxLength)
	if (maxLength == undefined) {
		if (!regex.test(elementValue)) {
			return false;
		}
	}
	var str_length = elementValue.length; // 전체길이
	try {
		if (maxLength > 0) {
			if (Number(str_length) > Number(maxLength)) {
				console.log(returnStr);
				return false;
			}

			console.log(" maxlength is defined " + maxLength + " : " + elementValue.length)
			// regex = new RegExp('^[a-z]([-a-z0-9]*[a-z0-9]){' + maxLength+'}$') ;
			//regex = new RegExp('^[a-z]([-a-z0-9]*[a-z0-9]){ 5,' + maxLength+'}$') ;

			regex = new RegExp("^[a-z]([-a-z0-9]*[a-z0-9]){3," + maxLength + "}$", "g");

			if (!regex.test(elementValue)) {
				console.log("return val " + elementValue)
				return false;
			}
		}
	} catch (e) {
		return false;
		// console.log(e);
	}
	console.log("validate return")
	return true;
	//return charsPattern.test(elementValue);
}


// 해당 table 의 limit를 초과하면 scroll이 생기도록
// width는 colgroup이 없는 채로 ht, td 에 width class를 추가한다.
function setTableHeightForScroll(tableId, limitHeight) {
	var tableHeight = $("#" + tableId).height();
	if (tableHeight > limitHeight) {
		$("#" + tableId).css({ height: limitHeight });
	}
}

// 비어있으면 false, 안비어있으면 true
function checkEmptyString(stringVal) {
	if (stringVal == null ||
		stringVal == undefined ||
		stringVal == 0) {
		return false;
	}
	return true;
}

// plus 버튼을 추가하는 script
function getPlusVm(){
	var append = "";
	append = append + '<li id="plusVmIcon" >';
	append = append + '<div class="server server_add" onClick="displayNewServerForm()">';
	append = append + '</div>';
	append = append + '</li>';
	return append;
}
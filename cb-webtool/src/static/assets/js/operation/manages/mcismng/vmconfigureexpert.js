$(document).ready(function () {

  //tab menu Server / OS_HW /	Network / Security / Other 위치 표시
  $(".create_tab .nav a").click(function () {
    var idx = $(".create_tab .nav a").index(this);
    for (i = 0; i < $(".create_tab .nav a").length; i++) {
      if (i == idx) {
        $('.config_location > div').removeClass('on');
        $('.config_location > div > span').eq(idx).parent().addClass('on');
      }
    }
  });
  //tab 내용 다음
  $(".create_tab .btn_next").click(function (e) {
    var $active = $('.create_tab .nav li > .active');
    $active.parent().next().find('.nav-link').removeClass('disabled');
    nextTab($active);
  });

  //tab 내용 이전
  $(".create_tab .btn_prev").click(function (e) {
    var $active = $('.create_tab .nav li > a.active');
    prevTab($active);
  });

  //Deployment Target table scrollbar
  $('.btn_assist').on('click', function () {
    console.log("btn_assist click--")
    $("#Deployment_box").modal();
    $('.dtbox.scrollbar-inner').scrollbar();
  });

  //Server Configuration clear
  $(".btn_clear").click(function () {
    //$('.svc_ipbox').find('input, textarea').val('');
    $("#es_name").val('');
    $("#es_description").val('');
    $("#es_regConnectionName").val('');
    setConnectionValue("");

    osHardwareClear();
    vnetClear();
    vmSecurityClear();
  });

  //OS_HW - Clear
  $("#OS_HW .btn_clear").click(function () {
    // $('#OS_HW .tab_ipbox').find('input, textarea').val('');
    osHardwareClear();
  });
  //Network - Clear
  $("#Network .btn_clear").click(function () {
    // $('#Network .tab_ipbox').find('input, textarea').val('');
    vnetClear();
  });
  //Security - Clear
  $("#Security .btn_clear").click(function () {
    // $('#Security .tab_ipbox').find('input, textarea').val('');
    vmSecurityClear();
  });
  //Other - Clear
  $("#Other .btn_clear").click(function () {
    $('#Other .tab_ipbox').find('input, textarea').val('');
  });

});

$(document).on('change', 'input[name="securityGroup_chk"]',
  function () {
    console.log("in checkbox change")
    var chkIdArr = $(this).attr('id').split("_");// 0번째와 2번째를 합치면 id 추출가능  ex) securityGroup_Raw_0
    console.log("chkIdArr " + chkIdArr + " : " + $(this).is(":checked"))
    if ($(this).is(":checked")) {

      $("#assistSelectedIndex").val(chkIdArr[2]).val();
      // 해당 securityGroup의 connection과 form의 connection이 다르면초기화 후 set
      // 같으면 securityGroup set

      var selectedId = $("#" + chkIdArr[0] + "_id_" + chkIdArr[2]).val()//id="securityGroup_id_{{$securityGroupIndex}}"
      var selectedInfo = $("#" + chkIdArr[0] + "_info_" + chkIdArr[2]).val()
      var selectedConnectionName = $("#" + chkIdArr[0] + "_connectionName_" + chkIdArr[2]).val()
      if ($("#e_connectionName").val() != "" && $("#e_connectionName").val() != selectedConnectionName) {

        var targetTabOjbInfo = "tab_securityGroupInfo";
        var targetTabObjConnectionName = "tab_securityGroupConnectionName";

        $("#" + targetTabOjbInfo).val(selectedInfo);
        $("#" + targetTabObjConnectionName).val(selectedConnectionName);
        $("#t_connectionName").val(selectedConnectionName);// 임시 connectionName set
        rollbackObjArr[0] = targetTabOjbInfo;
        rollbackObjArr[1] = targetTabObjConnectionName;

        commonConfirmOpen("DifferentConnectionAtSecurityGroup");
        // TODO : commonConfirmOpen 해서 OK면 초기화  :
        // securityGroupTable 및 display securityGroup ifno 를 현재 connection Name으로 Set.
        // 다른 table은 1개만 선택하므로 display input box 들 초기화.
        //  connection정보로 선택된 항목의 connectionName 비교 후 초기화 function 만들 것.

      } else {
        $("#e_connectionName").val(selectedConnectionName);
        setMuipleValueToFormObj('securityGroup_chk', 'tab_securityGroupInfo', 'e_securityGroupIds')
      }
    } else {//Uncheck event
      //   // alert("B " + $(this).is(":checked") + " : " +  securityGroupId);
      $("#t_connectionName").val(selectedConnectionName).val("");
      setMuipleValueToFormObj('securityGroup_chk', 'tab_securityGroupInfo', 'e_securityGroupIds')
    }
  }
)
//////// Tab Clear //////
function osHardwareClear() {
  $("#tab_vmImageInfo").val("");
  $("#e_imageId").val("");

  $("#tab_vmSpecInfo").val("");
  $("#e_specId").val("");
}

function vnetClear() {
  $("#tab_vNetInfo").val("");
  $("#e_vNetId").val("");
  $("#e_subnetId").val("");
}

function vmSecurityClear() {
  clearCheckbox("securityGroup_chk");
  $("#tab_securityGroupInfo").val("");
  $("#e_securityGroupIds").val("");

  $("#tab_sshKeyInfo").val("");
  $("#e_sshKeyId").val("");

  $("#es_vmUserAccount").val("");
  $("#es_vmUserPassword").val("");

}



// multi 선택되는 checkbox를 한번에 초기화
function clearCheckbox(chkboxName) {
  $('input:checkbox[name="' + chkboxName + '"]').prop('checked', false);
}

// 1. tempConnection에는 값이 set되어 있어야 함.
// 2. 돌면서 connectionName이 다르면 uncheck
function uncheckDifferentConnectionAtSecurityGroup() {
  var tempConnectionNameValue = $("#t_connectionName").val();
  clearCheckboxByConnectionName('securityGroup_chk', tempConnectionNameValue);
  // 체크 된 securityGroup으로 재 설정
  setMuipleValueToFormObj('securityGroup_chk', 'tab_securityGroupInfo', 'e_securityGroupIds');
}

// 체크 된 박스에서 connectionName이 다르면 체크해제
//
function clearCheckboxByConnectionName(chkboxName, connectionName) {

  $('input:checkbox[name="' + chkboxName + '"]').each(function () {
    if (this.checked) {//checked 처리된 항목의 값
      var chkIdArr = $(this).attr('id').split("_");// 0번째와 2번째를 합치면 id 추출가능  ex) securityGroup_Raw_0
      console.log("clearCheckboxByConnectionName = ");
      var securityGroupId = $("#" + chkIdArr[0] + "_id_" + chkIdArr[2]).val()//id="securityGroup_id_{{$securityGroupIndex}}"
      securityGroupConnectionName = $("#" + chkIdArr[0] + "_connectionName_" + chkIdArr[2]).val()
      if (securityGroupConnectionName != connectionName) {
        console.log("체크 해제하자 = ");
        this.checked = false;// 체크 해제
        console.log("체크 해제완료 = ");
      }
    }
  });
}
/////////////////
$(document).ready(function () {
  //table 스크롤바 제한
  $(window).on("load resize", function () {
    var vpwidth = $(window).width();
    if (vpwidth > 768 && vpwidth < 1800) {
      $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
      $(".dataTable.scrollbar-inner").scrollbar();
    } else {
      $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
    }
  });
});


//next
function nextTab(elem) {
  $(elem).parent().next().find('a[data-toggle="tab"]').click();
}
//prev
function prevTab(elem) {
  $(elem).parent().prev().find('a[data-toggle="tab"]').click();
}

// 조회결과 table hide 시키기
function hideFilterResultTable(targetObjId, hideButtonObj) {
  document.getElementById(targetObjId).style.display = "none";
  if (hideButtonObj != undefined) {
    document.getElementById(hideButtonObj).style.display = "none";
  }
}

// TODO : util.js로 옮길 것
// select box의 option text에 compareText가 있으면 show 없으면 hide
function selectBoxFilterByText(targetObject, compareText) {
  $('#' + targetObject + ' option').filter(function () {
    if (this.value == "") return;
    console.log(this.text + " : " + compareText)
    console.log(this.text.indexOf(compareText) > -1)
    this.text.indexOf(compareText) > -1 ? $(this).show() : $(this).hide();
  });
}

// TODO : util.js로 옮길 것
// select box의 option text에 compareText1 && compareText2가 모두 있으면 show 없으면 hide
function selectBoxFilterBy2Texts(targetObject, compareText1, compareText2) {
  $('#' + targetObject + ' option').filter(function () {
    if (this.value == "") return;
    console.log(this.text + " : " + compareText1)
    console.log(this.text.indexOf(compareText1) > -1)
    if (this.text.indexOf(compareText1) > -1 && this.text.indexOf(compareText2) > -1) {
      $(this).show()
    } else {
      $(this).hide();
    }
  });
}

// Expert Mode=on 상태에서 Cloud Provider 를 변경했을 때, 해당 Provider의 region목록 조회 => 실제로는 조회되어 있으므로 filter
// 추가로 connection 정보도 조회하라고 호출
function getRegionListFilterForSelectbox(provider, targetRegionObj, targetConnectionObj) {

  // region 목록 filter
  selectBoxFilterByText(targetRegionObj, provider)
  $("#" + targetRegionObj + " option:eq(0)").attr("selected", "selected");

  // connection 목록 filter
  selectBoxFilterByText(targetConnectionObj, provider)
  $("#" + targetConnectionObj + " option:eq(0)").attr("selected", "selected");
}

// region변경시 connection 정보 filter
function getConnectionListFilterForSelectbox(regionValue, referenceObj, targetConnectionObj) {
  var referenceVal = $('#' + referenceObj).val();
  //var regionValue = region.substring(region.indexOf("]") ).trim();  
  // console.log(region + ", regionValue = " + regionValue);
  if (referenceVal == "") {
    selectBoxFilterByText(targetConnectionObj, regionValue)
  } else {
    selectBoxFilterBy2Texts(targetConnectionObj, referenceVal, regionValue)
  }

  // $("#" + targetConnectionObj + " option:eq(0)").attr("selected", "selected");
  $("#es_regConnectionName").val("");
  setConnectionValue("");// val("")을 했을 때 자동으로 설정이 안되어서 setConnectionValue("")으로 값 set.

}

// TODO : filter 기능 check
// provider, region, connection은 먼저 선택이 필수가 아닐 수 있음.
// 그래도 하위에서 일단 선택되면 변경시 알려줘야할 듯.
// 1. provider 선택시 -> 
// 2. region 선택시
// 3. OS Platform(Image) 선택 시
// 4. HW Spec 선택시
// 5. Vnet 선택시
// 6. SecurityGroup 선택시
// 7. sshKey 선택시
// 8. subnet 선택시??

//e_imageID

// Asist를 클릭했을 때 나타나는 popup에서 provider 변경 시 region selectbox와 connection table을 filter
function popProviderChange(providerObj, regionObj, targetTableObj) {
  var providerVal = $("#" + providerObj).val();
  console.log("popProviderChange " + providerVal);
  selectBoxFilterByText(regionObj, providerVal)

  $("#" + regionObj + " option:eq(0)").attr("selected", "selected");

  // table filter
  getConnectionListFilterForTable(providerObj, regionObj, targetTableObj)
}

function getConnectionListFilterForTable(providerObj, regionObj, targetTableObj) {
  var providerVal = $("#" + providerObj).val();
  var regionVal = $("#" + regionObj).val();

  $("#" + targetTableObj + " > tbody >  tr").filter(function () {
    console.log("filter table " + $(this).text());
    var compareText = $(this).text().toLowerCase()
    var toggleStatus = true;
    if (providerVal == "" && regionVal == "") {

      //return;
    } else if (providerVal == "" && compareText.indexOf(regionVal.toLowerCase()) > -1) {
      toggleStatus = true
    } else if (regionVal == "" && compareText.indexOf(providerVal.toLowerCase()) > -1) {
      toggleStatus = true
    } else if (compareText.indexOf(providerVal.toLowerCase()) > -1 && compareText.indexOf(regionVal.toLowerCase()) > -1) {
      toggleStatus = true
    } else {
      toggleStatus = false
    }
    //$(this).toggle(toggleStatus)
    if (toggleStatus) {
      $(this).show();
    } else {
      $(this).hide();
    }
  });

}
// Expert Mode=on 상태에서 Popup의 Cloud Provider 를 변경했을 때, 해당 Provider의 region목록 조회. 
// getRegionListForSelectbox() 와 동작방식은 동일
function getRegionListForPopSelectbox(provider, targetRegionObj, targetConnectionObj) {
  $('#' + targetRegionObj + ' option').filter(function () {
    if (this.value == "") return;

    return this.text.indexOf(provider) > -1 ? $(this).show() : $(this).hide();
  });

  // connection filter
}

// region 변경시, 해당 provider, region으로 connection 목록 조회
function getConnectionListByRegionForSelectbox(region, targetProviderObj, targetConnectionObj) {

}


const Expert_Server_Config_Arr = new Array();
var expert_data_cnt = 0
const expertServerCloneObj = obj => JSON.parse(JSON.stringify(obj))
function expertDone_btn() {
  console.log("expert Done")
  // TODO : 원래는 같은 VM 을 여러개 만들 때 vmGroupSize를 set 하는 것 같은데... for문으로 돌리고 있음.... 고칠까?
  // $("#e_vmGroupSize").val( $("#es_vm_add_cnt").val() )
  // validation check 
  if ($("#e_name").val() == "") { commonAlert("VM Name is required"); return; }
  if ($("#e_connectionName").val() == "") { commonAlert("Connection is required"); return; }
  if ($("#e_vNetId").val() == "") { commonAlert("vNet is required"); return; }
  if ($("#e_subnetId").val() == "") { commonAlert("Subnet is required"); return; }
  if ($("#e_securityGroupIds").val() == "") { commonAlert("SecurityGroup is required"); return; }
  if ($("#e_sshKeyId").val() == "") { commonAlert("SSH Key is required"); return; }
  if ($("#e_imageId").val() == "") { commonAlert("VM Image is required"); return; }
  if ($("#e_specId").val() == "") { commonAlert("VM Spec is required"); return; }


  $("#e_vm_add_cnt").val($("#es_vm_add_cnt").val());// 추가수량 값을 form에 추가.

  // expertForm에는 vm생성에 필요한 값들만 있음.
  var expert_form = $("#expert_form").serializeObject()
  var server_name = expert_form.name
  var server_cnt = parseInt(expert_form.vmAddCount) // expert
  console.log('server_cnt : ', server_cnt)
  var add_server_html = "";

  if (server_cnt > 1) {
    for (var i = 1; i <= server_cnt; i++) {
      var new_vm_name = server_name + "-" + i;
      var object = cloneObj(expert_form)
      object.name = new_vm_name

      add_server_html += '<li onclick="view_simple(\'' + expert_data_cnt + '\')">'
        + '<div class="server server_on bgbox_b">'
        + '<div class="icon"></div>'
        + '<div class="txt">' + new_vm_name + '</div>'
        + '</div>'
        + '</li>';
      Expert_Server_Config_Arr.push(object)
      console.log(i + "번째 Simple form data 입니다. : ", object);
    }
  } else {
    Expert_Server_Config_Arr.push(expert_form)
    add_server_html += '<li onclick="view_simple(\'' + expert_data_cnt + '\')">'
      + '<div class="server server_on bgbox_b">'
      + '<div class="icon"></div>'
      + '<div class="txt">' + server_name + '</div>'
      + '</div>'
      + '</li>';

  }
  $(".expert_servers_config").removeClass("active");
  $("#mcis_server_list").prepend(add_server_html)
  // $("#mcis_server_list").append(add_server_html)
  $("#plusVmIcon").remove();
  $("#mcis_server_list").prepend(getPlusVm());
  console.log("expert btn click and expert form data : ", expert_form)
  console.log("expert data array : ", Expert_Server_Config_Arr);
  expert_data_cnt++;
  $("#expert_form").each(function () {
    this.reset();
  })
}



//////////////////// filterling 기능 ///////////////

// 조회조건인 connection 변경시 호출.
// TODO : 실제 connectionVal이 바뀌는것이 아니라. 다른 것들의 조회 조건 filter의 기본값으로 set.
function setConnectionValue(connName) {
  console.log(" connection change")
  var connectionObj = $("#e_connectionName");
  var tempConnectionObj = $("#t_connectionName");
  if (connName == "") {
    connectionObj.val(connName);
    tempConnectionObj.val(connName);
    return;
  }

  if (connectionObj.val() == "") {// 비어있으면 그냥 set
    console.log(" initial connName")
    connectionObj.val(connName);
  } else if (connectionObj.val() != connName) {
    console.log(" diff connName " + connName + " : " + connectionObj.val())
    tempConnectionObj.val(connName);
    commonConfirmOpen("DifferentConnection")
  } else {

  }
}

// 다른 connectinName으로 set 할 때 기존에 있던 것들 중 connectionName이 다른 것들은 초기화
function setAndClearByDifferentConnectionName(caller) {
  var tempConnectionName = $("#t_connectionName").val();
  console.log("setAndClearByDifferentConnectionName " + tempConnectionName);
  //$("#expert_form").reset();// 이거하면 싹 날아가므로 connectionName이 다른 항목들만 초기화.

  $("#e_connectionName").val(tempConnectionName);

  if ($("#tab_vmImageConnectionName").val() != tempConnectionName) {
    console.log("clear tab_vmImageConnectionName " + $("#tab_vmImageConnectionName").val());
    $("#e_imageId").val("");
    $("#tab_vmImageInfo").val("")
    $("#tab_vmImageConnectionName").val("")
  }
  if ($("#tab_vmSpecConnectionName").val() != tempConnectionName) {
    console.log("clear tab_vmSpecConnectionName " + $("#tab_vmSpecConnectionName").val());
    $("#e_specId").val("");
    $("#tab_vmSpecInfo").val("");
    $("#tab_vmSpecConnectionName").val("")
  }
  // vnet
  // if( $("#tab_securityGroupConnectionName").val() != tempConnectionName ){
  //   console.log("clear tab_vmImageConnectionName " + $("#tab_securityGroupConnectionName").val());
  //   $("#e_securityGroupIds").val("");
  //   $("#tab_securityGroupConnectionName").val("")
  // }

  if ($("#tab_securityGroupConnectionName").val() != tempConnectionName) {
    console.log("clear tab_securityGroupConnectionName " + $("#tab_securityGroupConnectionName").val());
    $("#e_securityGroupIds").val("");
    $("#tab_securityGroupInfo").val("")
    $("#tab_securityGroupConnectionName").val("")
  }
  if ($("#tab_sshKeyConnectionName").val() != tempConnectionName) {
    console.log("clear tab_sshKeyConnectionName " + $("#tab_sshKeyConnectionName").val());
    $("#e_sshKeyId").val("");
    $("#tab_sshKeyInfo").val("")
    $("#tab_sshKeyConnectionName").val("")
  }

  // e_vNetId<input type="text" name="vNetId" id="e_vNetId" />								
  // 	e_subnetId<input type="text" name="subnetId" id="e_subnetId" />								

}

function hardwareSpecFilterByEnter(targetObjId, keyword) {
  console.log(event.KeyCode + " : " + keyword);
  if (event.KeyCode == 13) {
    selectBoxFilterByText(targetObjId, keyword)
  }
}

function hardwareSpecFilter(targetObjId, keywordObjId) {
  var keyword = $("#" + keywordObjId).val();
  selectBoxFilterByText(targetObjId, keyword)
}


// hidden에 들어있는 값을 기준으로 filter
function filterEnterToHidden(keywordObjId, filterColumnName, targetTableId) {
  if (window.event.keyCode == 13) {
    filterToHidden(keywordObjId, filterColumnName, targetTableId)
  }
}
// vm image 를 filter. 입력한 단어를 전체에서 찾기
// hidden에 들어있는 값을 기준으로 filter
function filterToHidden(keywordObjId, filterColumnName, tableId) {
  $("#" + keywordObjId).modal();
  var keyword = $("#" + keywordObjId).val();
  if (keyword == '') {
    // commonAlert("검색할 단어를 입력하세요")
    // return;
    keyword = "ALL";
  }

  var selectedConnectionName = $("#es_regConnectionName option:selected").val()
  if (selectedConnectionName != "") {// 선택한 connectionName이 있으면
    var connectionNameColumn = filterColumnName.split("_")[0] + "_connectionName";
    let multipleColumnMap = new Map();
    // 
    multipleColumnMap.set(connectionNameColumn, selectedConnectionName);

    multipleColumnMap.set(filterColumnName, keyword);
    filterTableByMultipleHiddenColumn(tableId, multipleColumnMap);
  } else {
    // hidden field 의 Data 를 기준으로 filter.
    // 해당 table의 column을 지정하여 filter하는데... hidden을 filterling하는 법 찾자
    filterTableByHiddenColumn(tableId, filterColumnName, keyword)
  }
  //

  // target Table이 display:none이면 보이도록
  document.getElementById(tableId).style.display = "";
  if (document.getElementById(tableId + "Hide")) {// Hide button이 있는 경우
    document.getElementById(tableId + "Hide").style.display = "";
  }

}

function filterEnterVnetGroupToHidden(keywordObjId) {
  if (window.event.keyCode == 13) {
    filterVnetGroupToHidden(keywordObjId)
  }
}

// vnet 을 한번에 filter하기위한 function
// 이름에 따라 keyword 추가
function filterVnetGroupToHidden(keywordObjId) {
  var keyword = $("#" + keywordObjId).val();
  if (keyword == '') {
    // commonAlert("검색할 단어를 입력하세요");
    // return;
    keyword = "ALL";
  }

  let multipleColumnMap = new Map();
  if ($("#filter_networkName").val() != "") {
    multipleColumnMap.set("vNet_name", $("#filter_networkName").val());
  }
  if ($("#filter_networkDescription").val() != "") {
    multipleColumnMap.set("vNet_description", $("#filter_networkDescription").val());
  }
  if ($("#filter_networkCidrBlock").val() != "") {
    multipleColumnMap.set("vNet_cidrBlock", $("#filter_networkCidrBlock").val());
  }
  if ($("#filter_networkSubnetname").val() != "") {
    multipleColumnMap.set("vNet_subnetInfos", $("#filter_networkSubnetname").val());
  }

  // connection 조건을 선택한 경우 해당 column도 함께 filterling 조건에 추가.
  var connectionNameColumn = "vNet_connectionName";
  var selectedConnectionName = $("#es_regConnectionName option:selected").val();
  multipleColumnMap.set(connectionNameColumn, selectedConnectionName);

  var targetTableId = "es_vNetList";
  filterTableByMultipleHiddenColumn(targetTableId, multipleColumnMap);

  // target Table이 display:none이면 보이도록
  document.getElementById(targetTableId).style.display = "";
}


// ConnectionName이 조건에 있는 경우 connectinName에 해당 하는 것 들 중에서 keyword가 맞는 애들만
// column 당 keyword는 1개.(map으로 처리하므로 )
function filterTableByMultipleHiddenColumn(tableId, multipleColumnMap) {
  var trs = $('#' + tableId + ' tr');
  // console.log(trs + " : " + trs.size());
  // console.log(multipleColumnMap);
  for (var i = 1; i < trs.size(); i++) {
    var isShow = true;
    multipleColumnMap.forEach(function (filterKeyword, columnName, arr) {
      console.log(columnName + ":" + filterKeyword)
      //console.log(arr)
      try {
        if (isShow == false) { return true; }

        console.log(trs.eq(i).find('input:hidden[name="' + columnName + '"]').val())
        var hiddenval = trs.eq(i).find('input:hidden[name="' + columnName + '"]').val();
        if (filterKeyword == "ALL") {
          // return true;// continue
        } else if (hiddenval.toUpperCase().indexOf(filterKeyword.toUpperCase()) > -1) {
          isShow = true;
          // return false;// break
        } else {
          isShow = false;
          return false;
        }
      } catch (e) {
        // compare failed
        console.log("filter compare failed");
        console.log(e);
      }

    });
    console.log("show hide " + isShow)
    if (isShow == true) {
      trs.eq(i).css("display", "");
    } else {
      trs.eq(i).css("display", "none");
    }
  }

  $.each(multipleColumnMap, function (key, value) {
    console.log(key + ': ' + value);
  });
}

// vNet의 경우 value set이 다른 tab과 다름(setValueToFormObj 와 기능은 같음)
// subnet 처리를 위한 로직 추가됨
function setVnetValueToFormObj(prefixTargetTabName, prefixVnetId, prefixName, selectedIndex, targetObjId) {
  setValueToFormObj(prefixTargetTabName, prefixName, selectedIndex, targetObjId);

  // subnet set 추가
  //<input type="hidden" name="vNet_subnet_{{$vNetItem.ID}}" id="vNet_subnet_{{$vNetItem.ID}}_{{$subnetIndex}}" value="{{$subnetItem.IID}}"/>
  // var selectedSubnetIndex = 0;// default : 0
  // var selectedSubnetId = $("#" + prefixName + "_subnet_" + prefixVnetId + "_" + selectedSubnetIndex).val();
  var selectedSubnetIds = $("#vNet_subnetInfos_" + selectedIndex).val();
  console.log("selectedSubnetId = " + selectedSubnetIds);

  $("#e_subnetId").val(selectedSubnetIds);// subnetId가 1개 이상이지만 securityGroup종속임.
}

// 대상 table, 선택한 tr의 index, set할 값, 대상 form의 obj 지정하여
// tr의 check시 해당 값이 obj에 저장
// TODO : 초기화는??	
function setValueToFormObj(prefixTargetTabName, prefixName, selectedIndex, targetObjId) {
  console.log("prefixTargetTabName=" + prefixTargetTabName + ",prefixName=" + prefixName + ", selectedIndex=" + selectedIndex + ", targetObjId=" + targetObjId)
  var selectedId = $("#" + prefixName + "_id_" + selectedIndex).val();
  var selectedConnectionName = $("#" + prefixName + "_connectionName_" + selectedIndex).val();

  var econnectionName = $("#e_connectionName").val();
  console.log(econnectionName + " : " + selectedConnectionName);
  var targetTabObjId = prefixTargetTabName + "Info";
  var targetTabConnectionNameObjId = prefixTargetTabName + "ConnectionName";
  if (econnectionName != "" && econnectionName != selectedConnectionName) {
    $("#" + targetTabConnectionNameObjId).val(selectedConnectionName);
    $("#" + targetObjId).val(selectedId);

    rollbackObjArr[0] = targetTabObjId;
    rollbackObjArr[1] = targetObjId;
    $("#t_connectionName").val(selectedConnectionName);// confirm을 통해서 form에 set 되므로 임시(t_connectionName)로 저장.
    commonConfirmOpen("DifferentConnection");
  } else {
    console.log("setValueToFormObj=" + targetTabObjId);
    $("#" + targetTabConnectionNameObjId).val(selectedConnectionName);
    $("#" + targetObjId).val(selectedId);

    var esSelectedConnectionName = $("#es_regConnectionName option:selected").val()
    if (esSelectedConnectionName == "") {// 선택한 connectionName이 없으면 set
      $("#es_regConnectionName").val(selectedConnectionName);
    }
    $("#e_connectionName").val(selectedConnectionName);
  }
}

// assist Popup에서 항목 선택 시 임시로 set. -> Apply버튼 클릭 시 화면으로 set할 것임.
function setAssistValue(index) {
  $("#assistSelectedIndex").val(index);
}
// assist Popup에서 apply버튼 클릭 시 assestSelectedIndex의 값으로 선택된 값들을 main으로 set
// orgPrefix : assist창의 object, targetPrefix : 메인창의 object
function applyAssistValues(caller) {
  //
  var selectedIndex = $("#assistSelectedIndex").val();
  var applyConnectionName = "";
  console.log(caller + " : " + selectedIndex);

  if (caller == "vmImageAssist") {
    var orgPrefix = "vmImageAssist_";
    var targetPrefix = "tab_vmImage_";
    console.log("vmImageAssist----")
    console.log(targetPrefix + "cspImageId")
    console.log(orgPrefix + "cspImageId_" + selectedIndex)
    console.log($("#" + orgPrefix + "cspImageId_" + selectedIndex).val())
    $("#" + targetPrefix + "id").val($("#" + orgPrefix + "id_" + selectedIndex).val());
    $("#" + targetPrefix + "name").val($("#" + orgPrefix + "name_" + selectedIndex).val());
    $("#" + targetPrefix + "cspImageId").val($("#" + orgPrefix + "cspImageId_" + selectedIndex).val());
    $("#" + targetPrefix + "cspImageName").val($("#" + orgPrefix + "cspImageName_" + selectedIndex).val());
    $("#" + targetPrefix + "guestOS").val($("#" + orgPrefix + "guestOS_" + selectedIndex).val());
    $("#" + targetPrefix + "description").val($("#" + orgPrefix + "description_" + selectedIndex).val());
    $("#" + targetPrefix + "connectionName").val($("#" + orgPrefix + "connectionName_" + selectedIndex).val());

    applyConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val()

    $("#imageAssist").modal("hide");
  } else if (caller == "vmSpecAssist") {
    var orgPrefix = "vmSpecAssist_";
    var targetPrefix = "tab_vmSpec_";

    $("#" + targetPrefix + "id").val($("#" + orgPrefix + "id_" + selectedIndex).val());
    $("#" + targetPrefix + "name").val($("#" + orgPrefix + "name_" + selectedIndex).val());
    // $("#" + targetPrefix + "cspSpecId").val($("#" + orgPrefix + "cspSpecId_" + selectedIndex).val());
    $("#" + targetPrefix + "cspSpecName").val($("#" + orgPrefix + "cspSpecName_" + selectedIndex).val());
    $("#" + targetPrefix + "memGiB").val($("#" + orgPrefix + "memGiB_" + selectedIndex).val());
    $("#" + targetPrefix + "numvCPU").val($("#" + orgPrefix + "numvCPU_" + selectedIndex).val());
    $("#" + targetPrefix + "numGpu").val($("#" + orgPrefix + "numGpu_" + selectedIndex).val());

    applyConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val()
    $("#specAssist").modal("hide");
  } else if (caller == "networkAssist") {
    var orgPrefix = "vNetAssist_";
    var targetPrefix = "tab_vNet_";

    $("#" + targetPrefix + "id").val($("#" + orgPrefix + "id_" + selectedIndex).val());
    $("#" + targetPrefix + "name").val($("#" + orgPrefix + "name_" + selectedIndex).val());
    $("#" + targetPrefix + "connectionName").val($("#" + orgPrefix + "connectionName_" + selectedIndex).val());
    $("#" + targetPrefix + "description").val($("#" + orgPrefix + "description_" + selectedIndex).val());
    $("#" + targetPrefix + "cidrBlock").val($("#" + orgPrefix + "cidrBlock_" + selectedIndex).val());
    $("#" + targetPrefix + "cspVnetName").val($("#" + orgPrefix + "cspVnetName_" + selectedIndex).val());
    $("#" + targetPrefix + "subnetId").val($("#" + orgPrefix + "subnetId_" + selectedIndex).val());
    $("#" + targetPrefix + "subnetName").val($("#" + orgPrefix + "subnetName_" + selectedIndex).val());

    applyConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val()
    $("#networkAssist").modal("hide");
  } else if (caller == "securityGroupAssist") {

    var orgPrefix = "securityGroupAssist_";
    var firewallRules = "firewallRules_";
    var targetPrefix = "tab_securityGroup_";

    var securityGroupIds = "";
    var count = 0;
    $("input[name='securityGroupAssist_chk']:checked").each(function () {
      var sgId = $(this).attr("id")
      var sgEleArr = sgId.split("_");
      var sgIndex = sgEleArr[sgEleArr.length - 1];

      if (count > 0) { securityGroupIds += "," }

      securityGroupIds += $("#" + orgPrefix + "id_" + sgIndex).val()
      $("#" + targetPrefix + "name").val($("#" + orgPrefix + "name_" + sgIndex).val());
      $("#" + targetPrefix + "vpcId").val($("#" + orgPrefix + "vNetId_" + sgIndex).val());
      $("#" + targetPrefix + "connectionName").val($("#" + orgPrefix + "connectionName_" + sgIndex).val());
      $("#" + targetPrefix + "description").val($("#" + orgPrefix + "description_" + sgIndex).val());

      $("#" + targetPrefix + firewallRules + "cidr").val($("#" + orgPrefix + firewallRules + "cidr_" + sgIndex).val());
      $("#" + targetPrefix + firewallRules + "direction").val($("#" + orgPrefix + firewallRules + "direction_" + sgIndex).val());
      $("#" + targetPrefix + firewallRules + "fromPort").val($("#" + orgPrefix + firewallRules + "fromPort_" + sgIndex).val());
      $("#" + targetPrefix + firewallRules + "toPort").val($("#" + orgPrefix + firewallRules + "toPort_" + sgIndex).val());
      $("#" + targetPrefix + firewallRules + "ipProtocol").val($("#" + orgPrefix + firewallRules + "ipProtocol_" + sgIndex).val());

      applyConnectionName = $("#" + orgPrefix + "connectionName_" + sgIndex).val();
      count++;
    });

    $("#" + targetPrefix + "id").val(securityGroupIds);

    $("#securityGroupAssist").modal("hide");

  } else if (caller == "sshKeyAssist") {

    var orgPrefix = "sshKeyAssist_";
    var targetPrefix = "tab_sshKey_";

    $("#" + targetPrefix + "id").val($("#" + orgPrefix + "id_" + selectedIndex).val());
    $("#" + targetPrefix + "name").val($("#" + orgPrefix + "name_" + selectedIndex).val());
    $("#" + targetPrefix + "connectionName").val($("#" + orgPrefix + "connectionName_" + selectedIndex).val());
    $("#" + targetPrefix + "description").val($("#" + orgPrefix + "description_" + selectedIndex).val());

    applyConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val()
    $("#sshKeyAssist").modal("hide");
  }

  console.log($("#e_connectionName").val())
  console.log("applyConnectionName = " + applyConnectionName)
  //선택된 connection과 기존 connection이 다른 tab의 data는 초기화하고 set한다
  if ($("#e_connectionName").val() != "" && $("#e_connectionName").val() != applyConnectionName) {
    setAndClearByDifferentConnectionName(caller);
  }

  if (caller == "vmImageAssist") {
    var orgPrefix = "vmImageAssist_";

    $("#e_imageId").val($("#" + orgPrefix + "id_" + selectedIndex).val());
  } else if (caller == "vmSpecAssist") {
    var orgPrefix = "vmSpecAssist_";

    $("#e_specId").val($("#" + orgPrefix + "id_" + selectedIndex).val());
  } else if (caller == "networkAssist") {
    var orgPrefix = "vNetAssist_";

    $("#e_vNetId").val($("#" + orgPrefix + "id_" + selectedIndex).val());
    $("#e_subnetId").val($("#" + orgPrefix + "subnetId_" + selectedIndex).val());

  } else if (caller == "securityGroupAssist") {
    var orgPrefix = "securityGroupAssist_";
    var securityGroupIds = "";
    var count = 0;
    $("input[name='securityGroupAssist_chk']:checked").each(function () {
      var sgId = $(this).attr("id")
      var sgEleArr = sgId.split("_");
      var sgIndex = sgEleArr[sgEleArr.length - 1];

      if (count > 0) { securityGroupIds += "," }

      securityGroupIds += $("#" + orgPrefix + "id_" + sgIndex).val();
      count++;
    });

    $("#e_securityGroupIds").val(securityGroupIds);
  } else if (caller == "sshKeyAssist") {
    var orgPrefix = "sshKeyAssist_";

    $("#e_sshKeyId").val($("#" + orgPrefix + "id_" + selectedIndex).val());
  }

  console.log("apply connection =" + applyConnectionName);
  $("#e_connectionName").val(applyConnectionName);
  $("#es_regConnectionName").val(applyConnectionName);
}
// diffConnectionCheck
// 1. 다르면 -> confirm에서 OK면 -> applyAssistValues(caller) 호출하여 set되도록
// 2. 같으면 바로 applyAssistValues(caller) 호출하여 set
function applyAssistValidCheck(caller) {
  var selectedIndex = $("#assistSelectedIndex").val();

  // 선택한 connection check : 이미 선택된 connection이 있을 때 비교하여 다른 connection이면 confirm을 띄우고 OK면 초기화 시키고 set
  var selectedConnectionName = "";
  if (caller == "vmImageAssist") {
    var orgPrefix = "vmImageAssist_";
    selectedConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val();
  } else if (caller == "vmSpecAssist") {
    var orgPrefix = "vmSpecAssist_";
    selectedConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val();
  } else if (caller == "networkAssist") {
    var orgPrefix = "vNetAssist_";
    selectedConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val();
  } else if (caller == "securityGroupAssist") {
    var orgPrefix = "securityGroupAssist_";
    var tempConnectionName = "";
    var isSameConnection = true;
    $("input[name='securityGroupAssist_chk']:checked").each(function () {
      var sgId = $(this).attr("id")
      var sgEleArr = sgId.split("_");
      var sgIndex = sgEleArr[sgEleArr.length - 1];

      var currentConnectionName = $("#" + orgPrefix + "connectionName_" + sgIndex).val();
      if (tempConnectionName == "") {
        tempConnectionName = currentConnectionName
      } else if (tempConnectionName != currentConnectionName) {
        isSameConnection = false;
        return;
      }
    });

    if (!isSameConnection) {
      commonAlert("서로다른 ConnectionName이 선택되어 있습니다.");
      return;
    }
    selectedConnectionName = tempConnectionName;

  } else if (caller == "sshKeyAssist") {
    var orgPrefix = "sshKeyAssist_";
    selectedConnectionName = $("#" + orgPrefix + "connectionName_" + selectedIndex).val();
  }

  console.log("caller=" + caller)
  console.log($("#e_connectionName").val())
  console.log("selectedConnectionName=" + selectedConnectionName)
  $("#t_connectionName").val(selectedConnectionName);
  if ($("#e_connectionName").val() != "" && $("#e_connectionName").val() != selectedConnectionName) {
    //commonConfirmOpen("DifferentConnectionAtSecurityGroup");
    commonConfirmOpen("DifferentConnectionAtAssistPopup", caller)
  } else {
    applyAssistValues(caller)
  }
}

// Table에서 connection 선택시 hidden에 connection정보 set.
function setConnectionByAssist(providerName, regionName, configName) {
  $("#assistConnectionConfigName").val(configName);
}

// assist에서 apply버튼 클릭 시
function applyConnectionConfigName() {
  var selectedConnectionName = $("#assistConnectionConfigName").val();
  console.log("applyConnectionConfigName = " + selectedConnectionName);
  $("#es_regConnectionName").val(selectedConnectionName);

  $("#popProvider").val('');
  $("#popRegion").val('');
  $("#assistConnectionConfigName").val('');

  $("#Deployment_box").modal("hide");

  getConnectionListFilterForTable('popProvider', 'popRegion', 'popCloudConnection');
  console.log("getConnectionListFilterForTable = ");
}

// Textbox 값이 변경 된 경우 해당 값을 form obj에 set
function setTextValueToFormObj(setValue, targetObjId) {
  $("#" + targetObjId).val(setValue);
}

// TODO 1: Table의 Check Box를 여러개 체크하여 해당 값을 securityGroupIds 에 넣는 function : securityGroup에서 사용
// 동일한 Name을 가져야 함.
// 체크된 목록을 보여줄 objId
// form에 set할 objId
function setMuipleValueToFormObj(chkboxName, targetObjId, formObjId) {
  var checkedIds = "";
  var idandConnectionName = "";
  var securityGroupConnectionName = "";

  var tempConnectionNameValue = $("#t_connectionName").val();
  $('input:checkbox[name="' + chkboxName + '"]').each(function () {
    if (this.checked) {//checked 처리된 항목의 값
      var chkIdArr = $(this).attr('id').split("_");// 0번째와 2번째를 합치면 id 추출가능  ex) securityGroup_Raw_0
      console.log("setMuipleValueToFormObj = " + formObjId);
      var securityGroupId = $("#" + chkIdArr[0] + "_id_" + chkIdArr[2]).val()//id="securityGroup_id_{{$securityGroupIndex}}"
      securityGroupConnectionName = $("#" + chkIdArr[0] + "_connectionName_" + chkIdArr[2]).val()
      idandConnectionName += securityGroupId + "(" + securityGroupConnectionName + ")" + ",";
      checkedIds += securityGroupId + ",";
    }
  });

  idandConnectionName = idandConnectionName.substr(0, idandConnectionName.length - 1);
  checkedIds = checkedIds.substr(0, checkedIds.length - 1);

  $("#" + targetObjId).val(idandConnectionName);// 선택항목 display
  tempConnectionNameValue
  $("#" + formObjId).val(checkedIds);// set
  $("#e_connectionName").val(tempConnectionNameValue);
}
// TODO 2: 여러 filter 조건으로 table filter 하는 function 만들 것.
//         TODO3에서 사용하는 것으로 hidden obj에 값을 세팅하고
//         검색시 해당값과 &조건으로 filterling  
// TODO 3: connection  정보가 필수인데, 상단부분을 선택할 수도 있고 Tab에서 선택한 항목에서 connection정보를 가져올 수도 있다.
//       상단에서 변경한 경우 hidden connection에 값이 없으면 set
//       상단에서 변경한 경우 hidden connection에 값이 있으면
//          값을 비교하여 동일하면 continue
//          값디 다르면 confirm : 선택한 connection정보와 기존에 설정된 connection정보가 다릅니다. 
//              선택한 connection으로 할 경우 설정된 값들은 초기화 됩니다.
//              OK 누르면 값들 초기화 -> Tab에서 다시선택해야함
//       Tab에서 항목 조회시 hidden connection이 있으면
//          값을 비교하여 동일하면 continue
//          값이 다르면 confirm : 선택한 connection정보와 기존에 설정된 connection정보가 다릅니다. 
//              선택한 connection으로 할 경우 설정된 값들은 초기화 됩니다.
//              OK 누르면 값들 초기화 -> Tab에서 다시 선택해야함.

//////////////////// filterling 기능 //////////////
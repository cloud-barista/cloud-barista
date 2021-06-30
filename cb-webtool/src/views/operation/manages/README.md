
##
Expert Mode로 VM 추가하기


실제 vm을 생성은 form 안에 hidden으로 있음.

각 TAB( OS/HW, Network, Security, Other)은 각 조건들의 세부항목을 조회하여
원하는 ID를 조회하는 것임.

ex) Spec table에서 여러 조건을 filtering하여 알맞은 spec을 체크하면
  form 내 hiddend으로 정의된 spec에 값이 set.

Done 버튼을 누르면 deploy 대상으로 추가.


souce 
1. table
 . table id : prefix_이름List.  ex) id="es_imageList"
 . tr onclick event : setValueToFromObj( tableId, tr의 checkboxName, index, set할 값, set할 object Id)
   ex) <tr onclick="setValueToFormObj('es_imageList', 'vmImage_chk', '{{$vmInageIndex}}', '{{$vmImageItem.ID}}', 'e_imageId');">
 . td : checkbox, key값(ID), 전체값
        checkbox : name=이름_chk, id=이름Raw_index
        id obj : id=이름_id_index, value=해당 value
        info obj : id=이름_info_index, value는 | 를 구분자로 넣고자 하는 값들 모두

   ex)  <td class="overlay hidden" data-th="">
            <input type="checkbox" name="vmImage_chk" id="vmImageRaw_{{$vmInageIndex}}" title="" />
            <input type="hidden" id="vmImage_id_{{$vmInageIndex}}" value="{{$vmImageItem.ID}}"/>
            <input type="hidden" name="vmImageInfo" id="vmImage_info_{{$vmInageIndex}}" value="{{$vmImageItem.ID}}|{{$vmImageItem.Name}}|{{$vmImageItem.ConnectionName}}|{{$vmImageItem.CspImageId}}|{{$vmImageItem.CspImageName}}|{{$vmImageItem.GuestOS}}|{{$vmImageItem.Description}}"/>
            <label for="td_ch1"></label> <span class="ov off"></span>
        </td>
. table은 첫번째 Row는 th 로 header 이므로 로직에서는 0번째를 제외한 1번째부터 계산

3. 조회조건 정의
 . text obj : id=filter_이름, onkeydown="filterEnterToHidden( this.id, hidden obj id, 대상 table)
        hidden으로 정의된 객체의 값을 대상으로 filterling.
        엔터를 쳤을 때 작업 수행
 . 돋보기 버튼 : text obj의 단어로 대상 table을 filterling
        클릭 시 filterToHidden 호출
        onclick=filterToHidden( keyword를 가진 객체, hidden 객체 id, 대상 table)

<span class="sbox">
    <input type="text" name="" value="" placeholder="Filter Items" class="pline ip_1 search_ip" id="filter_image" onkeydown="filterEnterToHidden(this.id, 'vmImageInfo', 'es_imageList');"/>
    <input type="submit" name="" class="btn_search" value="" title="" onclick="filterToHidden('filter_image', 'vmImageInfo', 'es_imageList')"/>
</span>


* 일반적으로 table filter는 대상테이블에 TH 이름 칼럼을 기준으로 동작함
* expert mode의 table filter는 hidden obj에 넣고 싶은 값을 모두 넣고 이를 기준으로 동작시킴 ( util.filterTableByHiddenColumn)
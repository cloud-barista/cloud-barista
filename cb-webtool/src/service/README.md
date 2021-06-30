API 호출기준 CRUD 에 따른 명명규칙

1. 조회(목록) : GetXXXList
2. 조회(항목) : GetXXXData
3. 등록 : RegXXX
4. 삭제 : DelXXX

CommonHandler

NameSpaceHandler

- GetNameSpaceList
- RegNameSpace
- CreateDefaultNameSpace : Namespace가 없는경우 기본으로 1개의 namespace를 자동으로 생성
- DelNameSpace

CloudConnectionHandler
ResourceHandler
McisHandler

---

handler의 return 값중 두번째 인자는 model.WebStatus 로 한다.
WebStatus.Status 는 code를, Message에는 message를
error가 났을 때 Status = 500, Message에는 error의 값을

정상적으로 호출했으나 해당 내용이 Error일 때는 최종 수신단(UI)에서 StatusCode에 따라 결정한다.

---

TODO

1. TB : lookup, search 등 method 호출 테스트 필요(UI등 에서 어떻게 사용될 지)
2. TB : lifecycle 변경의 경우 TB API명세에 없는데 호출 됨.
3. TB : lifecycle 호출하는 handlerMethod 명 변경해야하나?? Get, Reg, Del 외에...

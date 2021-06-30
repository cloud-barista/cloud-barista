DB model

- 대상 system에 따라 folder 로 분류 : dragonfly, ladybug, spider, tumblebug
- 기본 struct : 업무Info 등
- 조회(request용) struct : 업무reqInfo 등 기본 struct에서 요청에 필요한 것들만 추출

- list의 경우 [] 를 사용한다.
- 첫번째 자리 : 변수명
- 두번째 자리 : 자료형
- 세번째 자리 : json객체명
  type InspectResourcesResponse struct {
  ResourcesOnCsp ResourcesOnCsp `json:"resourcesOnCsp"`
  ResourcesOnSpider ResourcesOnSpider `json:"resourcesOnSpider"`
  ResourcesOnTumblebug []ResourcesOnTumblebug `json:"resourcesOnTumblebug"`
  }

. Renderling 할 때에는 변수명을 사용하나
axios 등 통신을 할 경우에는 json객체명을 사용

. 변수명의 경우 첫글자 대문자
. json 인 경우 첫글자 소문자

# Entity

## Key
```
/ns/{namespace}/clusters/{cluster}
```

## Value
```
  {
    kind: "Cluster",
    name: "",
    status: {
      phase: "",
      reason: "",
      message: ""
    },
    mcis: "",
    namespace: "",
    clusterConfig: "",
    cpLeader: "",
    networkCni: "",
    label: "",
    installMonAgent: "",
    description: "",
    createdTime: "",
    nodes: [
      {
        name: "",
        credential: "",
        publicIp: "",
        role: "control-plane",
        spec: "",
        csp: "",
        createdTime: "",
        cspLabel: "",
        regionLabel: "",
        zoneLabel: "",
      },
      {
        name: "",
        credential: "",
        publicIp: "",
        role: "worker",
        spec: "",
        csp: "",
        createdTime: "",
        cspLabel: "",
        regionLabel: "",
        zoneLabel: "",
      },
      ...
    ]
  }
```

---
## Cluster
> 클러스터 정보

|속성               |이름                          |타입   |비고                                   |
|---                |---                         |---    |---                                  |
|kind               |종류                         |string |Cluster                              |
|name               |클러스터 명                    |string |                                     |
|status             |프로비저닝 상태                 |object |                                     |
|status.phase       |프로비저닝 단계                 |string |아래 "ClusterPhase" 참조               |
|status.reason      |프로비저닝 오류                 |string |아래 "ClusterReason" 참조              |
|status.message     |프로비저닝 오류 메시지            |string |                                     |
|mcis               |MCIS 명                      |string |                                     |
|namespace          |MCIS 네임스페이스               |string |                                     |
|clusterConfig      |클러스터 연결정보                |string |Kubernetes 인 경우 kubeconfig.yaml     |
|cpLeader           |control plane leader 노드명   |string |                                     |
|networkCni         |network CNI 정보             |string |                                     |
|label              |label                       |string |                                     |
|installMonAgent    |모니터링 에이전트 설치 여부        |string | yes/no (no가 아니면 설치)              |
|description        |description                 |string |                                     |
|createdTime        |생성일자                      |string |                                     |

### ClusterPhase
> 프로비저닝 단계

* Pending
* Provisioning
* Provisioned
* Failed
* Deleting

### ClusterReason
> 프로비저닝 오류 원인 (Phase == Failed 경우)

* GetMCISFailedReason : MCIS 조회 실패
* CreateMCISFailedReason : MCIS 생성 실패
* AlreadyExistMCISFailedReason : 이미 존재하는 MCIS 
* GetControlPlaneConnectionInfoFailedReason : ControlPlane 노드 생성위한 클라우드 연결정보 조회 실패
* GetWorkerConnectionInfoFailedReason : Worker 노드 생성위한 클라우드 연결정보 조회 실패
* CreateVpcFailedReason : VCP 생성 실패
* CreateSecurityGroupFailedReason : Security Group 생성 실패
* CreateSSHKeyFailedReason : SSH 키 생성 실패
* CreateVmImageFailedReason : VM 이미지 생성 실패
* CreateVmSpecFailedReason : VM 타입 생성 실패
* SetupBoostrapFailedReason : OS 기본 패키지 설치 실패
* SetupHaproxyFailedReason : HAProxy 설치 실패
* InitControlPlaneFailedReason : ControlPlane Init. 실패
* SetupNetworkCNIFailedReason : Network CNI 설치 실패
* JoinControlPlaneFailedReason : ControlPlane join 실패
* JoinWorkerFailedReason : Worker 노드 join 실패

## Node
> 클러스터의 노드 정보

|속성           |이름               |타입   |비고                 |
|---            |---                |---    |---                  |
|kind           |종류               |string |Node                 |
|name           |노드명             |string |mcis vm 이름과 동일  |
|credential     |private key        |string |                     |
|publicIp       |공인 IP            |string |                     |
|role           |역할               |string |control-plane/worker |
|spec           |spec               |string |                     |
|csp            |csp 정보           |string |                     |
|createdTime    |생성일자            |string |                    |
|cspLabel       |CSP Label         |string |<label_key>=<label_value> |
|regionLabel    |Region Label      |string |<label_key>=<label_value> |
|zoneLabel      |Zone Label        |string |<label_key>=<label_value> |

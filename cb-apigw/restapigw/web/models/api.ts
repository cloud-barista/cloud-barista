import Util from "@/utils";

/****************************************************************
 * Define models for APIs
 ****************************************************************/

export class HostConfig {
  host: string = ""; // 호스트 도메인
  weight: number = 0; // 가중치 (정수)
}

export class HealthCheckConfig {
  url: string = ""; // 상태 검증용 URL
  timeout: any = "1s"; // 처리 제한 시간
  public AdjustSendValues() {
    this.timeout = Util.timeParser.ToDuration(this.timeout);
  }
  public AdjustReceiveValues() {
    this.timeout = Util.timeParser.FromDuration(this.timeout, "s");
  }
}

export class BackendConfig {
  hosts: Array<HostConfig> = []; // 백엔드 서비스 도메인 리스트 (Defintion에서 전역으로 설정한 경우는 생략 가능)
  timeout: any = "3s"; // 처리 제한 시간 (Defintion에서 전역으로 설정한 경우는 생략 가능)
  method: string = "GET"; // 호출 메서드 (Definition에서 전역으로 설정한 경우는 생략 가능)
  url_pattern: string = ""; // 서비스 URL (Domain 제외)
  encoding: string = "json"; // 데이터 인코딩 (xml, json : Definition에서 전역으로 설정한 경우는 생략 가능)
  group: string = ""; // 결과를 묶어서 처리할 그룹 명
  blacklist: Array<string> = []; // 결과에서 제외할 필드명 리스트 (flatmap 적용방식 "." 오퍼레이션 가능)
  whitelist: Array<string> = []; // 결과에서 포함할 필드명 리스트 (flatmap 적용방식 "." 오퍼레이션 가능)
  mapping: Map<string, string> = new Map(); // 결과에서 필드명을 변경할 리스트
  is_collection: boolean = false; // 결과가 Collection인지 여부
  wrap_collection_to_json: boolean = false; // 결과 Collection을 JSON의 "collection = []" 방식으로 처리할지 여부, false명 Array 상태로 처리
  target: string = ""; // 결과 중 특정한 필드만 처리할 경우의 필드명
  middleware?: any = {};
  disable_host_sanitize: boolean = false; // host 정보를 정제 작업할지 여부
  lb_mode: string = ""; // 백엔드 로드밸런싱 모드 ("rr", "wrr")

  public AdjustSendValues(def: ApiDefinition) {
    // Host 값 조정
    if (this.hosts.length === 0) {
      this.hosts = def.hosts;
    }
    // Method 값 조정
    if (this.method === "") {
      this.method = def.method;
    }
    // Encoding 값 조정
    if (this.encoding === "") {
      this.encoding = "json";
    }
    this.timeout = Util.timeParser.ToDuration(this.timeout);
  }
  public AdjustReceiveValues(_: ApiDefinition) {
    this.timeout = Util.timeParser.FromDuration(this.timeout, "s");
  }
  public Validate(def: ApiDefinition): string {
    if (this.url_pattern === "")
      return "Backend 서비스 URL을 지정하셔야 합니다.";
    if (this.encoding === "")
      return "Backend 서비스 Encoding을 지정하셔야 합니다.";
    if (!this.hosts || this.hosts.length === 0) {
      return `Backend에 Host 정보를 설정하셔야 합니다. (${def.name} - URL: ${this.url_pattern})`;
    }
    return "";
  }
}

export class ApiDefinition {
  name: string = ""; // 엔드포인트 식별 명
  active: boolean = false; // 엔드포인트 활성화 여부
  endpoint: string = ""; // 엔드포인트 URL
  hosts: Array<HostConfig> = []; // 백엔드 서비스 도메인 리스트
  method: string = "GET"; // 호출 메서드
  timeout: any = "1m"; // 처리 제한 시간
  cache_ttl: any = "3600s"; // 캐시 TTL
  output_encoding: string = "json"; // 데이터 인코딩 (xml, json)
  except_querystrings: Array<string> = []; // 벡엔드로 전달하지 않을 Query String 리스트
  except_headers: Array<string> = []; // 벡엔드로 전달하지 않을 Header 리스트
  middleware?: any = {};
  health_check: HealthCheckConfig = new HealthCheckConfig(); // 헬스 검증용 설정
  backend: Array<BackendConfig> = []; // 백엔드 설정

  public AdjustSendValues() {
    this.timeout = Util.timeParser.ToDuration(this.timeout);
    this.cache_ttl = Util.timeParser.ToDuration(this.cache_ttl);
    this.health_check.AdjustSendValues();
    this.backend.forEach(b => b.AdjustSendValues(this));
    if (this.output_encoding === "") this.output_encoding = "json";
  }
  public AdjustReceiveValues() {
    this.timeout = Util.timeParser.FromDuration(this.timeout, "s");
    this.cache_ttl = Util.timeParser.FromDuration(this.cache_ttl);
    this.backend.forEach(b => b.AdjustReceiveValues(this));
    this.health_check.AdjustReceiveValues();
  }
  public Validate(): string {
    if (this.name === "") return "API Definition 이름을 지정하셔야 합니다.";
    if (this.endpoint === "")
      return "API Definition의 Endpoint를 지정하셔야 합니다.";
    if (this.method === "")
      return "API Definition의 호출 메서드를 지정하셔야 합니다.";
    if (this.backend.length === 0)
      return "API Definition의 Backend 서비스 정보를 등록하셔야 합니다.";
    for (let i = 0; i < this.backend.length; i++) {
      const error = this.backend[i].Validate(this);
      if (error !== "") return error;
    }

    return "";
  }
}

export class ApiGroup {
  name: string = "";
  definitions: Array<ApiDefinition> = [];

  public AdjustSendValues() {
    this.definitions.forEach(d => d.AdjustSendValues());
  }
  public AdjustReceiveValues() {
    this.definitions.forEach(d => d.AdjustReceiveValues());
  }
  public Validate(): string {
    if (this.name === "") return "API Group 이름을 지정하셔야 합니다.";
    if (this.definitions.length === 0)
      return "API Definition 정보가 존재하지 않습니다.";
    for (let i = 0; i < this.definitions.length; i++) {
      const error = this.definitions[i].Validate();
      if (error !== "") return error;
    }
    return "";
  }
}

export function deserializeGroupFromJSON(
  val: any,
  isReceived: boolean = false
) {
  const group: ApiGroup = Object.assign(new ApiGroup(), val);
  group.definitions = group.definitions.map(d => {
    const def: ApiDefinition = Object.assign(new ApiDefinition(), d);
    def.health_check = Object.assign(new HealthCheckConfig(), def.health_check);
    def.backend = def.backend.map(b => Object.assign(new BackendConfig(), b));
    return def;
  });
  if (isReceived) group.AdjustReceiveValues();
  else group.AdjustSendValues();

  return group;
}

export function deserializeDefinitionFromJSON(
  val: any,
  isReceived: boolean = false
) {
  const def: ApiDefinition = Object.assign(new ApiDefinition(), val);
  def.health_check = Object.assign(new HealthCheckConfig(), def.health_check);
  def.backend = def.backend.map(b => Object.assign(new BackendConfig(), b));
  if (isReceived) def.AdjustReceiveValues();
  else def.AdjustSendValues();

  return def;
}

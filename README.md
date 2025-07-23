# Cloud-Barista

[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cloud-barista/releases/latest)
[![Pre Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=brightgreen&include_prereleases&label=release%28dev%29)](https://github.com/cloud-barista/cloud-barista/releases)
[![License](https://img.shields.io/github/license/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cb-tumblebug/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-Cloud--Barista-brightgreen)](https://join.slack.com/t/cloud-barista/shared_invite/zt-bda8zhkg-tlOCr7_TdQGE_oUSz4mlkA)

*The Cloud-Barista is a Multi-Cloud Service Platform SW.*

Cloud-Barista consists of multiple frameworks (sub-systems) to accommodate microservice-like architecture.

Please take a look [Cloud-Barista Website](https://cloud-barista.github.io/technology/) for a detail decription.

<details>
<summary>Note for developing and using Cloud-Barista</summary>

#### Development stage of Cloud-Barista

```
Cloud-Barista is currently under development. (not v1.0 yet)
We welcome any new suggestions, issues, opinions, and controbutors !
Please note that the functionalities of Cloud-Barista are not stable and secure yet.
Becareful if you plan to use the current release in production.
If you have any difficulties in using Cloud-Barista, please let us know.
(Open an issue or Join the Cloud-Barista Slack)
```

</details>

---
*This repository is an integrated archive for repository of major frameworks.* These repositories are included and listed in the root directory. This repo reflects the latest release only.

The main frameworks or tools are as follows (the release version of each repository may vary),

- **CB-Tumblebug** (manages multi-cloud infrastructures)
  - Upstream repo: <https://github.com/cloud-barista/cb-tumblebug>
    - Use CB-TB to run all basic Cloud-Barista components
- **CB-Spider** (connects all clouds in a single interface)
  - Upstream repo: <https://github.com/cloud-barista/cb-spider>
- CB-Bridge/cb-log (provides log library to Cloud-Barista system)
  - Upstream repo: <https://github.com/cloud-barista/cb-log>  
- CB-Bridge/cb-store (provides an unified DB interface for meta info of Cloud-Barista)
  - Upstream repo: <https://github.com/cloud-barista/cb-store>  
- CB-Bridge/cb-operator (operation tool for Cloud-Barista system runtime)
  - Upstream repo: <https://github.com/cloud-barista/cb-operator>
  - Note: As the components of Cloud-Barista are currently simplified, using cb-operator might require additional effort for users. It is recommended to use cb-tumblebug directly.
  - Note: not updated, possible to be removed  
- (Deprecated) CB-Dragonfly (monitors multi-cloud services)
  - Upstream repo: <https://github.com/cloud-barista/cb-dragonfly>
  - Note: not updated since v0.8.0, possible to be removed

***

## [Installation and Execution]

- Cloud-Barista Platform Execution 
  - CB-Tumblebug Installation and Execution (Docker Compose based)
    - [Refer to cloud-barista/cb-tumblebug README for configuration and installation](https://github.com/cloud-barista/cb-tumblebug)
    - Run CB-Tumblebug (CB-Spider runs simultaneously)


- Cloud-Barista Platform Individual Framework Source Download and Installation

  - CB-Tumblebug Installation and Execution
    - [Refer to cloud-barista/cb-tumblebug README for configuration and installation](https://github.com/cloud-barista/cb-tumblebug)
      - Configure CB-Spider API server address in conf/setup.env
      - Configure CB-Dragonfly API server address in conf/setup.env
    - Run CB-Tumblebug

  - CB-Spider Installation and Execution
    - [Refer to cloud-barista/cb-spider README for configuration and installation](https://github.com/cloud-barista/cb-spider)
    - Run CB-Spider


***

## [Notable Information]

- Development Stage: Function development priority stage (stabilization and supplementation required for commercial use)
- CSP Integration Verification Status
  - CSPs tested and completed based on CB-Spider: Refer to [Link 1](https://github.com/cloud-barista/cb-spider#3-제공-자원) and [Link 2](https://github.com/cloud-barista/cb-spider/wiki/Supported-CloudOS)
  - CSPs tested and completed based on CB-Tumblebug: Refer to [Link](https://github.com/cloud-barista/cb-tumblebug?tab=readme-ov-file#cb-tb-)
  - Currently in development stage, so functional stability may be low (bug reports and contributions are welcome!)

***

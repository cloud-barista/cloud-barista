# Cloud-Barista

[![Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cloud-barista/releases/latest)
[![Pre Release Version](https://img.shields.io/github/v/release/cloud-barista/cb-tumblebug?color=brightgreen&include_prereleases&label=release%28dev%29)](https://github.com/cloud-barista/cloud-barista/releases)
[![License](https://img.shields.io/github/license/cloud-barista/cb-tumblebug?color=blue)](https://github.com/cloud-barista/cb-tumblebug/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-Cloud--Barista-brightgreen)](https://join.slack.com/t/cloud-barista/shared_invite/zt-bda8zhkg-tlOCr7_TdQGE_oUSz4mlkA)

*The Cloud-Barista is a Multi-Cloud Service Platform SW.*

Cloud-Barista is an open source multi-cloud platform designed to make cross cloud operations simple and consistent.
It provides common APIs and core components to connect to different cloud providers and manage infrastructure in a unified way.
This helps teams build portable cloud environments and reduce provider specific complexity.

Cloud-Barista currently consists of the following main components:
 - CB-Spider: A multi-cloud interface that connects to various cloud providers.
 - CB-Tumblebug: A multi-cloud orchestration that manages multi-cloud infrastructures

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
- (Deprecated) CB-Bridge/cb-operator (operation tool for Cloud-Barista system runtime)
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

***

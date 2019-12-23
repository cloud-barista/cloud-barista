package util

import (
	"fmt"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/vm-ssh"
)

func RunCommand(vmIP string, userName string, privateKey string, cmd string) (*string, error) {

	// VM SSH 접속정보 설정 (외부 연결 정보, 사용자 아이디, Private Key)
	serverEndpoint := fmt.Sprintf("%s:22", vmIP)
	sshInfo := sshrun.SSHInfo{
		ServerPort: serverEndpoint,
		UserName:   userName,
		PrivateKey: []byte(privateKey),
	}

	// VM SSH 명령어 실행
	if result, err := sshrun.SSHRun(sshInfo, cmd); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

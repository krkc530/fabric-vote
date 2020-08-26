# fabric-vote

## 1. Network 실행
`cd network`

`./start.sh`

네트워크 구성 정보는 connection.json 참조

## 2. Chaincode 설치 및 배포, 테스트
`cd network`

`./cc.sh`

## 3. NodeJS App
디렉토리 이동

`cd application`

### 업로드
`node vote upload -k key -f filePath -t tag(Optional)`

태그 옵션이 없으면 태그가 default로 설정됩니다. 

### 다운로드
`node vote download -k key -f filePath`

키에 해당하는 파일을 filePath에 저장합니다.

### 리스트(검색)
`node vote list -t tag(Optional)`

태그 옵션이 없으면 모든 파일의 리스트를 가져옵니다.

### 파일 내용 확인
`node vote show -k key`

리스트와 기능이 동일하지만 키를 통해서 단일 파일을 호출합니다.

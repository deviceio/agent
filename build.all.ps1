$env:GOOS='windows';$env:GOARCH='amd64';go install std; go build -o $env:GOPATH\bin\deviceio-agent.windows.amd64.exe github.com/deviceio/agent 
$env:GOOS='windows';$env:GOARCH='386';go install std; go build -o $env:GOPATH\bin\deviceio-agent.windows.386.exe github.com/deviceio/agent 

$env:GOOS='linux';$env:GOARCH='amd64';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.amd64 github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='386';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.386 github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='ppc64';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.ppc64 github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='ppc64le';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.ppc64le github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='mips';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.mips github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='mipsle';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.mipsle github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='mips64';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.mips64 github.com/deviceio/agent
$env:GOOS='linux';$env:GOARCH='mips64le';go install std; go build -o $env:GOPATH\bin\deviceio-agent.linux.mips64le github.com/deviceio/agent

$env:GOOS='darwin';$env:GOARCH='amd64';go install std; go build -o $env:GOPATH\bin\deviceio-agent.darwin.amd64 github.com/deviceio/agent
$env:GOOS='darwin';$env:GOARCH='386';go install std; go build -o $env:GOPATH\bin\deviceio-agent.darwin.386 github.com/deviceio/agent
windows: build_deploy_windows
freebsd: build_deploy_freebsd
build_deploy_windows:
	@bash -c "echo Building for Windows && GOOS=windows GOARCH=amd64 go build && ssh -t windows 'TASKKILL /IM server.exe /F || true' 2>/dev/null && echo Transferring && scp -q ./server.exe windows: && echo Starting && ssh -t windows './server &' 2> /dev/null && echo Deployed"
build_deploy_freebsd:
	@bash -c "echo Building for FreeBSD && GOOS=freebsd GOARCH=386 go build && ssh -t freebsd 'sudo pkill -f \"sudo ./server\" || true' 2>/dev/null && echo Transferring && scp -q ./server freebsd: && scp -q -r static freebsd: && echo Starting && ssh -t freebsd 'nohup sudo ./server &' 2>/dev/null && echo Deployed"

set -e
alias acbuild="/home/faraz/acbuild"
acbuild begin
trap "{ export EXT=$?; acbuild --debug end && exit $EXT;   }" EXIT
acbuild set-name server
acbuild --debug copy-to-dir server /srv
acbuild --debug set-exec /srv/server
acbuild label add version 0.0.1
acbuild annotation add authors "Faraz Fazli <farazfazli@gmail.com>"
acbuild write --overwrite server.aci
acbuild end

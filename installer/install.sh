#!/bin/bash

cp -v -f ./go_ws_sh-tls.service /etc/systemd/system/

cp -v -f ./go_ws_sh.service /etc/systemd/system/
mkdir -p -v /root/opt/go_ws_sh/server/
cp -v -f ./loop.sh /root/opt/go_ws_sh/server/

cp -v -f ./loop-tls.sh /root/opt/go_ws_sh/server/




chmod +x /root/opt/go_ws_sh/server/loop.sh
chmod +x /root/opt/go_ws_sh/server/loop-tls.sh


cp -v -f ./server-config.json /root/opt/go_ws_sh/server/


cp -v -f ./server-config-tls.json /root/opt/go_ws_sh/server/

cp -v -f ./main /root/opt/go_ws_sh/server/

chmod +x /root/opt/go_ws_sh/server/main



cp -v -f ./*.crt /root/opt/go_ws_sh/server/


cp -v -f ./*.key /root/opt/go_ws_sh/server/

systemctl daemon-reload

systemctl enable go_ws_sh.service

systemctl enable go_ws_sh-tls.service

systemctl restart go_ws_sh.service

systemctl restart go_ws_sh-tls.service
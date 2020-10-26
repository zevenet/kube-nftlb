s/\" : \"/\": \"/g
s/  / /g
s/\("virtual-addr": \)\(.*\)/\1"IP",/g
s/\("ip-addr": \)\(.*\)/\1"IP",/g
s/\("mark": \)\(.*\)/\1"MARK",/g
s/\(^[[:space:]]\{20\}"name": \)\(.*\)/\1"NAME",/g

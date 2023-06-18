create user

curl -X POST -d '{"username":"binbomb","password":"secret","full_name":"binbomb","email":"binbomb@gmail.com"}' http://127.0.0.1:8888/users

user login to create account

curl -X POST  -d '{"username":"binbomb","password":"secret"}' http://127.0.0.1:8888/users/login

create account with currency

curl -H "Authorization: Bearer v2.local.mZ1y9Rhs5ZcfCpqsiM_7noAS2ztyjC1yDoj1oko95_6v2g8dUikpdjSNgM-VblzBleR1uyOgmtt01pNnun-RAhJRIyLX_ZnCf5_5ah63HXEphZkFBa-w3sbLO_7tpvUZMQ62ZRZTfZU3PhQqO9Y7GeB70tT26FzFhzuoaQ0mxxJcP7RVE9CZ3sN4Whbc.bnVsbA" -d '{"currency":"USD"}' "http://127.0.0.1:8888/accounts" -vvv


get account
curl -H "Authorization: Bearer v2.local.mZ1y9Rhs5ZcfCpqsiM_7noAS2ztyjC1yDoj1oko95_6v2g8dUikpdjSNgM-VblzBleR1uyOgmtt01pNnun-RAhJRIyLX_ZnCf5_5ah63HXEphZkFBa-w3sbLO_7tpvUZMQ62ZRZTfZU3PhQqO9Y7GeB70tT26FzFhzuoaQ0mxxJcP7RVE9CZ3sN4Whbc.bnVsbA"  "http://127.0.0.1:8888/accounts/1" -vvv

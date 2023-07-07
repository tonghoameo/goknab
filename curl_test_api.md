create user

curl -X POST -d '{"username":"balance4","password":"secret","full_name":"balance","email":"balance4@gmail.com"}' http://127.0.0.1:8888/v1/create_user
curl -X POST -d '{"username":"duymeo","password":"secret","full_name":"duymeo","email":"duymeo@gmail.com"}' http://127.0.0.1:8888/v1/create_user
curl -X POST -d '{"username":"tonhoa","password":"secret","full_name":"tonhoa","email":"tonhoa@gmail.com"}' http://127.0.0.1:8888/v1/create_user
user login to create account

curl -X POST  -d '{"username":"balance4","password":"secret"}' http://127.0.0.1:8888/users/login
curl -X POST  -d '{"username":"duymeo","password":"secret"}' http://127.0.0.1:8888/users/v1/login_user



create account with currency

curl -H "Authorization: Bearer v2.local.WGfFYKs0A4Y_QasAOXB88xSksBKuLVoNTChpdKQCb_vsxYzcT-PL5BhiJXhUyP6e0jIAcFhUfEHN3qu5lTv-5CQA-n0d15_lyQSndSS0GvslB-iIN3jPP_XsTZhw7FePRTAwN3WTfMQaSVeAiRrHIbYFVwB3BNn8DOwjbvMImop-MORnoKnw6FlwKqE.bnVsbA" -d '{"currency":"USD"}' "http://127.0.0.1:8888/accounts" -vvv

tonho
v2.local.c_wxId-iZXkXfyZsFz6bw6A_bsm8rwwLrdyqSg_JadqA25jJdQiEmYRY-6-9wOtj8SwezuKQ0wPnipYiu-kFiggVDVAskoA67izmqm41L2zKoTCmr7CwGDdQ2fRdbm-ZaBLWxayoYJl7-P5nj2b-zbr6Rz7-60SiO-U5D1jyEAe29TAs4Uu3PZtG9Q.bnVsbA
duymeo  -H "Authorization: Bearer v2.local.LEa4eHz8ZTQ56PEtik1dcqbCHiIsYtfPun1Kx-0Fg_BqmPOH371ZQQsUpb1oG2g5Od1UtPp_lrX_O80GA2IJSvP8Ba155N2apDETc6BCO09riighdJ1-oud259Km6G4tRLFHvuzwsf-QO-0Pm9zxP0TpGNl-yyG5l2cs-i-_Q_ib6nnPaIJzEkDNk54.bnVsbA"

balance4 

v2.local.RBlAhCitWV-F7ZWE-skrgO8rLG08VaxVuhu68sdwvH1eT2N3rE8nmZ4MXHlReF0T1w9Zzs_tTnIeIBzmUJjyaQZEMJZuRoBqXFxOb0u0u_-W00oOTlK2o7M_lzDH5BrY7UGI1YXGWopNghVQQDP2nlJBS9j-hA3bCRNeXEVtWmkRgJDbDU6mOSscbS8LMA.bnVsbA

get account
curl -H "Authorization: Bearer v2.local.WGfFYKs0A4Y_QasAOXB88xSksBKuLVoNTChpdKQCb_vsxYzcT-PL5BhiJXhUyP6e0jIAcFhUfEHN3qu5lTv-5CQA-n0d15_lyQSndSS0GvslB-iIN3jPP_XsTZhw7FePRTAwN3WTfMQaSVeAiRrHIbYFVwB3BNn8DOwjbvMImop-MORnoKnw6FlwKqE.bnVsbA"  "http://127.0.0.1:8888/accounts/3" -vvv


curl -X POST  -d '{"refresh_token":"v2.local.H8ML25BEvM-6ztTWUFC-8xTuhLx6-8qDm30l_R14ULpLelFdFdgZJXsYRwdG92H0u2njbKV3cvFXKbGMUlc4yFdodDPfmTbFS1D3V05FgnP0Z2uBberKX9VuI7UJpv4htrdIWlQJ3I6nDUQ0HQVZ_6Fmt-RmJnPwI_L7ajeJd9z1wJjYy_8u3lUOul4.bnVsbA"}' http://127.0.0.1:8888/tokens/renew_access
"Authorization: Bearer v2.local.MyRfcecpGnzJPmw_LTUYwkTbZKvtW3BdsEmagW6w4KA6NZXDUHfjk6Cp9VqaPpVN-tjdRokFBZMqqq6nOfABH_OxprQvnLcotsDWZxTSE7G0h8Vinm_Ja0xqOI1ftoFUwKRJfvbfnxnuZ9bHo4Ow-xzZL6-_MBGZBdqBQYCQC1BJCYjtmcFUYmJg9tA.bnVsbA"  

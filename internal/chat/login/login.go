/*
Mkae the login and start the session 


The key and the field -- 

| Redis Key             | Data Structure | Fields (example)                                                         |
| --------------------- | -------------- | ------------------------------------------------------------------------ |
| `session:abc123:Abhi` | Hash           | `ws_connected: 1` <br> `last_seen: 2025-07-08T20:00:00` <br> `notify: 0` |
| `session:abc123:Anny` | Hash           | `ws_connected: 0` <br> `last_seen: 2025-07-08T19:55:00` <br> `notify: 1` |
| `session:def456:Bob`  | Hash           | `ws_connected: 1` <br> `last_seen: 2025-07-08T20:01:00` <br> `notify: 0` |
| `session:def456:Cara` | Hash           | `ws_connected: 0` <br> `last_seen: 2025-07-08T19:50:00` <br> `notify: 1` |


HSET session:abc123:Abhi chat_id abc123
HSET session:abc123:Abhi user Abhi
HSET session:abc123:Abhi last_seen 2025-07-08T20:00:00
HSET session:abc123:Abhi ws_connected 1
HSET session:abc123:Abhi notify 0

when login store the session data into the redis
and start hearbeat protocol until logout or endchat

*/

/*
func LoginUser 

params --> hash , user 

hash comes from the link

match the user in the pgsql
create the login
store the session data into the redis 
and start the hearbeat protocol to check the session 
and keep updating the status
pass a success message if user is valid to front-end 

*/




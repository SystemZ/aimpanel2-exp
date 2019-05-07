# Master

## API Handler Errors

|Endpoint|Code Range|
|---|---|
|auth|1000-1999|
|user|2000-2999|
|host|3000-3999|
|game|4000-4999|
|game_server|5000-5999|

## 5xx
|HTTP|Code|Msg|
|---|---|---|
|503|20003|DB is fucked|
|503|20004|CPU is grilling|
|502|20003|Admin commited suicide|
|500|20002|host exploded|

## 4xx
|HTTP|Code|Msg|
|---|---|---|
|420|10003|customer is mad, enhace your calm https://httpstatusdogs.com/420-enhance-your-calm|

# Slave

## CLI errors
|Code|Msg|
|---|---|
|112|Data on disk is too fat|


# Auth
* 1001 - Invalid body.
* 1002 - Passwords do not match.
* 1003 - Emails do not match
* 1004 - User with this username already exist.
* 1005 - Something went wrong.
* 1006 - Something went wrong.
* 1007 - Invalid body.
* 1008 - Something went wrong.
* 1009 - Wrong username or password.

# User
* 2001 - Invalid body.
* 2002 - Passwords do not match.
* 2003 - Current password is wrong.
* 2004 - Invalid body.
* 2005 - Emails do not match.
* 2006 - Current email is wrong.
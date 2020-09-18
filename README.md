# firebase-messaging-notifications-flutter-golang

Confirm that you have your Credential configuration from firebase
===================================================================
Expected output
```
YYYY/MM/DD hh:mm:ss &{map[]]  <PROJECT-ID>   
 [][]https://www.googleapis.com/auth/cloud-platform 
 https://www.googleapis.com/auth/datastore 
 https://www.googleapis.com/auth/devstorage.full_control 
 https://www.googleapis.com/auth/firebase 
 https://www.googleapis.com/auth/identitytoolkit 
 https://www.googleapis.com/auth/userinfo.email]]}}
```

Execute golang server
===================== 
```
$ cd firebase-messaging-notifications-flutter-golang/golang
$ go run main.go
```
```
USAGE:

        curl localhost:8080/sendToToken -d '{
            "data": {"k": "v"},
            "notification": {"title": "t", "body": "b"},
            "token":"<YOUR TOKEN>"
            }'

        curl localhost:8080/sendToTokens -d '{
            "data": {"k": "v"},
            "notification": {"title": "t", "body": "b"},
            "tokens": ["<TOKEN-0>", "<TOKEN-1>", ...]
            }'

Starting Firebase Cloud Messaging Notification Generator server on port: 8080
127.0.0.1:8080
```


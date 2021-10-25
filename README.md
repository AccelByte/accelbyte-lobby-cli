# Justice Lobby CLI

## Overview
This is a lobby CLI, allow Accelbyte client to simulate all the lobby features

## Run requirement
* Install GO, please refer to this url https://golang.org/doc/install

## How to run
* this project need some environment variables, you can define it directly in the Operation System or in `.env` file, 
  please check the ```.envExample``` for the sample
* go to the main directory and run ```go run *.go```

### Environment Variables
| Name                       | Description                                   | Example Value                                                        |
|----------------------------|-----------------------------------------------|----------------------------------------------------------------------|
| LOBBY_BASE_URL             | Lobby websocket url                           | wss://the-lobby-url                            |
| IAM_BASE_URL               | Redis lobby address port                      | https://the-iam-url                                                                 |
| QOS_BASE_URL               | Redis lobby address password                  | https://the-qos-url                                                       |
| IAM_CLIENT_ID              | Redis matchmaking address host                | 2ca0636b03154050ac85f771e978e44c                            |
| IAM_CLIENT_SECRET          | Redis matchmaking address port                | client-secret-if-any                                                 |

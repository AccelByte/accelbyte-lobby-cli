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
| LOBBY_BASE_URL             | Lobby websocket url                           | wss://demo.accelbyte.io/lobby/                                       |
| IAM_BASE_URL               | IAM base url                                  | https://demo.accelbyte.io/iam                                        |
| QOS_BASE_URL               | QOS base url                                  | https://demo.accelbyte.io/qosm                                       |
| IAM_CLIENT_ID              | IAM client ID                                 | 2ca0636b03154050ac85f771e978e44c                                     |
| IAM_CLIENT_SECRET          | IAM client secret                             | client-secret-if-any                                                 |

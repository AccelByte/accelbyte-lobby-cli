// Copyright (c) 2021 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	infoCmd                       = "1"
	createCmd                     = "2"
	leaveCmd                      = "3"
	inviteCmd                     = "4"
	joinCmd                       = "5"
	kickCmd                       = "6"
	chatCmd                       = "7"
	partyChatCmd                  = "8"
	getOnlineFriendsCmd           = "9"
	getFriendsCmd                 = "11"
	setUserStatusCmd              = "12"
	startMatchmakingCmd           = "13"
	cancelMatchmakingCmd          = "23"
	setReadyConsentMatchmakingCmd = "24"
	getLatenciesCmd               = "25"
	sendMessageToDSCmd            = "26"
	requestDSCmd                  = "33"
	requestFriendsCmd             = "14"
	listIncomingFriendsCmd        = "15"
	listOutgoingFriendsCmd        = "16"
	acceptFriendsCmd              = "17"
	rejectFriendsCmd              = "18"
	cancelFriendsRequestCmd       = "19"
	unfriendCmd                   = "20"
	listOfFriendsCmd              = "21"
	getFriendshipStatusCmd        = "22"
	personalChatHistoryCmd        = "30"
	joinDefaultChatChannelCmd     = "35"
	sendChatChannelCmd            = "36"
	rejectCmd                     = "37"
	blockCmd                      = "38"
	unblockCmd                    = "39"
	promoteLeaderCmd              = "27"
	generatePartyCodeCmd          = "28"
	getPartyCodeCmd               = "29"
	deletePartyCodeCmd            = "31"
	joinViaPartyCodeCmd           = "32"
	userMetricCmd                 = "40"
	setSessionAttributeCmd        = "41"
	getSessionAttributeCmd        = "42"
	getAllSessionAttributeCmd     = "43"
	sendPartyNotifCmd             = "34"
	refreshTokenCmd               = "44"
	getBlockedPlayersCmd          = "45"
	quitCmd                       = "99"
)

var (
	reader *bufio.Reader
	conn   *websocket.Conn
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
}

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	config := getConfig()

	reader = bufio.NewReader(os.Stdin)

	fmt.Println("Input lobby URL: [" + config.LobbyBaseURL + "]")
	lobbyURL := getInput()
	if lobbyURL == "" {
		lobbyURL = config.LobbyBaseURL
	}

	fmt.Println("Input IAM URL: [" + config.IAMBaseURL + "]")
	iamURL := getInput()
	if iamURL == "" {
		iamURL = config.IAMBaseURL
	}

	fmt.Println("Input Client ID: [" + config.IAMClientID + "]")
	clientID := getInput()
	if clientID == "" {
		clientID = config.IAMClientID
	}

	fmt.Println("Input Client Secret: [" + config.IAMClientSecret + "]")
	clientSecret := getInput()
	if clientSecret == "" {
		clientSecret = config.IAMClientSecret
	}

	fmt.Println("Input username/email:")
	username := getInput()

	fmt.Println("Input password:")
	password := getInput()

	fmt.Println("Input entitlement token:")
	entToken := getInput()

	req := gorequest.New()
	req.Debug = true

	var tokenResp TokenResponse
	fmt.Println(fmt.Sprintf("iam: %s", iamURL))
	resp, _, errs := req.Post(iamURL+"/oauth/token").
		SetBasicAuth(clientID, clientSecret).
		Type("form").
		Send(fmt.Sprintf(
			`{ "grant_type": "password", "username": "%s", "password": "%s" }`,
			username, password)).
		EndStruct(&tokenResp)

	if errs != nil {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(-1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		os.Exit(-1)
	}

	authHeader := "Bearer " + tokenResp.AccessToken

	fmt.Println(fmt.Sprintf("djay lobby: %s", lobbyURL))
	conn = connect(lobbyURL, authHeader, entToken)
	defer conn.Close()

	done := make(chan struct{})
	go heartbeat(done)
	go receive(done, true)

	serve(config)

	logrus.Debug("Done")
}

func connect(lobbyURL, authHeader, entitlement string) *websocket.Conn {
	logrus.Debug("Connecting user to lobby")

	req, err := http.NewRequest(http.MethodGet, lobbyURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Entitlement", entitlement)

	connection, res, err := websocket.DefaultDialer.Dial(req.URL.String(), req.Header)
	if err == websocket.ErrBadHandshake {
		b, e := ioutil.ReadAll(res.Body)
		if e == nil {
			logrus.Error("Bad handshake", res.Status, string(b))
		}
	}
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	connection.SetCloseHandler(func(code int, text string) error {
		logrus.Infof("handling close message, code: %d, message: %s\n", code, text)
		err := connection.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, text), time.Now().Add(time.Second))
		if err != nil {
			logrus.Error("error writing control message: ", err)
		}
		return nil
	})
	connection.SetPongHandler(func(text string) error {
		// logrus.Infof("handling pong message, message: %s\n", text)
		err := connection.SetReadDeadline(time.Now().Add(6 * time.Second))
		if err != nil {
			logrus.Error("error setting read deadline: ", err)
		}
		return nil
	})

	return connection
}

// nolint: gocyclo
func serve(config *config) {
	for {
		printHelp()
		command := getInput()
		switch command {
		case infoCmd:
			info()
		case createCmd:
			create()
		case leaveCmd:
			leave()
		case inviteCmd:
			invite()
		case joinCmd:
			join()
		case rejectCmd:
			reject()
		case kickCmd:
			kick()
		case chatCmd:
			chat()
		case blockCmd:
			block()
		case unblockCmd:
			unblock()
		case getBlockedPlayersCmd:
			getBlocked()
		case partyChatCmd:
			partyChat()
		case getOnlineFriendsCmd:
			getOnlineFriends()
		case getFriendsCmd:
			getFriendsStatus()
		case setUserStatusCmd:
			setUserStatus()
		case startMatchmakingCmd:
			startMatchmaking()
		case cancelMatchmakingCmd:
			cancelMatchmaking()
		case requestDSCmd:
			requestDS()
		case requestFriendsCmd:
			requestFriends()
		case listIncomingFriendsCmd:
			listIncomingFriends()
		case listOutgoingFriendsCmd:
			listOutgoingFriends()
		case acceptFriendsCmd:
			acceptFriends()
		case rejectFriendsCmd:
			rejectFriends()
		case cancelFriendsRequestCmd:
			cancelFriendsRequest()
		case unfriendCmd:
			unfriend()
		case listOfFriendsCmd:
			listOfFriends()
		case getFriendshipStatusCmd:
			getFriendshipStatus()
		case personalChatHistoryCmd:
			personalChatHistory()
		case joinDefaultChatChannelCmd:
			joinDefaultChatChannel()
		case sendChatChannelCmd:
			sendChannelChat()
		case setReadyConsentMatchmakingCmd:
			setReadyConsent()
		case getLatenciesCmd:
			getPingLatencies(config)
		case sendMessageToDSCmd:
			sendMessageToDS()
		case promoteLeaderCmd:
			promoteLeader()
		case generatePartyCodeCmd:
			generatePartyCode()
		case getPartyCodeCmd:
			getPartyCode()
		case deletePartyCodeCmd:
			deletePartyCode()
		case joinViaPartyCodeCmd:
			joinViaPartyCode()
		case userMetricCmd:
			userMetric()
		case setSessionAttributeCmd:
			setSessionAttribute()
		case getSessionAttributeCmd:
			getSessionAttribute()
		case getAllSessionAttributeCmd:
			getAllSessionAttribute()
		case sendPartyNotifCmd:
			sendPartyNotif()
		case refreshTokenCmd:
			refreshToken()
		case quitCmd:
			logrus.Print("disconnect message: ")
			msg := getInput()
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg))
			return
		}
	}
}

func receive(done chan struct{}, dump bool) {
	for {
		_, msg, subErr := conn.ReadMessage()
		if subErr != nil {
			logrus.Info("read message failed: ", subErr)
			close(done)
			return
		}
		if dump {
			logrus.Printf("\nReceived %d bytes, message: %s\n", len(msg), string(msg))
		}
	}
}

func create() {
	logrus.Debug("Sending create message")
	text := fmt.Sprintf("type: %s\n%s", TypeCreateRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func info() {
	logrus.Debug("Sending info message")
	text := fmt.Sprintf("type: %s\n%s", TypeInfoRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func leave() {
	logrus.Debug("Sending leave message")
	text := fmt.Sprintf("type: %s\n%s", TypeLeaveRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func invite() {
	fmt.Println("Friend ID:")
	id := getInput()
	logrus.Debug("Sending invite message")
	text := fmt.Sprintf("type: %s\n%s\nfriendID: %s", TypeInviteRequest, generateMessageID(), id)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func join() {
	fmt.Println("Party ID:")
	partyID := getInput()
	fmt.Println("Invitation token:")
	token := getInput()
	logrus.Debug("Sending join message")
	text := fmt.Sprintf("type: %s\n%s\npartyID: %s\ninvitationToken: %s", TypeJoinRequest, generateMessageID(), partyID, token)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func kick() {
	fmt.Println("Member ID:")
	id := getInput()
	logrus.Debug("Sending kick message")
	text := fmt.Sprintf("type: %s\n%s\nmemberID: %s", TypeKickRequest, generateMessageID(), id)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func reject() {
	fmt.Println("Party ID:")
	partyID := getInput()
	fmt.Println("Invitation token:")
	token := getInput()
	logrus.Debug("Sending reject message")
	text := fmt.Sprintf("type: %s\n%s\npartyID: %s\ninvitationToken: %s", TypeRejectRequest, generateMessageID(), partyID, token)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func block() {
	fmt.Println("Namespace:")
	namespace := getInput()
	fmt.Println("BlockedUserId:")
	blockedUserID := getInput()
	logrus.Debug("Sending block message")
	text := fmt.Sprintf("type: %s\n%s\nnamespace: %s\nblockedUserId: %s", TypeBlockPlayerRequest, generateMessageID(), namespace, blockedUserID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func unblock() {
	fmt.Println("Namespace:")
	namespace := getInput()
	fmt.Println("UnblockedUserId:")
	unblockedUserID := getInput()
	logrus.Debug("Sending unblock message")
	text := fmt.Sprintf("type: %s\n%s\nnamespace: %s\nunblockedUserId: %s", TypeUnblockPlayerRequest, generateMessageID(), namespace, unblockedUserID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getBlocked() {
	logrus.Debug("Sending get blocked message")
	text := fmt.Sprintf("type: %s\n%s", TypeGetBlockedPlayerRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func chat() {
	fmt.Println("Friend ID:")
	friendID := getInput()
	fmt.Println("Message:")
	content := getInput()
	id := generateID()
	logrus.Debug("Sending chat message")
	text := fmt.Sprintf("type: %s\nid: %s\nto: %s\npayload: %s", TypePersonalChatRequest, id, friendID, content)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func partyChat() {
	fmt.Println("Message:")
	content := getInput()
	id := generateID()
	logrus.Debug("Sending party chat message")
	text := fmt.Sprintf("type: %s\nid: %s\npayload: %s", TypePartyChatRequest, id, content)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getOnlineFriends() {
	logrus.Debug("get online friends")
	text := fmt.Sprintf("type: %s\n%s", TypeFriendsPresenceRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getFriendsStatus() {
	fmt.Println("get friends status:")
	text := fmt.Sprintf("type: %s\n%s", TypeFriendsPresenceRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func setUserStatus() {
	fmt.Println("Availability:")
	availability := getInput()
	fmt.Println("Activity:")
	activity := getInput()
	logrus.Debug("set users status")
	text := fmt.Sprintf("type: %s\n%s\navailability: %s\nactivity: %s", TypeSetUserStatusRequest, generateMessageID(), availability, activity)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func startMatchmaking() {
	fmt.Println("game mode:")
	gameMode := getInput()
	fmt.Println("DS name:")
	dsName := getInput()
	fmt.Println("Client version:")
	clientVersion := getInput()
	fmt.Println("Latencies: (hint: use get ping latencies command output here)")
	latencies := getInput()
	fmt.Println(`Party Attributes: format json, example: {"key":"value"}`)
	partyAttribute := getInput()
	fmt.Println(`Temp Party: (comma-separated userIDs, e.g. userA,userB)`)
	tempParty := getInput()
	fmt.Println(`Extra attributes: (comma-separated attributes, e.g. attrA,attrB)`)
	extraAttr := getInput()
	textCommand := fmt.Sprintf(
		"type: startMatchmakingRequest\n"+
			"%s\n"+
			"gameMode: %s\n"+
			"serverName: %s\n"+
			"clientVersion: %s\n"+
			"latencies: %s\n"+
			"partyAttributes: %s\n"+
			"tempParty: %s\n"+
			"extraAttributes: %s",
		generateMessageID(),
		gameMode,
		dsName,
		clientVersion,
		latencies,
		partyAttribute,
		tempParty,
		extraAttr,
	)
	err := conn.WriteMessage(websocket.TextMessage, []byte(textCommand))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func cancelMatchmaking() {
	fmt.Println("game mode:")
	gameMode := getInput()
	textCommand := fmt.Sprintf("type: cancelMatchmakingRequest\n%s\ngameMode: %s", generateMessageID(), gameMode)
	err := conn.WriteMessage(websocket.TextMessage, []byte(textCommand))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func requestDS() {
	fmt.Println("match ID:")
	matchID := getInput()
	fmt.Println("game mode:")
	gameMode := getInput()
	fmt.Println("DS name:")
	dsName := getInput()
	fmt.Println("client version:")
	clientVersion := getInput()
	fmt.Println("region:")
	region := getInput()
	fmt.Println("deployment:")
	deployment := getInput()
	textCommand := fmt.Sprintf(
		"type: createDSRequest\n"+
			"%s\n"+
			"matchId: %s\n"+
			"gameMode: %s\n"+
			"serverName: %s\n"+
			"clientVersion: %s\n"+
			"region: %s\n"+
			"deployment: %s\n",
		generateMessageID(),
		matchID,
		gameMode,
		dsName,
		clientVersion,
		region,
		deployment,
	)
	err := conn.WriteMessage(websocket.TextMessage, []byte(textCommand))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func requestFriends() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("friends request")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeRequestFriendsRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func listIncomingFriends() {
	logrus.Debug("get list of incoming friends")
	text := fmt.Sprintf("type: %s\n%s", TypeListIncomingFriendsRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func listOutgoingFriends() {
	logrus.Debug("get list of outgoing friends")
	text := fmt.Sprintf("type: %s\n%s", TypeListOutgoingFriendsRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func acceptFriends() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("accept friends")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeAcceptFriendsRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func rejectFriends() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("reject friends")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeRejectFriendsRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func cancelFriendsRequest() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("cancel friends")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeCancelFriendsRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func unfriend() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("unfriend")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeUnfriendRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func listOfFriends() {
	logrus.Debug("get list of friends")
	text := fmt.Sprintf("type: %s\n%s", TypeListOfFriendsRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getFriendshipStatus() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("get friendship status")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypeGetFriendshipStatusRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func personalChatHistory() {
	fmt.Println("Friends UserID:")
	friendID := getInput()
	logrus.Debug("load personal chat history")
	text := fmt.Sprintf("type: %s\n%s\nfriendId: %s", TypePersonalChatHistoryRequest, generateMessageID(), friendID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func joinDefaultChatChannel() {
	logrus.Debug("join default channel")
	text := fmt.Sprintf("type: %s\n%s", TypeJoinDefaultChannelRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func sendChannelChat() {
	fmt.Println("Channel Slug:")
	channelSlug := getInput()
	fmt.Println("Message:")
	payload := getInput()
	logrus.Debug("send channel chat")
	text := fmt.Sprintf("type: %s\n"+
		"%s\n"+
		"channelSlug: %s\n"+
		"payload: %s",
		TypeSendChannelChatRequest,
		generateMessageID(),
		channelSlug,
		payload)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func setReadyConsent() {
	fmt.Println("Match ID:")
	matchID := getInput()
	logrus.Debug("set ready consent")
	text := fmt.Sprintf("type: %s\n%s\nmatchId: %s", TypeSetReadyConsentRequest, generateMessageID(), matchID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func heartbeat(done chan struct{}) {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				logrus.Errorf("cannot write heartbeat: %v", err)
			}
		case <-done:
			logrus.Info("done signal received, stop heartbeat.")
			return
		}
	}
}

func printHelp() {
	fmt.Printf(`
Commands:
# Party
%s: User party info
%s: Create party
%s: Leave party
%s: Invite friend
%s: Join party
%s: Kick member
%s: Promote new leader
%s: Generate new party code
%s: Get party code
%s: Delete party code
%s: Join via party code
%s: Send party notif

# Chat
%s: 1to1 Chat
%s: Party Chat
%s: Load Personal Chat History
%s: Join Default Chat Channel
%s: Send Channel Chat

# Presence
%s: Get online friends
%s: Get Friends Status
%s: Set User Status

# Session Attribute
%s: Set user's session attribute
%s: Get user's session attribute
%s: Get user's all session attributes


# Matchmaking
%s: Start matchmaking
%s: Cancel matchmaking
%s: Set Ready Consent matchmaking
%s: Get ping latencies
%s: Send message to DS
%s: request DS

# Friends
%s: Request Friends
%s: List of Incoming Friends
%s: List of Outgoing Friends
%s: Accept Friends
%s: Reject Friends
%s: Cancel Friends Request
%s: Unfriend
%s: List of Friends
%s: Get Friendship Status

# Block
%s: Block Player
%s: Unblock Player
%s: Get blocked players

# Lobby
%s: Get Connected Player Count
%s: Refresh user token
%s: Quit

`,
		// party
		infoCmd,
		createCmd,
		leaveCmd,
		inviteCmd,
		joinCmd,
		kickCmd,
		promoteLeaderCmd,
		generatePartyCodeCmd,
		getPartyCodeCmd,
		deletePartyCodeCmd,
		joinViaPartyCodeCmd,
		sendPartyNotifCmd,

		// chat
		chatCmd,
		partyChatCmd,
		personalChatHistoryCmd,
		joinDefaultChatChannelCmd,
		sendChatChannelCmd,

		// presence
		getOnlineFriendsCmd,
		getFriendsCmd,
		setUserStatusCmd,

		// session attribute
		setSessionAttributeCmd,
		getSessionAttributeCmd,
		getAllSessionAttributeCmd,

		// matchmaking
		startMatchmakingCmd,
		cancelMatchmakingCmd,
		setReadyConsentMatchmakingCmd,
		getLatenciesCmd,
		sendMessageToDSCmd,
		requestDSCmd,

		// friends
		requestFriendsCmd,
		listIncomingFriendsCmd,
		listOutgoingFriendsCmd,
		acceptFriendsCmd,
		rejectFriendsCmd,
		cancelFriendsRequestCmd,
		unfriendCmd,
		listOfFriendsCmd,
		getFriendshipStatusCmd,

		// block
		blockCmd,
		unblockCmd,
		getBlockedPlayersCmd,

		// lobby
		userMetricCmd,
		refreshTokenCmd,
		quitCmd)
}

func getInput() string {
	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return text
}

func generateMessageID() string {
	return "id: " + generateID()
}

// nolint
// generateMockAccessToken is the function to give option for user to use mock access token to access lobby
func generateMockAccessToken() string {
	userID := generateID()
	logrus.Infof("creating mock user access token, user id: %v", userID)
	payload := fmt.Sprintf(`{"namespace":"mock","display_name":"LobbyTestUser","roles":["c0e5281daf05b2251438839e948d783e"],
						"permissions":[],"bans":[],"jflgs":0,"aud":"4a1ab8b880ca007cdd607f28ea98f3be","exp":1634406946,"iat":1534385419,
						"sub":"%s","jti":"19987012-4087-4ff8-8a8d-755a21b932f0"}`, userID)
	return fmt.Sprintf("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.%s.mock", base64.StdEncoding.EncodeToString([]byte(payload)))
}

func getPingLatencies(config *config) {
	type qoSServer struct {
		Region     string    `json:"region"`
		IP         string    `json:"ip"`
		Port       int       `json:"port"`
		LastUpdate time.Time `json:"last_update"`
	}

	type listQoSServerResponse struct {
		Servers []qoSServer `json:"servers"`
	}

	var serverList listQoSServerResponse

	fmt.Println("getting QoS server list...")

	resp, _, errs := gorequest.New().SetDebug(true).
		Get(config.QOSBaseURL + "/public/qos").
		EndStruct(&serverList)

	if errs != nil {
		logrus.WithField("error", errs).Error("unable to request list of QoS servers")
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("status", resp.StatusCode).Error("QoS Manager returned non-OK")
		return
	}

	fmt.Printf("QoS Servers: %+v\n\n", serverList)

	latencies := make(map[string]int)

	for _, server := range serverList.Servers {
		fmt.Print("pinging ", server, "...  ")
		now := time.Now().UTC()
		err := ping(server.IP, server.Port)
		if err != nil {
			logrus.WithField("error", err).Error("unable to ping server")
			continue
		}
		elapsed := time.Since(now)
		fmt.Println(elapsed)
		latencies[server.Region] = int(elapsed.Milliseconds())
	}

	b, err := json.Marshal(latencies)
	if err != nil {
		logrus.WithField("error", err).Error("unable to marshal latencies")
		return
	}

	fmt.Println("\nlatencies:")
	fmt.Println(string(b))
}

func ping(host string, port int) error {
	out, err := sendUDPMessage(host, fmt.Sprintf("%d", port), "PING")
	if err != nil {
		return err
	}
	if out != "PONG" {
		return errors.New("unexpected response: " + out)
	}
	return nil
}

func sendUDPMessage(host string, port string, msg string) (string, error) {
	c, err := net.Dial("udp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = c.Close()
	}()

	err = c.SetDeadline(time.Now().UTC().Add(5 * time.Second))
	if err != nil {
		return "", err
	}

	_, err = fmt.Fprintf(c, msg)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}

	s := string(buf[:n])
	return s, nil
}

func sendMessageToDS() {
	fmt.Println("IP:")
	ip := getInput()
	fmt.Println("Port:")
	port := getInput()
	fmt.Println("Message:")
	msg := getInput()
	ret, err := sendUDPMessage(ip, port, msg)
	if err != nil {
		logrus.WithField("error", err).Error("cannot send UDP message")
	}
	fmt.Println("response: ", ret)
}

func promoteLeader() {
	fmt.Println("New Leader User ID:")
	userID := getInput()
	text := fmt.Sprintf("type: %s\n%s\nnewLeaderUserId: %s", TypePromoteLeaderRequest, generateMessageID(), userID)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getPartyCode() {
	logrus.Debug("Sending get party code message")
	text := fmt.Sprintf("type: %s\n%s", TypeGetPartyCodeRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func generatePartyCode() {
	logrus.Debug("Sending generate party code message")
	text := fmt.Sprintf("type: %s\n%s", TypeGeneratePartyCodeRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func deletePartyCode() {
	logrus.Debug("Sending delete party code message")
	text := fmt.Sprintf("type: %s\n%s", TypeDeletePartyCodeRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func joinViaPartyCode() {
	fmt.Println("Party code:")
	partyCode := getInput()
	logrus.Debug("Sending join via code message")
	text := fmt.Sprintf("type: %s\n%s\npartyCode: %s", TypeJoinViaPartyCodeRequest, generateMessageID(), partyCode)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func setSessionAttribute() {
	fmt.Println("key:")
	key := getInput()
	fmt.Println("value:")
	value := getInput()
	logrus.Debug("Sending set session attribute")
	text := fmt.Sprintf("type: %s\n%s\nkey: %s\nvalue: %s", TypeSetSessionAttributeRequest, generateMessageID(), key, value)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getSessionAttribute() {
	fmt.Println("key:")
	key := getInput()
	logrus.Debug("Sending get session attribute")
	text := fmt.Sprintf("type: %s\n%s\nkey: %s", TypeGetSessionAttributeRequest, generateMessageID(), key)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func getAllSessionAttribute() {
	logrus.Debug("Sending get all session attribute")
	text := fmt.Sprintf("type: %s\n%s", TypeGetAllSessionAttributeRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func userMetric() {
	fmt.Println("Get User Metric:")
	text := fmt.Sprintf("type: %s\n%s", TypeUserMetricRequest, generateMessageID())
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func sendPartyNotif() {
	fmt.Println("Send party notif:")
	fmt.Println("topic:")
	topic := getInput()
	fmt.Println("payload:")
	payload := getInput()
	text := fmt.Sprintf("type: %s\n%s\ntopic: %s\npayload: %s", TypeSendPartyNotifRequest, generateMessageID(), topic, payload)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func refreshToken() {
	fmt.Println("Refresh token:")
	fmt.Println("token:")
	token := getInput()
	text := fmt.Sprintf("type: %s\n%s\ntoken: %s", TypeRefreshTokenRequest, generateMessageID(), token)
	err := conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		logrus.Fatalf("cannot write message: %v", err)
	}
}

func generateID() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}

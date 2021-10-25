package main

// Lobby messaging protocol type
const (
	// Party
	TypeInfoRequest                 = "partyInfoRequest"
	TypeInfoResponse                = "partyInfoResponse"
	TypeCreateRequest               = "partyCreateRequest"
	TypeCreateResponse              = "partyCreateResponse"
	TypeLeaveRequest                = "partyLeaveRequest"
	TypeLeaveResponse               = "partyLeaveResponse"
	TypeLeaveNotif                  = "partyLeaveNotif"
	TypeInviteRequest               = "partyInviteRequest"
	TypeInviteResponse              = "partyInviteResponse"
	TypeInviteNotif                 = "partyInviteNotif"
	TypeGetInvitatedNotif           = "partyGetInvitedNotif"
	TypeJoinRequest                 = "partyJoinRequest"
	TypeJoinResponse                = "partyJoinResponse"
	TypeJoinNotif                   = "partyJoinNotif"
	TypeRejectRequest               = "partyRejectRequest"
	TypeRejectResponse              = "partyRejectResponse"
	TypeRejectNotif                 = "partyRejectNotif"
	TypeKickRequest                 = "partyKickRequest"
	TypeKickResponse                = "partyKickResponse"
	TypeKickNotice                  = "partyKickNotif"
	TypePersonalChatRequest         = "personalChatRequest"
	TypePersonalChatResponse        = "personalChatResponse"
	TypePersonalChatNotif           = "personalChatNotif"
	TypePartyChatRequest            = "partyChatRequest"
	TypePartyChatResponse           = "partyChatResponse"
	TypePartyChatNotif              = "partyChatNotif"
	TypePersonalChatHistoryRequest  = "personalChatHistoryRequest"
	TypePersonalChatHistoryResponse = "personalChatHistoryResponse"
	TypePromoteLeaderRequest        = "partyPromoteLeaderRequest"
	TypePromoteLeaderResponse       = "partyPromoteLeaderResponse"
	TypeGeneratePartyCodeRequest    = "partyGenerateCodeRequest"
	TypeGeneratePartyCodeResponse   = "partyGenerateCodeResponse"
	TypeGetPartyCodeRequest         = "partyGetCodeRequest"
	TypeGetPartyCodeResponse        = "partyGetCodeResponse"
	TypeDeletePartyCodeRequest      = "partyDeleteCodeRequest"
	TypeDeletePartyCodeResponse     = "partyDeleteCodeResponse"
	TypeJoinViaPartyCodeRequest     = "partyJoinViaCodeRequest"
	TypeJoinViaPartyCodeResponse    = "partyJoinViaCodeResponse"
	TypeSendPartyNotifRequest       = "partySendNotifRequest"
	TypeSendPartyNotifResponse      = "partySendNotifResponse"
	TypePartyNotif                  = "partyNotif"

	// Presence
	TypeFriendsPresenceRequest  = "friendsStatusRequest"
	TypeFriendsPresenceResponse = "friendsStatusResponse"
	TypeSetUserStatusRequest    = "setUserStatusRequest"
	TypeSetUserStatusResponse   = "setUserStatusResponse"
	TypeUserStatusNotif         = "userStatusNotif"

	// TypeClientResetRequest is request from clienthandler to lobby
	TypeClientResetRequest = "clientResetRequest"

	// Notification
	TypeNotificationMessage = "messageNotif"

	// Matchmaking
	TypeStartMatchmakingRequest   = "startMatchmakingRequest"
	TypeStartMatchmakingResponse  = "startMatchmakingResponse"
	TypeCancelMatchmakingRequest  = "cancelMatchmakingRequest"
	TypeCancelMatchmakingResponse = "cancelMatchmakingResponse"
	TypeMatchmakingNotif          = "matchmakingNotif"
	TypeSetReadyConsentRequest    = "setReadyConsentRequest"
	TypeSetReadyConsentResponse   = "setReadyConsentResponse"
	TypeSetReadyConsentNotif      = "setReadyConsentNotif"
	TypeRematchmakingNotif        = "rematchmakingNotif"

	// Friends
	TypeRequestFriendsRequest       = "requestFriendsRequest"
	TypeRequestFriendsResponse      = "requestFriendsResponse"
	TypeRequestFriendsNotif         = "requestFriendsNotif"
	TypeListIncomingFriendsRequest  = "listIncomingFriendsRequest"
	TypeListIncomingFriendsResponse = "listIncomingFriendsResponse"
	TypeListOutgoingFriendsRequest  = "listOutgoingFriendsRequest"
	TypeListOutgoingFriendsResponse = "listOutgoingFriendsResponse"
	TypeAcceptFriendsRequest        = "acceptFriendsRequest"
	TypeAcceptFriendsResponse       = "acceptFriendsResponse"
	TypeAcceptFriendsNotif          = "acceptFriendsNotif"
	TypeRejectFriendsRequest        = "rejectFriendsRequest"
	TypeRejectFriendsNotif          = "rejectFriendsNotif"
	TypeRejectFriendsResponse       = "rejectFriendsResponse"
	TypeCancelFriendsRequest        = "cancelFriendsRequest"
	TypeCancelFriendsResponse       = "cancelFriendsResponse"
	TypeCancelFriendsNotif          = "cancelFriendsNotif"
	TypeUnfriendRequest             = "unfriendRequest"
	TypeUnfriendResponse            = "unfriendResponse"
	TypeUnfriendNotif               = "unfriendNotif"
	TypeListOfFriendsRequest        = "listOfFriendsRequest"
	TypeListOfFriendsResponse       = "listOfFriendsResponse"
	TypeGetFriendshipStatusRequest  = "getFriendshipStatusRequest"
	TypeGetFriendshipStatusResponse = "getFriendshipStatusResponse"
	TypeBlockPlayerRequest          = "blockPlayerRequest"
	TypeBlockPlayerResponse         = "blockPlayerResponse"
	TypeBlockPlayerNotif            = "blockPlayerNotif"
	TypeUnblockPlayerRequest        = "unblockPlayerRequest"
	TypeUnblockPlayerResponse       = "unblockPlayerResponse"
	TypeUnblockPlayerNotif          = "unblockPlayerNotif"
	TypeGetBlockedPlayerRequest     = "getBlockedPlayerRequest"
	TypeGetBlockedPlayerResponse    = "getBlockedPlayerResponse"

	// DSM
	TypeCreateDSRequest  = "createDSRequest"
	TypeCreateDSResponse = "createDSResponse"
	TypeWaitForDSRequest = "waitForDSRequest"
	TypeClaimDSRequest   = "claimDSRequest"
	TypeNotifyDSRequest  = "notifyDSRequest"
	TypeDSNotif          = "dsNotif"

	// System events
	TypeSystemComponentsStatus = "systemComponentsStatus"

	// Party
	TypePartyDataUpdateNotif = "partyDataUpdateNotif"

	// Channel Chat
	TypeJoinDefaultChannelRequest  = "joinDefaultChannelRequest"
	TypeJoinDefaultChannelResponse = "joinDefaultChannelResponse"
	TypeUserBannedNotification     = "userBannedNotification"
	TypeUserUnbannedNotification   = "userUnbannedNotification"
	TypeExitAllChannel             = "exitAllChannel"
	TypeSendChannelChatRequest     = "sendChannelChatRequest"
	TypeSendChannelChatResponse    = "sendChannelChatResponse"
	TypeChannelChatNotif           = "channelChatNotif"

	// session attribute
	TypeSetSessionAttributeRequest     = "setSessionAttributeRequest"
	TypeSetSessionAttributeResponse    = "setSessionAttributeResponse"
	TypeGetSessionAttributeRequest     = "getSessionAttributeRequest"
	TypeGetSessionAttributeResponse    = "getSessionAttributeResponse"
	TypeGetAllSessionAttributeRequest  = "getAllSessionAttributeRequest"
	TypeGetAllSessionAttributeResponse = "getAllSessionAttributeResponse"

	// Signaling
	TypeSignalingP2PNotif = "signalingP2PNotif"

	// User Metric
	TypeUserMetricRequest  = "userMetricRequest"
	TypeUserMetricResponse = "userMetricResponse"

	// System Metric
	TypeSendMatchResultNotif = "sendMatchResultNotif"
)

// Message types enum
const (
	TypeUnknown = "unknown"

	// Server messages
	TypeConnected            = "connectNotif"
	TypeDisconnected         = "disconnectNotif"
	TypeError                = "errorNotif"
	TypeShutdown             = "shutdownNotif"
	TypeHeartbeat            = "heartbeat"
	TypeRefreshTokenRequest  = "refreshTokenRequest"
	TypeRefreshTokenResponse = "refreshTokenResponse"

	// SuccessCode
	SuccessCode = 0

	// Matchmaking status
	MMStatusStart   = "start"
	MMStatusDone    = "done"
	MMStatusCancel  = "cancel"
	MMStatusTimeout = "timeout"
	MMStatusBanned  = "banned"

	// System Components
	SystemComponentChat = "chat"
)
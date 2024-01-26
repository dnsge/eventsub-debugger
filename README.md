# eventsub-debugger

A simple Twitch EventSub Websocket client for testing and debugging a single subscription.

## Usage 
```
Usage of eventsub-debugger:
  -client-id string
    	Twitch Application ClientID
  -condition string
    	JSON Condition to pass
  -server string
    	Twitch EventSub WebSocket endpoint (default "wss://eventsub.wss.twitch.tv/ws")
  -sub-type string
    	EventSub Subscription Type
  -sub-version string
    	EventSub Subscription Version
  -token string
    	Twitch OAuth Access Token
```

## Example (channel.chat.message)

First, generate an access token to read the chat: `twitch token -u -s "user:read:chat"`. Then, specify the chat to connect to and your user id in the condition:

```
$ eventsub-debugger \
      --client-id "<your client id>" \
      --token "<your access token>" \
      --sub-type "channel.chat.message" \
      --sub-version "1" \
      --condition '{ "broadcaster_user_id": "22484632", "user_id": "<your user id>" }'
```

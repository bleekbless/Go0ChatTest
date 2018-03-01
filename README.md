## Introduction
If you don't have few hundreds users for load-testing your alpha-version of Facebook chat-bot application, the application might be helpful for you. Especially for kind of applications like quiz, poll, etc.

## How it works ?
The application can imitate Facebook in very simple way.

Your chat-bot app sends requests to the app.

And if you need, the app sends requests in response.

So you need:
1. write `config.json` like this [config.sample.json](https://github.com/ss-dev/fb-chat-emulator/config.sample.json)
2. run binary `$ fb-chat-emulator` for your config
3. change Facebook URL in your app to the app URL, for example:
```python
# URL_BASE = “https://graph.facebook.com/v2.6/me/”
URL_BASE = “http://127.0.0.1:8877/”
```
4. create test users into your app
5. and test it

## Configuration
In a config file you can set:

`AppURL` - base URL to your chat-bot app.

Facebook spends some time to work our requests so you can set the time between `RequestTimeMin` and `RequestTimeMax` ms.

Users spend some time to respond - from `ResponsePauseMin` to `ResponsePauseMax` ms.

By using `Rules` you can describe some reaction.
For example, chat-bot app sends a question for a user (say during some quiz):

```
POST https://graph.facebook.com/v2.6/me/messages
{"message": {"attachment": {"type": "template", ... "title": "What is AppURL parameter above?"}}, "recipient": {"id": "***************"}}
```

so we can describe it as

```javascript
"Request": {
    "Name": "question#1",
    "URL": "/messages",
    "BodySegment": "What is AppURL parameter above?"
}
```

there `Name` need only for readable statistic.

And if right answer is `It is address of your app server` our `Response` would be like this:

```javascript
"Response": {
    "URL": "/fb-bot/",
    "Body": "{'entry': [{'messaging': [{'timestamp': [timestamp], 'postback': {'payload': 'question:10', 'title': 'It is address of your app server'}, 'recipient': {'id': 1234567890}, 'sender': {'id': [RecipientId]}}]}]}"
}
```

and it's a full rule:

```javascript
{
    "Request": {
        "Name": "question#1",
        "URL": "/messages",
        "BodySegment": "What is AppURL parameter above?"
    },
    "Response": {
        "URL": "/fb-bot/",
        "Body": "{'entry': [{'messaging': [{'timestamp': [timestamp], 'postback': {'payload': 'question:10', 'title': 'It is address of your app server'}, 'recipient': {'id': 1234567890}, 'sender': {'id': [RecipientId]}}]}]}"
    }
}
```

You can write empty `Response` then the app won't send any response but will write statistic for the requests.

## Templates
As you could see in examples above you can use some template-variables in config.

`[timestamp]` - it's replaced by current timestamp for response body

`[RecipientId]` - it gets from request body

## Statistics
You can watch realtime statistic during testing process in your browser `http://127.0.0.1:8877/stat`.

And reset it between tests `http://127.0.0.1:8877/reset`.

It looks like this:
```
+------------+-----------+------------+--------------+-----------------+---------------+--------------+----------------+---------------+
|    NAME    | REQ COUNT | RESP COUNT | RESP NET ERR | RESP STATUS ERR | FIRST REQUEST | LAST REQUEST | FIRST RESPONSE | LAST RESPONSE |
+------------+-----------+------------+--------------+-----------------+---------------+--------------+----------------+---------------+
| greetings  |       100 |          0 |            0 |               0 | 12:42:00.513  | 12:42:01.977 | 00:00:00.000   | 00:00:00.000  |
| question#1 |       100 |        100 |            0 |               0 | 12:42:26.805  | 12:42:28.285 | 12:42:28.261   | 12:42:32.915  |
| question#2 |       100 |        100 |            0 |               0 | 12:42:48.154  | 12:42:49.408 | 12:42:49.501   | 12:42:54.147  |
| question#3 |       100 |        100 |            0 |               0 | 12:43:11.119  | 12:43:12.327 | 12:43:12.592   | 12:43:16.970  |
| OTHERS     |      2525 |          0 |            0 |               0 | 12:42:09.003  | 12:43:49.505 | 00:00:00.000   | 00:00:00.000  |
+------------+-----------+------------+--------------+-----------------+---------------+--------------+----------------+---------------+
```

By the way, thanks [olekukonko](https://github.com/olekukonko) for the great library [tablewriter](https://github.com/olekukonko/tablewriter) !

## Flags
If you want to change app port or path to config or say to enable debug mode you should use `flags`.

By the way, Debug mode can help you to find out what exactly do you send to Facebook.

Use `$ fb-chat-emulator --help` for more information.

## Releases

You can build your own binary by using `go tools` or use prepared builds there [https://github.com/ss-dev/fb-chat-emulator/releases](https://github.com/ss-dev/fb-chat-emulator/releases)

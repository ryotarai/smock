# Smock

Slack mock

## Installation

```
$ go get -u github.com/ryotarai/smock
```

## Usage with Bolt

First, make Slack API URL configurable:

```js
let config = {
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
};

if (process.env.SLACK_API_URL !== undefined) {
  config.clientOptions = {
    slackApiUrl: process.env.SLACK_API_URL,
  };
}

const app = new App(config);
```

Start Smock and keep it running:

```
$ smock start --event-url=http://localhost:3000/slack/events --listen=localhost:3100 --external-url=http://localhost:3100
>>> 
```

Start your Bolt app in another terminal:

```
$ SLACK_BOT_TOKEN=dummy SLACK_SIGNING_SECRET=dummy SLACK_API_URL=http://localhost:3100/api node app.js
```

Type your commands or messages to Smock prompt:

```
>>> /your-slash-command
<<< from your app
>>> your message
<<< from your app
```

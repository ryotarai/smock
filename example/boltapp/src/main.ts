import * as bolt from "@slack/bolt";
import * as webapi from "@slack/web-api";
import { ReceiverMultipleAckError } from "@slack/bolt";

const app = new bolt.App({
  token: process.env.SLACK_BOT_TOKEN,
  signingSecret: process.env.SLACK_SIGNING_SECRET,
  logLevel: bolt.LogLevel.DEBUG,
  clientOptions: {
    slackApiUrl: 'http://localhost:3002/api',
  },
});

app.command("/echo", async ({ack, command, respond}: bolt.SlackCommandMiddlewareArgs) => {
  console.log(command);
  await ack({
    text: `ack ${command.text}`,
  });
  await respond({
    text: "Hello respond",
  });
});

app.message(':wave:', async ({ message, say }) => {
  console.log(message);
  await say(`Hello, <@${message.user}>`);
});

(async () => {
  // Start your app
  await app.start(process.env.PORT || 3000);

  console.log('⚡️ Bolt app is running!');
})();

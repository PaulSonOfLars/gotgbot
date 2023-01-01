# Sample bots

This is a short description of all the sample bots present in this directory. These are intended to be a source of
inspiration and learning for anyone looking to use this library.


## samples/callbackqueryBot

This bot demonstrates some example interactions with telegram callback queries.
It has a basic start command which contains a button. When pressed, this button is edited via a callback query.

## samples/commandBot

This bot demonstrates some example interactions with commands on telegram.
It has a basic start command with a bot intro.
It also has a source command, which sends the bot sourcecode, as a file.

## samples/conversationBot

This bot demonstrates a basic conversation handler. Conversation handlers help you to track states across multiple
messages.
In this case, the bot has commands to start a conversation, which then causes it to ask for your name and your age.

## samples/echoBot

This bot is as basic as it gets - it simply repeats everything you say.

## samples/echoMultiBot

This bot demonstrates how to create echo bot that works with multiple bot instances at once.
It also shows how to stop the bot gracefully using the Updater.Stop() mechanism.
It has options to use either polling or webhooks.

## samples/echoWebhookBot

This bot repeats everything you say - but it uses webhooks instead of long polling.
Webhooks are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
Simply install ngrok, make an account on the website, and run:
`ngrok http 8080`
Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
from the samples/echoWebhookBot directory:
`TOKEN="<your_token_here>" WEBHOOK_DOMAIN="<your_domain_here>"  WEBHOOK_SECRET="<random_string_here>" go run .`
Then, simply send /start to your bot; if it replies, you've successfully set up webhooks!

## samples/middlewareBot

This bot shows how to effectively use middlewares to modify and intercept HTTP requests to the bot API server.
In this example, the middleware sets the allow_sending_without_reply to certain methods, as well as make sure to log all error messages.

## samples/webappBot

This bot shows how to use this library to server a webapp.
Webapps are slightly more complex to run, since they require a running webserver, as well as an HTTPS domain.
For development purposes, we recommend running this with a tool such as ngrok (https://ngrok.com/).
Simply install ngrok, make an account on the website, and run:
`ngrok http 8080`
Then, copy-paste the HTTPS URL obtained from ngrok (changes every time you run it), and run the following command
from the samples/webappBot directory:
`URL="<your_url_here>" TOKEN="<your_token_here>" go run .`
Then, simply send /start to your bot, and enjoy your webapp demo.

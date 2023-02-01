# technews-bot

Aggregate tech news articles from multiple news sources, filter them based on
your favorite subjects and send them to a chat platform of your choice.

## Invite link

[Technews Bot invite link](https://discord.com/oauth2/authorize?client_id=1020502388462854165&scope=bot&permissions=8)

## Supported news sources

-   [https://lobste.rs/](https://lobste.rs/)
-   [https://news.ycombinator.com/news](https://news.ycombinator.com/news)

## Configuration

To run this bot, you must provide at least a discord token and the channel ID
where to dump all of the interresting articles via environment variables. You
can also use a `.env` file to set the environment variables. The env file looks like this:

```sh
DISCORD_TOKEN=<Your token here>
DISCORD_CHANNEL=<Your text channel id here>
```

## Project vision

-   This project is not configurable as of right now, aside from the discord token
    and channel id, but will aim to be fully configurable via text commands in the
    future.
-   I'm not sure if this bot should only support Discord for now. It would be
    interresting to add support for many chat services such as Slack or Matrix.

# Contribute

This project aims to aggregate news from _tons_ of different news sources
into one place. If you would like to have a news source added, feel free to
either contribute support for it or [open an
issue](https://github.com/notarock/technews-bot/issues/new) explaining why
the news source should be added.

# Telegram News Bot

Bot for Telegram that gets and posts news to a channel.

# Features

- Fetching articles from RSS feeds
- Article summaries powered by GPT-3.5 or llama3
- Admin commands for managing sources

# Configuration

## Environment variables

- `TG_BOT_TOKEN` — token for Telegram Bot API
- `TG_CHANNEL_ID` — ID of the channel to post to, can be obtained via [@JsonDumpBot](https://t.me/JsonDumpBot)
- `DB_DSN` — PostgreSQL connection string
- `FETCH_INTERVAL` — the interval of checking for new articles, default `10m`
- `NOTIFICATION_INTERVAL` — the interval of delivering new articles to Telegram channel, default `1m`
- `FILTER_KEYWORDS` — comma separated list of words to skip articles containing these words
- `OPENAI_KEY` — token for OpenAI API
- `OPENAI_PROMPT` — prompt for GPT-3.5 Turbo to generate summary

## HCL

News Feed Bot can be configured with HCL config file. The service is looking for config file in following locations:

- `./config.hcl`
- `./config.local.hcl`
- `$HOME/.config/telegram-news-bot/config.hcl`

The names of parameters are the same except that there is no prefix and names are in lower case instead of upper case.

# Nice to have features (backlog)

- [ ] More types of resources — not only RSS
- [x] Summary for the article
- [ ] Dynamic source priority (based on 👍 and 👎 reactions) — currently blocked by Telegram Bot API
- [ ] Article types: text, video, audio
- [ ] De-duplication — filter articles with the same title and author
- [ ] Low quality articles filter — need research
    - Ban by author?
    - Check article length — not working with audio/video posts, but it will be fixed after article type implementation

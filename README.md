# Pinger bot

The purpose: add **mention all** functionality to telegram groups. This is native feature
in slack, rocket-chat and many other apps. But not in telegram.

> By the way, we will write a simple bot framework

## Usage (WIP)

First and main usage: add **@all** to your message and bot will mention every person in
the chat.

## Future ideas

Add categories. For example: **@all/backend**
Add users ability to signup for desired categories (topics).

## TODO

- [x] Basic polling
- [ ] Dockerfile & CI/CD
- [ ] Logging
- [ ] Basic target logic (add to chat, get members list, react on **@all** keyword)
- [ ] Administrative API (get chats and some other stats and metrics)
- [ ] Process each chat concurrently
- [ ] Nice docs page
- [ ] Webhooks
- [ ] Topics


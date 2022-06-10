# Mr.Pinger - notify them all!
---
## Introduction
#### Initial problem

My team of developers uses telegram to communicate on daily basis. We really like it's simplicity and speed.
The only feature that we miss is mentioning all members of a group (something like Slack's **@all** or **@channel**). It's a pitty that telegram doesn't have native support for this.

_Bots to the resque!_

#### One more bot?
It would be naive to think that i'm a pioneer in solving this problem. There are some solutions out there, and we have tried some of them. But there were always some issues (downtime, just stopping forever, etc.).

> Becides, i'm learning Go at the moment. And I never had a real pet-project which someone else could use. So.. why not?

#### Development and final goal

I'm going to avoid using external libraries in the beginning (in educational purposes). But than I will add some packages that can really improve the code.

One more thing to mention is that I'm going stick with **Bot API** (not more complex Telegram API).

The development process will be iterative and what I'm going to achieve is something like this:
- working application with functionality described below
- automatic deployment to cloud (SaaS model)
- good **go** code using all best practices (idiomatic solutions, linting, docs, tests etc.)

In the end the plan is to write a post about the journey and bot itself on habr or medium.

#### Target functionality

- [x] Ping all members of the group with **/ping** command.
- [ ] Add members to categories (eg.: #backend, #frontend or any others)
- [ ] Ping users by categories with **/ping#**_{tag}_ command

Bot should definitely work in **privacy mode** so any team could use it safely.

> Maybe some other functions will be suggested by my teammates or other users. Any contribution is welcomed!
___
## Current state (v0.1.0)

#### Privacy
BEWARE! At the moment privacy mode **is disabled**. It will be fixed in the next minor version.

#### Limitations

Standard Telegram Bot API doesn't allow bots to list all users in group. At least I couldn't find this option. That's why users should be registered before they can receive notifications. There are 2 ways:
- user is added to "known list" when he joins the group (and removed when leaves)
- user is added when he writes a message

Second option is **unsecure** (needs privacy mode to be disabled) and **will be removed** in next version. Instead new command will be implemented: **/add** and **/addme**

#### Usage

Add **@eonae_pinger_bot** to you telegram group.

Then type **/ping <your message>**. It will trigger bot's reply to your message, mentioning all known members of the group.
___
## TODO

- [x] Basic long polling
- [x] Logging
- [x] Sending messages
- [x] Groups registry, known members registry
- [x] Basic target logic (add to chat, get members list, react on **/ping** keyword)
- [ ] Implement /add and /addme
- [ ] Enable privacy mode
- [ ] Implement tags/topics (#backend,#frontend)
- [ ] CI/CD (not sure how exactly yet)
- [ ] Monitoring with grafana and maybe prometheus
- [ ] Process each chat concurrently
- [ ] Unit tests
- [ ] Linters (golangci-lint)
- [ ] Nice docs page
- [ ] Webhooks

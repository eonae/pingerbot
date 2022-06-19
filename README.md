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
- [x] Add members to categories (eg.: #backend, #frontend or any others)
- [x] Ping users by categories with **/ping#**_{tag}_ command

Bot should definitely work in **privacy mode** so any team could use it safely.

> Maybe some other functions will be suggested by my teammates or other users. Any contribution is welcome!
___
## Current state (v0.2.0)

#### Privacy
Unlike in v0.1.0, now privacy mode **is enabled**.

#### Usage

Add **@eonae_pinger_bot** to you telegram group.

Now you can send him some commands.

##### Registration

Standard _Telegram Bot API_ doesn't allow bots to list all users in group. At least I couldn't find this option. That's why users should be registered before they can receive notifications.

Everyone can add himself this way: 
> **/addme**

To add someone else use **/add** For example:
> **/add** @jessy, @jimmy ...

NOTE: anyone can add users, not only admin

##### Tags

Since 0.2.0 you user can be added with tag. Eg:

> **/addme** _#backend_ (you can specify more tags)

or

> **/add** @jimmy @jessy #pm #important

NOTE: Adding without tag means adding with #all (created automatically)

##### List

You can get a list of users grouped by tags:

> **/ls** [tags]

If tags are provided, only users added to these tags will be shown. Output will be something like this:
```
#all
-> jessy
-> jimmy
-> vince
-> kate
#backend
-> kate
-> vince
#pm
-> jessy
-> jimmy
```

##### Ping

If you type **/ping** in a message, it will trigger bot to reply to it, mentioning all members he know in the group. Like this:

> **/ping** Hey, guys! Time for daily!
```
> [replyto] Hey guys! Time for daily!
@jimmy, @jessy, @kate, @vince
```

You can specify tags to mention not all of the users:

> **/ping** _#backend_ Backenders, come on!

##### Unregister

Commands **/remove** and **/removeme** are just opposite to /add and /addme and can be used with tags too:

> **/removeme** _#backend_
> **/remove** _@vince_ _#backend_

___
## TODO

- [x] Basic long polling
- [x] Logging
- [x] Sending messages
- [x] Groups registry, known members registry
- [x] Basic target logic (add to chat, get members list, react on **/ping** keyword)
- [x] Implement /add and /addme
- [x] Enable privacy mode
- [x] Implement tags/topics (#backend,#frontend)
- [ ] CI/CD (not sure how exactly yet)
- [ ] Monitoring with grafana and maybe prometheus
- [ ] Process each chat concurrently
- [ ] Unit tests
- [ ] Linters (golangci-lint)
- [ ] Nice docs page
- [ ] Webhooks

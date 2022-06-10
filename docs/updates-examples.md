
### Examples
Two updates happen when bot is added to group. First is of my_chat_member.
``` json
{
    "update_id": 31227837,
    "my_chat_member": {
        "chat": {
            "id": -758526669,
            "title": "testing-bot",
            "type": "group",
            "all_members_are_administrators": true
        },
        "from": {
            "id": 388218616,
            "is_bot": false,
            "first_name": "Sergey",
            "username": "eonae_white",
            "language_code": "en"
        },
        "date": 1654294588,
        "old_chat_member": {
            "user": {
                "id": 2030512582,
                "is_bot": true,
                "first_name": "Pinger",
                "username": "eonae_pinger_bot"
            },
            "status": "left"
        },
        "new_chat_member": {
            "user": {
                "id": 2030512582,
                "is_bot": true,
                "first_name": "Pinger",
                "username": "eonae_pinger_bot"
            },
            "status": "member"
        }
    }
}
```
Then a message is emitted that bot was added to chat.
``` json
{
    "update_id": 31227838,
    "message": {
        "message_id": 29,
        "from": {
            "id": 388218616,
            "is_bot": false,
            "first_name": "Sergey",
            "username": "eonae_white",
            "language_code": "en"
        },
        "chat": {
            "id": -758526669,
            "title": "testing-bot",
            "type": "group",
            "all_members_are_administrators": true
        },
        "date": 1654294588,
        "new_chat_participant": {
            "id": 2030512582,
            "is_bot": true,
            "first_name": "Pinger",
            "username": "eonae_pinger_bot"
        },
        "new_chat_member": {
            "id": 2030512582,
            "is_bot": true,
            "first_name": "Pinger",
            "username": "eonae_pinger_bot"
        },
        "new_chat_members": [
            {
                "id": 2030512582,
                "is_bot": true,
                "first_name": "Pinger",
                "username": "eonae_pinger_bot"
            }
        ]
    }
}

```
- [ ] Remove bot from group
``` json
{
    "update_id": 31227839,
    "my_chat_member": {
        "chat": {
            "id": -758526669,
            "title": "testing-bot",
            "type": "group",
            "all_members_are_administrators": true
        },
        "from": {
            "id": 388218616,
            "is_bot": false,
            "first_name": "Sergey",
            "username": "eonae_white",
            "language_code": "en"
        },
        "date": 1654294598,
        "old_chat_member": {
            "user": {
                "id": 2030512582,
                "is_bot": true,
                "first_name": "Pinger",
                "username": "eonae_pinger_bot"
            },
            "status": "member"
        },
        "new_chat_member": {
            "user": {
                "id": 2030512582,
                "is_bot": true,
                "first_name": "Pinger",
                "username": "eonae_pinger_bot"
            },
            "status": "left"
        }
    }
}
```
- [ ] Member added to group
``` json
{
    "update_id": 31227843,
    "message": {
        "message_id": 32,
        "from": {
            "id": 388218616,
            "is_bot": false,
            "first_name": "Sergey",
            "username": "eonae_white",
            "language_code": "en"
        },
        "chat": {
            "id": -758526669,
            "title": "testing-bot",
            "type": "group",
            "all_members_are_administrators": true
        },
        "date": 1654295229,
        "new_chat_participant": {
            "id": 58894736,
            "is_bot": false,
            "first_name": "Vladimir",
            "last_name": "Korennoy",
            "username": "icewind666"
        },
        "new_chat_member": {
            "id": 58894736,
            "is_bot": false,
            "first_name": "Vladimir",
            "last_name": "Korennoy",
            "username": "icewind666"
        },
        "new_chat_members": [
            {
                "id": 58894736,
                "is_bot": false,
                "first_name": "Vladimir",
                "last_name": "Korennoy",
                "username": "icewind666"
            }
        ]
    }
}
```
- [ ] Member removed from group
``` json
{
    "update_id": 31227834,
    "message": {
        "message_id": 27,
        "from": {
            "id": 388218616,
            "is_bot": false,
            "first_name": "Sergey",
            "username": "eonae_white",
            "language_code": "en"
        },
        "chat": {
            "id": -758526669,
            "title": "testing-bot",
            "type": "group",
            "all_members_are_administrators": true
        },
        "date": 1654294163,
        "left_chat_participant": {
            "id": 58894736,
            "is_bot": false,
            "first_name": "Vladimir",
            "last_name": "Korennoy",
            "username": "icewind666"
        },
        "left_chat_member": {
            "id": 58894736,
            "is_bot": false,
            "first_name": "Vladimir",
            "last_name": "Korennoy",
            "username": "icewind666"
        }
    }
}
```
- [ ] Message with **/ping** command published into group chat.
``` json
{
    "message_id": 37,
    "from": {
        "id": 388218616,
        "is_bot": false,
        "first_name": "Sergey",
        "username": "eonae_white",
        "language_code": "en"
    },
    "chat": {
        "id": -758526669,
        "title": "testing-bot",
        "type": "group",
        "all_members_are_administrators": true
    },
    "date": 1654296261,
    "text": "/ping Some message",
    "entities": [
        {
            "offset": 0,
            "length": 5,
            "type": "bot_command"
        }
    ]
}

```

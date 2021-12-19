# STOMP Chat

Chat client to operate through a STOMP pub/sub server.

Operation:
- Message frames look like this:
    ```
    MESSAGE
    destination:/channel/main
    username:imamod

    this is a chat message on the channel main
    ```
That's it, it's that simple.

Prompt:
```
username : #channel>
```

Switch channel:
```
username : #channel>/chan other-channel
```

This creates that channel if it does not exist.

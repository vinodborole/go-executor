## Worker

It is responsible to consume the messages from the redis queue and perform respective tasks.

It also reports tasks status on another redis queue which manager consumes
## Worker

It is responsible to consume the messages from the redis queue and perform respective tasks.

It also reports tasks status on another redis queue which manager consumes and also can receive messages related to stopping job execution in between from the manager
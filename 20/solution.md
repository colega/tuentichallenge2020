This solution is incomplete.

Create exploit runing `python3 exploit.py` and then send it to the app.

Although it works against my local instance:

    root@b9563459fc59:/data# netcat localhost 6666 < exploit
    Current secret handshake: 29.
    > Welcome agent.
    SECRET:This_flag_is_only_for_prepro
    Agent options:
    1: Read company motto.
    2: Send task status.
    3: Request emergency evac.
    0: Exit.

It doesn't work for the remote one:

    root@b9563459fc59:/data# netcat 52.49.91.111 32666 < exploit
    Current secret handshake: 29.
    > Welcome agent.

It was fun anyway.

Budget chat

- tcp based chat room so protocol is tcp
- msg will be single line of ascii text by `\n`
- clients can send multiple msgs per connection, doesnt that mean websocket ?
- optionally strip trailing whitespace and carriage return chars like `\r`

Upon connection 
- ask user for username `Welcome to budgetchat! What shall I call you?`
    - first msg is username
    - username must contain 1 < username <= 16 chars
    - consist entirely of alphanumeric characters (uppercase, lowercase, and digits)
    - may choose to either allow or reject duplicate names
- optionally error out on illegal name and DISCONNECT the session

User joins
- announce to all other players when a user joins
    - * bob has entered the room

- announce to new user, the list of all ACTIVE users [ even when room is empty ]
    - * The room contains: bob, charlie, dave 

- broadcast a users msg to all
    open square bracket character
    the sender's name
    close square bracket character
    space character
    the sender's message

    If "bob" sends "hello", other users would receive "[bob] hello".

- limit on message len <= 1000

User leaves (only counts if it joined after connecting)
* bob has entered the room

Make sure you support at least 10 simultaneous clients.
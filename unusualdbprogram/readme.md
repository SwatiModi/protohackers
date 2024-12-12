- db is accessed over UDP
- no retransmission of dropped packets
-  two types of request: insert and retrieve. 
    - insert allows a client to insert a value for a key 
    - retrieve allows a client to retrieve the value for a key.

- insert request will have an "=" in request packet
    if key already exists, value must be updated
- rectreive will not have an "=" in request packet
- An insert request does not yield a response
- response should be served by ip and port which received it and should respond to ip and port which sent the request
- for key that doesnt exist, either "key=" or no response

- version reporting
    - special key "version"
    - attempts to update this key must be ignored

- all reqs and resps should be smaller than 1000 bytes

- ignore all udp related issues, it must assume it works fine


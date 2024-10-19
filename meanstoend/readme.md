
Means to an end

9 bytes long msg
 0 - > char
 1 -> 5 int32
 6 -> 9 int32


 I -> Insert 
 Q -> Query

 two signed two's complement 32-bit integers in network byte order

each client has its own data, basically ip+port has it own data, clients can not query for data they did not insert, conn.RemoteAddr() gives us this info
LocalAddr(): Info about your own server (your IP and port).
RemoteAddr(): Info about the client (their IP and port).

49    00 00 30 39    00 00 00 65
hex ^^

convert to decimal
73 | 0 0 48 57 | 0 0 0 101

first char to ascii
next 4 + 4 to big endian

I | 48 x 256 + 57 | 101 x 1

I | 12345 | 101


----

hex to decimal conversion

hex has base 16
decimal has base 10

0 -> 0 
1 -> 1
2 -> 2
3 -> 3
4 -> 4
5 -> 5
6 -> 6 
7 -> 7
8 -> 8
9 -> 9
A -> 10
B -> 
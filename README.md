# pgsqlxx

This is simply an attempt to have Pgx and Sqlx interop together directly, rather then through pgx/stdlib. Unfortunately some of sqlx needs to be reimplemented in order to achieve this, but ideally this will be kept to a minimum.

## Why?

sqlx and pgx are both great libraries, and having some of the sqlx functionality work directly with pgx makes using pgx specific funtionality possible, as well as leading to improved performance (hopefully?)

## Reimplementation of sqlx

Whenever I reimplement (read: pretty much copy paste) sqlx functionality in this library, 
I make a point of commenting the exact line and commit the code block is coming from. This is mostly caused by certain 
struct attributes / functions being private to the sqlx package. If they are ever made public, or have some way to 
publicly interop with them, these reimplementations wouldn't be necessary.

## Status

VERY much a work in progress. 

## Thanks

BIG THANKS to the authors of both pgx and sqlx. This would not be possible without them.
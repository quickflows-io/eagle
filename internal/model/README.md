# Model

The Model layer, or `Entity`, the entity layer, is used to store our entity classes, which are basically consistent with the attribute values in the database.

The structures returned by http access are also placed here, and the structures are converted before output. Generally in the form of `XXXInfo`.
For example: convert `userModel` to struct `UserInfo` before returning the end user.

## database conventions

Here is used by default `MySQL` database，use as much as possible `InnoDB` as storage engine。

### Related tables use a uniform prefix

such as user-related，use `user_` as table prefix：

```bash
user_base       // User base table
user_follow     // User attention form
user_fans       // User fan list
user_stat       // User Statistics
```

### Uniform field name

Three fields that need to be included in a table: primary key (id), creation time (created_at), update time (updated_at)
If a user id is required, it is generally represented by `user_id`.
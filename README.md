# silo-server
Server side implementation for Silo. Realtime modular applications on Android: auto-updating, analytics, and more. I made the server and client side integration in 6 days, so it's currently beta work.

For Android client side implementation, see [here] (https://github.com/farazfazli/silo/tree/master).

# Core Technology

It's built in Go, on top of TCP, and data is stored in Redis. Files are pretty self-explanatory.

# Adding Modules

It's easy to add modules. See [Blaise module] (https://github.com/farazfazli/silo-server/blob/master/blaise.go) for an example. Silo is designed to be extensible, and each module acts upon the same set of data which is device identifiable, allowing for powerful segmentation.

# Bugs/Suggestions

Feel free to open a GitHub issue for any bugs you find, or suggestions you have.

# License

Apache

# Contact

Email me at farazfazli@gmail.com if you'd like to contact me. And if you're using Silo, I'd love to hear from you and try out your app!

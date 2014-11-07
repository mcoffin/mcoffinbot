function help(target, command) {
    irc("PRIVMSG " + target + " :" + command + " message");
}

function quit(command, source, target) {
    var message = "Quitting";
    if (arguments.length > 3) {
        message = arguments[3];
    }
    irc("QUIT " + message);
}

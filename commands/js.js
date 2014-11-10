function js(command, source, target, suffix) {
    var reconstructed = suffix;
    var message = "";
    try {
        message = eval(reconstructed);
    } catch(err) {
        message = err;
    }
    irc("PRIVMSG " + target + " :" + message);
}

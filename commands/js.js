function js(command, source, target) {
    var reconstructed = "";
    for (var i = 3; i < arguments.length; i++) {
        reconstructed += arguments[i];
    }
    var message = "";
    try {
        message = eval(reconstructed);
    } catch(err) {
        message = err;
    }
    irc("PRIVMSG " + target + " :" + message);
}

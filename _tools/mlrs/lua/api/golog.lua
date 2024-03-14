go_log = require("go_log")
log = {}
log.debug = function(msg)
    go_log.debug(msg)
end
log.info = function(msg)
    go_log.info(msg)
end

log.warn = function(msg)
    go_log.warn(msg)
end

log.error = function(msg)
    go_log.error(msg)
end
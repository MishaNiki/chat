$(function() {
    var conn

    $('.message-submit').click(function() {
        sendMessage(conn)
	});

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://"+ document.location.host +"/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            // Здесь вывод сообщения
            var msg = evt.data
            message = JSON.parse(msg)
            $('<div class="message-cutor"><p class="message-personal-auth">'+ message.auth + ' seys:</p>'+ message.body + '</div>').appendTo($('.messages-content')).addClass('new');
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
});

function sendMessage(conn) {
	msg = $('.message-input').val();
	auth = $('.auth').val();
	if ($.trim(msg) == '') {
		return false;
	}
	if ($.trim(auth) == '') {
		auth = "guest";
    }
    if (!conn) {
        return false;
    }
    conn.send(toJString(auth, msg));
    $('.message-input').val(null);
}

function toJString(auth, msg) {
    return '{"auth" : "' + auth + '", "body":"'+ msg +'"}';
}
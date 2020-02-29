$(function() {

	$('.message-submit').click(function() {
		insertMessage();
	});

	$(window).on('keydown', function(e) {
		if (e.which == 13) {
			insertMessage();
			return false;
  		}
	});
});

var $messages = $('.messages-content'),
	d, h, m;

function setDate(){
	d = new Date()
	if (m != d.getMinutes()) {
		m = d.getMinutes();
		$('<div class="timestamp">' + d.getHours() + ':' + d.getMinutes() + '</div>').appendTo($('.message:last'));
	}
}

function cutorMessage() {
	msg = 'my affairs are very good';
	if ($.trim(msg) == '') {
		return false;
	}
	$('<div class="message-cutor"><b>Server seys:</b><br>' + msg + '</div>').appendTo($('.messages-content')).addClass('new');

}

function insertMessage() {
	msg = $('.message-input').val();
	auth = $('.auth').val();
	if ($.trim(msg) == '') {
		return false;
	}
	if ($.trim(auth) == '') {
		auth = "guest";
	}
    $('<div class="message-personal"><p class="message-personal-auth">'+ auth + ' seys:</p>'+ msg + '</div>').appendTo($('.messages-content')).addClass('new');
	setDate();  // дописать
	$('.message-input').val(null);
	cutorMessage(); // дописать
}

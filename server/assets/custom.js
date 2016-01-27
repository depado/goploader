var loader = $("#upload-loader");
var form = $("#upload-form");
var result = $("#upload-result");
var upurl = $("#upload-url");
var uperror = $("#upload-error")
$('#upload-btn').click(function($e) {
    $e.preventDefault();
    var data = new FormData();
    data.append('file', $("#upload-file")[0].files[0])
    form.fadeOut(400, function() {
        loader.fadeIn();
        loader.promise().done(function() {
            var req = $.ajax({
                url: '/',
                data: data,
                cache: false,
                contentType: false,
                processData: false,
                type: 'POST'
            });
            req.done(function(data) {
                upurl.html("Here is your file :<br /><a href='" + data + "'>" + data + "</a>");
                loader.fadeOut(400, function() {
                    result.fadeIn();
                });
            });
            req.fail(function(jqxhr, statusmsg) {
                loader.fadeOut(400, function() {
                    uperror.html("The file is too big or an error occured on the server.").show();
                    form.fadeIn();
                });
            });
        });
    });
});
$("#upload-again").click(function($e) {
    result.fadeOut(400, function() {
        form.fadeIn();
    });
});
'use strict';

;
(function(document, window, index) {
    var inputs = document.querySelectorAll('.inputfile');
    Array.prototype.forEach.call(inputs, function(input) {
        var label = input.nextElementSibling,
            labelVal = label.innerHTML;

        input.addEventListener('change', function(e) {
            var fileName = '';
            if (this.files && this.files.length > 1)
                fileName = (this.getAttribute('data-multiple-caption') || '').replace('{count}', this.files.length);
            else
                fileName = e.target.value.split('\\').pop();

            if (fileName)
                label.querySelector('span').innerHTML = fileName;
            else
                label.innerHTML = labelVal;
        });

        // Firefox bug fix
        input.addEventListener('focus', function() {
            input.classList.add('has-focus');
        });
        input.addEventListener('blur', function() {
            input.classList.remove('has-focus');
        });
    });
}(document, window, 0));

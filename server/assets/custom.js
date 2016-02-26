'use strict';

var loader = $("#upload-loader");
var form = $("#upload-form");
var result = $("#upload-result");
var upurl = $("#upload-url");
var uperror = $("#upload-error");

$('#upload-btn').click(function($e) {
    $e.preventDefault();
    var data = new FormData();
    data.append('file', $("#upload-file")[0].files[0]);
    data.append('duration', $("#duration").val());
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

var active = $("#upload");
var pages = ["introduction", "client", "curl", "server"];

$("a[id^='toggle-']").click(function(evt) {
    var toggleid = $(this).attr('id').split('-')[1];
    active.fadeOut(400, function() {
        window.scrollTo(0, 0);
        $("#" + toggleid).fadeIn();
    });
    if (window.location.href.indexOf("#") > -1) {
        if (toggleid == "upload") {
            window.location = window.location.href.substring(0, window.location.href.indexOf("#")) + "#";
        } else {
            window.location = window.location.href.substring(0, window.location.href.indexOf("#")) + "#" + toggleid;
        }
    } else if (toggleid != "upload") {
        window.location = window.location + "#" + toggleid;
    }
    active = $("#" + toggleid);
    evt.preventDefault();
});

$(document).ready(function() {
    for (var i = 0; i < pages.length; i++) {
        var page = pages[i];
        if (window.location.href.indexOf("#"+page) > -1) {
            var current = window.location.href.substring(window.location.href.indexOf("#"+page))
            active.fadeOut(400, function() {
                $("#"+page).fadeIn(400, function () {
                    $("html, body").animate({scrollTop: $(current).offset().top }, 200);
                });
            });
            active = $("#"+page);
            return;
        }
    }
});

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

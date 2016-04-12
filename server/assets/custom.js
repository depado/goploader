'use strict';

var loader = $("#upload-loader");
var form = $("#upload-form");
var result = $("#upload-result");
var upurl = $("#upload-url");
var upclipboard = $("#upload-clipboard");
var uperror = $("#upload-error");
var oneviewlabel = $("label[for=one-view]");
var sourcelabel = $("label[for=source]");
var filelabel = $("label[for=upload-file]");
var uploadsum = $("#upload-summary");

var currentfile = "";
var active = $("#upload");
var pages = ["introduction", "client", "curl", "server"];
toastr.options = {
    "progressBar": true,
    "positionClass": "toast-top-right",
    "showDuration": "300",
    "hideDuration": "300",
    "timeOut": "2000",
    "extendedTimeOut": "1000",
    "showEasing": "swing",
    "hideEasing": "linear",
    "showMethod": "fadeIn",
    "hideMethod": "fadeOut"
};

var clipboard = new Clipboard('#upload-clipboard');
clipboard.on('success', function(e) {
    toastr.success('Copied to Clipboard');
    e.clearSelection();
});

function getsum() {
    var sum = "Your file will live for " + $("#duration option:selected").text() + " and will be visible ";
    if ($('#one-view').is(":checked")) {
        sum += "only once.";
    } else {
        sum += "without restrictions.";
    }
    return sum
}

$('#one-view').change(function() {
    if ($('#one-view').is(":checked")) {
        oneviewlabel.text("One Download");
    } else {
        oneviewlabel.text("No Restriction");
    }
    uploadsum.text(getsum());
});

$('#source').change(function() {
    if ($('#source').is(":checked")) {
        sourcelabel.text("Text");
        filelabel.fadeOut(200, function() {
            uploadsum.text(getsum());
            $('#upload-text').fadeIn(200);
        });
    } else {
        sourcelabel.text("File");
        $('#upload-text').fadeOut(200, function() {
            uploadsum.text(getsum());
            filelabel.fadeIn(200);
        });
    }
});

$("#duration").change(function() {
    uploadsum.text(getsum());
});

$("#upload-again").click(function($e) {
    result.fadeOut(400, function() {
        form.fadeIn();
    });
});

$('#upload-btn').click(function($e) {
    $e.preventDefault();
    var data = new FormData();
    if ($('#source').is(":checked")) {
        if ($("#upload-text").val() == "") {
            toastr.success('Please paste some text')
            return
        }
        data.append('file', new File([new Blob([$("#upload-text").val()])], "stdin"));
    } else {
        if ($("#upload-file")[0].files.length != 1) {
            toastr.success('Please select a file');
            return
        }
        data.append('file', $("#upload-file")[0].files[0]);
    }
    data.append('duration', $("#duration").val());
    if ($('#one-view').is(":checked")) {
        data.append('once', 'true');
    }
    toastr.success('File transfer in progress');
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
                upurl.html("Here is your file :<br /><a href='" + data + "' target='_blank'>" + data + "</a>");
                upclipboard.attr("data-clipboard-text", data);
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
        if (window.location.href.indexOf("#" + page) > -1) {
            var current = window.location.href.substring(window.location.href.indexOf("#" + page))
            active.fadeOut(400, function() {
                $("#" + page).fadeIn(400, function() {
                    $("html, body").animate({
                        scrollTop: $(current).offset().top
                    }, 200);
                });
            });
            active = $("#" + page);
            return;
        }
    }
    if ($('#one-view').is(":checked")) {
        oneviewlabel.text("One Download");
    } else {
        oneviewlabel.text("No Restriction");
    }
    if ($('#source').is(":checked")) {
        sourcelabel.text("Text");
        filelabel.fadeOut(200, function() {
            $('#upload-text').fadeIn(200);
        });
    } else {
        sourcelabel.text("File");
        $('#upload-text').fadeOut(200, function() {
            filelabel.fadeIn(200);
        });
    }
    uploadsum.text(getsum());
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

            if (fileName) {
                currentfile = fileName;
                label.querySelector('span').innerHTML = fileName;
                buildsummary();
            } else {
                label.innerHTML = labelVal;
            }
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

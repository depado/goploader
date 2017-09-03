'use strict';

var finalurl;

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

var lineslabel = $("label[for=lines]");
var themelabel = $("label[for=theme]");

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

function updateQueryStringParameter(uri, key, value) {
    var re = new RegExp("([?&])" + key + "=.*?(&|$)", "i");
    var separator = uri.indexOf('?') !== -1 ? "&" : "?";
    if (uri.match(re)) {
        return uri.replace(re, '$1' + key + "=" + value + '$2');
    }
    else {
        return uri + separator + key + "=" + value;
    }
}

function removeURLParameter(url, parameter) {
    var urlparts= url.split('?');   
    if (urlparts.length>=2) {
        var prefix= encodeURIComponent(parameter)+'=';
        var pars= urlparts[1].split(/[&;]/g);
        for (var i= pars.length; i-- > 0;) {    
            if (pars[i].lastIndexOf(prefix, 0) !== -1) {  
                pars.splice(i, 1);
            }
        }
        url= urlparts[0] + (pars.length > 0 ? '?' + pars.join('&') : "");
        return url;
    } else {
        return url;
    }
}

function getsum() {
    var sum = "Your file will live for " + $("#duration option:selected").text() + " and will be visible ";
    if ($('#one-view').is(":checked")) {
        sum += "only once.";
    } else {
        sum += "without restrictions.";
    }
    return sum
}

$('#language').change(function() {
    var lang = $('#language').val();
    var segments = finalurl.pathname.split("/").length - 1;
    if (segments > 3) {
        if (lang=="none") {
            var value = finalurl.pathname.substring(finalurl.pathname.lastIndexOf('/'));
            finalurl.pathname = finalurl.pathname.replace(value, "");
        } else {
            var value = finalurl.pathname.substring(finalurl.pathname.lastIndexOf('/') + 1);
            finalurl.pathname = finalurl.pathname.replace(value, lang);
        }
    } else {
        if (lang!="none") {
            finalurl.pathname = finalurl.pathname + "/" +lang;
        }
    }
    $('#final-url').attr('href', finalurl.toString());
    $('#final-url').text(finalurl.toString());
});

$('#lines').change(function() {
    var url = finalurl.toString();
    if ($('#lines').is(":checked")) {
        lineslabel.text("With Lines");
        finalurl = new URL(updateQueryStringParameter(url, "lines", "true"))
    } else {
        lineslabel.text("No Lines");
        finalurl = new URL(removeURLParameter(url, "lines"))
    }
    $('#final-url').attr('href', finalurl.toString());
    $('#final-url').text(finalurl.toString());
});

$('#theme').change(function() {
    var url = finalurl.toString();
    if ($('#theme').is(":checked")) {
        themelabel.text("Light");
        finalurl = new URL(updateQueryStringParameter(url, "theme", "light"))
    } else {
        themelabel.text("Dark");
        finalurl = new URL(removeURLParameter(url, "theme"))
    }
    $('#final-url').attr('href', finalurl.toString());
    $('#final-url').text(finalurl.toString());
});

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
    $('#language').prop('selectedIndex', 0);
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
                xhr: function() {
                    var percentComplete = 0;
                    var xhr = new window.XMLHttpRequest();

                    xhr.upload.addEventListener("progress", function(evt) {
                        if (evt.lengthComputable) {
                            var percentComplete = evt.loaded / evt.total;
                            percentComplete = parseInt(percentComplete * 100);
                            $(".progress>div").css("width", percentComplete+"%");
                        }
                    }, false);

                    return xhr;
                },
                url: '/',
                data: data,
                cache: false,
                contentType: false,
                processData: false,
                type: 'POST'
            });
            req.done(function(data) {
                finalurl = new URL(data)
                upurl.html("Here is your file :<br /><a id='final-url' href='" + data + "' target='_blank'>" + data + "</a>");
                upclipboard.attr("data-clipboard-text", data);
                loader.fadeOut(400, function() {
                    result.fadeIn(400, function() {
                        $(".progress>div").css("width", "0%");
                    });
                });
                $("label[for=upload-file] span").text("Choose a file…");
                $("#upload-file").replaceWith($("#upload-file").val('').clone(true));
                $("#upload-text").val('');
            });
            req.fail(function(jqxhr, statusmsg) {
                loader.fadeOut(400, function() {
                    uperror.html("The file is too big or an error occured on the server.").show();
                    form.fadeIn(400, function() {
                        $(".progress>div").css("width", "0%");
                    });
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

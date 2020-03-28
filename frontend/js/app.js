"use strict";

function initUserPage(username){
    function doPollForUserStatus(){
        $.get('/api/@'+username+'/status', function(data) {
            if($('#progress-identifier').hasClass('dn')) {
                $('#progress-identifier').removeClass('dn').addClass('dib')
            } else if($('#progress-identifier').hasClass('dib')) {
                $('#progress-identifier').removeClass('dib').addClass('dn')
            }
            if (data.goodreads_connected || data.books_count > 0 ) {
                window.location.reload();
            } else {
                setTimeout(doPollForUserStatus,2000);
            }
        });
    }
    doPollForUserStatus();
}

function initUserBookPage(bookEditionID) {
    //buttons
    $.fn.editableform.buttons = '<button type="submit" class="editable-submit hover-gray pointer input-reset bn bg-transparent sans-serif">ok</button>' +
        '<button type="button" class="editable-cancel hover-gray pointer input-reset bn bg-transparent sans-serif">cancel</button>';
    $.fn.editable.defaults.mode = 'inline';
    $.fn.editable.defaults.ajaxOptions = {
        type: 'patch',
        dataType: 'json',
        contentType: 'application/json'
    };
    $.fn.combodate.defaults = {
        minYear: 1950,
        maxYear: 2020,
        yearDescending: true
    };
    $('#book_review').editable({
        emptytext: 'leave review',
        showbuttons: 'bottom',
        params: function (params) {
            params["edition"] = bookEditionID;
            return JSON.stringify(params);
        },
        display: function (value, response) {
            if (response) {
                $(this).html(response.content_html);
            }
        }
    });

    $('.reading-note').editable({
        emptytext: 'add note',
        params: function (params) {
            params["edition"] = bookEditionID;
            return JSON.stringify(params);
        }
    });

    $('.reading-start-date').editable({
        params: function (params) {
            params["edition"] = bookEditionID;
            return JSON.stringify(params);
        },
        display: function (value, response) {
            if (response) {
                $(this).html(response.start_date);
                $(this).parent().parent().find(".reading-duration").html(response.duration);
            }
        }
    });
    $('.reading-finish-date').editable({
        params: function (params) {
            params["edition"] = bookEditionID;
            return JSON.stringify(params);
        },
        display: function (value, response) {
            if (response) {
                $(this).html(response.finish_date);
                $('form.status-change button').removeClass('underline');
                $('form.status-change.status-1 button').addClass('underline');
                $(this).parent().parent().find(".reading-duration").html(response.duration);
            }
        }
    });
}

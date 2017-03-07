function err_message_quietly(msg, f) {
    $.layer({
        title: false,
        closeBtn: false,
        time: 2,
        dialog: {
            msg: msg
        },
        end: f
    });
}

function ok_message_quietly(msg, f) {
    $.layer({
        title: false,
        closeBtn: false,
        time: 1,
        dialog: {
            msg: msg,
            type: 1
        },
        end: f
    });
}

function my_confirm(msg, btns, yes_func, no_func) {
    $.layer({
        shade: [ 0 ],
        area: [ 'auto', 'auto' ],
        dialog: {
            msg: msg,
            btns: 2,
            type: 4,
            btn: btns,
            yes: yes_func,
            no: no_func
        }
    });
}

function handle_quietly(json, f) {
    if (json.msg.length > 0) {
        err_message_quietly(json.msg);
    } else {
        ok_message_quietly("successfully:-)", f);
    }
}

// - business function -
function all_select() {
    var boxes = $("input[type=checkbox]");
    for (var i = 0; i < boxes.length; i++) {
        boxes[i].checked="checked";
    }
}

function reverse_select() {
    var boxes = $("input[type=checkbox]");
    for (var i = 0; i < boxes.length; i++) {
        if (boxes[i].checked) {
            boxes[i].checked=""
        } else {
            boxes[i].checked="checked";
        }
    }
}

function batch_solve() {
    var boxes = $("input[type=checkbox]");
    var ids = []
    for (var i = 0; i < boxes.length; i++) {
        if (boxes[i].checked) {
            ids.push($(boxes[i]).attr("alarm"))
        }
    }

    $.post("/event/solve", {"ids": ids.join(',,')}, function(msg){
        if (msg=="") {
            location.reload();
        } else {
            alert(msg);
        }
    });
}

function solve(id) {
    $.post("/event/solve", {"ids": id}, function(msg){
        if (msg=="") {
            location.reload();
        } else {
            alert(msg);
        }
    });
}

// Gets data from provided url and updates DOM element.
function generate_os_data(url, element) {
    $.get(url, function (d) {
        if (d.msg == "success") {
            $(element).text(d.data);
        } else {
            $(element).text(d.msg);
        }
    }, "json");
}

// If dataTable with provided ID exists, destroy it.
function destroy_dataTable(table_id) {
    var table = $("#" + table_id);
    var ex = document.getElementById(table_id);
    if ($.fn.DataTable.fnIsDataTable(ex)) {
        table.hide().dataTable().fnClearTable();
        table.dataTable().fnDestroy();
    }
}

//DataTables
//Sort file size data.
jQuery.extend(jQuery.fn.dataTableExt.oSort, {
    "file-size-units": {
        K: 1024,
        M: Math.pow(1024, 2),
        G: Math.pow(1024, 3),
        T: Math.pow(1024, 4),
        P: Math.pow(1024, 5),
        E: Math.pow(1024, 6)
    },

    "file-size-pre": function (a) {
        var x = a.substring(0, a.length - 1);
        var x_unit = a.substring(a.length - 1, a.length);
        if (jQuery.fn.dataTableExt.oSort['file-size-units'][x_unit]) {
            return parseInt(x * jQuery.fn.dataTableExt.oSort['file-size-units'][x_unit], 10);
        }
        else {
            return parseInt(x + x_unit, 10);
        }
    },

    "file-size-asc": function (a, b) {
        return ((a < b) ? -1 : ((a > b) ? 1 : 0));
    },

    "file-size-desc": function (a, b) {
        return ((a < b) ? 1 : ((a > b) ? -1 : 0));
    }
});

//DataTables
//Sort numeric data which has a percent sign with it.
jQuery.extend(jQuery.fn.dataTableExt.oSort, {
    "percent-pre": function (a) {
        var x = (a === "-") ? 0 : a.replace(/%/, "");
        return parseFloat(x);
    },

    "percent-asc": function (a, b) {
        return ((a < b) ? -1 : ((a > b) ? 1 : 0));
    },

    "percent-desc": function (a, b) {
        return ((a < b) ? 1 : ((a > b) ? -1 : 0));
    }
});

//DataTables
//Sort IP addresses
jQuery.extend(jQuery.fn.dataTableExt.oSort, {
    "ip-address-pre": function (a) {
        // split the address into octets
        //
        var x = a.split('.');

        // pad each of the octets to three digits in length
        //
        function zeroPad(num, places) {
            var zero = places - num.toString().length + 1;
            return Array(+(zero > 0 && zero)).join("0") + num;
        }

        // build the resulting IP
        var r = '';
        for (var i = 0; i < x.length; i++)
            r = r + zeroPad(x[i], 3);

        // return the formatted IP address
        //
        return r;
    },

    "ip-address-asc": function (a, b) {
        return ((a < b) ? -1 : ((a > b) ? 1 : 0));
    },

    "ip-address-desc": function (a, b) {
        return ((a < b) ? 1 : ((a > b) ? -1 : 0));
    }
});

/*******************************
 Data Call Functions
 *******************************/

var dashboard = {};

dashboard.getRam = function () {
    $.get("/page/memory", function (data) {
        if (data.msg == "success") {
            var ram_total = data.data[0];
            var ram_used = Math.round((data.data[1] / ram_total) * 100);
            var ram_free = Math.round((data.data[2] / ram_total) * 100);

            $("#ram-total").text(ram_total);
            $("#ram-used").text(data.data[1]);
            $("#ram-free").text(data.data[2]);

            $("#ram-free-per").text(ram_free);
            $("#ram-used-per").text(ram_used);
        }
    }, "json");
}

dashboard.getDf = function () {
    $.get("/page/df", function (data) {
        var table = $("#df_dashboard");
        var ex = document.getElementById("df_dashboard");
        if ($.fn.DataTable.fnIsDataTable(ex)) {
            table.hide().dataTable().fnClearTable();
            table.dataTable().fnDestroy();
        }

        table.dataTable({
            aaData: data.data,
            aoColumns: [
                { sTitle: "Filesystem" },
                { sTitle: "BTotal", sType: "file-size" },
                { sTitle: "BUsed", sType: "file-size" },
                { sTitle: "BFree", sType: "file-size" },
                { sTitle: "BUse%", sType: "percent" },
                { sTitle: "Mounted" },
                { sTitle: "ITotal", sType: "file-size" },
                { sTitle: "IUsed", sType: "file-size" },
                { sTitle: "IFree", sType: "file-size" },
                { sTitle: "IUse%", sType: "percent" },
                { sTitle: "Vfstype" }
            ],
            bPaginate: false,
            bFilter: false,
            bAutoWidth: true,
            bInfo: false
        }).fadeIn();
    }, "json");
}

dashboard.getCpu = function () {
    $.get("/page/cpu/usage", function (data) {
        if (data.msg != "success") {
            return
        }
        
        var table = $("#cpu_dashboard");
        var ex = document.getElementById("cpu_dashboard");
        if ($.fn.DataTable.fnIsDataTable(ex)) {
            table.hide().dataTable().fnClearTable();
            table.dataTable().fnDestroy();
        }

        table.dataTable({
            aaData: data.data,
            aoColumns: [
                { sTitle: "idle"},
                { sTitle: "busy"},
                { sTitle: "user"},
                { sTitle: "nice"},
                { sTitle: "system"},
                { sTitle: "iowait"},
                { sTitle: "irq"},
                { sTitle: "softirq"},
                { sTitle: "steal"},
                { sTitle: "guest"}
            ],
            bPaginate: false,
            bFilter: false,
            bAutoWidth: true,
            bInfo: false
        }).fadeIn();
    }, "json");
}

dashboard.getDiskstats = function () {
    $.get("/page/diskio", function (data) {
        var table = $("#diskstats_dashboard");
        var ex = document.getElementById("diskstats_dashboard");
        if ($.fn.DataTable.fnIsDataTable(ex)) {
            table.hide().dataTable().fnClearTable();
            table.dataTable().fnDestroy();
        }

        table.dataTable({
            aaData: data.data,
            aoColumns: [
                { sTitle: "Device"},
                { sTitle: "rrqm/s"},
                { sTitle: "wrqm/s"},
                { sTitle: "r/s"},
                { sTitle: "w/s"},
                { sTitle: "rkB/s"},
                { sTitle: "wkB/s"},
                { sTitle: "avgrq-sz"},
                { sTitle: "avgqu-sz"},
                { sTitle: "await"},
                { sTitle: "svctm"},
                { sTitle: "%util"},
            ],
            bPaginate: false,
            bFilter: false,
            bAutoWidth: true,
            bInfo: false
        }).fadeIn();
    }, "json");
}

dashboard.getOs = function () {
    generate_os_data("/proc/kernel/version", "#os-info");
    generate_os_data("/proc/kernel/hostname", "#os-hostname");
    generate_os_data("/system/date", "#os-time");
    generate_os_data("/page/system/uptime", "#os-uptime");

    $.get("/version", function(d){
        $("#agent-version").text(d);
    });
}

dashboard.getLoadAverage = function () {
    $.get("/page/system/loadavg", function (d) {
        if (d.msg != "success") {
            return
        }
        $("#cpu-1min").text(d.data[0][0]);
        $("#cpu-5min").text(d.data[1][0]);
        $("#cpu-15min").text(d.data[2][0]);
        $("#cpu-1min-per").text(d.data[0][1]);
        $("#cpu-5min-per").text(d.data[1][1]);
        $("#cpu-15min-per").text(d.data[2][1]);
    }, "json");
    generate_os_data("/proc/cpu/num", "#core-number");
}

/**
 * Refreshes all widgets. Does not call itself recursively.
 */
dashboard.getAll = function () {
    for (var item in dashboard.fnMap) {
        if (dashboard.fnMap.hasOwnProperty(item) && item !== "all") {
            dashboard.fnMap[item]();
        }
    }
}

dashboard.fnMap = {
    all: dashboard.getAll,
    ram: dashboard.getRam,
    df: dashboard.getDf,
    os: dashboard.getOs,
    load: dashboard.getLoadAverage,
    cpu: dashboard.getCpu,
    diskstats: dashboard.getDiskstats,
};

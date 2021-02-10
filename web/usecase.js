let attributes;
let attributeHtml;
let attributeTypeMap;
$(document).ready(function () {
    $('input[type=datetime-local]').val(new Date().toJSON().slice(0, 16));

    let coll = document.getElementsByClassName("collapsible");
    for (let i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function () {
            this.classList.toggle("active");
            let content = this.nextElementSibling;
            if (content.style.display === "block") {
                content.style.display = "none";
            } else {
                content.style.display = "block";
            }
        });
    }

    populateAttribute();
});

$(document).ajaxStart(function () {
    $("#loader").css("display", "block");
}).ajaxSuccess(function () {
    $("#loader").css("display", "none");
}).ajaxError(function () {
    $("#loader").css("display", "none");
});

function copyToClipboard(e) {
    let rightTdIndex = $(e).parent().index() + 1;
    let rightTd = $(e).closest("tr").find("td:eq(" + rightTdIndex + ")");
    execCopyToClipboard($(rightTd).text());
}

function execCopyToClipboard(txt) {
    let el = document.createElement('textarea');
    el.value = txt;
    el.setAttribute('readonly', '');
    el.style.position = 'absolute';
    el.style.left = '-9999px';
    document.body.appendChild(el);
    el.select();
    document.execCommand('copy');
    document.body.removeChild(el);
}

$(function () {
    $("#option-type-input").change(function () {
        populateAttribute();
        $('#dynamicTable').empty();
    });

    $("#option-service-input").change(function () {
        populateAttribute();
        $('#dynamicTable').empty();
    });

    $("#btn-search").click(function (e) {
        e.preventDefault();
        let error = validateForm();
        if (error === '') {
            $('#error').html('<div class=""></div>');
            let query = buildQuery();
            if (query !== "") {
                search(query);
            }
        } else {
            $('#error').html('<div class="alert alert-danger">' + error + '</div>');
        }
    });

    $("#btn-generate-query").click(function (e) {
        e.preventDefault();
        let error = validateForm();
        if (error === '') {
            $('#error').html('<div class=""></div>');
            let query = buildQuery();
            if (query !== "") {
                generateQuery(query);
            }
        } else {
            $('#error').html('<div class="alert alert-danger">' + error + '</div>');
        }
    });

    $("#btn-new-group").click(function (e) {
        e.preventDefault();
        let $table = $('<table class="table item_table" id="item_table[]" style="margin-top: 5px">');
        let $tbody = $table.append('<tbody />').children('tbody');
        let $tfoot = $table.append('<tfoot />').children('tfoot');

        $tbody
            .append("<th><select name=\"item_group_logical_operator[]\" class=\"form-control form-control-sm item_group_logical_operator\"><option value=\"&&\"  selected=\"selected\">AND</option>\\n' +\n" +
                "'<option value=\"||\">OR</option></select></th>")
            .append("<th>Logical</th>")
            .append("<th>Attribute</th>")
            .append("<th>Operator</th>")
            .append("<th>Value</th>")
            .append("<th></th>");

        $tfoot
            .append("<th></th>")
            .append("<th></th>")
            .append("<th colspan=\"3\"><button type=\"button\" name=\"add\" class=\"btn btn-success btn-sm add btn-block\"><span class=\"glyphicon glyphicon-plus\"></span> Add attribute</button></th>")
            .append("<th><button type=\"button\" name=\"remove_group\" class=\"btn btn-warning btn-sm remove_group btn-block\" data-toggle=\"tooltip\" data-placement=\"top\" title=\"Delete context group\"><span class=\"btn-inner--icon\"></span>Remove context group</button></th>");


        $table.on('click', '.add', function () {
            let display = "none";
            if ($tbody[0].rows.length > 0) {
                display = "block";
            }

            let html = '';
            html += '<tr>';
            html += '<td></td>';
            html += '<td><select name="item_logical_operator[]" class="form-control form-control-sm item_logical_operator" style="display:' + display + '"><option value="&&"  selected="selected">AND</option>\n' +
                '<option value="||">OR</option></select></td>';
            html += '<td><select name="item_attribute[]" class="form-control form-control-sm item_attribute">' + attributeHtml + '</select></td>';
            html += '<td><select name="item_operator[]" class="form-control form-control-sm item_operator">' +
                '<option value="=" selected="selected">Equal</option>' +
                '<option value=">">Greater than</option>' +
                '<option value=">=">Greater than or eq</option>' +
                '<option value="<">Less than</option>' +
                '<option value="<=">Less than or eq</option></select></td>';
            html += '<td><input type="text" name="item_value[]" class="form-control form-control-sm item_value" /></td>';
            html += '<td><button type="button" name="remove" class="btn btn-danger btn-sm remove"><span class="btn-inner--icon"><i class="ni ni-fat-remove"></i></span></button></td></tr>';
            $table.append(html);
        });

        $table.on('click', '.remove', function () {
            $(this).closest('tr').remove();
        });

        $table.on('click', '.remove_group', function () {
            $('#error').html('<div class=""></div>');
            $table.remove();
        });

        $table.appendTo('#dynamicTable');
    });
});

function populateAttribute() {
    attributes = getAttribute();
    attributeHtml = "";
    attributeTypeMap = new Map();
    $.each(attributes, function (key, data) {
        attributeHtml += "<option value=\"" + data[0] + "\">" + data[1] + "</option>"
        attributeTypeMap.set(data[0], data[2]);
    });
}

function buildQuery() {
    let startTimestamp = $('#start-datetime-input').val();
    let endTimestamp = $('#end-datetime-input').val();
    let starts = startTimestamp.split("T");
    let ends = endTimestamp.split("T");

    let startDate, endDate;
    let startTime, endTime;
    if (starts.length >= 2) {
        startDate = starts[0];
        startTime = parseInt(starts[1].split(":")[0]);
    }
    if (ends.length >= 2) {
        endDate = ends[0];
        endTime = parseInt(ends[1].split(":")[0]);
    }

    if (startDate !== endDate) {
        swal("Warning", "Date must be same", "warning");
        return "";
    }
    if (startTime > endTime) {
        swal("Warning", "Start time can't greater than end time", "warning");
        return "";
    }

    let contextQueries = [];
    let isFirstDateKey = true;
    for (let i = startTime; i <= endTime; i++) {
        if (i === startTime) {
            contextQueries.push("(");
        }
        if (!isFirstDateKey) {
            contextQueries.push("||");
        }
        let strTime = "" + i;
        if (i < 10) {
            strTime = "0" + i;
        }
        contextQueries.push("datekey=\"" + startDate.replaceAll("-", "") + strTime + "\"");
        if (i === endTime) {
            contextQueries.push(")");
        }
        isFirstDateKey = false;
    }

    contextQueries.push("&& timestamp>=\"" + startTimestamp + "\"");
    contextQueries.push("&& timestamp<=\"" + endTimestamp + "\"");

    $('.item_table').each(function () {
        let groupLogicalOperator = $(this).find('th .item_group_logical_operator').val();
        let isContain = false;
        let isFirst = true;
        let combination = "";
        $(this).find('tr').each(function (i, row) {
            let itemLogicalOperator = $(row).find('.item_logical_operator').val();
            let itemAttribute = $(row).find('.item_attribute').val();
            let itemOperator = $(row).find('.item_operator').val();
            let itemValue = $(row).find('.item_value').val().trim();
            if (!isFirst) {
                combination += " " + itemLogicalOperator;
            }
            if (attributeTypeMap.get(itemAttribute) === "alphanumeric") {
                itemValue = "\"" + itemValue + "\"";
            }
            combination += " " + itemAttribute + itemOperator + itemValue;
            isContain = true;
            isFirst = false;
        });

        if (isContain) {
            contextQueries.push(groupLogicalOperator);
            contextQueries.push("(", combination, ")");
        }
    });

    let query = contextQueries.join(" ");
    console.log("QUERY = " + query);
    return query;
}

function validateForm() {
    let error = '';
    $('.item_logical_operator').each(function (i, data) {
        if ($(this).val() === '') {
            error += "<p>Select logical operator at " + (i + 1) + " row</p>";
            return false;
        }
    });

    $('.item_attribute').each(function (i, data) {
        if ($(this).val() === '') {
            error += "<p>Select attribute at " + (i + 1) + " row</p>";
            return false;
        }
    });

    $('.item_value').each(function (i, data) {
        if ($(this).val() === '') {
            error += "<p>Enter value at " + (i + 1) + " row</p>";
            return false;
        }
    });
    return error;
}

function getAttribute() {
    let service = $('#option-service-input').val();
    let source = $('#option-type-input').val();
    let arrayReturn = [];
    $.ajax({
        url: "attributes?service=" + service + "&source=" + source,
        async: false,
        type: "GET",
        dataType: 'json',
        success: function (response) {
            if (response !== null) {
                if (response.success === true) {
                    if (response.data != null) {
                        if (response.data.length > 0) {
                            response.data.forEach(function (object, index) {
                                arrayReturn.push([
                                    object.attribute,
                                    object.name,
                                    object.type
                                ]);
                                return false;
                            });
                        }
                    }
                }
            }
        }
    });
    return arrayReturn;
}

function generateQuery(contextQuery) {
    if ($('#start-datetime-input').val() === "" || $('#end-datetime-input').val() === "" || $('#option-service-input').val() === "") {
        swal("Warning", "You must complete the data", "warning");
    }
    let limit = parseInt($('#option-limit-input').val());
    let param = {};
    param.service = $('#option-service-input').val();
    param.message_query = $('#text-message-input').val().trim();
    param.context_query = contextQuery;
    param.limit = limit;
    param.type = $('#option-type-input').val();

    $('#generated-query-input').val("");
    $.ajax({
        url: "generate_query",
        async: false,
        type: "POST",
        data: JSON.stringify(param),
        contentType: "application/json; charset=utf-8",
        dataType: 'json',
        success: function (response) {
            if (response !== null) {
                console.log(response);
                if (response.success) {
                    $("#modal-default .modal-title").text("Generated query");
                    $("#modal-default .modal-body").html("<code>" + response.data + "</code>");
                    $('#modal-default').modal();
                } else {
                    swal("Oops", response.error, "error");
                }
            } else {
                swal("Oops", "Something went wrong!", "error");
            }
        },
        error: function (xhr, status, error) {
            var err = JSON.parse(xhr.responseText);
            swal("Oops", err.error.message, "error");
        }
    });
}

function search(contextQuery) {
    if ($('#start-datetime-input').val() === "" || $('#end-datetime-input').val() === "" || $('#option-service-input').val() === "") {
        swal("Warning", "You must complete the data", "warning");
    }
    let limit = parseInt($('#option-limit-input').val());
    let param = {};
    param.service = $('#option-service-input').val();
    param.message_query = $('#text-message-input').val().trim();
    param.context_query = contextQuery;
    param.limit = limit;
    param.type = $('#option-type-input').val();

    $.ajax({
        url: "query",
        async: true,
        type: "POST",
        data: JSON.stringify(param),
        contentType: "application/json; charset=utf-8",
        dataType: 'json',
        success: function (response) {
            if (response !== null) {
                console.log(response);
                if (response.success) {
                    setDataTable(response.data);
                } else {
                    swal("Oops", response.error, "error");
                }
            } else {
                swal("Oops", "Something went wrong!", "error");
            }
        },
        error: function (xhr, status, error) {
            let err = JSON.parse(xhr.responseText);
            swal("Oops", err.error.message, "error");
        }
    });
}

function setDataTable(data) {
    let htmlResult = [];
    $.each(data, function (i, data) {
        data.message = data.message.replaceAll(/<br\/>|<p>|<\/p>|<br>/g, " ");
        let row = $(document.createElement("tr"));
        row.append($('<td></td>').append(data.timestamp));
        row.append($('<td></td>').append(data.hostname));
        row.append($('<td></td>').append(data.flowid));
        row.append($('<td></td>').append(data.type));
        row.append("<td><button type=\"button\" name=\"btn-show-context\" class=\"btn btn-secondary btn-sm btn-show-context \"><span class=\"btn-inner--icon\"></span>Show context</button>" +
            "<button type=\"button\" class=\"btn btn-info btn-sm\" onclick=\"copyToClipboard(this)\">Copy</button></td>");
        row.append($('<td style="white-space:nowrap;overflow: hidden;text-overflow: ellipsis;"><div></td>').append(data.message));

        row.on('click', '.btn-show-context', function () {
            $("#modal-default .modal-title").text("Available context");
            $("#modal-default .modal-body").html(data.context);
            $('#modal-default').modal();
        });

        htmlResult.push(row);
    });

    $('#tableResult tbody')
        .empty()
        .append(...htmlResult);

}
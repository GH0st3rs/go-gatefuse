function createGateTableRow(item) {
    var $row = $('<tr>', { 'id': item.UUID });
    // Checkbox column
    var $checkbox = $('<input>', {
        'type': 'checkbox',
        'class': 'form-check-input'
    });
    if (item.Active) {
        $checkbox.prop('checked', true);
    }
    var $checkboxDiv = $('<div>', { 'class': 'form-switch ai-c' }).append($checkbox);
    var $checkboxColumn = $('<th>', { 'scope': 'row' }).append($checkboxDiv);

    // Append columns to the row
    $row.append(
        $checkboxColumn,
        // Other columns
        $('<td>').text(item.DstAddr),
        $('<td>').text(item.DstPort),
        $('<td>').text(item.SrcAddr),
        $('<td>').text(item.SrcPort),
        $('<td>').text(item.Protocol),
        $('<td>').text(item.Comment)
    );
    return $row;
}

function clickOnRow(event) {
    var row = event.target.parentElement;
    var rowId = row.getAttribute('id');
    $('#createNewRuleForm #uuidNew').val(rowId);
    $('#createNewRuleForm #activeNew').prop('checked', $('#' + rowId + ' th input').prop('checked'));

    var dstAddrNew = $('#' + rowId + ' td:eq(0)').text()
    $('#createNewRuleForm #dstAddrNew').val(dstAddrNew);

    var dstPortNew = $('#' + rowId + ' td:eq(1)').text()
    $('#createNewRuleForm #dstPortNew').val(dstPortNew);

    var srcAddrNew = $('#' + rowId + ' td:eq(2)').text()
    $('#createNewRuleForm #srcAddrNew').val(srcAddrNew);

    var srcPortNew = $('#' + rowId + ' td:eq(3)').text()
    $('#createNewRuleForm #srcPortNew').val(srcPortNew);

    var protoNew = $('#' + rowId + ' td:eq(4)').text()
    $('#createNewRuleForm #protoNew').val(protoNew);

    var commentNew = $('#' + rowId + ' td:eq(5)').text()
    $('#createNewRuleForm #commentNew').val(commentNew);

    $('#ruleSettingsModal').modal('show');
}

function toggleRule(event) {
    var $row = $(this).parent().parent().parent()
    var formData = JSON.stringify({ "UUID": $row.prop('id'), "Active": $(this).prop('checked') });
    $.ajax({
        type: "POST",
        url: "/toggle",
        dataType: "json",
        contentType: "application/json",
        data: formData
    }).done(function (data) {
        if (!data.status) {
            $(this).prop('checked', false);
        }
    });
}

function loadGateTable() {
    $.ajax({
        type: "GET",
        url: "/list",
        dataType: "json"
    }).done(function (data) {
        if (data.status) {
            var $dataTable = $('#gateTable tbody');
            $dataTable.empty();
            data.routes.forEach(function (item) {
                var $row = createGateTableRow(item);
                // Append the row to the table
                $dataTable.append($row);
            });
            // gateTable Add click events to every column in the table
            $('#gateTable tbody td').click(clickOnRow);
            $('#gateTable input[type="checkbox"]').change(toggleRule);
        } else {
            alert(data.error);
        }
    });
}

function validateNewRuleForm() {
    var form = $('#createNewRuleForm')[0];
    if (form.checkValidity() === false) {
        return false;
    }
    form.classList.add('was-validated');
    return true;
}

function closeNewRuleForm() {
    var $form = $('#createNewRuleForm');
    // Remove validation
    $form.removeClass("was-validated");
    // Clean the form after sending
    $form[0].reset();
    $('#uuidNew').val('');
    // Hide this modal window
    $('#ruleSettingsModal').modal('hide');
}

function createRequest() {
    // Serialize the form data
    var formData = $('#createNewRuleForm').serialize();
    $.ajax({
        type: "POST",
        url: "/create",
        data: formData,
        dataType: "json"
    }).done(function (data) {
        if (data.status) {
            var $dataTable = $('#gateTable tbody');
            var $row = createGateTableRow(data.response);
            // Append the row to the table
            $dataTable.append($row);
        } else {
            alert(data.error);
        }
    });
}

function updateRequest() {
    // Serialize the form data
    var formData = $('#createNewRuleForm').serialize();
    $.ajax({
        type: "POST",
        url: "/update",
        data: formData,
        dataType: "json"
    }).done(function (data) {
        if (data.status) {
            var $row = $('#' + data.response.UUID);
            $row = createGateTableRow(data.response);
        } else {
            alert(data.error);
        }
    });
}

document.addEventListener('DOMContentLoaded', function () {
    // ruleSettingsModal save settings event
    $('#createNewGateRule').click(function (event) {
        //Validate form and decline any scripts
        if (!validateNewRuleForm()) {
            event.preventDefault();
            event.stopPropagation();
            return
        }
        // Submit data
        if (!$('#uuidNew').val()) {
            createRequest();
        }
        else {
            updateRequest();
        }
        closeNewRuleForm();
    });

    // ruleSettingsModal change protocol event
    $("#protoNew").change(function () {
        switch ($(this).val()) {
            case 'tcp':
                $('#generateSubDomain').prop('disabled', false);
                $('#dstAddrNew').prop('disabled', false);
                if (!$('#dstPortNew').val()) $('#dstPortNew').val(80);
                break;
            case 'udp':
                $('#dstAddrNew').val(settings.MainDomain).prop('disabled', true);
                $('#generateSubDomain').prop('disabled', true);
                break;
        }
    });

    // ruleSettingsModal generate sub domain event
    $('#generateSubDomain').click(function () {
        $.ajax({
            url: "/generate_domain",
            method: "GET",
        }).done(function (res) {
            if (res.status) {
                $('#dstAddrNew').val(res.response);
            }
        });
    });

    // ruleSettingsModal save settings event
    $('#deleteGateRule').click(function (event) {
        // Submit data
        if (!$('#uuidNew').val()) { return }
        // Serialize the data
        var jsonData = JSON.stringify({ "uuid": $('#uuidNew').val() });
        $.ajax({
            type: "POST",
            url: "/delete",
            data: jsonData,
            contentType: "application/json",  // This is important to set for JSON payload
            dataType: "json",  // The type of data you're expecting to receive
        }).done(function (data) {
            if (data.status) {
                var $row = $('#' + data.response);
                $row.remove();
            } else {
                alert(data.error);
            }
        });
        closeNewRuleForm();
    });

    // ruleSettingsModal show event
    $('#ruleSettingsModal').on('show.bs.modal', function () {
        if (!$('#uuidNew').val()) {
            $('#ruleSettingsModalLabel').text("Add new rule");
        } else {
            $('#ruleSettingsModalLabel').text("Edit this rule");
        }
    });
    $('#ruleSettingsModal').on('hidden.bs.modal', closeNewRuleForm);

    // gateTable Load table with rules
    loadGateTable();
});
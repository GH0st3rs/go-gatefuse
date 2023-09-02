function createGateTableRow(item) {
    var $row = $('<tr>', {
        'data-bs-toggle': 'modal',
        'data-bs-target': '#exampleModal',
        'id': item.UUID
    });
    // Checkbox column
    var $checkbox = $('<input>', {
        'type': 'checkbox',
        'class': 'form-check-input'
    });
    if (item.Active) {
        $checkbox.prop('checked', true);
    }
    var $checkboxDiv = $('<div>', {
        'class': 'form-switch ai-c'
    }).append($checkbox);
    var $checkboxColumn = $('<th>', {
        'scope': 'row'
    }).append($checkboxDiv);

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
    // Remove validation
    $('#createNewRuleForm').removeClass("was-validated");

    var $modal = $('#ruleSettingsModal');
    // Clean the form after sending
    $modal[0].reset();
    // Hide this modal window
    $modal.modal('hide');
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
        if (!$('#uuidNew').val()) { createRequest(); }
        else { updateRequest(); }
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


});
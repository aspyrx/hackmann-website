window.onload = function() {
    animateBits();
};

document.getElementById('navbar-register-button').onclick = function() {
    window.location.href = '#scroll-register';
}

formSubmit = function(e) {
    e.preventDefault();

    var form = e.target;
    if (!form.checkValidity()) {
        return;
    }

    var formData = {}
    for (var id in form.elements) {
        var input = form.elements[id];
        if (input instanceof HTMLElement && input.name.length > 0) {
            formData[input.name] = input.value;
        }
    }

    var formName = e.target.id.substring(0, e.target.id.indexOf('Form'));
    statusClassList = document.getElementById(formName + '-form-status').classList;
    statusTitle = document.getElementById(formName + '-form-status-title');
    statusMessage = document.getElementById(formName + '-form-status-message');
    statusMessageSuccessText = 'Thanks for signing up for HackMANN 2015! Look out soon for an email from us for more details.';
    if (formName === 'sponsor') {
        statusMessageSuccessText = 'Thanks for your interest in sponsoring HackMANN 2015! Feel free to email us; otherwise, we will get in touch with you as soon as possible.';
    }
    jQuery.ajax(form.action, {
        method: 'POST',
        data: formData,
        success: function(data, textStatus, jqXHR) {
            statusClassList.remove('hidden');
            statusClassList.remove('alert-danger');
            statusClassList.add("alert-success");
            statusTitle.innerHTML = 'Success!';
            statusMessage.innerHTML = statusMessageSuccessText;
        },
        error: function(jqXHR, textStatus, errorThrown) {
            statusClassList.remove('hidden');
            statusClassList.remove('alert-success');
            statusClassList.add('alert-danger');
            statusTitle.innerHTML = 'Uh oh...';
            statusMessage.innerHTML = 'Something went wrong! Please try again or contact us via email.';
        }
    });
}

document.getElementById('registerForm').onsubmit = formSubmit;
document.getElementById('mentorForm').onsubmit = formSubmit;
document.getElementById('sponsorForm').onsubmit = formSubmit;

function animateBits() {
    const BITS_ID = 'bits';
    const BITS_TEXT = 'bits';
    const BITS_ANIMATE_INTERVAL = 50.0;
    const BITS_ANIMATE_DURATION = 1500.0;
    var bitsTextTicks = 0;
    var animateBITS_ID = window.setInterval(function() {
        var len = Math.floor(bitsTextTicks / (BITS_ANIMATE_DURATION / BITS_ANIMATE_INTERVAL) * 4);
        var text = BITS_TEXT;
        for (i = len; i < 4; i++) {
            text = text.substr(0, i) + (Math.random() > 0.5 ? '1' : '0') + text.substr(i + 1);
        }
        document.getElementById(BITS_ID).innerHTML = text;
        bitsTextTicks++;
    }, BITS_ANIMATE_INTERVAL);
    window.setTimeout(function() {
        window.clearInterval(animateBITS_ID);
        document.getElementById(BITS_ID).innerHTML = BITS_TEXT;
    }, BITS_ANIMATE_DURATION);
}


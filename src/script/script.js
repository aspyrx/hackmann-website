window.onload = function() {
    animateBits();
};

document.getElementById('navbar-register-button').onload = function(e) {
    window.location.href = '#scroll-register';
}

document.getElementById('registerForm').onsubmit = function(e) {
    e.preventDefault();

    var form = e.target;
    if (!form.checkValidity()) {
        console.log('Form is invalid');
        return;
    }

    var formData = {}
    for (var id in form.elements) {
        var input = form.elements[id];
        if (input instanceof HTMLElement && input.name.length > 0) {
            formData[input.name] = input.value;
        }
    }
    console.log(formData);

    jQuery.ajax(form.action, {
        method: 'POST',
        data: formData,
        success: function(data, textStatus, jqXHR) {
            console.log("Successfully submitted");
            console.log(data);
        },
        error: function(jqXHR, textStatus, errorThrown) {
            console.log(jqXHR, textStatus, errorThrown);
        }
    });
}

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


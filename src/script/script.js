window.onload = function() {
    animateBits();
};

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

function submitRegisterForm(form) {
    if (!form.checkValidity()) {
        return;
    }

    var req = new XMLHttpRequest();
    req.onload = function() {
        if (this.status != 200) {
            console.log('Error submitting registerForm');
        } else {
            console.log('Successfully submitted registerForm');
        }
    }
    req.open('post', form.action, true);
    req.send(new FormData(form));
}

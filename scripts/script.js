window.onload = function() {
    animateBits();
}

var bitsId = 'bits';
var bitsAnimateInterval = 50.0;
var bitsAnimateDuration = 1500.0;
function animateBits() {
    var animateBitsId = window.setInterval(changeBitsText, bitsAnimateInterval);
    window.setTimeout(function() {
        window.clearInterval(animateBitsId);
        document.getElementById(bitsId).innerHTML = 'bits';
    }, bitsAnimateDuration);
}

var bitsTextTicks = 0;
function changeBitsText() {
    var len = Math.floor(bitsTextTicks / (bitsAnimateDuration / bitsAnimateInterval) * 4);
    var text = 'bits';
    for (i = len; i < 4; i++) {
        text = text.substr(0, i) + (Math.random() > 0.5 ? '1' : '0') + text.substr(i + 1);
    }
    document.getElementById(bitsId).innerHTML = text;
    bitsTextTicks++;
}

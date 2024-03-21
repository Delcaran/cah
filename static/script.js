function showLang(element)
{
    var lang_divs = document.querySelectorAll('div.lang_sets');
    for (var i = 0; i < lang_divs.length; i++) {
        lang_divs[i].style.display = "none";
    }
    document.querySelector('#' + element.value).style.display = "block";
}

function onCheckBoxChange()
{
    var checked = document.querySelectorAll('input[type=checkbox]:checked').lenght;
    var min = document.getElementById('min_checked').value;
    document.getElementById("submit").disabled = checked < min;
}

window.onload = function () {
    var conn;
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function logMessage(msg) {
        var item = document.createElement("div");
        item.innerHTML = msg;
        appendLog(item);
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (document.getElementById("czar").checked) {
            // send czar response to players
            return false;
        } else {
            // send players selected cards as json to the czar
            var checkedCards = document.querySelectorAll('input[name=cards]:checked');
            if (checkedCards.length > 0) {
                let msg = {};
                msg.player_id = document.getElementById("player_id").value;
                msg.cards = {}
                for (var i = 0; i < checkedCards.length; i++) {
                    index = checkedCards[i].value;
                    msg.cards[index] = document.getElementById("lbl_" + index).innerText;
                }
                conn.send(JSON.stringify(msg));
                logMessage("<b>Wait for Czar selection.</b>")
            }
            return false;
        }
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            logMessage("<b>Connection closed.</b>");
        };

        conn.onmessage = function (evt) {
            if (document.getElementById("czar").checked) {
                // parse players selections
                msg = document.createElement("div");
                msg.innerText = evt.data;
                appendLog(msg);
                const obj = JSON.parse(evt.data);
                var chk = document.createElement("input");
                chk.required = true;
                chk.type = "radio";
                chk.name = "players";
                chk.id = "player_" + obj.player_id;
                chk.value = obj.player_id;
                document.getElementById("white_cards").appendChild(chk)
                for (let card_index in obj.cards) {
                    var lbl = document.createElement("label")
                    lbl.htmlFor = chk.id
                    lbl.innerText = obj.cards[card_index]
                    document.getElementById("white_cards").appendChild(lbl)
                }
            } else {
                // parse czar selections
            }
        };

    } else {
        logMessage("<b>Your browser does not support WebSockets.</b>");
    }
};
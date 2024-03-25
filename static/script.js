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
            var selectedPlayer = document.querySelectorAll('input[name=players]:checked');
            if (selectedPlayer.length == 1) {
                // only one winner
                let msg = {};
                msg.kind = "winner"
                msg.payload = selectedPlayer[0].payload
                // message is already stored in value
                conn.send(JSON.stringify(msg));
            }
            return false;
        } else {
            // hide submit button (will be replaced with another button when Czar decision is received)
            document.getElementById("submit").style.display = "none";
            // send players selected cards as json to the czar
            var selectedPlayer = document.querySelectorAll('input[name=cards]:checked');
            if (selectedPlayer.length > 0) {
                let msg = {};
                msg.kind = "submission"
                msg.payload = {}
                msg.payload.player_id = document.getElementById("player_id").value;
                msg.payload.cards = {}
                for (var i = 0; i < selectedPlayer.length; i++) {
                    index = selectedPlayer[i].value;
                    msg.payload.cards[index] = document.getElementById("lbl_" + index).innerText;
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
            logMessage(evt.data)
            const obj = JSON.parse(evt.data);     
            switch(obj.kind) {
                case 'submission':
                    // parse players selections
                    var chk = document.createElement("input");
                    chk.required = true;
                    chk.type = "radio";
                    chk.name = "players";
                    chk.id = "player_" + obj.payload.player_id;
                    chk.value = obj.payload.player_id;
                    chk.payload = obj.payload; // (ab)use custom DOM property
                    document.getElementById("white_cards").appendChild(chk)
                    for (let card_index in obj.payload.cards) {
                        var lbl = document.createElement("label")
                        lbl.id = "lbl_" + obj.payload.player_id + "_" + card_index
                        lbl.name = "lbl_"+ obj.payload.player_id
                        lbl.htmlFor = chk.id
                        lbl.innerText = obj.payload.cards[card_index]
                        document.getElementById("white_cards").appendChild(lbl)
                    }
                    break;
                case 'winner':
                    // parse czar selections
                    break;
                default:
                    logMessage("ERROR")
            }      
        };

    } else {
        logMessage("<b>Your browser does not support WebSockets.</b>");
    }
};
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

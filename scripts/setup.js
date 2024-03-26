function showLang(element)
{
    var lang_divs = document.querySelectorAll('div.lang_sets');
    for (var i = 0; i < lang_divs.length; i++) {
        lang_divs[i].style.display = "none";
    }
    document.querySelector('div#' + element.value).style.display = "block";
}

function onCheckBoxChange()
{
    var is_disabled = true
    var num_checked = document.querySelectorAll('input[name=sets]:checked').length;
    if (num_checked) {
        is_disabled = num_checked < 1
    }
    document.getElementById("submit").disabled = is_disabled;
}

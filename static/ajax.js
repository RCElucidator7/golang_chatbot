const form = $("#input-e");
const list = $("#list");

form.keypress(function(event){
    if(event.keyCode != 13){
        return;
    }

    event.preventDefault();

    const userInput = form.val();
    form.val("");
    
    list.append("<li class='list-group-item  text-left list-group-item-success'>"
    + "User : " + userInput + "</li>");

    const prompt = {"input" : userInput }
    $.get("/chatbot", prompt)
        .done(response => {
            const output = "<li class='list-group-item  text-right list-group-item-warning'>" + response + " : Frank! " + "</li>";
            $("html, body").scrollTop($("body").height());
            setTimeout(function(){list.append(output)}, 2000);
        }).fail(() => {
            const output = "<li class='list-group-item list-group-item-danger'>Unavailable.</li>";
            list.append(output);
        });
});
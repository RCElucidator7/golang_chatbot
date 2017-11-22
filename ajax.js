const form = $("#input");
const list = $("#list");

form.keypress(function(event){
    if(event.keyCode != 13){
        return;
    }

    event.preventDefault();

    const userInput = form.val();
    form.val(" ");
    
    list.append("<li class='list-group-item  text-left list-group-item-success' id='leftList'>"
    + "User : " + userInput + "</li>");

    const queryParams = {"input" : userInput }
    $.get("/chat", queryParams)
    
        .done(function(resp){
            setTimeout(function(){   
                const newItem = userInput         
                list.append(newItem)
            }, 2000);
        }).fail(function(){
            const newItem = "<li class='list-group-item list-group-item-danger' >Sorry I'm not home right now.</li>";
            list.append(newItem);
        });
});
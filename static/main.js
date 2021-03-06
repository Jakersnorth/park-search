$(function(){
    $.get("/ping", function(data){
        if(data.error == "true"){
            $("#results").prepend("<div class='alert alert-danger'><strong>Error!</strong> "+ data.message +"</div>");
        }
    }, "json")
    
    $.get("/query1", function(data){
        $("#firstQuery").append(data);
    }, "html")

    $.get("/query2", function(data){
        $("#secondQuery").append(data);
    }, "html")
/*
    $.get("/query3", function(data){
        $("#thirdQuery").append(data);
    }, "html")
*/
    $("#submit").click(function(){
      $.post("/submit", {description: $("#description").val(), rating: $("#rating").val(), warnings: $("#warnings").val()})
        .done(function(data){
            $("#result").text("Review submitted");
        });
    });

    $("#search").click(function(){
      $("#currRes").remove();
      $.post("/search", {searchTerm: $("#searchTerm").val()})
        .done(function(data){
            $("#searchResults").append(data);
        });
    });

})

"use strict";


// Notification Setup
$(document).ready(function() {
    $(".toast").toast('show')
})


window.onload = function() {
    console.info("Hello")

    
    $("#studentAddSaveButton").click(function() {
        console.info("Button: Student Add Save")

        // TODO Add some sort of input validation

        var data = {}
        data["firstName"] = $("#studentAddFirstNameInput").val()
        data["lastName"] = $("#studentAddLastNameInput").val()
        data["graduationYear"] = $("#studentAddGradutationYearInput").val()
        data["id"] = $("#studentAddIdInput").val()
        data["schoolId"] = $("#studentAddSchoolIdInput").val()

        var request = JSON.stringify(data)

        console.info("Request: " + request)
        
        var url = "/api/students"
        var postReq = new XMLHttpRequest(); 
        postReq.open("POST", url)
        postReq.setRequestHeader('Content-Type', 'application/json');
        postReq.send(request);

        postReq.onreadystatechange = function() {
            console.log("POST Readystate: " + postReq.readyState)
            if(postReq.readyState == 4) {
                if (postReq.status == 200) {
                    // Success
                    console.info("Requst succedded")
                    $.toast({
                        title: 'Student Add',
                        // subtitle: '11 mins ago',
                        content: 'Request Success',
                        type: 'success',
                        delay: 1000
                    })
                } else {
                    // Failed
                    console.error("Request failed")
                    $.toast({
                        title: 'Student Add',
                        // subtitle: '11 mins ago',
                        content: 'Request Failed',
                        type: 'success',
                        delay: 1000
                    })
                }
            }
        }

    })
}
---
category: AlarmManager
apiurl: '/api/v1/fault'
title: "Create Fault"
type: 'POST'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required
* Params in request body:
    * title: Title of fault. The field can not be empty.
    * note: Note of fault.
    * owner: The people who is responsible for the fault. The field can not be empty.
    * tags: Be used to count how many faults a service has.
* Params in response body:
    * Code: The statud code of response. 201 is returned if successful, 400 or 500 is returned if failed.
    * Data: The info of fault in detail. 
        * Creator: The people who creates the fault will follow the fault by default.
        * State: The processing progress of fault.
        * Events: Event in the fault.
        * Followers: The people who follows the fault.
        * Comments: Be used to communicate by people who cares about the fault.
    * Message: Operational prompt information.

### Request
```{ "title":"mipush service down", "note":"test", "owner":"Bob", "tags":["miphone","miui"] }```

### Response

```Status:200```
```{
    "Code": 200,
    "Data": {
        "Id": 93,
        "CreatedAt": "2018-11-12T17:52:56+08:00",
        "Title": "mipush service down",
        "Note": "test",
        "Creator": "root",
        "Owner": "Bob",
        "State": "PROCESSING",
        "Tags": [
            "miphone",
            "miui"
        ],
        "Events": [],
        "Followers": [
            "root"
        ],
        "Comments": []
    },
    "Message": "fault create succeed"
}```


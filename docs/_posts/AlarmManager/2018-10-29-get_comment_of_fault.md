---
category: AlarmManager
apiurl: '/api/v1/fault/:id/comment'
title: "Get comment of fault"
type: 'GET'
sample_doc: 'alarm_manager.html'
layout: default
---

* [Session](#/authentication) Required.

### Response

```Status:200```
```{
    "Code": 200,
    "Data": [
        {
            "CreatedAt": "2018-11-13T14:14:42+08:00",
            "Creator": "root",
            "Comment": "\"test comment\""
        }
    ],
    "Message": "Get comment of fault successfully"
}```


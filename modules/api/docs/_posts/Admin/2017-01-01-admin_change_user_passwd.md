---
category: Admin
apiurl: '/api/v1/admin/change_user_passwd'
title: "Change User's Password"
type: 'PUT'
sample_doc: 'admin.html'
layout: default
---

* [Session](#/authentication) Required
* `Admin` usage

### Request
```{"user_id": 14, "password": "newpasswd"}```

### Response

```Status: 200```
```{"message":"password updated!"}```

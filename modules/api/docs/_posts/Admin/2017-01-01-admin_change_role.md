---
category: Admin
apiurl: '/api/v1/admin/change_user_role'
title: "Change User's role"
type: 'PUT'
sample_doc: 'admin.html'
layout: default
---

* [Session](#/authentication) Required
* `Admin` usage
* admin:
  * accept option:
    * yes
    * no

### Request
```{"user_id": 14, "admin": "yes"}```

### Response

```Status: 200```
```{"message":"user role update sccuessful, affect row: 1"}```

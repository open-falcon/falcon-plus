---
category: Admin
apiurl: '/api/v1/admin/change_user_profile'
title: "Change User's Profile"
type: 'PUT'
sample_doc: 'admin.html'
layout: default
---

* [Session](#/authentication) Required
* `Admin` usage

### Request
```
{
"user_id": 14,
"cnname": "翱鶚Test",
"email": "root123@cepave.com",
"im": "44955834958",
"phone": "99999999999",
"qq": "904394234239"
}
```
### Response

```Status: 200```
```{"message":"user profile updated!"}```

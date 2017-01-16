---
category: User
apiurl: '/api/v1/user/cgpasswd'
title: 'Change Password'
type: 'PUT'
sample_doc: 'user.html'
layout: default
---

* [Session](#/authentication) Required

### Request
```{
  "new_password": "test1",
  "old_password": "test1"
}```

### Response

```Status: 200```
```{"message":"password updated!"}```

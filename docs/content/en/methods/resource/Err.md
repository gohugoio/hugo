---
title: Err
description: Applicable to resources returned by the resources.GetRemote function, returns an error message if the HTTP request fails, else nil.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: resource.resourceError
    signatures: [RESOURCE.Err]
expiryDate: 2027-01-16 # deprecated 2025-01-16 in v0.141.0
---

{{< deprecated-in 0.141.0 >}}
Use the `try` statement instead. See [example].

[example]: /functions/go-template/try/#example
{{< /deprecated-in >}}

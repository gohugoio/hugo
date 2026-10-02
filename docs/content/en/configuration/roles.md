---
title: Configure roles
linkTitle: Roles
description: Configure roles.
categories: []
keywords: []
---

{{< new-in 0.153.0 />}}

This is the default configuration:

{{< code-toggle config=roles />}}

## Base settings

Configure the following base settings:

{{< code-toggle file=hugo >}}
defaultContentRole = 'guest'
defaultContentRoleInSubdir = false
{{< /code-toggle >}}

`defaultContentRole`
: (`string`) The project's [default role](g). When one or more roles are defined, this value must match one of the defined role names.

`defaultContentRoleInSubdir`
: (`bool`) Whether to publish the default content role to a subdirectory matching the [`defaultContentRole`](#defaultcontentrole). Default is `false`.

## Role settings

Use the following setting to define how Hugo orders roles.

`weight`
: (`int`) The role [weight](g). When set to a non-zero value, this is the primary sort criteria for this role.

## Sort order

Hugo sorts roles by weight in ascending order, then lexicographically in ascending order. This affects build order and complement selection.

## Example

The following configuration demonstrates how to define multiple roles with specific weights.

{{< code-toggle >}}
[roles.guest]
weight = 20
[roles.member]
weight = 10
{{< /code-toggle >}}

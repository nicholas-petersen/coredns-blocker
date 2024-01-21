# Virtualhost

## Name

*blocker* - Redirect banned hosts to 0.0.0.0

## Description

*blocker* redirect hosts ip addresses to 0.0.0.0 instead of resolve the related ip address.

## Syntax

```
blocker hosts
```

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metrics are exported:
* `coredns_blocker_{domain} - Count of the blocked hosts`

## Examples

Hosts file example [Steven Black](https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts)

```
example.com {
    blocker hosts
}
```


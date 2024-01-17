# Virtualhost

## Name

*blocker* - 

## Description


## Syntax

```
blocker
```

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metrics are exported:
* `coredns_blocker_{domain} - Counter of blocked domains `

## Examples

```
example.com {
    blocker
}
```

